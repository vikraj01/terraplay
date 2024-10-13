package controllers

import (
	"fmt"
	"log"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
)

func DestroySession(c *gin.Context) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not initialize database. Please try again later.",
		})
		return
	}

	sessionId := strings.TrimSpace(c.Param("session_id"))
	log.Print(sessionId)
	if sessionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Usage: Provide a valid session_id as a path parameter",
		})
		return
	}

	workspaceID, err := game.DestroyGameSession(sessionId, dynamodbService)
	if err != nil {
		log.Printf("Error during game session destruction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	message := fmt.Sprintf(
		"🛠️ **Destruction of the game session has been initiated!**\n\n"+
			"**Session ID:** `%s`\n"+
			"**Workspace ID:** `%s`\n\n"+
			"Your game session is now in the process of being stopped. "+
			"Thank you for using the Terraplay service!",
		sessionId, workspaceID,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
