package server

import (
	"github.com/gin-gonic/gin"
	"github.com/vikraj01/terraplay/internals/server/controllers"
	"github.com/vikraj01/terraplay/internals/server/middleware"
)

func RegisterRoutes(router *gin.Engine) {

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/game/create", controllers.CreateGame)
		protected.POST("/game/stop", controllers.StopGame)
		protected.POST("/game/restart", controllers.RestartGame)
		protected.GET("/game/sessions", controllers.ListSessions)
	}

	router.GET("/game/list", controllers.ListGames)
	router.GET("/auth/discord/initiate", controllers.InitiateDiscordOAuth)
	router.GET("/auth/discord/callback", controllers.DiscordOAuthCallback)
	router.GET("/auth/discord/status", controllers.CheckAuthStatus)
}
