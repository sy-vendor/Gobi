package services

import (
	"gobi/internal/models"
	"time"
)

// CacheService defines the interface for cache operations
type CacheService interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Flush()
	GetStats() map[string]interface{}
}

// EncryptionService defines the interface for encryption operations
type EncryptionService interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	GenerateJWT(userID uint, role string) (string, error)
	ValidateJWT(token string) (uint, string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

// PermissionService defines the interface for permission checks
type PermissionService interface {
	CanAccess(userID uint, resourceID uint, resourceType string, isAdmin bool) bool
	CheckOwnership(userID uint, resourceID uint, resourceType string) bool
	IsAdmin(userID uint) bool
}

// ValidationService defines the interface for validation operations
type ValidationService interface {
	ValidateSQL(sql string) error
	ValidateChartType(chartType string) error
	ValidateDataSource(ds *models.DataSource) error
	ValidateChartConfig(config string) error
	ValidateChartData(data string) error
}

// SQLExecutionService defines the interface for SQL execution
type SQLExecutionService interface {
	ExecuteSQL(ds models.DataSource, sql string) ([]map[string]interface{}, error)
	ExecuteSQLWithTimeout(ds models.DataSource, sql string, timeout time.Duration) ([]map[string]interface{}, error)
	ExecuteSQLWithLimit(ds models.DataSource, sql string, limit int) ([]map[string]interface{}, error)
}

// ReportGeneratorService defines the interface for report generation
type ReportGeneratorService interface {
	GenerateExcelFromTemplate(data string, template []byte, filename string) ([]byte, error)
	GeneratePDFFromTemplate(data string, template []byte, filename string) ([]byte, error)
}

// WebhookTriggerService defines the interface for webhook operations
type WebhookTriggerService interface {
	TriggerWebhook(url string, payload interface{}) error
	ValidateWebhookURL(url string) error
}
