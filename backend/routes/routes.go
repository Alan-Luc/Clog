package routes

import (
	"github.com/Alan-Luc/VertiLog/backend/handlers"
	"github.com/gin-gonic/gin"
)

// Route grouping
func authRoutes(router *gin.RouterGroup) {
	router.POST("/register", handlers.UserRegisterHandler)
	router.POST("/login", handlers.UserLoginHandler)
}

func sessionRoutes(router *gin.RouterGroup) {
	router.GET("/sessions", handlers.SessionGetAllHandler)
	router.GET("/sessions/:id", handlers.SessionGetByIDHandler)
	router.GET("/sessions/summaries", handlers.SessionGetSummariesByDateHandler)
}

func climbRoutes(router *gin.RouterGroup) {
	router.GET("/climbs", handlers.ClimbGetAllHandler)
	router.GET("/climbs/:id", handlers.ClimbGetByIDHandler)
	router.POST("/logClimb", handlers.ClimbLogHandler)
}

func userRoutes(router *gin.RouterGroup) {
	router.PATCH("/profile", handlers.UserProfileHandler)
}
