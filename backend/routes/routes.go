package routes

import (
	"net/http"

	"github.com/Alan-Luc/VertiLog/backend/handlers"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/Alan-Luc/VertiLog/backend/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// public routes
	authRoutes(router.Group("/"))

	// protected routes
	protected := router.Group("/app")
	protected.Use(auth.JWTAuthMiddleWare())
	sessionRoutes(protected)
	climbRoutes(protected)

	return router
}

func StartServer() *http.Server {
	router := SetupRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Failed to start server:", zap.Error(err))
		}
	}()

	return server
}

// Route grouping
func authRoutes(router *gin.RouterGroup) {
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginUser)
}

func sessionRoutes(router *gin.RouterGroup) {
	router.GET("/sessions/:id", handlers.GetSessionByIDHandler)
	router.GET("/sessions", handlers.GetAllSessionsHandler)
}

func climbRoutes(router *gin.RouterGroup) {
	router.POST("/logClimb", handlers.LogClimbHandler)
	router.GET("/climbs/:id", handlers.GetClimbByIDHandler)
	router.GET("/climbs", handlers.GetAllClimbsHandler)
}
