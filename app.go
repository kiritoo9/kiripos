package main

import (
	"log"
	"os"

	"github.com/alexsasharegan/dotenv"

	"kiripos/src/configs"
	"kiripos/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	routes.Init(router)
	configs.Connect()

	// SERVE
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	router.Run(":" + os.Getenv("APP_PORT"))
}
