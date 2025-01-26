package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"irm_backend/internal/config"
	"irm_backend/internal/handlers"
	"irm_backend/internal/middleware"
	"irm_backend/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Elasticsearch
	db, err := repository.NewElasticsearchDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
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

	// Initialize handlers
	h := handlers.NewHandler(db, cfg)

	// Public routes
	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/forgot-password", h.ForgotPassword)
			auth.POST("/reset-password", h.ResetPassword)
		}
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// User routes
		user := protected.Group("/users")
		{
			user.GET("/me", h.GetProfile)
			user.PUT("/me", h.UpdateProfile)
		}

		// Social accounts routes
		social := protected.Group("/social")
		{
			social.GET("/accounts", h.GetSocialAccounts)
			social.POST("/accounts/:platform/connect", h.ConnectSocialAccount)
			social.DELETE("/accounts/:id", h.DisconnectSocialAccount)
		}

		// Posts routes
		posts := protected.Group("/posts")
		{
			posts.GET("", h.GetPosts)
			posts.POST("", h.CreatePost)
			posts.GET("/:id", h.GetPost)
			posts.PUT("/:id", h.UpdatePost)
			posts.DELETE("/:id", h.DeletePost)
			posts.POST("/:id/publish", h.PublishPost)
			posts.POST("/:id/schedule", h.SchedulePost)
		}

		// Media routes
		media := protected.Group("/media")
		{
			media.POST("/upload", h.UploadMedia)
			media.GET("", h.GetMedia)
			media.DELETE("/:id", h.DeleteMedia)
		}

		// Analytics routes
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/overview", h.GetAnalyticsOverview)
			analytics.GET("/:platform", h.GetPlatformAnalytics)
			analytics.GET("/posts/:id", h.GetPostAnalytics)
		}
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(cfg), middleware.AdminMiddleware())
	{
		admin.GET("/users", h.GetUsers)
		admin.GET("/analytics/global", h.GetGlobalAnalytics)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
