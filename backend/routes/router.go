package routes

import (
	"log"
	"net/http"

	"github.com/Alan-Luc/clog/backend/utils/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Set-Cookie")
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// public routes
	authRoutes(router.Group("/"))

	// protected routes
	protected := router.Group("/app")
	protected.Use(auth.JWTAuthMiddleWare())
	sessionRoutes(protected)
	climbRoutes(protected)
	userRoutes(protected)

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
			log.Printf("Failed to start server: %+v", err)
		}
	}()

	return server
}
