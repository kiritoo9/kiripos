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

	var datas []models.Branches
	err := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(location) LIKE ?", "%"+keywords+"%").
		Order("name ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&datas).Error
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
		"data":       datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func BranchDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data *models.Branches
	res_branch := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("id = ?", id).
		First(&data)
	if res_branch.RowsAffected <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data is not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    data,
	})
}

func BranchInsert(c *gin.Context) {
	var body *models.Branches_Form

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.IsMain {
		configs.DB.Model(&models.Branches{}).Update("is_main", false)
	}

	branch := models.Branches{
		Id:          uuid.New(),
		Code:        body.Code,
		Name:        body.Name,
		Location:    body.Location,
		Phone:       body.Phone,
		Email:       body.Email,
		IsMain:      body.IsMain,
		IsActive:    body.IsActive,
		CreatedDate: time.Now(),
	}
	err_branch := configs.DB.Create(&branch).Error
	if err_branch != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_branch.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Data inserted",
		"data_inserted": branch,
	})
}

func BranchUpdate(c *gin.Context) {
	var body *models.Branches_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.IsMain {
		configs.DB.Model(&models.Branches{}).Update("is_main", false)
	}

	branch_update := map[string]interface{}{
		"code":      body.Code,
		"name":      body.Name,
		"location":  body.Location,
		"phone":     body.Phone,
		"email":     body.Email,
		"is_main":   body.IsMain,
		"is_active": body.IsActive,
	}

	err_branch := configs.DB.Model(&models.Branches{}).
		Where("id = ?", body.Id).
		Updates(branch_update).Error
	if err_branch != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_branch.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data updated",
		"data_updated": branch_update,
	})
}

func BranchDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_branch := configs.DB.Model(&models.Branches{}).Where("id = ?", id).Update("deleted", true).Error
	if err_branch != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_branch.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data deleted",
	})
}
