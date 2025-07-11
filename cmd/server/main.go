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
	"gobi/pkg/errors"
	"gobi/pkg/utils"
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

// ConsoleAlertChannel 控制台告警通道
type ConsoleAlertChannel struct{}

func (c *ConsoleAlertChannel) SendAlert(alert *errors.Alert) error {
	utils.Logger.WithFields(map[string]interface{}{
		"severity": alert.Severity,
		"type":     alert.Type,
		"message":  alert.Message,
		"details":  alert.Details,
	}).Error("ALERT")
	return nil
}

// setupServer 设置和配置服务器
func setupServer(cfg *config.Config) (*http.Server, error) {
	// 初始化错误监控
	errorMonitor := errors.GetGlobalMonitor()
	defer errorMonitor.Stop()

	// 添加控制台告警通道
	consoleAlertChannel := &ConsoleAlertChannel{}
	errorMonitor.AddAlertChannel(consoleAlertChannel)

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		return nil, errors.WrapError(err, "Failed to initialize database")
	}
	db := database.GetDB()

	// 初始化连接管理器
	database.InitConnectionManager(cfg)
	defer database.CloseAllConnections()

	// 初始化智能缓存
	utils.InitQueryCache(cfg)
	utils.InitReportGenerator()
	defer utils.StopReportGenerator()

	// Create infrastructure services
	cacheService := infrastructure.NewCacheService(cfg)
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
		AllowOrigins:     cfg.Security.CORSOrigins,
		AllowMethods:     cfg.Security.AllowedMethods,
		AllowHeaders:     cfg.Security.AllowedHeaders,
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
	r.Use(middleware.ErrorLogger())
	r.Use(middleware.SQLSecurityMiddleware())
	r.Use(gin.Logger())

	// Public routes
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/register", h.CreateUser)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware(cfg, h.UserService))
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
		authorized.GET("/system/error-stats", h.ErrorStats)

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

	return srv, nil
}

// gracefulShutdown 优雅关闭服务器
func gracefulShutdown(srv *http.Server, timeout time.Duration) error {
	utils.Logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Logger.WithError(err).Error("Server forced to shutdown")
		return err
	}

	utils.Logger.Info("Server exited gracefully")
	return nil
}

func main() {
	_ = godotenv.Load()

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		utils.Logger.WithError(err).Fatal("Failed to load config")
		os.Exit(1)
	}
	cfg := config.GetConfig()

	// 设置服务器
	srv, err := setupServer(cfg)
	if err != nil {
		utils.Logger.WithError(err).Fatal("Failed to setup server")
		os.Exit(1)
	}

	// 启动服务器
	go func() {
		utils.Logger.WithField("port", cfg.Server.Port).Info("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.WithError(err).Fatal("Server failed to start")
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	if err := gracefulShutdown(srv, 5*time.Second); err != nil {
		utils.Logger.WithError(err).Error("Failed to shutdown gracefully")
		os.Exit(1)
	}
}
