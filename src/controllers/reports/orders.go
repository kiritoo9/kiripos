package reports

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReportOrder(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	branchId := c.Query("branch_id")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var datas []models.Trx
	tx := configs.DB.Table("trx").
		Select("trx.id", "trx.code", "trx.total_qty", "trx.total_price", "trx.discount", "trx.discount_desc", "trx.grand_total", "trx.status", "trx.note", "trx.created_date", "customers.name AS customer_name", "users.fullname AS user_name", "branches.name AS branch_name").
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN branches ON branches.id = trx.branch_id").
		Order("trx.created_date DESC").
		Where("trx.deleted = ?", false)
	if branchId != "" {
		tx.Where("trx.branch_id = ?", branchId)
	}
	if startDate != "" && endDate != "" {
		tx.Where("LEFT(trx.created_date::TEXT, 10) BETWEEN ? AND ?", startDate, endDate)
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
			"code":          datas[i].Code,
			"total_qty":     datas[i].TotalQty,
			"total_price":   datas[i].TotalPrice,
			"discount":      datas[i].Discount,
			"discount_desc": datas[i].DiscountDesc,
			"grand_total":   datas[i].GrandTotal,
			"status":        datas[i].Status,
			"note":          datas[i].Note,
			"created_date":  datas[i].CreatedDate,
			"user_name":     datas[i].UserName,
			"customer_name": datas[i].CustomerName,
			"branch_name":   datas[i].BranchName,
		})
	}

	var count int64
	var totalPage float64 = 1
	txcount := configs.DB.Table("trx").Distinct("trx.id").Where("trx.deleted = ?", false).
		Joins("LEFT JOIN customers ON customers.id = trx.customer_id").
		Joins("LEFT JOIN users ON users.id = trx.user_id").
		Joins("LEFT JOIN branches ON branches.id = trx.branch_id")
	if branchId != "" {
		txcount.Where("trx.branch_id = ?", branchId)
	}
	if startDate != "" && endDate != "" {
		txcount.Where("LEFT(trx.created_date::TEXT, 10) BETWEEN ? AND ?", startDate, endDate)
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
