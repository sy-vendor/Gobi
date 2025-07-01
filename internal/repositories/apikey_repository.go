package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// 只保留实现部分，接口定义在 interfaces.go

type APIKeyRepositoryImpl struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &APIKeyRepositoryImpl{db: db}
}

func (r *APIKeyRepositoryImpl) Create(key *models.APIKey) error {
	return r.db.Create(key).Error
}

func (r *APIKeyRepositoryImpl) FindByID(id uint) (*models.APIKey, error) {
	var key models.APIKey
	if err := r.db.First(&key, id).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &key, nil
}

func (r *APIKeyRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.APIKey, error) {
	var keys []models.APIKey
	query := r.db
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&keys).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch API keys")
	}
	return keys, nil
}

func (r *APIKeyRepositoryImpl) FindByPrefix(prefix string) (*models.APIKey, error) {
	var key models.APIKey
	if err := r.db.Where("prefix = ?", prefix).First(&key).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &key, nil
}

func (r *APIKeyRepositoryImpl) Revoke(keyID uint, userID uint, isAdmin bool) error {
	query := r.db.Model(&models.APIKey{}).Where("id = ?", keyID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	return query.Update("revoked", true).Error
}

func (r *APIKeyRepositoryImpl) Update(key *models.APIKey) error {
	return r.db.Save(key).Error
}

func (r *APIKeyRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.APIKey{}, id).Error
}
