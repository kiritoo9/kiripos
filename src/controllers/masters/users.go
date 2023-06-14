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
	var totalPage float64 = 1
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
		Select("id", "fullname", "email", "is_active", "created_date").
		Where("deleted = ?", false).
		Where("id = ?", id).
		Find(&data)

	if res.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data is not found",
		})
		return
	}

	var role_name string
	var user_role models.UserRoles
	res_role := configs.DB.Unscoped().
		Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
		Select("roles.name AS role_name").
		Where("user_roles.deleted = ?", false).
		Where("user_roles.user_id = ?", data.Id).
		First(&user_role)
	if res_role.RowsAffected > 0 {
		role_name = user_role.RoleName
	}

	output := map[string]interface{}{
		"user":      data,
		"role_name": role_name,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request success",
		"data":    output,
	})
}

func UserInsert(c *gin.Context) {
	var body models.Users_Form

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// VALIDATE DATA
	res := configs.DB.Unscoped().
		Where("deleted = ?", false).
		Where("email = ?", body.Email).
		First(&models.Users{})
	if res.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is already registered!",
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
	err_user := configs.DB.Create(&user).Error
	if err_user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_user.Error(),
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
		"message":       "Data inserted",
		"data_inserted": user,
	})
}

func UserUpdate(c *gin.Context) {
	var body models.Users_Form

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res_user := configs.DB.Unscoped().
		Where("id = ?", body.Id).
		Where("LOWER(email) = ?", strings.ToLower(body.Email)).
		Where("deleted = ?", false).
		First(&models.Users{})
	if res_user.RowsAffected <= 0 {
		res_user := configs.DB.Unscoped().
			Where("LOWER(email) = ?", strings.ToLower(body.Email)).
			Where("deleted = ?", false).
			First(&models.Users{})
		if res_user.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email is already registered",
			})
			return
		}
	}

	data_update := map[string]interface{}{
		"fullname":  body.Fullname,
		"email":     body.Email,
		"is_active": body.IsActive,
	}

	if body.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		body.Password = string(hash)
		data_update["password"] = body.Password
	}

	err := configs.DB.Model(&models.Users{}).Where("id = ?", body.Id).Updates(data_update).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_role := configs.DB.Model(&models.UserRoles{}).Where("user_id = ?", body.Id).Update("role_id", body.RoleId).Error
	if err_role != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_role.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data updated",
		"data_updated": data_update,
	})
}

func UserDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err_user := configs.DB.Model(&models.Users{}).
		Where("id = ?", id).
		Update("deleted", true).Error
	if err_user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_user.Error(),
		})
		return
	}

	err_user_role := configs.DB.Model(&models.UserRoles{}).
		Where("user_id = ?", id).
		Update("deleted", true).Error
	if err_user_role != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err_user_role.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data deleted",
	})
}
