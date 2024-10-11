package server

import (
	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/server/handlers"
	"github.com/vikraj01/terraplay/internals/server/middleware"
)

func RegisterRoutes(router *gin.Engine) {

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/game/create", handlers.CreateGame)
	}

	router.GET("/auth/discord/initiate", handlers.InitiateDiscordOAuth)
	router.GET("/auth/discord/callback", handlers.DiscordOAuthCallback)
	router.GET("/auth/discord/status", handlers.CheckAuthStatus)
}
