package handlers

import (
	"gobi/internal/models"
	"gobi/internal/services"
	"gobi/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	DB             *gorm.DB
	WebhookService *services.WebhookService
}

// NewWebhookHandler creates a new WebhookHandler instance
func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{
		DB:             db,
		WebhookService: services.NewWebhookService(db),
	}
}

// CreateWebhook handles webhook creation
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	var webhook models.Webhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook data", err))
		return
	}

	userID, _ := c.Get("userID")
	if err := h.WebhookService.CreateWebhook(&webhook, userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

// ListWebhooks handles webhook listing
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	webhooks, err := h.WebhookService.ListWebhooks(userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, webhooks)
}

// GetWebhook handles getting a specific webhook
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	id := c.Param("id")
	webhookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	webhook, err := h.WebhookService.GetWebhook(uint(webhookID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// UpdateWebhook handles webhook updates
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	id := c.Param("id")
	webhookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook ID", err))
		return
	}

	var req models.Webhook
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook data", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	webhook, err := h.WebhookService.UpdateWebhook(uint(webhookID), &req, userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// DeleteWebhook handles webhook deletion
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	id := c.Param("id")
	webhookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.WebhookService.DeleteWebhook(uint(webhookID), userID.(uint), isAdmin); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook deleted successfully"})
}

// ListWebhookDeliveries handles listing webhook delivery attempts
func (h *WebhookHandler) ListWebhookDeliveries(c *gin.Context) {
	id := c.Param("id")
	webhookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	deliveries, err := h.WebhookService.ListWebhookDeliveries(uint(webhookID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

// TestWebhook handles webhook testing
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	id := c.Param("id")
	webhookID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid webhook ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	webhook, err := h.WebhookService.GetWebhook(uint(webhookID), userID.(uint), isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	// Send test payload
	testPayload := map[string]interface{}{
		"event": "webhook.test",
		"data": map[string]interface{}{
			"message":    "This is a test webhook from Gobi",
			"timestamp":  webhook.CreatedAt,
			"webhook_id": webhook.ID,
		},
	}

	if err := h.WebhookService.TriggerWebhook("webhook.test", testPayload, userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test webhook sent successfully"})
}
