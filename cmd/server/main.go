package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/Cassandra-Labs-Foundation/core/internal/api/auth"
	"github.com/Cassandra-Labs-Foundation/core/internal/api/middleware"
	"github.com/Cassandra-Labs-Foundation/core/internal/config"
	authService "github.com/Cassandra-Labs-Foundation/core/internal/service/auth"
	"github.com/Cassandra-Labs-Foundation/core/pkg/jwt"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Create JWT service
	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.ExpiryMinutes)
	
	// Create auth service
	authSvc := authService.NewService(jwtService)
	
	// Create auth handler
	authHandler := auth.NewHandler(authSvc)
	
	// Create gin router
	r := gin.Default()
	
	// Define API routes
	api := r.Group("/api/v1")
	
	// Auth routes (no authentication required)
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}
	
	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(authSvc))
	{
		protected.GET("/auth/validate", authHandler.ValidateToken)
		
		// Add more protected routes here
		protected.GET("/hello", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"message": "Hello, authenticated user!",
				"userID":  userID,
			})
		})
	}
	
	// Start the server
	log.Printf("Starting server on port %s...\n", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}