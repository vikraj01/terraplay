package controllers


import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/game"
)

type GameConfig struct {
	Game     string `json:"game"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

func CreateGame(c *gin.Context) {
	var config GameConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if !game.IsValidGame(config.Game) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported game type"})
		return
	}

	session, workspace, err := game.CreateGameSession(config.UserID, config.UserName, config.Game)
	if err != nil {
		log.Printf("Failed to create game session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game session", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game session created successfully",
		"session_id": session.SessionId,
		"workspace":  workspace,
		"game":       session.GameName,
		"status":     session.Status,
	})
}
