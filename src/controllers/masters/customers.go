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

func CustomerList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 1
	}
	offset := (page * limit) - limit

	var datas []models.Customers
	err := configs.DB.Unscoped().
		Limit(int(limit)).
		Offset(int(offset)).
		Order("name ASC").
		Where("deleted = ?", false).
		Where("LOWER(code) LIKE ?", "%"+keywords+"%").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.Customers{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LOWER(code) LIKE ?", "%"+keywords+"%").
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
		"data":       datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func CustomerDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data *models.Customers
	res := configs.DB.Unscoped().
		Where("id = ?", id).
		Where("deleted = ?", false).
		First(&data)
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

func CustomerInsert(c *gin.Context) {
	var body *models.Customers

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var lastData int64
	var code string = "CS001"
	configs.DB.Model(&models.Customers{}).
		Where("deleted = ?", false).
		Count(&lastData)

	if lastData > 0 {
		count := strconv.FormatInt(lastData+1, 10)
		if len(count) == 1 {
			code = "CS00" + count
		} else if len(count) == 2 {
			code = "CS0" + count
		} else {
			code = "CS" + count
		}
	}

	var data_inserted = map[string]interface{}{
		"id":           uuid.New(),
		"code":         code,
		"name":         body.Name,
		"email":        body.Email,
		"phone":        body.Phone,
		"address":      body.Address,
		"created_date": time.Now(),
	}
	err_data := configs.DB.Model(&models.Customers{}).Create(&data_inserted).Error
	if err_data != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_data.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Data inserted",
		"data_inserted": data_inserted,
	})
}

func CustomerUpdate(c *gin.Context) {
	var body *models.Customers

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data_updated = map[string]interface{}{
		"name":    body.Name,
		"email":   body.Email,
		"phone":   body.Phone,
		"address": body.Address,
	}
	err_data := configs.DB.Model(&models.Customers{}).
		Where("id = ?", body.Id).
		Updates(&data_updated).Error
	if err_data != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_data.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data updated",
		"data_updated": data_updated,
	})
}

func CustomerDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_delete := configs.DB.Model(&models.Customers{}).
		Where("id = ?", id).
		Update("deleted", true).Error
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
