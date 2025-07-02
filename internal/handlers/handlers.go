package handlers

import (
	"encoding/json"
	"fmt"
	"gobi/config"
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/internal/services"
	"gobi/internal/services/infrastructure"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds all dependencies for API handlers.
type Handler struct {
	DB                *gorm.DB
	UserService       *services.UserService
	DataSourceService *services.DataSourceService
	QueryService      *services.QueryService
	ChartService      *services.ChartService
	ReportService     *services.ReportService
	TemplateService   *services.TemplateService
}

// NewHandler creates a new Handler instance
func NewHandler(db *gorm.DB) *Handler {
	// Create infrastructure services
	cacheService := infrastructure.NewCacheService(&config.AppConfig)
	validationService := infrastructure.NewValidationService()
	encryptionService := infrastructure.NewEncryptionService()
	authService := infrastructure.NewAuthService()

	// Create user repository for permission service
	userRepo := repositories.NewUserRepository(db)
	permissionService := infrastructure.NewPermissionService(userRepo)

	sqlExecutionService := infrastructure.NewSQLExecutionService()
	reportGeneratorService := infrastructure.NewReportGeneratorService()
	webhookTriggerService := infrastructure.NewWebhookTriggerService()

	// Create API key repository
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

	return &Handler{
		DB:                db,
		UserService:       serviceFactory.CreateUserService(),
		DataSourceService: serviceFactory.CreateDataSourceService(),
		QueryService:      serviceFactory.CreateQueryService(),
		ChartService:      serviceFactory.CreateChartService(),
		ReportService:     serviceFactory.CreateReportService(),
		TemplateService:   serviceFactory.CreateTemplateService(),
	}
}

// CreateUser handles new user registration
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(errors.NewBadRequestError("Invalid user data", err))
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user login and JWT token generation
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid login data", err))
		return
	}

	token, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe retrieves the authenticated user's information
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	user, err := h.UserService.GetUserByID(userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DataSource handlers
func (h *Handler) CreateDataSource(c *gin.Context) {
	var ds models.DataSource
	if err := c.ShouldBindJSON(&ds); err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source data", err))
		return
	}

	userID, _ := c.Get("userID")
	if err := h.DataSourceService.CreateDataSource(&ds, userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, ds)
}

func (h *Handler) ListDataSources(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	dataSources, err := h.DataSourceService.ListDataSources(userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dataSources)
}

func (h *Handler) GetDataSource(c *gin.Context) {
	id := c.Param("id")
	dsID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	ds, err := h.DataSourceService.GetDataSource(uint(dsID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ds)
}

func (h *Handler) UpdateDataSource(c *gin.Context) {
	id := c.Param("id")
	dsID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source ID", err))
		return
	}

	var req models.DataSource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source data", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	ds, err := h.DataSourceService.UpdateDataSource(uint(dsID), &req, userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ds)
}

func (h *Handler) DeleteDataSource(c *gin.Context) {
	id := c.Param("id")
	dsID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.DataSourceService.DeleteDataSource(uint(dsID), userID.(uint), isAdmin); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data source deleted successfully"})
}

// Query handlers
func (h *Handler) CreateQuery(c *gin.Context) {
	var query models.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(errors.NewBadRequestError("Invalid query data", err))
		return
	}

	userID, _ := c.Get("userID")
	if err := h.QueryService.CreateQuery(&query, userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, query)
}

func (h *Handler) ListQueries(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	queries, err := h.QueryService.ListQueries(userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, queries)
}

