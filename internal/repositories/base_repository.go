package repositories

import (
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// BaseRepository provides common repository operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new BaseRepository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// FindByID finds a record by ID with proper error handling
func (r *BaseRepository) FindByID(model interface{}, id uint) error {
	if err := r.db.First(model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "Could not find record")
	}
	return nil
}

// FindByUser finds records by user ID with admin permission check
func (r *BaseRepository) FindByUser(model interface{}, userID uint, isAdmin bool) error {
	query := r.db.Model(model)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(model).Error; err != nil {
		return errors.WrapError(err, "Could not find records")
	}
	return nil
}

// Create creates a new record with proper error handling
func (r *BaseRepository) Create(model interface{}) error {
	if err := r.db.Create(model).Error; err != nil {
		return errors.WrapError(err, "Could not create record")
	}
	return nil
}

// Update updates a record with proper error handling
func (r *BaseRepository) Update(model interface{}) error {
	if err := r.db.Save(model).Error; err != nil {
		return errors.WrapError(err, "Could not update record")
	}
	return nil
}

// Delete deletes a record by ID with proper error handling
func (r *BaseRepository) Delete(model interface{}, id uint) error {
	if err := r.db.Delete(model, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete record")
	}
	return nil
}

// CheckOwnership checks if a record belongs to a user
func (r *BaseRepository) CheckOwnership(model interface{}, id uint, userID uint) error {
	var count int64
	if err := r.db.Model(model).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
		return errors.WrapError(err, "Could not check ownership")
	}
	if count == 0 {
		return errors.ErrForbidden
	}
	return nil
}

// WrapDBError wraps database errors with consistent error handling
func (r *BaseRepository) WrapDBError(err error, operation string) error {
	if err == gorm.ErrRecordNotFound {
		return errors.ErrNotFound
	}
	return errors.WrapError(err, operation)
}
