package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
	})
}

func Detail(c *gin.Context) {

}

func Create(c *gin.Context) {

}

func Update(c *gin.Context) {

}

func Delete(c *gin.Context) {

}
