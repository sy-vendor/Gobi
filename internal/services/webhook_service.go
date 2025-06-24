package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/errors"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// WebhookService handles webhook-related business logic
type WebhookService struct {
	db *gorm.DB
}

// NewWebhookService creates a new WebhookService instance
func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{db: db}
}

// CreateWebhook creates a new webhook configuration
func (s *WebhookService) CreateWebhook(webhook *models.Webhook, userID uint) error {
	webhook.UserID = userID
	webhook.CreatedAt = time.Now()
	webhook.UpdatedAt = time.Now()

	// Generate secret for signature verification
	webhook.Secret = generateWebhookSecret()

	if err := s.db.Create(webhook).Error; err != nil {
		return errors.WrapError(err, "Could not create webhook")
	}
	return nil
}

// ListWebhooks lists all webhooks for a user (admin can list all)
func (s *WebhookService) ListWebhooks(userID uint, isAdmin bool) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	query := s.db
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&webhooks).Error; err != nil {
		return nil, errors.WrapError(err, "Could not list webhooks")
	}
	return webhooks, nil
}

// GetWebhook gets a specific webhook by ID
func (s *WebhookService) GetWebhook(webhookID uint, userID uint, isAdmin bool) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := s.db.First(&webhook, webhookID).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	if !isAdmin && webhook.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return &webhook, nil
}

// UpdateWebhook updates a webhook configuration
func (s *WebhookService) UpdateWebhook(webhookID uint, updates *models.Webhook, userID uint, isAdmin bool) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := s.db.First(&webhook, webhookID).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	if !isAdmin && webhook.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update allowed fields
	if updates.Name != "" {
		webhook.Name = updates.Name
	}
	if updates.URL != "" {
		webhook.URL = updates.URL
	}
	if updates.Events != "" {
		webhook.Events = updates.Events
	}
	if updates.Headers != "" {
		webhook.Headers = updates.Headers
	}
	webhook.Active = updates.Active
	webhook.UpdatedAt = time.Now()

	if err := s.db.Save(&webhook).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update webhook")
	}
	return &webhook, nil
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(webhookID uint, userID uint, isAdmin bool) error {
	var webhook models.Webhook
	if err := s.db.First(&webhook, webhookID).Error; err != nil {
		return errors.ErrNotFound
	}
	if !isAdmin && webhook.UserID != userID {
		return errors.ErrForbidden
	}
	return s.db.Delete(&webhook).Error
}

// TriggerWebhook sends a webhook notification for a specific event
func (s *WebhookService) TriggerWebhook(event string, payload interface{}, userID uint) error {
	// Get all active webhooks for the user that subscribe to this event
	var webhooks []models.Webhook
	if err := s.db.Where("user_id = ? AND active = ?", userID, true).Find(&webhooks).Error; err != nil {
		return errors.WrapError(err, "Could not fetch webhooks")
	}

	for _, webhook := range webhooks {
		// Check if webhook subscribes to this event
		if !s.webhookSubscribesToEvent(&webhook, event) {
			continue
		}

		// Create delivery record
		delivery := &models.WebhookDelivery{
			WebhookID: webhook.ID,
			Event:     event,
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		// Serialize payload
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			delivery.Status = "failed"
			delivery.Response = fmt.Sprintf("Failed to serialize payload: %v", err)
			s.db.Create(delivery)
			continue
		}
		delivery.Payload = string(payloadBytes)

		// Save delivery record
		if err := s.db.Create(delivery).Error; err != nil {
			continue
		}

		// Send webhook asynchronously
		go s.sendWebhook(&webhook, delivery)
	}

	return nil
}

// sendWebhook sends a webhook notification with retry logic
func (s *WebhookService) sendWebhook(webhook *models.Webhook, delivery *models.WebhookDelivery) {
	maxRetries := 3
	retryDelays := []time.Duration{5 * time.Second, 30 * time.Second, 5 * time.Minute}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		success := s.sendSingleWebhook(webhook, delivery)
		if success {
			return
		}

		if attempt < maxRetries {
			time.Sleep(retryDelays[attempt])
		}
	}

	// Final failure
	delivery.Status = "failed"
	delivery.Attempts = maxRetries + 1
	s.db.Save(delivery)
}

// sendSingleWebhook sends a single webhook request
func (s *WebhookService) sendSingleWebhook(webhook *models.Webhook, delivery *models.WebhookDelivery) bool {
	// Parse custom headers
	var headers map[string]string
	if webhook.Headers != "" {
		if err := json.Unmarshal([]byte(webhook.Headers), &headers); err != nil {
			headers = make(map[string]string)
		}
	} else {
		headers = make(map[string]string)
	}

	// Set default headers
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = "Gobi-Webhook/1.0"

	// Generate signature
	timestamp := time.Now().Unix()
	signature := s.generateSignature(webhook.Secret, delivery.Payload, timestamp)
	headers["X-Gobi-Signature"] = signature
	headers["X-Gobi-Timestamp"] = fmt.Sprintf("%d", timestamp)
	headers["X-Gobi-Event"] = delivery.Event

	// Create HTTP request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBufferString(delivery.Payload))
	if err != nil {
		delivery.Response = fmt.Sprintf("Failed to create request: %v", err)
		return false
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		delivery.Response = fmt.Sprintf("Request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		delivery.Response = fmt.Sprintf("Failed to read response: %v", err)
		return false
	}

	// Update delivery record
	delivery.Status = "success"
	if resp.StatusCode >= 400 {
		delivery.Status = "failed"
	}
	delivery.Response = string(body)
	delivery.Attempts++
	now := time.Now()
	delivery.SentAt = &now

	s.db.Save(delivery)

	return resp.StatusCode < 400
}

// webhookSubscribesToEvent checks if a webhook subscribes to a specific event
func (s *WebhookService) webhookSubscribesToEvent(webhook *models.Webhook, event string) bool {
	if webhook.Events == "" {
		return false
	}

	var events []string
	if err := json.Unmarshal([]byte(webhook.Events), &events); err != nil {
		return false
	}

	for _, e := range events {
		if e == event || e == "*" {
			return true
		}
	}
	return false
}

// generateSignature generates HMAC-SHA256 signature for webhook payload
func (s *WebhookService) generateSignature(secret, payload string, timestamp int64) string {
	message := fmt.Sprintf("%d.%s", timestamp, payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// generateWebhookSecret generates a random webhook secret
func generateWebhookSecret() string {
	// Simple implementation - in production, use crypto/rand
	return fmt.Sprintf("whsec_%d", time.Now().UnixNano())
}

// ListWebhookDeliveries lists webhook delivery attempts
func (s *WebhookService) ListWebhookDeliveries(webhookID uint, userID uint, isAdmin bool) ([]models.WebhookDelivery, error) {
	var deliveries []models.WebhookDelivery
	query := s.db.Preload("Webhook").Where("webhook_id = ?", webhookID)

	if !isAdmin {
		query = query.Joins("JOIN webhooks ON webhook_deliveries.webhook_id = webhooks.id").
			Where("webhooks.user_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Limit(100).Find(&deliveries).Error; err != nil {
		return nil, errors.WrapError(err, "Could not list webhook deliveries")
	}
	return deliveries, nil
}
