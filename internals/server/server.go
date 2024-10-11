package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func StartServer() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Falling back to environment variables.")
	}
	if os.Getenv("DISCORD_CLIENT_ID") == "" || os.Getenv("DISCORD_CLIENT_SECRET") == "" {
		log.Fatal("DISCORD_CLIENT_ID or DISCORD_CLIENT_SECRET not set. Please check your environment variables.")
	}
	
	router := gin.Default()
	RegisterRoutes(router)

	log.Println("Starting Zephyr server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
