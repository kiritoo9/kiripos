package helpers

import (
	"encoding/base64"
	"fmt"
	"kiripos/src/configs"
	"kiripos/src/models"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

func GenerateImage(dir string, base64Image string, unique string) string {
	var filename string
	if base64Image != "" {
		var path string = "./cdn/"

		// CREATE "CDN" FOLDER
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0700); err != nil {
				panic(err.Error())
			}
		}

		// CREATE SPECIFIC FOLDER
		path += dir + "/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0700); err != nil {
				panic(err.Error())
			}
		}

		filename = unique + ".png"
		b64data := base64Image[strings.IndexByte(base64Image, ',')+1:]
		dec, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			panic(err.Error())
		}
		f, err := os.Create(path + filename)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		if _, err := f.Write(dec); err != nil {
			panic(err.Error())
		}
		if err := f.Sync(); err != nil {
			panic(err.Error())
		}
	}
	return filename
}

func GettRealPath(c *gin.Context, path string) string {
	var protocol string = "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}
	return protocol + "://" + c.Request.Host + "/cdn/" + path
}

func RemoveFile(path string) {
	os.Remove("./cdn/" + path)
}

func GetToken(c *gin.Context) map[string]interface{} {
	bearerToken := c.Request.Header.Get("Authorization")
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
				return map[string]interface{}{
					"id":        claims["id"],
					"fullname":  claims["fullname"],
					"branch_id": claims["branch_id"],
				}
			}
		}

	}
	return nil
}

func GenerateCustomerCode() (code string) {
	var lastData int64
	code = "CS001"
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
	return
}
