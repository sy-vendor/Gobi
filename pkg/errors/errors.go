package errors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ErrorCode 定义错误码类型
type ErrorCode string

// 错误码常量
const (
	// 通用错误码
	ErrCodeSuccess            ErrorCode = "SUCCESS"
	ErrCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeConflict           ErrorCode = "CONFLICT"
	ErrCodeInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout            ErrorCode = "TIMEOUT"
	ErrCodeRateLimit          ErrorCode = "RATE_LIMIT_EXCEEDED"

	// 认证相关错误码
	ErrCodeInvalidToken       ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenNotValidYet   ErrorCode = "TOKEN_NOT_VALID_YET"
	ErrCodeTokenMissingClaims ErrorCode = "TOKEN_MISSING_CLAIMS"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeUserExists         ErrorCode = "USER_EXISTS"
	ErrCodeInvalidAPIKey      ErrorCode = "INVALID_API_KEY"
	ErrCodeAPIKeyExpired      ErrorCode = "API_KEY_EXPIRED"

	// 数据库相关错误码
	ErrCodeDatabaseConnection ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseQuery      ErrorCode = "DATABASE_QUERY_ERROR"
	ErrCodeDatabaseTimeout    ErrorCode = "DATABASE_TIMEOUT"
	ErrCodeDatabaseConstraint ErrorCode = "DATABASE_CONSTRAINT_VIOLATION"

	// 业务逻辑错误码
	ErrCodeInvalidChartType   ErrorCode = "INVALID_CHART_TYPE"
	ErrCodeInvalidChartConfig ErrorCode = "INVALID_CHART_CONFIG"
	ErrCodeInvalidChartData   ErrorCode = "INVALID_CHART_DATA"
	ErrCodeInvalidSQL         ErrorCode = "INVALID_SQL"
	ErrCodeSQLInjection       ErrorCode = "SQL_INJECTION_DETECTED"
	ErrCodeQueryTimeout       ErrorCode = "QUERY_TIMEOUT"
	ErrCodeQueryLimitExceeded ErrorCode = "QUERY_LIMIT_EXCEEDED"

	// 数据源相关错误码
	ErrCodeDataSourceNotFound         ErrorCode = "DATASOURCE_NOT_FOUND"
	ErrCodeDataSourceConnection       ErrorCode = "DATASOURCE_CONNECTION_ERROR"
	ErrCodeDataSourceInvalid          ErrorCode = "DATASOURCE_INVALID"
	ErrCodeDataSourceNameRequired     ErrorCode = "DATASOURCE_NAME_REQUIRED"
	ErrCodeDataSourceTypeRequired     ErrorCode = "DATASOURCE_TYPE_REQUIRED"
	ErrCodeDataSourceHostRequired     ErrorCode = "DATASOURCE_HOST_REQUIRED"
	ErrCodeDataSourcePortRequired     ErrorCode = "DATASOURCE_PORT_REQUIRED"
	ErrCodeDataSourceDatabaseRequired ErrorCode = "DATASOURCE_DATABASE_REQUIRED"

	// 缓存相关错误码
	ErrCodeCacheConnection ErrorCode = "CACHE_CONNECTION_ERROR"
	ErrCodeCacheTimeout    ErrorCode = "CACHE_TIMEOUT"
	ErrCodeCacheFull       ErrorCode = "CACHE_FULL"

	// 文件相关错误码
	ErrCodeFileNotFound    ErrorCode = "FILE_NOT_FOUND"
	ErrCodeFileTooLarge    ErrorCode = "FILE_TOO_LARGE"
	ErrCodeFileInvalid     ErrorCode = "FILE_INVALID"
	ErrCodeFileUploadError ErrorCode = "FILE_UPLOAD_ERROR"

	// Webhook相关错误码
	ErrCodeWebhookInvalid   ErrorCode = "WEBHOOK_INVALID"
	ErrCodeWebhookTimeout   ErrorCode = "WEBHOOK_TIMEOUT"
	ErrCodeWebhookDelivery  ErrorCode = "WEBHOOK_DELIVERY_ERROR"
	ErrCodeWebhookSignature ErrorCode = "WEBHOOK_SIGNATURE_INVALID"
)

