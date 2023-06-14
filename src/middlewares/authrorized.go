package middlewares

import (
	"fmt"
	"kiripos/src/configs"
	"kiripos/src/models"
	"net/http"
	"os"
	"strings"

	"github.com/alexsasharegan/dotenv"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

func Authroized() gin.HandlerFunc {
	return func(c *gin.Context) {
		errenv := dotenv.Load()
		if errenv != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Key is not valid, please contact admin",
			})
			c.Abort()
			return
		}

		contentType := c.Request.Header.Get("Content-Type")
		bearerToken := c.Request.Header.Get("Authroization")
		if len(contentType) > 0 {
			var allowedToken = false
			if len(strings.Split(bearerToken, "")) >= 2 {
				cleanToken := strings.Split(bearerToken, " ")[1]

				token, err := jwt.Parse(cleanToken, func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method")
					}
					return []byte(os.Getenv("APP_KEY")), nil
				})

				if err == nil {
					claims, ok := token.Claims.(jwt.MapClaims)
					if ok && token.Valid {
						var id = claims["id"]
						var user *models.Users
						data := configs.DB.Find(&user, "id = ? AND deleted = ?", id, false)
						if data.RowsAffected > 0 {
							allowedToken = true
						}
					}
				}

			}

			if !allowedToken {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Access token is not valid",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Header must be json",
			})
			c.Abort()
			return
		}
	}
}
