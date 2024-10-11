package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	
	RegisterRoutes(router)

	log.Println("Starting Zephyr server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
