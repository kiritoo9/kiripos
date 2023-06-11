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

func RoleList(c *gin.Context) {

	/**
	DECLARE QUERY PARAMS
	@var page int
	@var limit int
	@var offset int
	@var keywords string
	*/

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

	/**
	QUERY USING GORM
	@var err *gorm.Error
	*/

	var roles []models.Roles
	err := configs.DB.Unscoped().
		Limit(int(limit)).
		Offset(int(offset)).
		Order("name asc").
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Find(&roles, "deleted = ?", false).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64
	configs.DB.Model(&models.Roles{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Count(&count)

	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       roles,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}
