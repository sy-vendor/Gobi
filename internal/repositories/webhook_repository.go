package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// WebhookRepositoryImpl implements WebhookRepository interface
type WebhookRepositoryImpl struct {
	db *gorm.DB
}

// NewWebhookRepository creates a new WebhookRepository instance
func NewWebhookRepository(db *gorm.DB) WebhookRepository {
	return &WebhookRepositoryImpl{db: db}
}

// Create creates a new webhook
func (r *WebhookRepositoryImpl) Create(webhook *models.Webhook) error {
	if err := r.db.Create(webhook).Error; err != nil {
		return errors.WrapError(err, "Could not create webhook")
	}
	return nil
}

// FindByID finds a webhook by ID
func (r *WebhookRepositoryImpl) FindByID(id uint) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := r.db.First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not find webhook")
	}
	return &webhook, nil
}

// FindByUser finds webhooks by user ID
func (r *WebhookRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	query := r.db
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&webhooks).Error; err != nil {
		return nil, errors.WrapError(err, "Could not find webhooks")
	}
	return webhooks, nil
}

// Update updates a webhook
func (r *WebhookRepositoryImpl) Update(webhook *models.Webhook) error {
	if err := r.db.Save(webhook).Error; err != nil {
		return errors.WrapError(err, "Could not update webhook")
	}
	return nil
}

// Delete deletes a webhook
func (r *WebhookRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.Webhook{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete webhook")
	}
	return nil
}

// CreateDelivery creates a webhook delivery record
func (r *WebhookRepositoryImpl) CreateDelivery(delivery *models.WebhookDelivery) error {
	if err := r.db.Create(delivery).Error; err != nil {
		return errors.WrapError(err, "Could not create webhook delivery")
	}
	return nil
}

// UpdateDelivery updates a webhook delivery record
func (r *WebhookRepositoryImpl) UpdateDelivery(delivery *models.WebhookDelivery) error {
	if err := r.db.Save(delivery).Error; err != nil {
		return errors.WrapError(err, "Could not update webhook delivery")
	}
	return nil
}

// ListDeliveries lists webhook delivery attempts
func (r *WebhookRepositoryImpl) ListDeliveries(webhookID uint) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	if err := r.db.Preload("Webhook").Where("webhook_id = ?", webhookID).
		Order("created_at DESC").Limit(100).Find(&deliveries).Error; err != nil {
		return nil, errors.WrapError(err, "Could not list webhook deliveries")
	}
	return deliveries, nil
}
