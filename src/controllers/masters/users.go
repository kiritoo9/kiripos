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
	"golang.org/x/crypto/bcrypt"
)

func UserList(c *gin.Context) {
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

	var datas []models.Users
	err := configs.DB.Unscoped().
		Limit(int(limit)).
		Offset(int(offset)).
		Order("fullname asc").
		Where("LOWER(fullname) LIKE ?", "%"+keywords+"%").
		Where("deleted = ?", false).
		Or("LOWER(email) LIKE ?", "%"+keywords+"%").
		Where("deleted = ?", false).
		Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	var output []models.Users_Output
	for i := 0; i < len(datas); i++ {
		output = append(output, models.Users_Output{
			Id:          datas[i].Id,
			Fullname:    datas[i].Fullname,
			Email:       datas[i].Email,
			IsActive:    datas[i].IsActive,
			CreatedDate: datas[i].CreatedDate,
		})
	}

	var count int64
	var totalPage float64
	configs.DB.Model(&models.Users{}).Distinct("id").
		Where("deleted = ?", false).
		Where("LOWER(fullname) LIKE ?", "%"+keywords+"%").
		Or("deleted = ?", false).
		Where("LOWER(email) LIKE ?", "%"+keywords+"%").
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

func UserDetail(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data models.Users
	res := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("id = ?", id).
		Find(&data)

	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request success",
		"data":    data,
	})
}

func UserInsert(c *gin.Context) {
	var body models.Users_Form

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// VALIDATE DATA
	res := configs.DB.Unscoped().Where("email = ?", body.Email).First(&models.Users{})
	if res.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is exists!",
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	body.Password = string(hash)

	user := models.Users{
		Id:          uuid.New(),
		Email:       body.Email,
		Password:    body.Password,
		Fullname:    body.Fullname,
		IsActive:    body.IsActive,
		CreatedDate: time.Now(),
	}
	result := configs.DB.Create(&user).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error(),
		})
		return
	}

	urole := models.UserRoles{
		Id:     uuid.New(),
		UserId: user.Id,
		RoleId: body.RoleId,
	}
	configs.DB.Create(&urole)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data inserted",
		"data":    user,
	})
}
