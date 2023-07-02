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

func SupplierList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	var offset int64 = (page * limit) - limit

	var datas []models.Suppliers
	err := configs.DB.
		Limit(int(limit)).
		Offset(int(offset)).
		Order("name ASC").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(address) LIKE ?", "%"+keywords+"%").
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.Suppliers{}).
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(address) LIKE ?", "%"+keywords+"%").
		Count(&count)
	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"datas":      datas,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func SupplierDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data models.Suppliers
	res := configs.DB.Where("deleted = ?", false).Where("id = ?", id).First(&data)
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

func SupplierInsert(c *gin.Context) {
	var body models.Suppliers
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var code string = "S001"
	var count int64 = 0
	configs.DB.Model(&models.Suppliers{}).Distinct("id").Where("deleted = ?", false).Count(&count)
	if count > 0 {
		count += 1
		var countstr string = strconv.FormatInt(count, 10)
		if len(countstr) == 1 {
			code = "S00" + countstr
		} else if len(countstr) == 2 {
			code = "S0" + countstr
		} else {
			code = "S" + countstr
		}
	}

	data_inserted := models.Suppliers{
		Id:          uuid.New(),
		Code:        code,
		Name:        body.Name,
		Email:       body.Email,
		Address:     body.Address,
		Phone:       body.Phone,
		CreatedDate: time.Now(),
	}

	err_insert := configs.DB.Create(&data_inserted).Error
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

func SupplierUpdate(c *gin.Context) {
	var body models.Suppliers
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data_updated := models.Suppliers{
		Id:      body.Id,
		Name:    body.Name,
		Address: body.Address,
		Email:   body.Email,
		Phone:   body.Phone,
	}

	err_update := configs.DB.Where("id = ?", body.Id).Updates(&data_updated).Error
	if err_update != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_update.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data updated",
		"data_updated": data_updated,
	})
}

func SupplierDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_delete := configs.DB.Model(&models.Suppliers{}).Where("id = ?", id).Update("deleted", true).Error
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
