package main

import (
	"context"
	"gobi/config"
	"gobi/internal/handlers"
	"gobi/internal/middleware"
	"gobi/internal/repositories"
	"gobi/internal/services"
	"gobi/internal/services/infrastructure"
	"gobi/pkg/database"
	"gobi/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	limiterlib "github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	_ = godotenv.Load()

	config.LoadConfig()
	cfg := config.AppConfig

	if err := database.InitDB(&cfg); err != nil {
		utils.Logger.Fatalf("Failed to initialize database: %v", err)
	}
	db := database.GetDB()

	// 初始化连接管理器
	database.InitConnectionManager(&cfg)

	defer database.CloseAllConnections()

	// 初始化智能缓存
	utils.InitQueryCache(&cfg)
	utils.InitReportGenerator()
	defer utils.StopReportGenerator()

	// Create infrastructure services
	cacheService := infrastructure.NewCacheService(&cfg)
	validationService := infrastructure.NewValidationService()
	encryptionService := infrastructure.NewEncryptionService()
	authService := infrastructure.NewAuthService()

	// Create user repository for permission service
	userRepo := repositories.NewUserRepository(db)
	permissionService := infrastructure.NewPermissionService(userRepo)

	sqlExecutionService := infrastructure.NewSQLExecutionService()
	reportGeneratorService := infrastructure.NewReportGeneratorService()
	webhookTriggerService := infrastructure.NewWebhookTriggerService()

	apiKeyRepo := repositories.NewAPIKeyRepository(db)

	// Create service factory
	serviceFactory := services.NewServiceFactory(
		db,
		cacheService,
		validationService,
		encryptionService,
		authService,
		permissionService,
		sqlExecutionService,
		reportGeneratorService,
		webhookTriggerService,
		apiKeyRepo,
	)

	h := handlers.NewHandler(db)
	reportHandler := handlers.NewReportHandler(db, serviceFactory)
	webhookHandler := handlers.NewWebhookHandler(db, serviceFactory)

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	rate, _ := limiterlib.NewRateFromFormatted("10-S")
	store := memory.NewStore()
	instance := ginlimiter.NewMiddleware(limiterlib.New(store, rate))
	r.Use(instance)

	p := ginprometheus.NewPrometheus("gobi")
	p.Use(r)

	r.GET("/healthz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(500, gin.H{"status": "db error"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.SQLSecurityMiddleware())
	r.Use(gin.Logger())

	// Public routes
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/register", h.CreateUser)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware(&cfg, h.UserService))
	{
		// User routes
		authorized.GET("/me", h.GetMe)

		// API Key management
		authorized.POST("/apikeys", h.CreateAPIKey)
		authorized.GET("/apikeys", h.ListAPIKeys)
		authorized.DELETE("/apikeys/:id", h.RevokeAPIKey)

		// Webhook management
		authorized.POST("/webhooks", webhookHandler.CreateWebhook)
		authorized.GET("/webhooks", webhookHandler.ListWebhooks)
		authorized.GET("/webhooks/:id", webhookHandler.GetWebhook)
		authorized.PUT("/webhooks/:id", webhookHandler.UpdateWebhook)
		authorized.DELETE("/webhooks/:id", webhookHandler.DeleteWebhook)
		authorized.GET("/webhooks/:id/deliveries", webhookHandler.ListWebhookDeliveries)
		authorized.POST("/webhooks/:id/test", webhookHandler.TestWebhook)

		// Query routes
		authorized.POST("/queries", h.CreateQuery)
		authorized.GET("/queries", h.ListQueries)
		authorized.GET("/queries/:id", h.GetQuery)
		authorized.PUT("/queries/:id", h.UpdateQuery)
		authorized.DELETE("/queries/:id", h.DeleteQuery)
		authorized.POST("/queries/:id/execute", h.ExecuteQuery)

		// Data source routes
		authorized.POST("/datasources", h.CreateDataSource)
		authorized.GET("/datasources", h.ListDataSources)
		authorized.GET("/datasources/:id", h.GetDataSource)
		authorized.PUT("/datasources/:id", h.UpdateDataSource)
		authorized.DELETE("/datasources/:id", h.DeleteDataSource)
		authorized.POST("/datasources/test", h.TestDatabaseConnection)

		// Chart routes
		authorized.POST("/charts", h.CreateChart)
		authorized.GET("/charts", h.ListCharts)
		authorized.GET("/charts/:id", h.GetChart)
		authorized.PUT("/charts/:id", h.UpdateChart)
		authorized.DELETE("/charts/:id", h.DeleteChart)

		// Excel template routes
		authorized.POST("/templates", h.UploadTemplate)
		authorized.GET("/templates", h.ListTemplates)
		authorized.GET("/templates/:id", h.GetTemplate)
		authorized.PUT("/templates/:id", h.UpdateTemplate)
		authorized.GET("/templates/:id/download", h.DownloadTemplate)
		authorized.DELETE("/templates/:id", h.DeleteTemplate)

		// Cache clear (admin only)
		authorized.POST("/cache/clear", h.ClearCache)

		// Dashboard stats
		authorized.GET("/dashboard/stats", h.DashboardStats)

		// System monitoring
		authorized.GET("/system/stats", h.SystemStats)

		// User management (admin only)
		authorized.GET("/users", h.ListUsers)
		authorized.PUT("/users/:id", h.UpdateUser)
		authorized.DELETE("/users/:id", h.DeleteUser)
		authorized.POST("/users/:id/reset-password", h.ResetPassword)

		// Report routes
		authorized.POST("/reports", reportHandler.CreateReport)
		authorized.GET("/reports", reportHandler.ListReports)
		authorized.GET("/reports/:id", reportHandler.GetReport)
		authorized.PUT("/reports/:id", reportHandler.UpdateReport)
		authorized.DELETE("/reports/:id", reportHandler.DeleteReport)
		authorized.POST("/reports/:id/generate", reportHandler.GenerateReport)
		authorized.GET("/reports/:id/status", reportHandler.GetReportStatus)

		// Report schedule routes
		authorized.POST("/reports/schedules", reportHandler.CreateReportSchedule)
		authorized.GET("/reports/schedules", reportHandler.ListReportSchedules)
		authorized.GET("/reports/schedules/:id", reportHandler.GetReportSchedule)
		authorized.PUT("/reports/schedules/:id", reportHandler.UpdateReportSchedule)
		authorized.DELETE("/reports/schedules/:id", reportHandler.DeleteReportSchedule)

		// Legacy report routes (for backward compatibility)
		authorized.POST("/reports/generate/excel", reportHandler.GenerateExcelReport)
		authorized.GET("/reports/:id/download", reportHandler.DownloadReport)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
