package services

import (
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
)

// DataSourceService handles data source-related business logic
type DataSourceService struct {
	dsRepo            repositories.DataSourceRepository
	encryptionService EncryptionService
	validationService ValidationService
}

// NewDataSourceService creates a new DataSourceService instance
func NewDataSourceService(
	dsRepo repositories.DataSourceRepository,
	encryptionService EncryptionService,
	validationService ValidationService,
) *DataSourceService {
	return &DataSourceService{
		dsRepo:            dsRepo,
		encryptionService: encryptionService,
		validationService: validationService,
	}
}

// CreateDataSource creates a new data source with encrypted password
func (s *DataSourceService) CreateDataSource(ds *models.DataSource, userID uint) error {
	ds.UserID = userID

	// 验证必填字段
	if ds.Name == "" {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceNameRequired,
			"DataSource name is required",
			nil,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}
	if ds.Type == "" {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceTypeRequired,
			"DataSource type is required",
			nil,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}
	if ds.Host == "" {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceHostRequired,
			"DataSource host is required",
			nil,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}
	if ds.Port == 0 {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourcePortRequired,
			"DataSource port is required",
			nil,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}
	if ds.Database == "" {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceDatabaseRequired,
			"DataSource database is required",
			nil,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}

	if ds.Password != "" {
		encryptedPassword, err := s.encryptionService.Encrypt(ds.Password)
		if err != nil {
			return errors.NewErrorWithSeverity(
				errors.ErrCodeInternalServer,
				"Failed to encrypt password",
				err,
				errors.SeverityHigh,
				errors.CategorySecurity,
			)
		}
		ds.Password = encryptedPassword
	}

	if err := s.dsRepo.Create(ds); err != nil {
		return errors.NewDatabaseError("Could not create data source", err)
	}
	return nil
}

// ListDataSources retrieves data sources based on user permissions
func (s *DataSourceService) ListDataSources(userID uint, isAdmin bool) ([]models.DataSource, error) {
	dataSources, err := s.dsRepo.FindByUser(userID, isAdmin)
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch data sources")
	}
	for i := range dataSources {
		dataSources[i].Password = ""
	}
	return dataSources, nil
}

// GetDataSource retrieves a specific data source
func (s *DataSourceService) GetDataSource(dsID uint, userID uint, isAdmin bool) (*models.DataSource, error) {
	ds, err := s.dsRepo.FindByID(dsID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if !isAdmin && ds.UserID != userID && !ds.IsPublic {
		return nil, errors.ErrForbidden
	}
	ds.Password = ""
	return ds, nil
}

// UpdateDataSource updates a data source
func (s *DataSourceService) UpdateDataSource(dsID uint, updates *models.DataSource, userID uint, isAdmin bool) (*models.DataSource, error) {
	ds, err := s.dsRepo.FindByID(dsID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if !isAdmin && ds.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if updates.Name != "" {
		ds.Name = updates.Name
	}
	if updates.Type != "" {
		ds.Type = updates.Type
	}
	if updates.Host != "" {
		ds.Host = updates.Host
	}
	if updates.Port != 0 {
		ds.Port = updates.Port
	}
	if updates.Database != "" {
		ds.Database = updates.Database
	}
	if updates.Username != "" {
		ds.Username = updates.Username
	}
	if updates.Description != "" {
		ds.Description = updates.Description
	}
	ds.IsPublic = updates.IsPublic
	if updates.Password != "" {
		encryptedPassword, err := s.encryptionService.Encrypt(updates.Password)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to encrypt password")
		}
		ds.Password = encryptedPassword
	}
	if err := s.dsRepo.Update(ds); err != nil {
		return nil, errors.WrapError(err, "Could not update data source")
	}
	ds.Password = ""
	return ds, nil
}

// DeleteDataSource deletes a data source
func (s *DataSourceService) DeleteDataSource(dsID uint, userID uint, isAdmin bool) error {
	ds, err := s.dsRepo.FindByID(dsID)
	if err != nil {
		return errors.ErrNotFound
	}
	if !isAdmin && ds.UserID != userID {
		return errors.ErrForbidden
	}
	if err := s.dsRepo.Delete(dsID); err != nil {
		return errors.WrapError(err, "Could not delete data source")
	}
	return nil
}

// TestConnection tests the connection to a data source
func (s *DataSourceService) TestConnection(ds *models.DataSource) error {
	// Decrypt password
	if ds.Password != "" {
		decryptedPassword, err := utils.DecryptAES(ds.Password)
		if err != nil {
			return errors.NewErrorWithSeverity(
				errors.ErrCodeInternalServer,
				"Could not decrypt password",
				err,
				errors.SeverityHigh,
				errors.CategorySecurity,
			)
		}
		ds.Password = decryptedPassword
	}

	// Use the connection manager to test the connection
	db, err := database.GetConnection(ds)
	if err != nil {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceConnection,
			"Failed to create connection",
			err,
			errors.SeverityHigh,
			errors.CategoryDatabase,
		)
	}

	if err := db.Ping(); err != nil {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeDataSourceConnection,
			"Database connection test failed",
			err,
			errors.SeverityHigh,
			errors.CategoryDatabase,
		)
	}

	return nil
}
