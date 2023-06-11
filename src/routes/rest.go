package routes

import (
	welcome "kiritech/src/controllers"
	"kiritech/src/controllers/auth"
	"kiritech/src/controllers/masters"
	"kiritech/src/middlewares"

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
		authorized.GET("/users", masters.UserList)
	}
}
