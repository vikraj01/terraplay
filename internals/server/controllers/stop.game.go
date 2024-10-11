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

func StopGame(c *gin.Context) {
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
			"error":   "Failed to initialize dynamodb session",
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

	user := "ec2-user"
	sshConfig, err := utils.GetSSHConfig(details.ServerIP, user)
	if err != nil {
		log.Printf("Error establishing ssh connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error establishing ssh connection",
			"details": err.Error(),
		})
		return
	}

	gameConfig := config.FindGameConfig(details.GameName)
	backupFile := "/tmp/backup.tar.gz"
	backupPath := gameConfig.VolumePath
	s3Bucket := os.Getenv("AWS_GAME_BACKUP_BUCKET")
	err = game.BackupAndStopEC2(sshConfig, backupPath, s3Bucket, backupFile, details.ServerIP)
	if err != nil {
		log.Printf("Error executing backup and stop: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error executing backup and stop",
			"details": err.Error(),
		})
		return
	}

	err = dynamodbService.UpdateSessionStatusAndIP(body.SessionID, "halted", details.ServerIP)
	if err != nil {
		log.Printf("Error updating session status in DynamoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update session status",
			"details": err.Error(),
		})
		return
	}

	message := fmt.Sprintf(
		"✅ EC2 instance with IP `%s` has been backed up and stopped. Workspace: `%s`",
		details.ServerIP, details.Workspace)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
