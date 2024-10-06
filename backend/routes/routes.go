package routes

import (
	"github.com/Alan-Luc/VertiLog/backend/handlers"
	"github.com/gin-gonic/gin"
)

// Route grouping
func authRoutes(router *gin.RouterGroup) {
	router.POST("/register", handlers.RegisterUserHandler)
	router.POST("/login", handlers.LoginUserHandler)
}

func sessionRoutes(router *gin.RouterGroup) {
	router.GET("/sessions/:id", handlers.GetSessionByIDHandler)
	router.GET("/sessions", handlers.GetAllSessionsHandler)
	router.GET("/sessions/summaries", handlers.GetSessionSummariesByDateHandler)
}

func climbRoutes(router *gin.RouterGroup) {
	router.POST("/logClimb", handlers.LogClimbHandler)
	router.GET("/climbs/:id", handlers.GetClimbByIDHandler)
	router.GET("/climbs", handlers.GetAllClimbsHandler)
}
