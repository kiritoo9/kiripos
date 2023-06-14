package masters

import (
	"kiripos/src/configs"
	"kiripos/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func _checkBranch(c *gin.Context) string {
	id, err := uuid.Parse(c.Param("branch_id"))
	if err != nil {
		return err.Error()
	}

	res := configs.DB.Unscoped().Where("deleted = ?", false).Where("id = ?", id).Find(&models.Branches{})
	if res.RowsAffected <= 0 {
		return "Branch is not found"
	}

	return ""
}

func BranchUserList(c *gin.Context) {
	if callback := _checkBranch(c); callback != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": callback,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request Success",
	})
}

func BranchUserDetail(c *gin.Context) {

}

func BranchUserInsert(c *gin.Context) {

}

func BranchUserUpdate(c *gin.Context) {

}

func BranchUserDelete(c *gin.Context) {

}
