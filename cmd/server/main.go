package main

import (
	"gobi/config"
	"gobi/internal/handlers"
	"gobi/internal/middleware"
	"gobi/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.DefaultConfig

	// Initialize database
	if err := database.InitDB(&cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.POST("/api/auth/login", handlers.Login)
	r.POST("/api/auth/register", handlers.Register)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware(&cfg))
	{
		// Query routes
		authorized.POST("/queries", handlers.CreateQuery)
		authorized.GET("/queries", handlers.ListQueries)
		authorized.GET("/queries/:id", handlers.GetQuery)
		authorized.PUT("/queries/:id", handlers.UpdateQuery)
		authorized.DELETE("/queries/:id", handlers.DeleteQuery)

		// Chart routes
		authorized.POST("/charts", handlers.CreateChart)
		authorized.GET("/charts", handlers.ListCharts)
		authorized.GET("/charts/:id", handlers.GetChart)
		authorized.PUT("/charts/:id", handlers.UpdateChart)
		authorized.DELETE("/charts/:id", handlers.DeleteChart)

		// Excel template routes
		authorized.POST("/templates", handlers.CreateTemplate)
		authorized.GET("/templates", handlers.ListTemplates)
		authorized.GET("/templates/:id", handlers.GetTemplate)
		authorized.PUT("/templates/:id", handlers.UpdateTemplate)
		authorized.DELETE("/templates/:id", handlers.DeleteTemplate)
	}

	// Start server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
