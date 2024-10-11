package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func DiscordAuthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Placeholder: Initiating Discord OAuth flow.",
    })
}
