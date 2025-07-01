package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// ChartRepositoryImpl implements ChartRepository interface
type ChartRepositoryImpl struct {
	db *gorm.DB
}

// NewChartRepository creates a new ChartRepository instance
func NewChartRepository(db *gorm.DB) ChartRepository {
	return &ChartRepositoryImpl{db: db}
}

// Create creates a new chart
func (r *ChartRepositoryImpl) Create(chart *models.Chart) error {
	if err := r.db.Create(chart).Error; err != nil {
		return errors.WrapError(err, "Could not create chart")
	}
	return nil
}

// FindByID finds a chart by ID
func (r *ChartRepositoryImpl) FindByID(id uint) (*models.Chart, error) {
	var chart models.Chart
	if err := r.db.Preload("User").Preload("Query").First(&chart, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not find chart")
	}
	return &chart, nil
}

// FindByUser finds charts by user ID
func (r *ChartRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.Chart, error) {
	var charts []models.Chart
	query := r.db.Preload("User").Preload("Query").Model(&models.Chart{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&charts).Error; err != nil {
		return nil, errors.WrapError(err, "Could not find charts")
	}
	return charts, nil
}

// Update updates a chart
func (r *ChartRepositoryImpl) Update(chart *models.Chart) error {
	if err := r.db.Save(chart).Error; err != nil {
		return errors.WrapError(err, "Could not update chart")
	}
	return nil
}

// Delete deletes a chart
func (r *ChartRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.Chart{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete chart")
	}
	return nil
}
