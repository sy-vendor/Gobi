package main

import (
	"gobi/config"
	"gobi/internal/handlers"
	"gobi/internal/middleware"
	"gobi/pkg/database"
	"gobi/pkg/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.DefaultConfig

	// Initialize database
	if err := database.InitDB(&cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize query cache (default 5 min, cleanup 10 min)
	utils.InitQueryCache(5*time.Minute, 10*time.Minute)

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Logger())

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

		// Data source routes
		authorized.POST("/datasources", handlers.CreateDataSource)
		authorized.GET("/datasources", handlers.ListDataSources)
		authorized.GET("/datasources/:id", handlers.GetDataSource)
		authorized.PUT("/datasources/:id", handlers.UpdateDataSource)
		authorized.DELETE("/datasources/:id", handlers.DeleteDataSource)

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

		// Cache clear (admin only)
		authorized.POST("/cache/clear", handlers.ClearCache)

		// Dashboard stats
		authorized.GET("/dashboard/stats", handlers.DashboardStats)

		// User list
		authorized.GET("/users", handlers.ListUsers)
		// User update
		authorized.PUT("/users/:id", handlers.UpdateUser)
		// User reset password
		authorized.POST("/users/:id/reset-password", handlers.ResetUserPassword)
	}

	// Start server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
