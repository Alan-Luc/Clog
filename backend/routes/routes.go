package routes

import (
	"log"
	"net/http"

	"github.com/Alan-Luc/VertiLog/backend/handlers"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	router := gin.Default()

	// user routes
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginUser)

	// protected routes
	protected := router.Group("/app")
	protected.Use(auth.JWTAuthMiddleWare())
	{
		protected.GET("/sessions", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"pogs": "on dogs"})
		})
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
