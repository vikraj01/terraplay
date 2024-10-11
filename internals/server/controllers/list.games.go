package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/game"
)

func ListGames(c *gin.Context) {
	games := game.GetGames()
	if len(games) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No games available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Available games",
		"games":   games,
	})
}
