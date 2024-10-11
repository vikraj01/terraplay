package server

import (
    "github.com/gin-gonic/gin"
    "github.com/vikraj01/terraplay/internals/server/handlers"
)

func RegisterRoutes(router *gin.Engine) {
    router.GET("/auth/discord", handlers.DiscordAuthHandler)

}
