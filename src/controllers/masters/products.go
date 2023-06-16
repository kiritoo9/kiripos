package masters

import (
	"encoding/json"
	"kiripos/src/configs"
	"kiripos/src/helpers"
	"kiripos/src/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ProductList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 0, 0)
	limit, _ := strconv.ParseInt(c.Query("limit"), 0, 0)
	keywords := strings.ToLower(c.Query("keywords"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		page = 1
	}
	var offset = (page * limit) - limit

	var datas []models.Products
	err := configs.DB.Unscoped().
		Limit(int(limit)).
		Offset(int(offset)).
		Order("name ASC").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(code) LIKE ?", "%"+keywords+"%").
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := make([]map[string]interface{}, len(datas))

	for i := 0; i < len(datas); i++ {
		var imgs []string
		json.Unmarshal([]byte(datas[i].Images), &imgs)
		for j := 0; j < len(imgs); j++ {
			imgs[j] = helpers.GettRealPath(c, "products/"+imgs[j])
		}

		output = append(output, map[string]interface{}{
			"id":           datas[i].Id,
			"code":         datas[i].Code,
			"name":         datas[i].Name,
			"description":  datas[i].Description,
			"is_active":    datas[i].IsActive,
			"images":       imgs,
			"created_date": datas[i].CreatedDate,
		})
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.Products{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LOWER(name) = ?", "%"+keywords+"%").
		Where("deleted = ?", false).
		Where("LOWER(code) = ?", "%"+keywords+"%").
		Count(&count)

	if count > 0 && limit > 0 {
		var x float64 = float64(count)
		var y float64 = float64(limit)
		totalPage = math.Ceil(x / y)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Request Success",
		"data":       output,
		"pageActive": page,
		"totalPage":  totalPage,
	})
}

func ProductDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data *models.Products
	err_data := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("id = ?", id).
		First(&data).Error
	if err_data != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_data.Error(),
		})
		return
	}

	var imgs []string
	json.Unmarshal([]byte(data.Images), &imgs)
	for i := 0; i < len(imgs); i++ {
		imgs[i] = helpers.GettRealPath(c, "products/"+imgs[i])
	}

	output := map[string]interface{}{
		"id":           data.Id,
		"code":         data.Code,
		"name":         data.Name,
		"description":  data.Description,
		"is_active":    data.IsActive,
		"images":       imgs,
		"created_date": data.CreatedDate,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    output,
	})
}

func ProductInsert(c *gin.Context) {
	var body *models.Products_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("LOWER(code) = ?", strings.ToLower(body.Code)).
		First(&models.Products{})
	if res.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product code is already exists",
		})
		return
	}

	body.Id = uuid.New()
	if body.Images != "" {
		body.Images = helpers.GenerateImage("products", body.Images, body.Id.String()+"-0")
	}

	var arr_imgs []string
	arr_imgs = append(arr_imgs, body.Images)

	images, _ := json.Marshal(arr_imgs)
	product := models.Products{
		Id:          body.Id,
		Code:        strings.ToUpper(body.Code),
		Name:        body.Name,
		Description: body.Description,
		IsActive:    body.IsActive,
		Images:      string(images),
		CreatedDate: time.Now(),
	}

	err := configs.DB.Model(&models.Products{}).Create(&product).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Data inserted",
		"data_inserted": product,
	})
}

func ProductUpdate(c *gin.Context) {
	var body *models.Products_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("id = ?", body.Id).
		Where("LOWER(code) = ?", strings.ToLower(body.Code)).
		First(&models.Products{})
	if res.RowsAffected <= 0 {
		res := configs.DB.Unscoped().
			Where("deleted = ?", false).
			Where("LOWER(code) = ?", strings.ToLower(body.Code)).
			First(&models.Products{})
		if res.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product code is already exists",
			})
			return
		}
	}

	data_update := map[string]interface{}{
		"code":        body.Code,
		"name":        body.Name,
		"description": body.Description,
		"is_active":   body.IsActive,
	}

	if body.Images != "" {
		var filename string = body.Id.String() + "-0"
		helpers.RemoveFile("products/" + body.Id.String() + "-0.png")
		body.Images = helpers.GenerateImage("products", body.Images, filename)

		var arr_imgs []string
		arr_imgs = append(arr_imgs, body.Images)
		images, _ := json.Marshal(arr_imgs)
		data_update["images"] = string(images)
	}

	err_update := configs.DB.Model(&models.Products{}).
		Where("id = ?", body.Id).
		Updates(data_update).Error
	if err_update != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_update.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data updated",
		"data_updated": data_update,
	})
}

func ProductDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_delete := configs.DB.Model(&models.Products{}).
		Where("id = ?", id).
		Update("deleted", true).Error
	if err_delete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data deleted",
	})
}
