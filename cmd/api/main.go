package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/priince9381/irm_backend/internal/config"
	"github.com/priince9381/irm_backend/internal/handlers"
	"github.com/priince9381/irm_backend/internal/middleware"
	"github.com/priince9381/irm_backend/internal/repository"
	"github.com/priince9381/irm_backend/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Elasticsearch
	esDB, err := repository.NewElasticsearchDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Ensure upload directory exists
	uploadDir := "/home/administrator/test/irm_backend/uploads"
	if err := utils.EnsureDir(uploadDir); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static files
	router.Static("/media/uploads", uploadDir)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(esDB, cfg.JWTSecret)
	postHandler := handlers.NewHandler(esDB, cfg)

	// Public routes
	public := router.Group("/api/v1/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.POST("/posts", postHandler.CreatePost)
		protected.GET("/posts", postHandler.GetPosts)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
