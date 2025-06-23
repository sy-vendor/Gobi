package services

import (
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"gobi/pkg/utils"

	"gorm.io/gorm"
)

// DataSourceService handles data source-related business logic
type DataSourceService struct {
	db *gorm.DB
}

// NewDataSourceService creates a new DataSourceService instance
func NewDataSourceService(db *gorm.DB) *DataSourceService {
	return &DataSourceService{db: db}
}

// CreateDataSource creates a new data source with encrypted password
func (s *DataSourceService) CreateDataSource(ds *models.DataSource, userID uint) error {
	ds.UserID = userID

	// Encrypt password before saving
	if ds.Password != "" {
		encryptedPassword, err := utils.EncryptAES(ds.Password)
		if err != nil {
			return errors.WrapError(err, "Failed to encrypt password")
		}
		ds.Password = encryptedPassword
	}

	if err := s.db.Create(ds).Error; err != nil {
		return errors.WrapError(err, "Could not create data source")
	}

	return nil
}

// ListDataSources retrieves data sources based on user permissions
func (s *DataSourceService) ListDataSources(userID uint, isAdmin bool) ([]models.DataSource, error) {
	var dataSources []models.DataSource

	query := s.db.Model(&models.DataSource{})
	if !isAdmin {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&dataSources).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch data sources")
	}

	// Clear passwords before sending
	for i := range dataSources {
		dataSources[i].Password = ""
	}

	return dataSources, nil
}

// GetDataSource retrieves a specific data source
func (s *DataSourceService) GetDataSource(dsID uint, userID uint, isAdmin bool) (*models.DataSource, error) {
	var ds models.DataSource
	if err := s.db.First(&ds, dsID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && ds.UserID != userID && !ds.IsPublic {
		return nil, errors.ErrForbidden
	}

	ds.Password = "" // Never return password
	return &ds, nil
}

// UpdateDataSource updates a data source
func (s *DataSourceService) UpdateDataSource(dsID uint, updates *models.DataSource, userID uint, isAdmin bool) (*models.DataSource, error) {
	var ds models.DataSource
	if err := s.db.First(&ds, dsID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && ds.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update allowed fields
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
	// Update IsPublic if provided
	ds.IsPublic = updates.IsPublic

	// Handle password update
	if updates.Password != "" {
		encryptedPassword, err := utils.EncryptAES(updates.Password)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to encrypt password")
		}
		ds.Password = encryptedPassword
	}

	if err := s.db.Save(&ds).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update data source")
	}

	ds.Password = "" // Never return password
	return &ds, nil
}

// DeleteDataSource deletes a data source
func (s *DataSourceService) DeleteDataSource(dsID uint, userID uint, isAdmin bool) error {
	var ds models.DataSource
	if err := s.db.First(&ds, dsID).Error; err != nil {
		return errors.ErrNotFound
	}

	if !isAdmin && ds.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&ds).Error; err != nil {
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
			return errors.WrapError(err, "Could not decrypt password")
		}
		ds.Password = decryptedPassword
	}

	// Use the connection manager to test the connection
	db, err := database.GetConnection(ds)
	if err != nil {
		return errors.WrapError(err, "Failed to create connection")
	}

	if err := db.Ping(); err != nil {
		return errors.WrapError(err, "Database connection test failed")
	}

	return nil
}
