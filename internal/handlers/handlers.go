package handlers

import (
	"encoding/json"
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Handler holds all dependencies for API handlers.
type Handler struct {
	DB *gorm.DB
}

// NewHandler creates a new handler with dependencies.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// CreateUser handles new user registration
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(errors.NewBadRequestError("Invalid user data", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(errors.WrapError(err, "Failed to hash password"))
		return
	}

	user.Password = string(hashedPassword)
	user.Role = "user" // Default role

	if err := h.DB.Create(&user).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create user"))
		return
	}

	user.Password = "" // Clear password before sending back
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

	var user models.User
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.Error(errors.ErrUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Error(errors.ErrUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.Error(errors.WrapError(err, "Could not generate token"))
		return
	}

	user.LastLogin = time.Now()
	h.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe retrieves the authenticated user's information
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	user.Password = "" // Never return password
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
	ds.UserID = userID.(uint)

	// Encrypt password before saving
	if ds.Password != "" {
		encryptedPassword, err := utils.EncryptAES(ds.Password)
		if err != nil {
			c.Error(errors.WrapError(err, "Failed to encrypt password"))
			return
		}
		ds.Password = encryptedPassword
	}

	if err := h.DB.Create(&ds).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create data source"))
		return
	}

	c.JSON(http.StatusCreated, ds)
}

func (h *Handler) ListDataSources(c *gin.Context) {
	var dataSources []models.DataSource
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := h.DB.Model(&models.DataSource{})
	if role.(string) != "admin" {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&dataSources).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not fetch data sources"))
		return
	}

	// Clear passwords before sending
	for i := range dataSources {
		dataSources[i].Password = ""
	}

	c.JSON(http.StatusOK, dataSources)
}

func (h *Handler) GetDataSource(c *gin.Context) {
	id := c.Param("id")
	var ds models.DataSource
	if err := h.DB.First(&ds, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && ds.UserID != userID.(uint) && !ds.IsPublic {
		c.Error(errors.ErrForbidden)
		return
	}

	ds.Password = "" // Never return password
	c.JSON(http.StatusOK, ds)
}

func (h *Handler) UpdateDataSource(c *gin.Context) {
	id := c.Param("id")
	var ds models.DataSource
	if err := h.DB.First(&ds, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && ds.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	var req models.DataSource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source data", err))
		return
	}

	// Preserve original fields
	req.ID = ds.ID
	req.UserID = ds.UserID
	req.CreatedAt = ds.CreatedAt

	// Encrypt password if it's being updated
	if req.Password != "" {
		encryptedPassword, err := utils.EncryptAES(req.Password)
		if err != nil {
			c.Error(errors.WrapError(err, "Failed to encrypt password"))
			return
		}
		req.Password = encryptedPassword
	} else {
		// If password is not provided in update, keep the old one
		req.Password = ds.Password
	}

	if err := h.DB.Save(&req).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update data source"))
		return
	}

	req.Password = ""
	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteDataSource(c *gin.Context) {
	id := c.Param("id")
	var ds models.DataSource
	if err := h.DB.First(&ds, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && ds.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := h.DB.Delete(&ds).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete data source"))
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
	query.UserID = userID.(uint)

	if err := h.DB.Create(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create query"))
		return
	}

	utils.QueryCache.Flush()

	c.JSON(http.StatusCreated, query)
}

func (h *Handler) ListQueries(c *gin.Context) {
	var queries []models.Query
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := h.DB.Preload("DataSource").Preload("User").Model(&models.Query{})
	if role.(string) != "admin" {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&queries).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not fetch queries"))
		return
	}

	c.JSON(http.StatusOK, queries)
}

func (h *Handler) GetQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := h.DB.Preload("DataSource").Preload("User").First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && query.UserID != userID.(uint) && !query.IsPublic {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, query)
}

func (h *Handler) ExecuteQuery(c *gin.Context) {
	id := c.Param("id")

	// Check cache first
	cacheKey := "query_result_" + id
	if result, found := utils.QueryCache.Get(cacheKey); found {
		c.JSON(http.StatusOK, gin.H{
			"success":       true,
			"data":          result,
			"source":        "cache",
			"executionTime": "0ms",
		})
		return
	}

	var query models.Query
	if err := h.DB.Preload("DataSource").First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	// Permission check
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && query.UserID != userID.(uint) && !query.IsPublic {
		c.Error(errors.ErrForbidden)
		return
	}

	// Decrypt password
	if query.DataSource.Password != "" {
		decryptedPassword, err := utils.DecryptAES(query.DataSource.Password)
		if err != nil {
			c.Error(errors.WrapError(err, "Could not decrypt password"))
			return
		}
		query.DataSource.Password = decryptedPassword
	}

	startTime := time.Now()
	results, err := utils.ExecuteSQL(query.DataSource, query.SQL)
	if err != nil {
		c.Error(errors.WrapError(err, "Failed to execute query"))
		return
	}
	executionTime := time.Since(startTime)

	// Update execution count
	query.ExecCount++
	h.DB.Save(&query)

	// Set cache
	utils.QueryCache.Set(cacheKey, results, 5*time.Minute)

	var columns []map[string]string
	if len(results) > 0 {
		for key := range results[0] {
			columns = append(columns, map[string]string{"name": key, "type": "unknown"})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"data":          results,
		"columns":       columns,
		"rowCount":      len(results),
		"executionTime": fmt.Sprintf("%.2fms", float64(executionTime.Nanoseconds())/1e6),
	})
}

