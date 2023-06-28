package orders

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"
	"strings"

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

	c.JSON(http.StatusCreated, gin.H{
		"message":       "data_inserted",
		"data_inserted": body,
	})
}

func OrderUpdate(c *gin.Context) {

}

func OrderDelete(c *gin.Context) {

}
