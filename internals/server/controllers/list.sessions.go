package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
)

func ListSessions(c *gin.Context) {
	var query struct {
		Status  string `form:"status" binding:"required"`
		UserID  string `form:"user_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required query parameters: status and user_id",
			"details": err.Error(),
		})
		return
	}

	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	statusInput := strings.ToLower(query.Status)
	sessions, err := game.ListSessions(query.UserID, statusInput, dynamoService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve sessions",
			"details": err.Error(),
		})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sessions found"})
		return
	}

	var sessionDetails []map[string]string
	for _, session := range sessions {
		sessionDetails = append(sessionDetails, map[string]string{
			"game":       session.GameName,
			"status":     session.Status,
			"run_id":     session.SessionId,
			"start_time": session.StartTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Listing sessions",
		"sessions": sessionDetails,
	})
}
