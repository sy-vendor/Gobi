package services

import (
	"gobi/internal/repositories"

	"gorm.io/gorm"
)

// ServiceFactory manages service dependencies
// 现在支持注入所有基础设施服务接口
// 便于测试和扩展
type ServiceFactory struct {
	db                  *gorm.DB
	cacheService        CacheService
	validationService   ValidationService
	encryptionService   EncryptionService
	authService         AuthService
	permissionService   PermissionService
	sqlExecutionService SQLExecutionService
	reportGenerator     ReportGeneratorService
	webhookTrigger      WebhookTriggerService
}

// NewServiceFactory creates a new ServiceFactory instance
func NewServiceFactory(
	db *gorm.DB,
	cache CacheService,
	validation ValidationService,
	encryption EncryptionService,
	auth AuthService,
	permission PermissionService,
	sqlExec SQLExecutionService,
	reportGen ReportGeneratorService,
	webhookTrig WebhookTriggerService,
) *ServiceFactory {
	return &ServiceFactory{
		db:                  db,
		cacheService:        cache,
		validationService:   validation,
		encryptionService:   encryption,
		authService:         auth,
		permissionService:   permission,
		sqlExecutionService: sqlExec,
		reportGenerator:     reportGen,
		webhookTrigger:      webhookTrig,
	}
}

// CreateUserService creates a UserService with all dependencies
func (f *ServiceFactory) CreateUserService() *UserService {
	userRepo := repositories.NewUserRepository(f.db)
	return NewUserService(
		userRepo,
		f.cacheService,
		f.authService,
	)
}

// CreateQueryService creates a QueryService with all dependencies
func (f *ServiceFactory) CreateQueryService() *QueryService {
	queryRepo := repositories.NewQueryRepository(f.db)
	return NewQueryService(
		queryRepo,
		f.cacheService,
		f.validationService,
		f.sqlExecutionService,
		f.encryptionService,
	)
}

// CreateDataSourceService creates a DataSourceService with all dependencies
func (f *ServiceFactory) CreateDataSourceService() *DataSourceService {
	dsRepo := repositories.NewDataSourceRepository(f.db)
	return NewDataSourceService(
		dsRepo,
		f.encryptionService,
		f.validationService,
	)
}

// CreateChartService creates a ChartService with all dependencies
func (f *ServiceFactory) CreateChartService() *ChartService {
	chartRepo := repositories.NewChartRepository(f.db)
	queryService := f.CreateQueryService()
	return NewChartService(
		chartRepo,
		*queryService,
		f.cacheService,
	)
}

// CreateReportService creates a ReportService with all dependencies
func (f *ServiceFactory) CreateReportService() *ReportService {
	reportRepo := repositories.NewReportRepository(f.db)
	return NewReportService(
		reportRepo,
		f.reportGenerator,
		f.permissionService,
	)
}

// CreateTemplateService creates a TemplateService with all dependencies
func (f *ServiceFactory) CreateTemplateService() *TemplateService {
	templateRepo := repositories.NewTemplateRepository(f.db)
	return NewTemplateService(
		templateRepo,
		f.permissionService,
	)
}

// CreateWebhookService creates a WebhookService with all dependencies
func (f *ServiceFactory) CreateWebhookService() *WebhookService {
	webhookRepo := repositories.NewWebhookRepository(f.db)
	return NewWebhookService(
		webhookRepo,
		f.webhookTrigger,
		f.permissionService,
	)
}

// 你可以继续为其他 Service 添加类似的 CreateXXXService 方法
