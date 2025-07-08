package middleware

import (
	"fmt"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)

		// 设置请求开始时间
		startTime := time.Now()
		c.Set("start_time", startTime)

		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err, requestID, startTime)
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := getRequestID(c)
				startTime := getStartTime(c)

				// 记录panic详情
				stack := debug.Stack()
				panicErr := fmt.Errorf("panic: %v\n%s", err, stack)

				handlePanic(c, panicErr, requestID, startTime)
			}
		}()
		c.Next()
	}
}

// handleError 处理错误
func handleError(c *gin.Context, err error, requestID string, startTime time.Time) {
	path := c.Request.URL.Path
	method := c.Request.Method
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	duration := time.Since(startTime)

	// 设置请求ID到错误中
	errors.SetRequestID(err, requestID)

	if customErr, ok := err.(*errors.CustomError); ok {
		// 记录自定义错误
		logError(customErr, path, method, userID, role, clientIP, userAgent, duration, requestID)

		// 设置响应头
		setErrorHeaders(c, customErr)

		// 返回错误响应
		c.JSON(getHTTPStatus(customErr.Code), createErrorResponse(customErr))
	} else {
		// 记录系统错误
		logSystemError(err, path, method, userID, role, clientIP, userAgent, duration, requestID)

		// 创建内部服务器错误
		internalErr := errors.NewErrorWithSeverity(
			errors.ErrCodeInternalServer,
			"Internal server error",
			err,
			errors.SeverityHigh,
			errors.CategorySystem,
		).WithContext(c.Request.Context())

		// 设置响应头
		setErrorHeaders(c, internalErr)

		// 返回错误响应
		c.JSON(http.StatusInternalServerError, createErrorResponse(internalErr))
	}

	c.Abort()
}

// handlePanic 处理panic
func handlePanic(c *gin.Context, err error, requestID string, startTime time.Time) {
	path := c.Request.URL.Path
	method := c.Request.Method
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	duration := time.Since(startTime)

	// 记录panic
	logPanic(err, path, method, userID, role, clientIP, userAgent, duration, requestID)

	// 创建panic错误
	panicErr := errors.NewErrorWithSeverity(
		errors.ErrCodeInternalServer,
		"Internal server error",
		err,
		errors.SeverityCritical,
		errors.CategorySystem,
	).WithContext(c.Request.Context())

	// 设置响应头
	setErrorHeaders(c, panicErr)

	// 返回错误响应
	c.JSON(http.StatusInternalServerError, createErrorResponse(panicErr))
	c.Abort()
}

// logError 记录自定义错误
func logError(err *errors.CustomError, path, method string, userID, role interface{}, clientIP, userAgent string, duration time.Duration, requestID string) {
	logFields := map[string]interface{}{
		"request_id":    requestID,
		"path":          path,
		"method":        method,
		"user_id":       userID,
		"role":          role,
		"client_ip":     clientIP,
		"user_agent":    userAgent,
		"duration_ms":   duration.Milliseconds(),
		"error_code":    string(err.Code),
		"error_message": err.Message,
		"severity":      string(err.Severity),
		"category":      string(err.Category),
		"timestamp":     err.Timestamp,
	}

	// 添加详细信息
	if err.Details != nil {
		logFields["details"] = err.Details
	}

	// 添加原始错误
	if err.Err != nil {
		logFields["original_error"] = err.Err.Error()
	}

	// 添加堆栈信息（仅在开发环境）
	if utils.IsDevelopment() && len(err.Stack) > 0 {
		logFields["stack_trace"] = err.Stack
	}

	// 根据严重程度选择日志级别
	switch err.Severity {
	case errors.SeverityCritical:
		utils.Logger.WithFields(logFields).Error("Critical error occurred")
	case errors.SeverityHigh:
		utils.Logger.WithFields(logFields).Error("High severity error occurred")
	case errors.SeverityMedium:
		utils.Logger.WithFields(logFields).Warn("Medium severity error occurred")
	case errors.SeverityLow:
		utils.Logger.WithFields(logFields).Info("Low severity error occurred")
	default:
		utils.Logger.WithFields(logFields).Error("Error occurred")
	}
}

// logSystemError 记录系统错误
func logSystemError(err error, path, method string, userID, role interface{}, clientIP, userAgent string, duration time.Duration, requestID string) {
	logFields := map[string]interface{}{
		"request_id":  requestID,
		"path":        path,
		"method":      method,
		"user_id":     userID,
		"role":        role,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"duration_ms": duration.Milliseconds(),
		"error":       err.Error(),
		"error_type":  fmt.Sprintf("%T", err),
		"severity":    string(errors.SeverityHigh),
		"category":    string(errors.CategorySystem),
	}

	utils.Logger.WithFields(logFields).Error("System error occurred")
}

// logPanic 记录panic
func logPanic(err error, path, method string, userID, role interface{}, clientIP, userAgent string, duration time.Duration, requestID string) {
	logFields := map[string]interface{}{
		"request_id":  requestID,
		"path":        path,
		"method":      method,
		"user_id":     userID,
		"role":        role,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"duration_ms": duration.Milliseconds(),
		"error":       err.Error(),
		"error_type":  "PANIC",
		"severity":    string(errors.SeverityCritical),
		"category":    string(errors.CategorySystem),
	}

	utils.Logger.WithFields(logFields).Error("Panic recovered")
}

