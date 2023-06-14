package masters

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BranchList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 1
	}
	var offset = (limit * page) - limit

	var data []models.Branches
	err := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(location) LIKE ?", "%"+keywords+"%").
		Order("name ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&data).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.Branches{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(location) LIKE ?", "%"+keywords+"%").
		Count(&count)

	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       data,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func BranchDetail(c *gin.Context) {

}

func BranchInsert(c *gin.Context) {

}

func BranchUpdate(c *gin.Context) {

}

func BranchDelete(c *gin.Context) {

}
