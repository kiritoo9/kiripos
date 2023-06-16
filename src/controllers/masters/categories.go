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

func CategoryList(c *gin.Context) {
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

	var datas []models.Categories
	err := configs.DB.Unscoped().
		Limit(int(limit)).
		Offset(int(offset)).
		Order("name ASC").
		Where("deleted = ?", false).
		Where("LOWER(name) LIKE ?", "%"+keywords+"%").
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := make([]map[string]interface{}, len(datas))

	for i := range datas {
		var arr_imgs []string
		json.Unmarshal([]byte(datas[i].Images), &arr_imgs)
		for j := range arr_imgs {
			arr_imgs[j] = helpers.GettRealPath(c, "categories/"+arr_imgs[j])
		}

		output = append(output, map[string]interface{}{
			"id":           datas[i].Id,
			"name":         datas[i].Name,
			"description":  datas[i].Description,
			"images":       arr_imgs,
			"created_date": datas[i].CreatedDate,
		})
	}

	var count int64
	var totalPage float64 = 1
	configs.DB.Model(&models.Categories{}).
		Distinct("id").
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
		"data":       output,
		"pageActive": page,
		"totalPage":  totalPage,
	})

}

func CategoryDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data *models.Categories
	res := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("id = ?", id).
		First(&data)
	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	var output map[string]interface{}
	if data.Images != "" {
		var arr_imgs []string
		json.Unmarshal([]byte(data.Images), &arr_imgs)
		for i := range arr_imgs {
			arr_imgs[i] = helpers.GettRealPath(c, "categories/"+arr_imgs[i])
		}
		output = map[string]interface{}{
			"id":           data.Id,
			"name":         data.Name,
			"description":  data.Description,
			"images":       arr_imgs,
			"created_date": data.CreatedDate,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
		"data":    output,
	})
}

func CategoryInsert(c *gin.Context) {
	var body *models.Categories_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	body.Id = uuid.New()
	if body.Images != "" {
		var arr_imgs []string
		filename := helpers.GenerateImage("categories", body.Images, body.Id.String()+"-0")
		arr_imgs = append(arr_imgs, filename)
		images, _ := json.Marshal(arr_imgs)
		body.Images = string(images)
	}

	data_inserted := models.Categories{
		Id:          body.Id,
		Name:        body.Name,
		Description: body.Description,
		Images:      body.Images,
		CreatedDate: time.Now(),
	}
	err_insert := configs.DB.Model(&models.Categories{}).Create(&data_inserted).Error
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

func CategoryUpdate(c *gin.Context) {
	var body *models.Categories_Form
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data_updated := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
	}

	if body.Images != "" {
		helpers.RemoveFile("categories/" + body.Id.String() + "-0.png")
		var arr_imgs []string
		filename := helpers.GenerateImage("categories", body.Images, body.Id.String()+"-0")
		arr_imgs = append(arr_imgs, filename)
		images, _ := json.Marshal(arr_imgs)
		data_updated["images"] = string(images)
	}

	err_update := configs.DB.Model(&models.Categories{}).Where("id = ?", body.Id).Updates(data_updated).Error
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

func CategoryDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_delete := configs.DB.Model(&models.Categories{}).Where("id = ?", id).Update("deleted", true).Error
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
