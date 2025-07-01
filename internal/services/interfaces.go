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

// ReportRepository defines the interface for report data operations
type ReportRepository interface {
	Create(report *models.Report) error
	FindByUser(userID uint, isAdmin bool) ([]models.Report, error)
	FindByID(reportID uint) (*models.Report, error)
	Update(report *models.Report) error
	Delete(reportID uint) error
	UpdateStatus(reportID uint, status string, error string) error
}

// TemplateRepository defines the interface for template data operations
type TemplateRepository interface {
	Create(template *models.ExcelTemplate) error
	FindByUser(userID uint, isAdmin bool) ([]models.ExcelTemplate, error)
	FindByID(templateID uint) (*models.ExcelTemplate, error)
	Update(template *models.ExcelTemplate) error
	Delete(templateID uint) error
	GetStats() (map[string]interface{}, error)
}

// WebhookRepository defines the interface for webhook data operations
type WebhookRepository interface {
	Create(webhook *models.Webhook) error
	FindByUser(userID uint, isAdmin bool) ([]models.Webhook, error)
	FindByID(webhookID uint) (*models.Webhook, error)
	Update(webhook *models.Webhook) error
	Delete(webhookID uint) error
	CreateDelivery(delivery *models.WebhookDelivery) error
	UpdateDelivery(delivery *models.WebhookDelivery) error
	ListDeliveries(webhookID uint) ([]models.WebhookDelivery, error)
}
