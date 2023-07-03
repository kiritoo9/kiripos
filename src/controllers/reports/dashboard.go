package reports

import (
	"kiripos/src/configs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
	var today time.Time = time.Now()
	t2 := today.Format("2006-01-2")
	var totalOrders int64 = 0
	var totalPurchases int64 = 0
	var totalProducts int64 = 0
	var totalSuppliers int64 = 0

	configs.DB.Table("trx").Where("deleted = ?", false).Where("LEFT(created_date::TEXT, 10) = ?", t2).Count(&totalOrders)
	configs.DB.Table("purchase_orders").Where("deleted = ?", false).Where("LEFT(created_date::TEXT, 10) = ?", t2).Count(&totalPurchases)
	configs.DB.Table("products").Where("deleted = ?", false).Count(&totalProducts)
	configs.DB.Table("suppliers").Where("deleted = ?", false).Count(&totalSuppliers)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Request Success",
		"totalOrders":    totalOrders,
		"totalPurchases": totalPurchases,
		"totalSuppliers": totalSuppliers,
		"totalProducts":  totalProducts,
	})

}
