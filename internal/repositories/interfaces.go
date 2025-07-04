package repositories

import (
	"gobi/internal/models"
)

// QueryRepository defines the interface for query data access
type QueryRepository interface {
	Create(query *models.Query) error
	FindByID(id uint) (*models.Query, error)
	FindByUser(userID uint, isAdmin bool) ([]models.Query, error)
	Update(query *models.Query) error
	Delete(id uint) error
	IncrementExecCount(id uint) error
}

// ChartRepository defines the interface for chart data access
type ChartRepository interface {
	Create(chart *models.Chart) error
	FindByID(id uint) (*models.Chart, error)
	FindByUser(userID uint, isAdmin bool) ([]models.Chart, error)
	Update(chart *models.Chart) error
	Delete(id uint) error
}

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	Create(report *models.Report) error
	FindByID(id uint) (*models.Report, error)
	FindByUser(userID uint, isAdmin bool) ([]models.Report, error)
	Update(report *models.Report) error
	Delete(id uint) error
	UpdateStatus(reportID uint, status string, error string) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	UpdateLastLogin(id uint) error
}

// DataSourceRepository defines the interface for datasource data access
type DataSourceRepository interface {
	Create(ds *models.DataSource) error
	FindByID(id uint) (*models.DataSource, error)
	FindByUser(userID uint, isAdmin bool) ([]models.DataSource, error)
	Update(ds *models.DataSource) error
	Delete(id uint) error
	TestConnection(ds *models.DataSource) error
}

// APIKeyRepository defines the interface for API key data access
type APIKeyRepository interface {
	Create(key *models.APIKey) error
	FindByID(id uint) (*models.APIKey, error)
	FindByPrefix(prefix string) (*models.APIKey, error)
	FindByUser(userID uint, isAdmin bool) ([]models.APIKey, error)
	Update(key *models.APIKey) error
	Delete(id uint) error
	Revoke(keyID uint, userID uint, isAdmin bool) error
}

// TemplateRepository defines the interface for template data access
type TemplateRepository interface {
	Create(template *models.ExcelTemplate) error
	FindByID(id uint) (*models.ExcelTemplate, error)
	FindByUser(userID uint, isAdmin bool) ([]models.ExcelTemplate, error)
	Update(template *models.ExcelTemplate) error
	Delete(id uint) error
	GetStats() (map[string]interface{}, error)
}

// WebhookRepository defines the interface for webhook data access
type WebhookRepository interface {
	Create(webhook *models.Webhook) error
	FindByID(id uint) (*models.Webhook, error)
	FindByUser(userID uint, isAdmin bool) ([]models.Webhook, error)
	Update(webhook *models.Webhook) error
	Delete(id uint) error
	CreateDelivery(delivery *models.WebhookDelivery) error
	UpdateDelivery(delivery *models.WebhookDelivery) error
	ListDeliveries(webhookID uint) ([]models.WebhookDelivery, error)
}
