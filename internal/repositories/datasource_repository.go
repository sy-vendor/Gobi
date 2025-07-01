package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// DataSourceRepositoryImpl implements DataSourceRepository interface
type DataSourceRepositoryImpl struct {
	db *gorm.DB
}

// NewDataSourceRepository creates a new DataSourceRepository instance
func NewDataSourceRepository(db *gorm.DB) DataSourceRepository {
	return &DataSourceRepositoryImpl{db: db}
}

// Create creates a new datasource
func (r *DataSourceRepositoryImpl) Create(ds *models.DataSource) error {
	if err := r.db.Create(ds).Error; err != nil {
		return errors.WrapError(err, "Could not create datasource")
	}
	return nil
}

// FindByID finds a datasource by ID
func (r *DataSourceRepositoryImpl) FindByID(id uint) (*models.DataSource, error) {
	var ds models.DataSource
	if err := r.db.Preload("User").First(&ds, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not find datasource")
	}
	return &ds, nil
}

// FindByUser finds datasources by user ID
func (r *DataSourceRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.DataSource, error) {
	var datasources []models.DataSource
	query := r.db.Preload("User").Model(&models.DataSource{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&datasources).Error; err != nil {
		return nil, errors.WrapError(err, "Could not find datasources")
	}
	return datasources, nil
}

// Update updates a datasource
func (r *DataSourceRepositoryImpl) Update(ds *models.DataSource) error {
	if err := r.db.Save(ds).Error; err != nil {
		return errors.WrapError(err, "Could not update datasource")
	}
	return nil
}

// Delete deletes a datasource
func (r *DataSourceRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.DataSource{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete datasource")
	}
	return nil
}

// TestConnection tests the connection to a datasource
func (r *DataSourceRepositoryImpl) TestConnection(ds *models.DataSource) error {
	// This would typically test the database connection
	// For now, we'll just return success
	return nil
}
