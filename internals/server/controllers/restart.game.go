package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
	"github.com/vikraj01/terraplay/internals/utils"
)

func RestartGame(c *gin.Context) {
	var body struct {
		SessionID string `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize DynamoDB session",
			"details": err.Error(),
		})
		return
	}

	details, err := dynamodbService.GetDetailsBySessionID(body.SessionID)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch session details",
			"details": err.Error(),
		})
		return
	}

	if details.InstanceId == "" {
		log.Println("⚠️ Error: Instance ID is missing in session details.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Instance ID is missing in the session data.",
		})
		return
	}

	newServerIP, err := game.RestartEC2(details.InstanceId)
	if err != nil {
		log.Printf("Error retrieving new server IP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error retrieving new server IP",
			"details": err.Error(),
		})
		return
	}

	user := "ec2-user"
	sshConfig, err := utils.GetSSHConfig(newServerIP, user)
	if err != nil {
		log.Printf("Error establishing SSH connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error establishing SSH connection",
			"details": err.Error(),
		})
		return
	}

	gameConfig := config.FindGameConfig(details.GameName)
	backupFile := "/tmp/backup.tar.gz"
	backupPath := gameConfig.VolumePath
	s3Bucket := os.Getenv("AWS_GAME_BACKUP_BUCKET")

	err = game.RestoreEC2(sshConfig, backupPath, s3Bucket, backupFile, newServerIP, details.InstanceId)
	if err != nil {
		log.Printf("Error restarting and restoring EC2: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error during restore process",
			"details": err.Error(),
		})
		return
	}

	err = dynamodbService.UpdateSessionStatusAndIP(body.SessionID, "running", newServerIP)
	if err != nil {
		log.Printf("Error updating session with new IP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update session status",
			"details": err.Error(),
		})
		return
	}

	message := fmt.Sprintf(
		"✅ EC2 instance with IP `%s` has been restarted and data restored. Workspace: `%s`",
		newServerIP, details.Workspace)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
