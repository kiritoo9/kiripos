package routes

import (
	welcome "kiripos/src/controllers"
	"kiripos/src/controllers/auth"
	"kiripos/src/controllers/masters"
	"kiripos/src/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.GET("/", welcome.Index)

	// AUTH
	_auth := router.Group("/auth")
	{
		_auth.POST("/login", auth.Login)
	}

	// V1
	authorized := router.Group("/v1")
	authorized.Use(middlewares.Authroized())
	{
		// MASTERS
		_roles := authorized.Group(("/roles"))
		{
			_roles.GET("/", masters.RoleList)
		}

		_users := authorized.Group(("/users"))
		{
			_users.GET("/", masters.UserList)
			_users.GET("/:id", masters.UserDetail)
			_users.POST("/", masters.UserInsert)
			_users.PUT("/", masters.UserUpdate)
			_users.DELETE("/:id", masters.UserDelete)
		}

		_branches := authorized.Group(("/branches"))
		{
			_branches.GET("/", masters.BranchList)
			_branches.GET("/:id", masters.BranchDetail)
			_branches.POST("/", masters.BranchInsert)
			_branches.PUT("/", masters.BranchUpdate)
			_branches.DELETE("/:id", masters.BranchDelete)
		}

		_branche_users := authorized.Group(("/branch_users"))
		{
			_branche_users.GET("/:branch_id", masters.BranchUserList)
			_branche_users.GET("/:branch_id/:id", masters.BranchUserDetail)
			_branche_users.POST("/:branch_id", masters.BranchUserInsert)
			_branche_users.PUT("/:branch_id", masters.BranchUserUpdate)
			_branche_users.DELETE("/:branch_id/:id", masters.BranchUserDelete)
		}
	}
}
