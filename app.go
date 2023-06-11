package main

import (
	"os"
	"log"
	"github.com/alexsasharegan/dotenv"
	
	"github.com/gin-gonic/gin"
	"kiritech/src/configs"
	"kiritech/src/routes"
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
	router.Run(":"+os.Getenv("APP_PORT"))
}
