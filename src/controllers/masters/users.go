package masters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "this is user",
	})
}

func UserInsert(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Insert user endpoint",
	})
}
