package reports

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReportPurchase(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	supplierId := c.Query("supplier_id")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var datas []models.PurchaseOrders
	tx := configs.DB.
		Select("purchase_orders.id", "purchase_orders.no_purchase", "purchase_orders.total_qty", "purchase_orders.total_price", "purchase_orders.discount", "purchase_orders.grand_total", "purchase_orders.status", "purchase_orders.note", "purchase_orders.created_date", "suppliers.name AS supplier_name").
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
		Order("purchase_orders.created_date DESC").
		Where("purchase_orders.deleted = ?", false)
	if supplierId != "" {
		tx.Where("purchase_orders.supplier_id = ?", supplierId)
	}
	if startDate != "" && endDate != "" {
		tx.Where("LEFT(purchase_orders.created_date::TEXT, 10) BETWEEN ? AND ?", startDate, endDate)
	}
	err := tx.Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var responses []map[string]interface{}
	for i := range datas {
		responses = append(responses, map[string]interface{}{
			"id":            datas[i].Id,
			"no_purchase":   datas[i].NoPurchase,
			"total_qty":     datas[i].TotalQty,
			"total_price":   datas[i].TotalPrice,
			"discount":      datas[i].Discount,
			"grand_total":   datas[i].GrandTotal,
			"status":        datas[i].Status,
			"note":          datas[i].Note,
			"created_date":  datas[i].CreatedDate,
			"supplier_name": datas[i].SupplierName,
		})
	}

	var count int64
	var totalPage float64 = 1
	txcount := configs.DB.Model(&models.PurchaseOrders{}).
		Distinct("purchase_orders.id").
		Where("purchase_orders.deleted = ?", false).
		Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id")
	if supplierId != "" {
		tx.Where("purchase_orders.supplier_id = ?", supplierId)
	}
	if startDate != "" && endDate != "" {
		txcount.Where("LEFT(purchase_orders.created_date::TEXT, 10) BETWEEN ? AND ?", startDate, endDate)
	}
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
		"data":       responses,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}
