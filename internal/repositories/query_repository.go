package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// QueryRepositoryImpl implements QueryRepository
type QueryRepositoryImpl struct {
	db *gorm.DB
}

// NewQueryRepository creates a new QueryRepository instance
func NewQueryRepository(db *gorm.DB) QueryRepository {
	return &QueryRepositoryImpl{db: db}
}

// Create creates a new query
func (r *QueryRepositoryImpl) Create(query *models.Query) error {
	if err := r.db.Create(query).Error; err != nil {
		return errors.WrapError(err, "Could not create query")
	}
	return nil
}

// FindByID finds a query by ID
func (r *QueryRepositoryImpl) FindByID(id uint) (*models.Query, error) {
	var query models.Query
	if err := r.db.Preload("DataSource").Preload("User").First(&query, id).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &query, nil
}

// FindByUser finds queries by user with permission check
func (r *QueryRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.Query, error) {
	var queries []models.Query

	query := r.db.Preload("DataSource").Preload("User").Model(&models.Query{})
	if !isAdmin {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&queries).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch queries")
	}

	return queries, nil
}

// Update updates a query
func (r *QueryRepositoryImpl) Update(query *models.Query) error {
	if err := r.db.Save(query).Error; err != nil {
		return errors.WrapError(err, "Could not update query")
	}
	return nil
}

// Delete deletes a query
func (r *QueryRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.Query{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete query")
	}
	return nil
}

// IncrementExecCount increments the execution count of a query
func (r *QueryRepositoryImpl) IncrementExecCount(id uint) error {
	if err := r.db.Model(&models.Query{}).Where("id = ?", id).UpdateColumn("exec_count", gorm.Expr("exec_count + ?", 1)).Error; err != nil {
		return errors.WrapError(err, "Could not increment exec count")
	}
	return nil
}