func (h *Handler) GetQuery(c *gin.Context) {
	id := c.Param("id")
	queryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid query ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	query, err := h.QueryService.GetQuery(uint(queryID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, query)
}

func (h *Handler) ExecuteQuery(c *gin.Context) {
	id := c.Param("id")
	queryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid query ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	result, err := h.QueryService.ExecuteQuery(uint(queryID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateQuery(c *gin.Context) {
	id := c.Param("id")
	queryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid query ID", err))
		return
	}

	var req models.Query
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid query data", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	query, err := h.QueryService.UpdateQuery(uint(queryID), &req, userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, query)
}

func (h *Handler) DeleteQuery(c *gin.Context) {
	id := c.Param("id")
	queryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid query ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.QueryService.DeleteQuery(uint(queryID), userID.(uint), isAdmin); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query deleted successfully"})
}

// Chart handlers
func (h *Handler) CreateChart(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		QueryID     uint   `json:"queryId"`
		Config      string `json:"config"`
		Data        string `json:"data"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid chart data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid chart data"))
		}
		return
	}

	// 验证图表类型
	validChartTypes := map[string]bool{
		"bar":               true,
		"line":              true,
		"pie":               true,
		"scatter":           true,
		"radar":             true,
		"heatmap":           true,
		"gauge":             true,
		"funnel":            true,
		"area":              true, // 面积图 - 显示数据趋势和变化量
		"3d-bar":            true,
		"3d-scatter":        true,
		"3d-surface":        true,
		"3d-bubble":         true,
		"treemap":           true,
		"sunburst":          true,
		"tree":              true,
		"boxplot":           true,
		"candlestick":       true,
		"wordcloud":         true,
		"graph":             true,
		"waterfall":         true,
		"polar":             true,
		"gantt":             true,
		"rose":              true,
		"geo":               true,
		"map":               true,
		"choropleth":        true,
		"progress":          true,
		"circular-progress": true,
	}

	if !validChartTypes[req.Type] {
		c.Error(errors.NewBadRequestError("Invalid chart type", nil))
		return
	}

	userID, _ := c.Get("userID")
	chart := models.Chart{
		Name:        req.Name,
		Type:        req.Type,
		QueryID:     req.QueryID,
		Config:      req.Config,
		Data:        req.Data,
		Description: req.Description,
		UserID:      userID.(uint),
	}

	if err := h.ChartService.CreateChart(&chart, userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, chart)
}

func (h *Handler) ListCharts(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	charts, err := h.ChartService.ListCharts(userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, charts)
}

func (h *Handler) GetChart(c *gin.Context) {
	id := c.Param("id")
	chartID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid chart ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	chart, err := h.ChartService.GetChart(uint(chartID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *Handler) UpdateChart(c *gin.Context) {
	id := c.Param("id")
	chartID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid chart ID", err))
		return
	}

	var req models.Chart
	if err := c.ShouldBindJSON(&req); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid chart data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid chart data"))
		}
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	chart, err := h.ChartService.UpdateChart(uint(chartID), &req, userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *Handler) DeleteChart(c *gin.Context) {
	id := c.Param("id")
	chartID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid chart ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.ChartService.DeleteChart(uint(chartID), userID.(uint), isAdmin); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chart deleted successfully"})
}

// Excel template handlers
func (h *Handler) UploadTemplate(c *gin.Context) {
	file, err := c.FormFile("template")
	if err != nil {
		c.Error(errors.NewBadRequestError("Template file is required", err))
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	if name == "" {
		c.Error(errors.NewBadRequestError("Template name is required", nil))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.Error(errors.WrapError(err, "Could not open template file"))
		return
	}
	defer src.Close()

	// Read file content into bytes
	var buf strings.Builder
	_, err = io.Copy(&buf, src)
	if err != nil {
		c.Error(errors.WrapError(err, "Could not read template file"))
		return
	}

	userID, _ := c.Get("userID")
	template := models.ExcelTemplate{
		Name:        name,
		Template:    []byte(buf.String()),
		Description: description,
		UserID:      userID.(uint),
	}

	if err := h.TemplateService.CreateTemplate(&template, userID.(uint)); err != nil {
		c.Error(errors.WrapError(err, "Could not save template"))
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	templates, err := h.TemplateService.ListTemplates(userID.(uint), isAdmin)
	if err != nil {
		c.Error(errors.WrapError(err, "Could not fetch templates"))
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	templateID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid template ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	template, err := h.TemplateService.GetTemplate(uint(templateID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	templateID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid template ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid template data", err))
		return
	}

	updates := &models.ExcelTemplate{
		Name:        req.Name,
		Description: req.Description,
	}

	template, err := h.TemplateService.UpdateTemplate(uint(templateID), updates, userID.(uint), isAdmin)
	if err != nil {
		c.Error(errors.WrapError(err, "Could not update template"))
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *Handler) DownloadTemplate(c *gin.Context) {
	id := c.Param("id")
	templateID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid template ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	template, err := h.TemplateService.DownloadTemplate(uint(templateID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	// 设置响应头，告诉浏览器这是一个文件下载
	c.Header("Content-Disposition", "attachment; filename="+template.Name)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(template.Template)))

	// 写入文件内容
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", template.Template)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	templateID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid template ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.TemplateService.DeleteTemplate(uint(templateID), userID.(uint), isAdmin); err != nil {
		c.Error(errors.WrapError(err, "Could not delete template"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// Admin handlers
func (h *Handler) DashboardStats(c *gin.Context) {
	stats, err := h.TemplateService.GetDashboardStats()
	if err != nil {
		c.Error(errors.WrapError(err, "Could not fetch dashboard stats"))
		return
	}

	c.JSON(http.StatusOK, stats)
}

// User management handlers (admin only)
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.UserService.ListUsers()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid user ID", err))
		return
	}

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid user data", err))
		return
	}

	user, err := h.UserService.UpdateUser(uint(userID), req.Email, req.Role)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid user ID", err))
		return
	}

	if err := h.UserService.DeleteUser(uint(userID)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Utility handlers
func (h *Handler) ClearCache(c *gin.Context) {
	var req struct {
		Type string `json:"type"` // "query" or "all"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid cache type", err))
		return
	}

	switch req.Type {
	case "query":
		utils.ClearCache()
	case "all":
		utils.ClearCache()
		// Add other caches here if they exist
	default:
		c.Error(errors.NewBadRequestError("Unsupported cache type", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Cache '%s' cleared", req.Type)})
}

func (h *Handler) TestDatabaseConnection(c *gin.Context) {
	var ds models.DataSource
	if err := c.ShouldBindJSON(&ds); err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source data", err))
		return
	}

	if err := h.DataSourceService.TestConnection(&ds); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid user ID", err))
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
		c.Error(errors.NewBadRequestError("Password is required", err))
		return
	}

	currentUserID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if role.(string) != "admin" && currentUserID.(uint) != uint(userID) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := h.UserService.ResetPassword(uint(userID), req.Password); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// CreateAPIKey handles API key creation for the authenticated user
func (h *Handler) CreateAPIKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Name      string     `json:"name"`
		ExpiresAt *time.Time `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid API key data", err))
		return
	}
	apiKey, plainKey, err := h.UserService.CreateAPIKey(userID.(uint), req.Name, req.ExpiresAt)
	if err != nil {
		c.Error(err)
		return
	}
	resp := gin.H{
		"api_key":    plainKey, // Only show once
		"prefix":     apiKey.Prefix,
		"name":       apiKey.Name,
		"expires_at": apiKey.ExpiresAt,
		"created_at": apiKey.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// ListAPIKeys lists all API keys for the user (admin can list all)
func (h *Handler) ListAPIKeys(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"
	keys, err := h.UserService.ListAPIKeys(userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, keys)
}

// RevokeAPIKey revokes an API key by ID
func (h *Handler) RevokeAPIKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"
	id := c.Param("id")
	keyID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid API key ID", err))
		return
	}
	if err := h.UserService.RevokeAPIKey(userID.(uint), uint(keyID), isAdmin); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}

// SystemStats returns system statistics
func (h *Handler) SystemStats(c *gin.Context) {
	dbStats := database.GetConnectionStats()

	cacheStats := utils.GetCacheStats()

	stats := map[string]interface{}{
		"database":  dbStats,
		"cache":     cacheStats,
		"timestamp": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, stats)
}
