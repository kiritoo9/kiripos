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

func _purchaseSearch(tx *gorm.DB, keywords string, startDate string, endDate string, supplierId string) *gorm.DB {
	fields := []string{"purchase_orders.no_purchase", "suppliers.name"}
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
		if supplierId != "" {
			tx.Where("purchase_orders.supplier_id = ?", supplierId)
		}
	}
	return tx
}

func PurchaseList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	supplierId := c.Query("supplier_id")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	var offset int64 = (page * limit) - limit
	var datas []models.PurchaseOrders
	tx := configs.DB.Limit(int(limit)).Offset(int(offset)).Order("created_date ASC").
		Select("purchase_orders.*", "suppliers.name AS supplier_name").
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id")
	tx = _purchaseSearch(tx, keywords, startDate, endDate, supplierId)
	err := tx.Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	txcount := configs.DB.Model(&models.PurchaseOrders{}).Distinct("purchase_orders.id").
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id")
	txcount = _purchaseSearch(txcount, keywords, startDate, endDate, supplierId)
	errcount := txcount.Count(&count).Error
	if errcount != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errcount.Error(),
		})
		return
	}

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data models.PurchaseOrders
	res := configs.DB.Select("purchase_orders.*", "suppliers.name AS supplier_name", "branches.name AS branch_name").
		Where("purchase_orders.deleted = ?", false).
		Where("purchase_orders.id = ?", id).
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
		Joins("LEFT JOIN branches ON branches.id = purchase_orders.branch_id").
		First(&data)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	var items []models.PurchaseOrderItems
	err_items := configs.DB.Select("purchase_order_items.*", "products.name AS product_name").
		Joins("LEFT JOIN products ON products.id = purchase_order_items.product_id").
		Where("purchase_order_items.purchase_order_id = ?", data.Id).Find(&items).Error
	if err_items != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_items.Error(),
		})
		return
	}

	type Response struct {
		Data  models.PurchaseOrders       `json:"data"`
		Items []models.PurchaseOrderItems `json:"items"`
	}
	var response Response = Response{
		Data:  data,
		Items: items,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    response,
	})
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
			totalPrice += body.Items[i].Qty * body.Items[i].Price
			details = append(details, models.PurchaseOrderItems{
				ProductId: body.Items[i].ProductId,
				Qty:       body.Items[i].Qty,
				Price:     body.Items[i].Price,
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
				return err
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
	var body models.PurchaseOrderForm

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var po models.PurchaseOrders
	res := configs.DB.Where("deleted = ?", false).Where("id = ?", body.Id).First(&po)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	if po.Status != "S1" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You can only update draft(S1) purchase_order",
		})
		return
	}

	var items []models.PurchaseOrderItems

	err_items := configs.DB.Select("purchase_order_items.*", "products.stock AS last_stock").
		Where("purchase_order_items.purchase_order_id = ?", po.Id).
		Joins("LEFT JOIN products ON products.id = purchase_order_items.product_id").Find(&items).Error
	if err_items != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_items.Error(),
		})
		return
	}

	var totalQty int = 0
	var totalPrice int = 0
	data_updated := map[string]interface{}{
		"supplier_id":   body.SupplierId,
		"purchase_date": body.PurchaseDate,
		"discount":      body.Discount,
		"status":        body.Status,
		"note":          body.Note,
		"total_qty":     0,
		"total_price":   0,
		"grand_total":   0,
	}
	err_update := configs.DB.Transaction(func(tx *gorm.DB) error {
		if strings.ToUpper(body.Status) == "S2" {
			for i := range items {
				var new_stock = items[i].LastStock - items[i].Qty
				err := configs.DB.Model(&models.Products{}).Where("id = ?", items[i].ProductId).Update("stock", new_stock).Error
				if err != nil {
					return err
				}
			}
		}

		err := configs.DB.Exec("DELETE FROM purchase_order_items WHERE purchase_order_id = '" + po.Id.String() + "'").Error
		if err != nil {
			return err
		}

		for i := range body.Items {
			detail := map[string]interface{}{
				"id":                uuid.New(),
				"purchase_order_id": po.Id,
				"product_id":        body.Items[i].ProductId,
				"qty":               body.Items[i].Qty,
				"price":             body.Items[i].Price,
			}
			err := configs.DB.Model(&models.PurchaseOrderItems{}).Create(&detail).Error
			if err != nil {
				return err
			}

			if strings.ToUpper(body.Status) == "S2" {
				var product models.Products
				err := configs.DB.Select("stock", "with_stock").Where("id = ?", body.Items[i].ProductId).First(&product).Error
				if err != nil {
					return err
				}
				if product.WithStock {
					errprod := configs.DB.Model(&models.Products{}).
						Where("id = ?", body.Items[i].ProductId).
						Update("stock", product.Stock+body.Items[i].Qty).Error
					if errprod != nil {
						return errprod
					}
					totalQty += body.Items[i].Qty
					totalPrice += body.Items[i].Price * body.Items[i].Qty
				}
			}
		}

		data_updated["total_qty"] = totalQty
		data_updated["total_price"] = totalPrice
		data_updated["grand_total"] = totalPrice - body.Discount
		errupdate := configs.DB.Model(&models.PurchaseOrders{}).Where("id = ?", po.Id).Updates(&data_updated).Error
		if errupdate != nil {
			return errupdate
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
		"message":      "Data updated",
		"data_updated": data_updated,
	})
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
		err := configs.DB.Model(&models.PurchaseOrders{}).Where("id = ?", id).Update("deleted", true).Error
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
				var new_stock int = product.Stock - items[i].Qty
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
