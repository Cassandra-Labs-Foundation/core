package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/Cassandra-Labs-Foundation/core/internal/api/auth"
	"github.com/Cassandra-Labs-Foundation/core/internal/api/middleware"
	personApi "github.com/Cassandra-Labs-Foundation/core/internal/api/person"
	"github.com/Cassandra-Labs-Foundation/core/internal/clients/supabase"
	"github.com/Cassandra-Labs-Foundation/core/internal/config"
	"github.com/Cassandra-Labs-Foundation/core/internal/repository"
	authService "github.com/Cassandra-Labs-Foundation/core/internal/service/auth"
	personService "github.com/Cassandra-Labs-Foundation/core/internal/service/person"
	"github.com/Cassandra-Labs-Foundation/core/pkg/jwt"
)

func main() {
	// Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
	
    // Load configuration
    cfg := config.Load()
    
    // Create Supabase client
    log.Printf("Connecting to Supabase at: %s", cfg.Supabase.URL)
    supabaseClient := supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.APIKey)
	
	// Create JWT service
	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.ExpiryMinutes)
	
	// Create auth service
	authSvc := authService.NewService(jwtService)
	
	// Create auth handler
	authHandler := auth.NewHandler(authSvc)
	
	// Create person repository and service using REST API
	personRepo := repository.NewPersonRestRepository(supabaseClient)
	personSvc := personService.NewService(personRepo)
	personHandler := personApi.NewHandler(personSvc)
	
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
		
		// Person entity routes
		personRoutes := protected.Group("/entities/person")
		{
			personRoutes.POST("", personHandler.Create)
			personRoutes.GET("", personHandler.List)
			personRoutes.GET("/:id", personHandler.Get)
			personRoutes.PATCH("/:id", personHandler.Update)
		}
		
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