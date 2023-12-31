package orders

import (
	"encoding/json"
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

func _handleSearch(tx *gorm.DB, keywords string, startDate string, endDate string, branchId string) *gorm.DB {
	searchFields := [3]string{"trx.code", "customers.name", "users.fullname"}
	for i := range searchFields {
		if i == 0 {
			tx.Where("trx.deleted = ?", false)
		} else {
			tx.Or("trx.deleted = ?", false)
		}
		tx.Where("LOWER("+searchFields[i]+") LIKE ?", "%"+keywords+"%")
		if startDate != "" && endDate != "" {
			tx.Where("trx.created_date::DATE BETWEEN ? AND ?", startDate, endDate)
		}

		branchId, errBranch := uuid.Parse(branchId)
		if errBranch == nil {
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
	branchId := c.Query("branch_id")
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
		Select("trx.id", "trx.user_id", "trx.customer_id", "trx.branch_id", "trx.code", "trx.total_qty", "trx.total_price", "trx.discount", "trx.discount_desc", "trx.status", "trx.created_date", "trx.note", "users.fullname AS user_name", "customers.name AS customer_name", "branches.name AS branch_name").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id").
		Joins("LEFT JOIN branches ON branches.id = trx.branch_id").
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data models.Trx
	res := configs.DB.
		Table("trx").
		Select("trx.id", "trx.user_id", "trx.customer_id", "trx.branch_id", "trx.code", "trx.total_qty", "trx.total_price", "trx.discount", "trx.discount_desc", "trx.status", "trx.created_date", "trx.note", "users.fullname AS user_name", "customers.name AS customer_name", "branches.name AS branch_name").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id").
		Joins("LEFT JOIN branches ON branches.id = trx.branch_id").
		Where("trx.deleted = ?", false).
		Where("trx.id = ?", id).
		First(&data)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	var orderDetail []models.TrxDetails
	configs.DB.Table("trx_items").
		Select("trx_items.*", "products.name AS product_name", "products.images AS product_images").
		Joins("LEFT JOIN products ON products.id = trx_items.product_id").
		Where("trx_items.trx_id = ?", data.Id).Find(&orderDetail)

	type details struct {
		Id            uuid.UUID `json:"id"`
		ProductId     uuid.UUID `json:"product_id"`
		Qty           int       `json:"qty"`
		Price         int       `json:"price"`
		ProductName   string    `json:"product_name"`
		ProductImages []string  `json:"product_images"`
	}
	var dataDetails []details
	for i := range orderDetail {
		var images []string
		json.Unmarshal([]byte(orderDetail[i].ProductImages), &images)
		for i := 0; i < len(images); i++ {
			images[i] = helpers.GettRealPath(c, "products/"+images[i])
		}

		d := details{
			Id:            orderDetail[i].Id,
			ProductId:     orderDetail[i].ProductId,
			Qty:           orderDetail[i].Qty,
			Price:         orderDetail[i].Price,
			ProductName:   orderDetail[i].ProductName,
			ProductImages: images,
		}
		dataDetails = append(dataDetails, d)
	}

	response := map[string]interface{}{
		"data":   data,
		"detail": dataDetails,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    response,
	})
}

func OrderCreate(c *gin.Context) {
	var body models.TrxForm
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Customer name cannot be empty",
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
				if validProducts[i]["with_stock"] == true {
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
	var body models.TrxForm
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if strings.ToUpper(body.Status) == "S2" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to update a paid order",
		})
		return
	}

	var order models.Trx
	var orderDetail []models.TrxDetails

	res := configs.DB.Table("trx").
		Where("deleted = ?", false).
		Where("id = ?", body.Id).First(&order)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	err_detail := configs.DB.Table("trx_items").
		Select("trx_items.*", "products.name AS product_name", "products.with_stock AS product_with_stock", "products.images AS product_images", "products.stock AS product_stock", "products.price AS product_price").
		Joins("LEFT JOIN products ON products.id = trx_items.product_id").
		Where("trx_items.trx_id = ?", order.Id).
		Find(&orderDetail).Error
	if err_detail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_detail.Error(),
		})
		return
	}

	/*
		Update items steps
		1. update product stock by last trx_items
		2. delete trx_items
		3. re-insert trx_items
	*/
	var data_updated map[string]interface{}
	err_update := configs.DB.Transaction(func(tx *gorm.DB) error {
		var totalQty int = 0
		var totalPrice int = 0
		for i := range orderDetail {
			if orderDetail[i].ProductWithStock {
				var new_stock int = orderDetail[i].ProductStock + orderDetail[i].Qty
				err_product := configs.DB.Model(&models.Products{}).
					Where("id = ?", orderDetail[i].ProductId).
					Update("stock", new_stock).Error
				if err_product != nil {
					return err_product
				}
			}
		}

		err_del_items := configs.DB.Exec("DELETE FROM trx_items WHERE trx_id = '" + order.Id.String() + "'").Error
		if err_del_items != nil {
			return err_del_items
		}

		for i := range body.Items {
			var product models.Products
			resprod := configs.DB.Select("id", "price", "with_stock", "stock").
				Where("deleted = ?", false).
				Where("id = ?", body.Items[i].ProductId).
				First(&product)
			if resprod.RowsAffected > 0 {
				_items := map[string]interface{}{
					"id":         uuid.New(),
					"trx_id":     order.Id,
					"product_id": product.Id,
					"qty":        body.Items[i].Qty,
					"price":      product.Price,
				}
				err_items := configs.DB.Table("trx_items").Create(&_items).Error
				if err_items != nil {
					return err_items
				} else {
					totalPrice += int(product.Price) * body.Items[i].Qty
					totalQty += body.Items[i].Qty
					if product.WithStock {
						var _stock int = product.Stock - body.Items[i].Qty
						err_update_product := configs.DB.Model(&models.Products{}).Where("id = ?", product.Id).Update("stock", _stock).Error
						if err_update_product != nil {
							return err_update_product
						}
					}
				}
			}
		}

		data_updated = map[string]interface{}{
			"discount":      body.Discount,
			"discount_desc": body.DiscountDesc,
			"status":        body.Status,
			"note":          body.Note,
			"total_qty":     totalQty,
			"total_price":   totalPrice,
			"grand_total":   totalPrice - body.Discount,
		}

		err_update := configs.DB.Table("trx").Where("id = ?", body.Id).Updates(&data_updated).Error
		if err_update != nil {
			return err_update
		}

		return nil
	})

	if err_update != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_update.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "data_updated",
		"data_updated": data_updated,
	})
}

func OrderDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var order models.Trx
	res_order := configs.DB.Table("trx").Where("deleted = ?", false).Where("id = ?", id).First(&order)
	if res_order.RowsAffected <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order is not found",
		})
		return
	}

	var orderDetail []models.TrxDetails
	err_detail := configs.DB.Table("trx_items").
		Select("trx_items.*", "products.name AS product_name").
		Joins("LEFT JOIN products ON products.id = trx_items.product_id").
		Where("trx_items.trx_id = ?", order.Id).Find(&orderDetail).Error
	if err_detail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_detail.Error(),
		})
		return
	}

	err_delete := configs.DB.Transaction(func(tx *gorm.DB) error {
		if err := configs.DB.Table("trx").Where("id = ?", id).Update("deleted", true).Error; err != nil {
			return err
		}

		if strings.ToUpper(order.Status) == "S1" {
			for i := range orderDetail {
				var product models.Products
				if err_product := configs.DB.Select("with_stock", "stock").Where("id = ?", orderDetail[i].ProductId).First(&product).Error; err_product != nil {
					return err_product
				}
				if product.WithStock {
					var new_stock int = product.Stock + orderDetail[i].Qty
					if err_update := configs.DB.Model(&models.Products{}).Where("id = ?", orderDetail[i].ProductId).Update("stock", new_stock).Error; err_update != nil {
						return err_update
					}
				}
			}
		}

		return nil
	})
	if err_delete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_delete.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data deleted",
	})
}
