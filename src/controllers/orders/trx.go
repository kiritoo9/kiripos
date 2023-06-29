package orders

import (
	"fmt"
	"kiripos/src/configs"
	"kiripos/src/helpers"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func _handleSearch(tx *gorm.DB, keywords string, startDate string, endDate string, branchId uuid.UUID) *gorm.DB {
	searchFields := [3]string{"trx.code", "customers.name", "users.fullname"}
	for i := range searchFields {
		if i == 0 {
			tx.Where("trx.deleted = ?", false)
		} else {
			tx.Or("trx.deleted = ?", false)
		}
		tx.Where("LOWER("+searchFields[i]+") = ?", "%"+keywords+"%")

		if branchId.String() != "" {
			tx.Where("trx.branch_id = ?", branchId)
		}
	}
	return tx
}

func OrderList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	branchId, _ := uuid.Parse(c.Query("branch_id"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var offset int64 = (page * limit) - limit
	var datas []models.Trx

	tx := configs.DB.
		Table("trx").
		Select("trx.id", "trx.code", "trx.total_qty", "trx.total_price", "trx.discount", "trx.discount_desc", "trx.status", "trx.created_date", "users.fullname AS user_name", "customers.name AS customer_name").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id").
		Limit(int(limit)).
		Offset(int(offset))
	tx = _handleSearch(tx, keywords, startDate, endDate, branchId)
	err := tx.Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	trxCount := configs.DB.Model(&models.Trx{}).Distinct("trx.id").
		Table("trx").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id")
	trxCount = _handleSearch(trxCount, keywords, startDate, endDate, branchId)
	trxCount.Count(&count)
	if count > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Floor(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func OrderDetail(c *gin.Context) {

}

func OrderCreate(c *gin.Context) {
	var body models.TrxForm
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var errProducts []string
	var validProducts []map[string]interface{}
	var totalQty int = 0
	var totalPrice int = 0

	for i := range body.Items {
		var d models.Products
		resprod := configs.DB.
			Where("deleted = ?", false).
			Where("id = ?", body.Items[i].ProductId).
			First(&d)
		if resprod.RowsAffected <= 0 {
			errProducts = append(errProducts, "Product with id : "+body.Items[i].ProductId.String()+" is not exists")
		} else {
			var allowed bool = true
			if d.WithStock {
				if d.Stock <= 0 {
					allowed = false
					errProducts = append(errProducts, "Product with id : "+body.Items[i].ProductId.String()+" have no stock")
				}
			}

			if allowed {
				totalQty += body.Items[i].Qty
				totalPrice += int(d.Price) * body.Items[i].Qty
				validProducts = append(validProducts, map[string]interface{}{
					"id":             uuid.New(),
					"product_id":     d.Id,
					"qty":            body.Items[i].Qty,
					"price":          d.Price,
					"with_stock":     d.WithStock,
					"existing_stock": d.Stock,
				})
			}
		}
	}

	if len(errProducts) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errProducts,
		})
		return
	}

	today := time.Now()
	tdate := today.Format("2006-01-02")
	var countToday int64 = 0
	configs.DB.Table("trx").
		Where("deleted = ?", false).
		Where("LEFT(created_date::TEXT, 10) = ?", tdate).
		Count(&countToday)
	countToday += 1

	order := map[string]interface{}{
		"id":            uuid.New(),
		"user_id":       nil,
		"customer_id":   nil,
		"branch_id":     nil,
		"code":          strconv.Itoa(int(countToday)),
		"total_qty":     totalQty,
		"total_price":   totalPrice,
		"discount":      body.Discount,
		"discount_desc": body.DiscountDesc,
		"grand_total":   totalPrice - body.Discount,
		"status":        body.Status,
		"note":          body.Note,
		"created_date":  today,
	}

	var customer models.Customers
	rescust := configs.DB.Where("deleted = ?", false).
		Where("LOWER(name) = ?", strings.ToLower(body.CustomerName)).
		First(&customer)
	if rescust.RowsAffected > 0 {
		order["customer_id"] = customer.Id
	} else {
		new_cust := models.Customers{
			Id:          uuid.New(),
			Code:        helpers.GenerateCustomerCode(),
			Name:        body.CustomerName,
			CreatedDate: time.Now(),
		}
		if err := configs.DB.Create(&new_cust).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			order["customer_id"] = new_cust.Id
		}
	}

	token := helpers.GetToken(c)
	if token["branch_id"] == nil || token["branch_id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing branch_id",
		})
		return
	}
	order["branch_id"] = token["branch_id"]
	order["user_id"] = token["id"]

	// INSERT TO DB
	err_insert := configs.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("trx").Create(&order).Error; err != nil {
			return err
		}

		for i := range validProducts {
			d := map[string]interface{}{
				"id":         validProducts[i]["id"],
				"trx_id":     order["id"],
				"product_id": validProducts[i]["product_id"],
				"qty":        validProducts[i]["qty"],
				"price":      validProducts[i]["price"],
			}
			if err := tx.Table("trx_items").Create(&d).Error; err != nil {
				return err
			} else {
				if strings.ToUpper(fmt.Sprint(order["status"])) == "S2" && validProducts[i]["with_stock"] == true {
					existingStock, _ := strconv.ParseInt(fmt.Sprint(validProducts[i]["existing_stock"]), 0, 0)
					qty, _ := strconv.ParseInt(fmt.Sprint(d["qty"]), 0, 0)
					var updatedStock int64 = existingStock - qty
					updated_data := map[string]interface{}{
						"stock": updatedStock,
					}
					if err := tx.Table("products").Where("id = ?", d["product_id"]).Updates(&updated_data).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err_insert != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_insert.Error(),
		})
		return
	}

	order["items"] = validProducts
	c.JSON(http.StatusCreated, gin.H{
		"message":       "data_inserted",
		"data_inserted": order,
	})
}

func OrderUpdate(c *gin.Context) {

}

func OrderDelete(c *gin.Context) {

}
