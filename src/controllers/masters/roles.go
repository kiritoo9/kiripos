package masters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Role list",
	})
}