func (h *Handler) UpdateQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := h.DB.First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && query.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(errors.NewBadRequestError("Invalid query data", err))
		return
	}

	if err := h.DB.Save(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update query"))
		return
	}

	utils.QueryCache.Flush()

	c.JSON(http.StatusOK, query)
}

func (h *Handler) DeleteQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := h.DB.First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && query.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := h.DB.Delete(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete query"))
		return
	}

	utils.QueryCache.Flush()

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
		"bar":        true,
		"line":       true,
		"pie":        true,
		"scatter":    true,
		"radar":      true,
		"heatmap":    true,
		"gauge":      true,
		"funnel":     true,
		"3d-bar":     true,
		"3d-scatter": true,
		"3d-surface": true,
		"3d-bubble":  true,
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

	if err := h.DB.Create(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create chart"))
		return
	}

	c.JSON(http.StatusCreated, chart)
}

func (h *Handler) ListCharts(c *gin.Context) {
	var charts []models.Chart
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := h.DB.Preload("Query").Preload("User").Model(&models.Chart{})
	if role.(string) != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&charts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch charts"})
		return
	}

	c.JSON(http.StatusOK, charts)
}

func (h *Handler) GetChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := h.DB.Preload("Query").Preload("User").First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *Handler) UpdateChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := h.DB.First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := c.ShouldBindJSON(&chart); err != nil {
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

	if err := h.DB.Save(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update chart"))
		return
	}

	c.JSON(http.StatusOK, chart)
}

func (h *Handler) DeleteChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := h.DB.First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := h.DB.Delete(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete chart"))
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

	buf := new(strings.Builder)
	_, err = io.Copy(buf, src) // This is incorrect, should use bytes.Buffer
	if err != nil {
		c.Error(errors.WrapError(err, "Could not read template file"))
		return
	}

	userID, _ := c.Get("userID")
	template := models.ExcelTemplate{
		Name:        name,
		Template:    []byte(buf.String()), // This is incorrect
		Description: description,
		UserID:      userID.(uint),
	}

	if err := h.DB.Create(&template).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not save template"))
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	var templates []models.ExcelTemplate
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := h.DB.Preload("User").Model(&models.ExcelTemplate{})
	if role.(string) != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&templates).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not fetch templates"))
		return
	}

	// Do not return template content in list view
	for i := range templates {
		templates[i].Template = nil
	}

	c.JSON(http.StatusOK, templates)
}

func (h *Handler) DownloadTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ExcelTemplate
	if err := h.DB.First(&template, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	// 权限检查：管理员可以下载所有模板，普通用户只能下载自己的模板
	if role.(string) != "admin" && template.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	// 设置响应头，告诉浏览器这是一个文件下载
	c.Header("Content-Disposition", "attachment; filename="+template.Name)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(template.Template)))

	// 写入文件内容
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", template.Template)
}

// Admin handlers
func (h *Handler) DashboardStats(c *gin.Context) {
	var totalQueries int64
	var totalCharts int64
	var totalUsers int64
	var todayQueries int64

	today := time.Now().Format("2006-01-02")
	h.DB.Model(&models.Query{}).Count(&totalQueries)
	h.DB.Model(&models.Chart{}).Count(&totalCharts)
	h.DB.Model(&models.User{}).Count(&totalUsers)
	h.DB.Model(&models.Query{}).Where("DATE(created_at) = ?", today).Count(&todayQueries)

	// 查询趋势（最近7天每天的查询数）
	queryTrends := []map[string]interface{}{}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		h.DB.Model(&models.Query{}).Where("DATE(created_at) = ?", date).Count(&count)
		queryTrends = append(queryTrends, map[string]interface{}{"date": date, "count": count})
	}

	// 热门查询（执行次数最多的前5个查询）
	type HotQuery struct {
		Name  string
		Count int64
	}
	hotQueries := []HotQuery{}
	h.DB.Table("queries").Select("name, exec_count as count").Order("exec_count desc").Limit(5).Scan(&hotQueries)

	c.JSON(http.StatusOK, gin.H{
		"totalQueries": totalQueries,
		"totalCharts":  totalCharts,
		"totalUsers":   totalUsers,
		"todayQueries": todayQueries,
		"queryTrends":  queryTrends,
		"hotQueries":   hotQueries,
	})
}

// List users handler
func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not fetch users"))
		return
	}

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser updates a user's role
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid role data", err))
		return
	}

	user.Role = req.Role
	if err := h.DB.Save(&user).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update user"))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ClearCache clears a specific cache type
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
		utils.QueryCache.Flush()
	case "all":
		utils.QueryCache.Flush()
		// Add other caches here if they exist
	default:
		c.Error(errors.NewBadRequestError("Unsupported cache type", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Cache '%s' cleared", req.Type)})
}

// TestDatabaseConnection tests the connection to a data source
func (h *Handler) TestDatabaseConnection(c *gin.Context) {
	var ds models.DataSource
	if err := c.ShouldBindJSON(&ds); err != nil {
		c.Error(errors.NewBadRequestError("Invalid data source data", err))
		return
	}

	// Decrypt password
	if ds.Password != "" {
		decryptedPassword, err := utils.DecryptAES(ds.Password)
		if err != nil {
			c.Error(errors.WrapError(err, "Could not decrypt password"))
			return
		}
		ds.Password = decryptedPassword
	}

	// Use the connection manager to test the connection
	db, err := database.GetConnection(&ds)
	if err != nil {
		c.Error(errors.WrapError(err, "Failed to create connection"))
		return
	}

	if err := db.Ping(); err != nil {
		c.Error(errors.WrapError(err, "Database connection test failed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database connection successful"})
}
