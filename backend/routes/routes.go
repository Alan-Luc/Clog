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

	// public routes
	authRoutes(router.Group("/"))

	// protected routes
	protected := router.Group("/app")
	protected.Use(auth.JWTAuthMiddleWare())
	sessionRoutes(protected)
	climbRoutes(protected)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Route grouping
func authRoutes(router *gin.RouterGroup) {
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginUser)
}

func sessionRoutes(router *gin.RouterGroup) {
	router.GET("/sessions/:id", handlers.GetSession)
}

func climbRoutes(router *gin.RouterGroup) {
	router.GET("/climbs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	router.POST("/logClimb", handlers.LogClimb)
}
