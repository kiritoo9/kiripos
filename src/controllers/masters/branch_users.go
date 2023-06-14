package masters

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func _checkBranch(c *gin.Context) string {
	id, err := uuid.Parse(c.Param("branch_id"))
	if err != nil {
		return err.Error()
	}

	res := configs.DB.Unscoped().Where("deleted = ?", false).Where("id = ?", id).Find(&models.Branches{})
	if res.RowsAffected <= 0 {
		return "Branch is not found"
	}

	return ""
}

func BranchUserList(c *gin.Context) {
	if callback := _checkBranch(c); callback != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": callback,
		})
		return
	}

	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 1
	}
	var offset = (page * limit) - limit

	var datas []models.BranchUsers
	err := configs.DB.Unscoped().
		Select("branch_users.*", "users.fullname AS user_name", "branches.name AS branch_name").
		Joins("LEFT JOIN users ON users.id = branch_users.user_id").
		Joins("LEFT JOIN branches ON branches.id = branch_users.branch_id").
		Limit(int(limit)).
		Offset(int(offset)).
		Order("users.fullname ASC").
		Where("branch_users.deleted = ?", false).
		Where("users.deleted = ?", false).
		Where("LOWER(users.fullname) LIKE ?", "%"+keywords+"%").
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.BranchUsers{}).Distinct("branch_users.id").
		Joins("LEFT JOIN users ON users.id = branch_users.user_id").
		Joins("LEFT JOIN branches ON branches.id = branch_users.branch_id").
		Where("branch_users.deleted = ?", false).
		Where("users.deleted = ?", false).
		Where("LOWER(users.fullname) LIKE ?", "%"+keywords+"%").
		Count(&count)

	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func BranchUserDetail(c *gin.Context) {
	if callback := _checkBranch(c); callback != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": callback,
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data *models.BranchUsers
	res := configs.DB.Unscoped().
		Select("branch_users.*", "users.fullname AS user_name", "branches.name AS branch_name").
		Joins("LEFT JOIN users ON users.id = branch_users.user_id").
		Joins("LEFT JOIN branches ON branches.id = branch_users.branch_id").
		Where("branch_users.deleted = ?", false).
		Where("users.deleted = ?", false).
		Where("branch_users.id = ?", id).
		Find(&data)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    data,
	})
}

func BranchUserInsert(c *gin.Context) {
	if callback := _checkBranch(c); callback != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": callback,
		})
		return
	}

	var body *models.BranchUsers_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	branch_id, _ := uuid.Parse(c.Param("branch_id"))
	data_inserted := map[string]interface{}{
		"id":           uuid.New(),
		"user_id":      body.UserId,
		"branch_id":    branch_id,
		"created_date": time.Now(),
	}
	err_insert := configs.DB.Model(&models.BranchUsers{}).Create(&data_inserted).Error
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

func BranchUserDelete(c *gin.Context) {
	if callback := _checkBranch(c); callback != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": callback,
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_delete := configs.DB.Model(&models.BranchUsers{}).Where("id = ?", id).Update("deleted", true).Error
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
