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
	}
}
