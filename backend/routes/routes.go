package routes

import (
	"github.com/Alan-Luc/VertiLog/backend/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// user routes
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login")

	// app routes
	app := router.Group("/app")
	{
		app.GET("/sessions")
	}

	return router
}
