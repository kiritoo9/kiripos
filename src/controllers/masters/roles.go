package masters

import (
	"kiripos/src/configs"
	"kiripos/src/models"
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
	@var err *gorm
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

	var totalPage = 0

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       roles,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}
