package services

import (
	"gobi/internal/repositories"
	"gobi/internal/services/infrastructure"

	"gorm.io/gorm"
)

// ServiceFactory manages service dependencies
type ServiceFactory struct {
	db *gorm.DB
}

// NewServiceFactory creates a new ServiceFactory instance
func NewServiceFactory(db *gorm.DB) *ServiceFactory {
	return &ServiceFactory{db: db}
}

// CreateQueryService creates a QueryService with all dependencies
func (f *ServiceFactory) CreateQueryService() *QueryService {
	// Create repositories
	queryRepo := repositories.NewQueryRepository(f.db)

	// Create infrastructure services
	cacheService := infrastructure.NewCacheService()
	validationService := infrastructure.NewValidationService()
	sqlExecutionService := infrastructure.NewSQLExecutionService()
	encryptionService := infrastructure.NewEncryptionService()

	// Create and return QueryService
	return NewQueryService(
		queryRepo,
		cacheService,
		validationService,
		sqlExecutionService,
		encryptionService,
	)
}
