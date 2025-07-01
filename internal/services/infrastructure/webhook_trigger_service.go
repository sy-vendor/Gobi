package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookTriggerServiceImpl implements WebhookTriggerService
type WebhookTriggerServiceImpl struct{}

// NewWebhookTriggerService creates a new WebhookTriggerService instance
func NewWebhookTriggerService() *WebhookTriggerServiceImpl {
	return &WebhookTriggerServiceImpl{}
}

// TriggerWebhook triggers a webhook
func (s *WebhookTriggerServiceImpl) TriggerWebhook(url string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Gobi-Webhook/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// ValidateWebhookURL validates a webhook URL
func (s *WebhookTriggerServiceImpl) ValidateWebhookURL(url string) error {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate webhook URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("webhook URL validation failed with status: %d", resp.StatusCode)
	}

	return nil
}