// ErrorSeverity 定义错误严重程度
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "LOW"
	SeverityMedium   ErrorSeverity = "MEDIUM"
	SeverityHigh     ErrorSeverity = "HIGH"
	SeverityCritical ErrorSeverity = "CRITICAL"
)

// ErrorCategory 定义错误分类
type ErrorCategory string

const (
	CategoryAuth       ErrorCategory = "AUTHENTICATION"
	CategoryAuthz      ErrorCategory = "AUTHORIZATION"
	CategoryDatabase   ErrorCategory = "DATABASE"
	CategoryBusiness   ErrorCategory = "BUSINESS_LOGIC"
	CategorySystem     ErrorCategory = "SYSTEM"
	CategoryNetwork    ErrorCategory = "NETWORK"
	CategorySecurity   ErrorCategory = "SECURITY"
	CategoryValidation ErrorCategory = "VALIDATION"
)

// ErrorResponse 定义错误响应结构
type ErrorResponse struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Error      string                 `json:"error,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Severity   ErrorSeverity          `json:"severity,omitempty"`
	Category   ErrorCategory          `json:"category,omitempty"`
	RetryAfter *int                   `json:"retry_after,omitempty"`
	HelpURL    string                 `json:"help_url,omitempty"`
}

// CustomError 自定义错误类型
type CustomError struct {
	Code       ErrorCode
	Message    string
	Err        error
	Severity   ErrorSeverity
	Category   ErrorCategory
	Details    map[string]interface{}
	RequestID  string
	Timestamp  time.Time
	RetryAfter *int
	HelpURL    string
	Stack      []string
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *CustomError) Unwrap() error {
	return e.Err
}

// Is 实现 errors.Is 接口
func (e *CustomError) Is(target error) bool {
	t, ok := target.(*CustomError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// As 实现 errors.As 接口
func (e *CustomError) As(target interface{}) bool {
	t, ok := target.(*CustomError)
	if !ok {
		return false
	}
	*t = *e
	return true
}

// 预定义错误
var (
	ErrInvalidRequest     = &CustomError{Code: ErrCodeInvalidRequest, Message: "Invalid request", Severity: SeverityMedium, Category: CategoryValidation}
	ErrUnauthorized       = &CustomError{Code: ErrCodeUnauthorized, Message: "Unauthorized", Severity: SeverityMedium, Category: CategoryAuth}
	ErrForbidden          = &CustomError{Code: ErrCodeForbidden, Message: "Access denied", Severity: SeverityMedium, Category: CategoryAuthz}
	ErrNotFound           = &CustomError{Code: ErrCodeNotFound, Message: "Resource not found", Severity: SeverityLow, Category: CategoryBusiness}
	ErrInternalServer     = &CustomError{Code: ErrCodeInternalServer, Message: "Internal server error", Severity: SeverityHigh, Category: CategorySystem}
	ErrDatabaseOperation  = &CustomError{Code: ErrCodeDatabaseQuery, Message: "Database operation failed", Severity: SeverityHigh, Category: CategoryDatabase}
	ErrInvalidToken       = &CustomError{Code: ErrCodeInvalidToken, Message: "Invalid token", Severity: SeverityMedium, Category: CategoryAuth}
	ErrTokenExpired       = &CustomError{Code: ErrCodeTokenExpired, Message: "Token expired", Severity: SeverityMedium, Category: CategoryAuth}
	ErrTokenNotValidYet   = &CustomError{Code: ErrCodeTokenNotValidYet, Message: "Token not valid yet", Severity: SeverityMedium, Category: CategoryAuth}
	ErrTokenMissingClaims = &CustomError{Code: ErrCodeTokenMissingClaims, Message: "Token missing required claims", Severity: SeverityMedium, Category: CategoryAuth}
	ErrUserExists         = &CustomError{Code: ErrCodeUserExists, Message: "User already exists", Severity: SeverityMedium, Category: CategoryBusiness}
	ErrInvalidCredentials = &CustomError{Code: ErrCodeInvalidCredentials, Message: "Invalid credentials", Severity: SeverityMedium, Category: CategoryAuth}
	ErrTimeout            = &CustomError{Code: ErrCodeTimeout, Message: "Request timeout", Severity: SeverityMedium, Category: CategorySystem}
	ErrRateLimit          = &CustomError{Code: ErrCodeRateLimit, Message: "Rate limit exceeded", Severity: SeverityMedium, Category: CategorySystem}
)

// NewError 创建新的自定义错误
func NewError(code ErrorCode, message string, err error) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewErrorWithDetails 创建带详细信息的错误
func NewErrorWithDetails(code ErrorCode, message string, err error, details map[string]interface{}) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Err:       err,
		Details:   details,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewErrorWithSeverity 创建带严重程度的错误
func NewErrorWithSeverity(code ErrorCode, message string, err error, severity ErrorSeverity, category ErrorCategory) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Err:       err,
		Severity:  severity,
		Category:  category,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// HandleError 处理错误并返回适当的响应
func HandleError(w http.ResponseWriter, err error) {
	var customErr *CustomError
	if e, ok := err.(*CustomError); ok {
		customErr = e
	} else {
		customErr = &CustomError{
			Code:      ErrCodeInternalServer,
			Message:   "Internal server error",
			Err:       err,
			Severity:  SeverityHigh,
			Category:  CategorySystem,
			Timestamp: time.Now(),
			Stack:     getStackTrace(),
		}
	}

	response := ErrorResponse{
		Code:       string(customErr.Code),
		Message:    customErr.Message,
		Details:    customErr.Details,
		RequestID:  customErr.RequestID,
		Timestamp:  customErr.Timestamp,
		Severity:   customErr.Severity,
		Category:   customErr.Category,
		RetryAfter: customErr.RetryAfter,
		HelpURL:    customErr.HelpURL,
	}

	if customErr.Err != nil {
		response.Error = customErr.Err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getHTTPStatus(customErr.Code))
	json.NewEncoder(w).Encode(response)
}

// WrapError 包装错误，添加上下文信息
func WrapError(err error, message string) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		return &CustomError{
			Code:      customErr.Code,
			Message:   message,
			Err:       customErr,
			Severity:  customErr.Severity,
			Category:  customErr.Category,
			Details:   customErr.Details,
			RequestID: customErr.RequestID,
			Timestamp: time.Now(),
			Stack:     getStackTrace(),
		}
	}
	return &CustomError{
		Code:      ErrCodeInternalServer,
		Message:   message,
		Err:       err,
		Severity:  SeverityHigh,
		Category:  CategorySystem,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// WrapErrorWithContext 包装错误并添加上下文
func WrapErrorWithContext(err error, message string, details map[string]interface{}) *CustomError {
	customErr := WrapError(err, message)
	customErr.Details = details
	return customErr
}

// NewBadRequestError 创建400错误
func NewBadRequestError(msg string, err error) *CustomError {
	return &CustomError{
		Code:      ErrCodeInvalidRequest,
		Message:   msg,
		Err:       err,
		Severity:  SeverityMedium,
		Category:  CategoryValidation,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewConflictError 创建409错误
func NewConflictError(msg string, err error) *CustomError {
	return &CustomError{
		Code:      ErrCodeConflict,
		Message:   msg,
		Err:       err,
		Severity:  SeverityMedium,
		Category:  CategoryBusiness,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(msg string, err error) *CustomError {
	return &CustomError{
		Code:      ErrCodeTimeout,
		Message:   msg,
		Err:       err,
		Severity:  SeverityMedium,
		Category:  CategorySystem,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(msg string, err error) *CustomError {
	return &CustomError{
		Code:      ErrCodeDatabaseQuery,
		Message:   msg,
		Err:       err,
		Severity:  SeverityHigh,
		Category:  CategoryDatabase,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewSecurityError 创建安全相关错误
func NewSecurityError(msg string, err error) *CustomError {
	return &CustomError{
		Code:      ErrCodeSQLInjection,
		Message:   msg,
		Err:       err,
		Severity:  SeverityCritical,
		Category:  CategorySecurity,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

// IsContentTypeError 检查是否为内容类型错误
func IsContentTypeError(err error) bool {
	return err != nil && (err.Error() == "request Content-Type isn't multipart/form-data")
}

// IsTimeoutError 检查是否为超时错误
func IsTimeoutError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code == ErrCodeTimeout
	}
	return strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded")
}

// IsRetryableError 检查是否为可重试错误
func IsRetryableError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		switch customErr.Code {
		case ErrCodeDatabaseConnection, ErrCodeDatabaseTimeout, ErrCodeCacheConnection, ErrCodeCacheTimeout, ErrCodeWebhookTimeout:
			return true
		}
	}
	return IsTimeoutError(err)
}

// GetRetryAfter 获取重试时间
func GetRetryAfter(err error) *int {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.RetryAfter
	}
	return nil
}

// SetRequestID 设置请求ID
func SetRequestID(err error, requestID string) {
	if customErr, ok := err.(*CustomError); ok {
		customErr.RequestID = requestID
	}
}

// SetRetryAfter 设置重试时间
func SetRetryAfter(err error, seconds int) {
	if customErr, ok := err.(*CustomError); ok {
		customErr.RetryAfter = &seconds
	}
}

// SetHelpURL 设置帮助URL
func SetHelpURL(err error, helpURL string) {
	if customErr, ok := err.(*CustomError); ok {
		customErr.HelpURL = helpURL
	}
}

// getHTTPStatus 根据错误码获取HTTP状态码
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeSuccess:
		return http.StatusOK
	case ErrCodeInvalidRequest, ErrCodeInvalidChartType, ErrCodeInvalidChartConfig, ErrCodeInvalidChartData, ErrCodeInvalidSQL:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired, ErrCodeTokenNotValidYet, ErrCodeTokenMissingClaims, ErrCodeInvalidCredentials, ErrCodeInvalidAPIKey, ErrCodeAPIKeyExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeDataSourceNotFound, ErrCodeFileNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeUserExists:
		return http.StatusConflict
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeTimeout, ErrCodeQueryTimeout, ErrCodeDatabaseTimeout, ErrCodeCacheTimeout, ErrCodeWebhookTimeout:
		return http.StatusRequestTimeout
	case ErrCodeServiceUnavailable, ErrCodeDatabaseConnection, ErrCodeCacheConnection:
		return http.StatusServiceUnavailable
	case ErrCodeInternalServer, ErrCodeDatabaseQuery, ErrCodeCacheFull, ErrCodeFileUploadError, ErrCodeWebhookDelivery:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// getStackTrace 获取堆栈跟踪
func getStackTrace() []string {
	var stack []string
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	return stack
}

// WithContext 添加上下文信息
func (e *CustomError) WithContext(ctx context.Context) *CustomError {
	if ctx != nil {
		if requestID, ok := ctx.Value("request_id").(string); ok {
			e.RequestID = requestID
		}
	}
	return e
}

// WithRetryAfter 添加重试时间
func (e *CustomError) WithRetryAfter(seconds int) *CustomError {
	e.RetryAfter = &seconds
	return e
}

// WithHelpURL 添加帮助URL
func (e *CustomError) WithHelpURL(helpURL string) *CustomError {
	e.HelpURL = helpURL
	return e
}

// WithDetails 添加详细信息
func (e *CustomError) WithDetails(details map[string]interface{}) *CustomError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// AddDetail 添加单个详细信息
func (e *CustomError) AddDetail(key string, value interface{}) *CustomError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// 兼容性函数
var As = errors.As

// Additional common errors
var (
	ErrInvalidChartType           = NewError(ErrCodeInvalidChartType, "Invalid chart type", nil)
	ErrInvalidChartConfig         = NewError(ErrCodeInvalidChartConfig, "Invalid chart configuration", nil)
	ErrInvalidChartData           = NewError(ErrCodeInvalidChartData, "Invalid chart data", nil)
	ErrDataSourceNameRequired     = NewError(ErrCodeDataSourceNameRequired, "DataSource name is required", nil)
	ErrDataSourceTypeRequired     = NewError(ErrCodeDataSourceTypeRequired, "DataSource type is required", nil)
	ErrDataSourceHostRequired     = NewError(ErrCodeDataSourceHostRequired, "DataSource host is required", nil)
	ErrDataSourcePortRequired     = NewError(ErrCodeDataSourcePortRequired, "DataSource port is required", nil)
	ErrDataSourceDatabaseRequired = NewError(ErrCodeDataSourceDatabaseRequired, "DataSource database is required", nil)
)