// setErrorHeaders 设置错误响应头
func setErrorHeaders(c *gin.Context, err *errors.CustomError) {
	// 设置请求ID
	if err.RequestID != "" {
		c.Header("X-Request-ID", err.RequestID)
	}

	// 设置重试时间
	if err.RetryAfter != nil {
		c.Header("Retry-After", fmt.Sprintf("%d", *err.RetryAfter))
	}

	// 设置错误类型
	c.Header("X-Error-Type", string(err.Category))
	c.Header("X-Error-Severity", string(err.Severity))

	// 设置帮助URL
	if err.HelpURL != "" {
		c.Header("X-Help-URL", err.HelpURL)
	}
}

// createErrorResponse 创建错误响应
func createErrorResponse(err *errors.CustomError) errors.ErrorResponse {
	response := errors.ErrorResponse{
		Code:       string(err.Code),
		Message:    err.Message,
		Details:    err.Details,
		RequestID:  err.RequestID,
		Timestamp:  err.Timestamp,
		Severity:   err.Severity,
		Category:   err.Category,
		RetryAfter: err.RetryAfter,
		HelpURL:    err.HelpURL,
	}

	// 添加原始错误信息（仅在开发环境）
	if utils.IsDevelopment() && err.Err != nil {
		response.Error = err.Err.Error()
	}

	return response
}

// getHTTPStatus 获取HTTP状态码
func getHTTPStatus(code errors.ErrorCode) int {
	switch code {
	case errors.ErrCodeSuccess:
		return http.StatusOK
	case errors.ErrCodeInvalidRequest, errors.ErrCodeInvalidChartType, errors.ErrCodeInvalidChartConfig, errors.ErrCodeInvalidChartData, errors.ErrCodeInvalidSQL:
		return http.StatusBadRequest
	case errors.ErrCodeUnauthorized, errors.ErrCodeInvalidToken, errors.ErrCodeTokenExpired, errors.ErrCodeTokenNotValidYet, errors.ErrCodeTokenMissingClaims, errors.ErrCodeInvalidCredentials, errors.ErrCodeInvalidAPIKey, errors.ErrCodeAPIKeyExpired:
		return http.StatusUnauthorized
	case errors.ErrCodeForbidden:
		return http.StatusForbidden
	case errors.ErrCodeNotFound, errors.ErrCodeDataSourceNotFound, errors.ErrCodeFileNotFound:
		return http.StatusNotFound
	case errors.ErrCodeConflict, errors.ErrCodeUserExists:
		return http.StatusConflict
	case errors.ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case errors.ErrCodeTimeout, errors.ErrCodeQueryTimeout, errors.ErrCodeDatabaseTimeout, errors.ErrCodeCacheTimeout, errors.ErrCodeWebhookTimeout:
		return http.StatusRequestTimeout
	case errors.ErrCodeServiceUnavailable, errors.ErrCodeDatabaseConnection, errors.ErrCodeCacheConnection:
		return http.StatusServiceUnavailable
	case errors.ErrCodeInternalServer, errors.ErrCodeDatabaseQuery, errors.ErrCodeCacheFull, errors.ErrCodeFileUploadError, errors.ErrCodeWebhookDelivery:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return uuid.New().String()
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return generateRequestID()
}

// getStartTime 获取请求开始时间
func getStartTime(c *gin.Context) time.Time {
	if startTime, exists := c.Get("start_time"); exists {
		if t, ok := startTime.(time.Time); ok {
			return t
		}
	}
	return time.Now()
}

// ErrorLogger 错误日志中间件
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始
		startTime := time.Now()
		requestID := getRequestID(c)

		utils.Logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request started")

		c.Next()

		// 记录请求结束
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		logFields := map[string]interface{}{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			utils.Logger.WithFields(logFields).Error("Request failed with server error")
		} else if statusCode >= 400 {
			utils.Logger.WithFields(logFields).Warn("Request failed with client error")
		} else {
			utils.Logger.WithFields(logFields).Info("Request completed")
		}
	}
}

// ValidationErrorHandler 验证错误处理中间件
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				if errors.IsValidationError(ginErr.Err) {
					// 处理验证错误
					validationErr := errors.NewErrorWithSeverity(
						errors.ErrCodeInvalidRequest,
						"Validation failed",
						ginErr.Err,
						errors.SeverityMedium,
						errors.CategoryValidation,
					).WithContext(c.Request.Context())

					// 添加验证详情
					if validationErrors, ok := ginErr.Err.(validator.ValidationErrors); ok {
						details := make(map[string]interface{})
						for _, fieldError := range validationErrors {
							details[fieldError.Field()] = map[string]string{
								"tag":   fieldError.Tag(),
								"value": fieldError.Value().(string),
							}
						}
						validationErr.WithDetails(details)
					}

					c.Error(validationErr)
					break
				}
			}
		}
	}
}

// TimeoutErrorHandler 超时错误处理中间件
func TimeoutErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				if errors.IsTimeoutError(ginErr.Err) {
					// 处理超时错误
					timeoutErr := errors.NewTimeoutError("Request timeout", ginErr.Err).
						WithRetryAfter(30). // 30秒后重试
						WithContext(c.Request.Context())

					c.Error(timeoutErr)
					break
				}
			}
		}
	}
}

// RetryableErrorHandler 可重试错误处理中间件
func RetryableErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				if errors.IsRetryableError(ginErr.Err) {
					// 处理可重试错误
					if customErr, ok := ginErr.Err.(*errors.CustomError); ok {
						// 设置重试时间
						if customErr.RetryAfter == nil {
							customErr.WithRetryAfter(5) // 默认5秒后重试
						}
					}
					break
				}
			}
		}
	}
}
