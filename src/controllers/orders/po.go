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

func _purchaseSearch(tx *gorm.DB, keywords string, startDate string, endDate string, branchId string) *gorm.DB {
	fields := []string{"purchase_orders.no_purchase", "suppliers.name", "branches.name"}
	for i := range fields {
		if i == 0 {
			tx.Where("purchase_orders.deleted = ?", false)
		} else {
			tx.Or("purchase_orders.deleted = ?", false)
		}

		tx.Where("LOWER("+fields[i]+") LIKE ?", "%"+keywords+"%")
		if startDate != "" && endDate != "" {
			tx.Where("LEFT(purchase_orders.purchase_date::TEXT, 10) BETWEEN ? AND ?", startDate, endDate)
		}
		if branchId != "" {
			tx.Where("purchase_orders.branch_id = ?", branchId)
		}
	}
	return tx
}

func PurchaseList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	branchId := c.Query("branchId")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	var offset int64 = (page * limit) - limit
	var datas []models.PurchaseOrders
	tx := configs.DB.Limit(int(limit)).Offset(int(offset)).Order("created_date ASC").
		Select("purchase_orders.*", "suppliers.name AS supplier_name", "branches.name AS branch_name").
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
		Joins("LEFT JOIN branches ON branches.id = purchase_orders.branch_id")
	tx = _purchaseSearch(tx, keywords, startDate, endDate, branchId)
	err := tx.Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.PurchaseOrders{}).Distinct("id").Count(&count)
	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"datas":      datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func PurchaseDetail(c *gin.Context) {

}

func PurchaseCreate(c *gin.Context) {
	var body models.PurchaseOrderForm

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var totalQty int = 0
	var totalPrice int = 0
	var details []models.PurchaseOrderItems

	for i := range body.Items {
		totalQty += body.Items[i].Qty
		var product models.Products
		resprod := configs.DB.
			Select("stock", "price").
			Where("deleted = ?", false).
			Where("with_stock = ?", true).
			Where("id = ?", body.Items[i].ProductId).
			First(&product)
		if resprod.RowsAffected <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product with id " + body.Items[i].ProductId.String() + " is not exists or not allowed to do this action",
			})
			return
		} else {
			totalPrice += body.Items[i].Qty * int(product.Price)
			details = append(details, models.PurchaseOrderItems{
				ProductId: body.Items[i].ProductId,
				Qty:       body.Items[i].Qty,
				Price:     int(product.Price),
				LastStock: product.Stock,
			})
		}
	}

	var no_purchase string = "PO001"
	today := time.Now()
	t2 := today.Format("2006-01")
	var countpo int64
	errcount := configs.DB.Model(&models.PurchaseOrders{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LEFT(created_date::TEXT, 7) = ?", t2).
		Count(&countpo).Error
	if errcount != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errcount.Error(),
		})
		return
	}
	countpostr := strconv.FormatInt((countpo + 1), 10)
	if len(countpostr) == 1 {
		no_purchase = "PO00" + countpostr
	} else if len(countpostr) == 2 {
		no_purchase = "PO0" + countpostr
	} else {
		no_purchase = "PO" + countpostr
	}

	token := helpers.GetToken(c)
	if token["branch_id"] == nil || token["branch_id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing branch_id OR user_id",
		})
		return
	}

	data_inserted := map[string]interface{}{
		"id":            uuid.New(),
		"user_id":       token["id"],
		"branch_id":     token["branch_id"],
		"supplier_id":   body.SupplierId,
		"no_purchase":   no_purchase,
		"purchase_date": body.PurchaseDate,
		"total_qty":     totalQty,
		"total_price":   totalPrice,
		"discount":      body.Discount,
		"grand_total":   totalPrice - body.Discount,
		"status":        body.Status,
		"note":          body.Note,
		"created_date":  time.Now(),
	}

	err_insert := configs.DB.Transaction(func(tx *gorm.DB) error {
		err := configs.DB.Model(&models.PurchaseOrders{}).Create(&data_inserted).Error
		if err != nil {
			return err
		}

		for i := range details {
			detail := map[string]interface{}{
				"id":                uuid.New(),
				"purchase_order_id": data_inserted["id"],
				"product_id":        details[i].ProductId,
				"qty":               details[i].Qty,
				"price":             details[i].Price,
			}
			err := configs.DB.Model(&models.PurchaseOrderItems{}).Create(&detail).Error
			if err != nil {
				return nil
			}

			if strings.ToUpper(fmt.Sprint(data_inserted["status"])) == "S2" {
				err := configs.DB.Model(&models.Products{}).Where("id = ?", details[i].ProductId).Update("stock", details[i].LastStock+details[i].Qty).Error
				if err != nil {
					return err
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

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Data inserted",
		"data_inserted": data_inserted,
	})
}

func PurchaseUpdate(c *gin.Context) {

}

func PurchaseDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var po models.PurchaseOrders
	res := configs.DB.Where("deleted = ?", false).Where("id = ?", id).First(&po)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	err_del := configs.DB.Transaction(func(tx *gorm.DB) error {
		err := configs.DB.Where("id = ?", id).Update("deleted", true).Error
		if err != nil {
			return err
		}

		if po.Status != "S1" {
			var items []models.PurchaseOrderItems
			err_items := configs.DB.Where("purchase_order_id", id).Find(&items).Error
			if err_items != nil {
				return err_items
			}

			for i := range items {
				var product models.Products
				err_prod := configs.DB.Where("id = ?", items[i].ProductId).First(&product).Error
				if err_prod != nil {
					return err_prod
				}
				var new_stock int = product.Stock + items[i].Qty
				err_update_product := configs.DB.Model(&models.Products{}).Where("id = ?", product.Id).Update("stock", new_stock).Error
				if err_update_product != nil {
					return err_update_product
				}
			}
		}

		return nil
	})

	if err_del != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_del.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data deleted",
	})
}
