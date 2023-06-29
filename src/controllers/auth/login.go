package auth

import (
	"log"
	"net/http"
	"os"

	"kiripos/src/configs"
	"kiripos/src/models"

	"github.com/alexsasharegan/dotenv"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var lm LoginModel
	errenv := dotenv.Load()
	if errenv != nil {
		log.Fatalf("Error loading .env file: %v", errenv)
	}

	if errbody := c.BindJSON(&lm); errbody != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Body is not valid",
		})
		return
	}

	var users models.Users
	result := configs.DB.Find(&users, "email = ? AND deleted = ?", lm.Email, false)
	if result.RowsAffected <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Email is not found",
		})
		return
	}

	hash := []byte(users.Password)
	plain := []byte(lm.Password)
	errbcrypt := bcrypt.CompareHashAndPassword(hash, plain)
	if errbcrypt != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password does not match",
		})
		return
	}

	// SETUP RESPONSE
	type Response struct {
		ID    string
		Email string
		Role  string
		Token string
	}
	var response = new(Response)
	response.ID = users.Id.String()
	response.Email = users.Email

	var role_name string = ""
	var user_role models.UserRoles
	res_role := configs.DB.Unscoped().
		Select("roles.name AS role_name").
		Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.deleted = ?", false).
		Where("user_roles.user_id = ?", users.Id).
		First(&user_role)
	if res_role.RowsAffected > 0 {
		role_name = user_role.RoleName
	}
	response.Role = role_name

	if role_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user have no role",
		})
		return
	}

	claims := jwt.MapClaims{}
	claims["id"] = users.Id.String()
	claims["fullname"] = users.Fullname

	var user_branch models.BranchUsers
	res_ubranch := configs.DB.Where("deleted = ?", false).Where("user_id = ?", users.Id).First(&user_branch)
	if res_ubranch.RowsAffected > 0 {
		claims["branch_id"] = user_branch.BranchId
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user have no branch",
		})
		return
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, errjwt := sign.SignedString([]byte(os.Getenv("APP_KEY")))
	if errjwt != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while generate your token",
		})
		return
	}
	response.Token = token

	c.JSON(http.StatusOK, gin.H{
		"message": "Request success",
		"data":    response,
	})
}
