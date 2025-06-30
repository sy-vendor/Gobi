package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse 定义错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// CustomError 自定义错误类型
type CustomError struct {
	Code    int
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// 预定义错误
var (
	ErrInvalidRequest     = &CustomError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrUnauthorized       = &CustomError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden          = &CustomError{Code: http.StatusForbidden, Message: "Access denied"}
	ErrNotFound           = &CustomError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInternalServer     = &CustomError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrDatabaseOperation  = &CustomError{Code: http.StatusInternalServerError, Message: "Database operation failed"}
	ErrInvalidToken       = &CustomError{Code: http.StatusUnauthorized, Message: "Invalid token"}
	ErrTokenExpired       = &CustomError{Code: http.StatusUnauthorized, Message: "Token expired"}
	ErrTokenNotValidYet   = &CustomError{Code: http.StatusUnauthorized, Message: "Token not valid yet"}
	ErrTokenMissingClaims = &CustomError{Code: http.StatusUnauthorized, Message: "Token missing required claims"}
	ErrUserExists         = &CustomError{Code: http.StatusConflict, Message: "User already exists"}
	ErrInvalidCredentials = &CustomError{Code: http.StatusUnauthorized, Message: "Invalid credentials"}
)

// NewError 创建新的自定义错误
func NewError(code int, message string, err error) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// HandleError 处理错误并返回适当的响应
func HandleError(w http.ResponseWriter, err error) {
	var customErr *CustomError
	if e, ok := err.(*CustomError); ok {
		customErr = e
	} else {
		customErr = &CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}

	response := ErrorResponse{
		Code:    customErr.Code,
		Message: customErr.Message,
	}

	if customErr.Err != nil {
		response.Error = customErr.Err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(customErr.Code)
	json.NewEncoder(w).Encode(response)
}

// WrapError 包装错误，添加上下文信息
func WrapError(err error, message string) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		return &CustomError{
			Code:    customErr.Code,
			Message: message,
			Err:     customErr,
		}
	}
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

func NewBadRequestError(msg string, err error) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: msg,
		Err:     err,
	}
}

func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func IsContentTypeError(err error) bool {
	return err != nil && (err.Error() == "request Content-Type isn't multipart/form-data")
}

func NewConflictError(msg string, err error) *CustomError {
	return &CustomError{
		Code:    http.StatusConflict,
		Message: msg,
		Err:     err,
	}
}

var As = errors.As

// Additional common errors
var (
	ErrInvalidChartType           = NewError(400, "Invalid chart type", nil)
	ErrInvalidChartConfig         = NewError(400, "Invalid chart configuration", nil)
	ErrInvalidChartData           = NewError(400, "Invalid chart data", nil)
	ErrDataSourceNameRequired     = NewError(400, "DataSource name is required", nil)
	ErrDataSourceTypeRequired     = NewError(400, "DataSource type is required", nil)
	ErrDataSourceHostRequired     = NewError(400, "DataSource host is required", nil)
	ErrDataSourcePortRequired     = NewError(400, "DataSource port is required", nil)
	ErrDataSourceDatabaseRequired = NewError(400, "DataSource database is required", nil)
)
