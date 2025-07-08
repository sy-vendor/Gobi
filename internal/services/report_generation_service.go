package services

import (
	"encoding/json"
	errs "errors"
	"fmt"
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ReportGenerationService handles report generation business logic
type ReportGenerationService struct {
	db             *gorm.DB
	webhookRepo    repositories.WebhookRepository
	webhookTrigger WebhookTriggerService
}

// NewReportGenerationService creates a new ReportGenerationService instance
func NewReportGenerationService(
	db *gorm.DB,
	webhookRepo repositories.WebhookRepository,
	webhookTrigger WebhookTriggerService,
) *ReportGenerationService {
	return &ReportGenerationService{
		db:             db,
		webhookRepo:    webhookRepo,
		webhookTrigger: webhookTrigger,
	}
}

// GenerateExcelReport generates an Excel report from chart data and template
func (s *ReportGenerationService) GenerateExcelReport(chartID uint, templateID uint, userID uint, isAdmin bool) ([]byte, error) {
	// Get chart
	var chart models.Chart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch chart")
	}

	// Permission check
	if !isAdmin && chart.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Get template
	var template models.ExcelTemplate
	if err := s.db.First(&template, templateID).Error; err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch template")
	}

	// Permission check for template
	if !isAdmin && template.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Generate Excel report
	return utils.GenerateExcelFromTemplate(chart.Data, template.Template, strconv.Itoa(int(chart.ID)))
}

// GeneratePDFReport generates a PDF report from chart data
func (s *ReportGenerationService) GeneratePDFReport(chartID uint, userID uint, isAdmin bool) ([]byte, error) {
	// Get chart
	var chart models.Chart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check
	if !isAdmin && chart.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// For now, return a simple PDF content
	// In a real implementation, you would use a PDF generation library
	pdfContent := fmt.Sprintf("PDF Report for Chart: %s\nGenerated at: %s\nData: %s",
		chart.Name, time.Now().Format("2006-01-02 15:04:05"), chart.Data)

	return []byte(pdfContent), nil
}

// GenerateScheduledReport generates a report based on a schedule
func (s *ReportGenerationService) GenerateScheduledReport(schedule *models.ReportSchedule) error {
	// Create a new report record
	report := models.Report{
		UserID:      schedule.UserID,
		Name:        schedule.Name,
		Type:        schedule.Type,
		Status:      "generating",
		GeneratedAt: time.Now(),
	}

	if err := s.db.Create(&report).Error; err != nil {
		return errors.WrapError(err, "Could not create report record")
	}

	// Parse queries, charts, and templates
	var queryIDs []uint
	var chartIDs []uint
	var templateIDs []uint

	if schedule.Queries != "" {
		if err := json.Unmarshal([]byte(schedule.Queries), &queryIDs); err != nil {
			report.Status = "failed"
			report.Error = "Invalid queries configuration"
			s.db.Save(&report)
			s.triggerReportWebhooks(schedule, &report, errors.NewError(errors.ErrCodeInvalidRequest, "Invalid queries configuration", err))
			return errors.WrapError(err, "Invalid queries configuration")
		}
	}

	if schedule.Charts != "" {
		if err := json.Unmarshal([]byte(schedule.Charts), &chartIDs); err != nil {
			report.Status = "failed"
			report.Error = "Invalid charts configuration"
			s.db.Save(&report)
			s.triggerReportWebhooks(schedule, &report, errors.NewError(errors.ErrCodeInvalidRequest, "Invalid charts configuration", err))
			return errors.WrapError(err, "Invalid charts configuration")
		}
	}

	if schedule.Templates != "" {
		if err := json.Unmarshal([]byte(schedule.Templates), &templateIDs); err != nil {
			report.Status = "failed"
			report.Error = "Invalid templates configuration"
			s.db.Save(&report)
			s.triggerReportWebhooks(schedule, &report, errors.NewError(errors.ErrCodeInvalidRequest, "Invalid templates configuration", err))
			return errors.WrapError(err, "Invalid templates configuration")
		}
	}

	// Generate report content based on configuration
	content, err := s.generateReportContent(queryIDs, chartIDs, templateIDs)
	if err != nil {
		report.Status = "failed"
		report.Error = err.Error()
		s.db.Save(&report)
		s.triggerReportWebhooks(schedule, &report, err)
		return err
	}

	// Update report with success status
	report.Status = "success"
	report.Content = content
	report.Error = ""
	if err := s.db.Save(&report).Error; err != nil {
		return errors.WrapError(err, "Could not update report status")
	}

	// Update schedule status
	schedule.LastRun = time.Now()
	s.db.Save(schedule)

	// Trigger webhook notifications
	s.triggerReportWebhooks(schedule, &report, nil)

	return nil
}

// generateReportContent generates the actual report content
func (s *ReportGenerationService) generateReportContent(queryIDs []uint, chartIDs []uint, templateIDs []uint) ([]byte, error) {
	// This is a simplified implementation
	// In a real scenario, you would:
	// 1. Execute queries and collect data
	// 2. Generate charts
	// 3. Apply templates
	// 4. Combine everything into a final report

	reportData := map[string]interface{}{
		"generated_at": time.Now().Format("2006-01-02 15:04:05"),
		"queries":      queryIDs,
		"charts":       chartIDs,
		"templates":    templateIDs,
	}

	// Convert to JSON for now
	content, err := json.Marshal(reportData)
	if err != nil {
		return nil, errors.WrapError(err, "Could not marshal report data")
	}

	return content, nil
}

// DownloadReport downloads a generated report
func (s *ReportGenerationService) DownloadReport(reportID uint, userID uint, isAdmin bool) (*models.Report, error) {
	var report models.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && report.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if report.Status != "success" {
		return nil, errors.NewBadRequestError("Report is not ready for download", nil)
	}

	return &report, nil
}

// triggerReportWebhooks sends webhook notifications for report events
func (s *ReportGenerationService) triggerReportWebhooks(schedule *models.ReportSchedule, report *models.Report, err error) {
	event := "report.generated"
	if err != nil {
		event = "report.failed"
	}

	// Get all webhooks for the user that are subscribed to report events
	webhooks, err := s.webhookRepo.FindByUser(schedule.UserID, false)
	if err != nil {
		utils.Logger.Errorf("Failed to fetch webhooks for user %d: %v", schedule.UserID, err)
		return
	}

	// Filter webhooks that are subscribed to this event type
	var eventWebhooks []models.Webhook
	for _, webhook := range webhooks {
		if !webhook.Active {
			continue
		}

		// Parse events from JSON string
		var events []string
		if err := json.Unmarshal([]byte(webhook.Events), &events); err != nil {
			utils.Logger.Errorf("Failed to parse webhook events for webhook %d: %v", webhook.ID, err)
			continue
		}

		// Check if webhook is subscribed to this event
		for _, subscribedEvent := range events {
			if subscribedEvent == event || subscribedEvent == "report.*" {
				eventWebhooks = append(eventWebhooks, webhook)
				break
			}
		}
	}

	// Prepare webhook payload
	payload := map[string]interface{}{
		"event":     event,
		"timestamp": time.Now().Unix(),
		"report": map[string]interface{}{
			"id":           report.ID,
			"name":         report.Name,
			"type":         report.Type,
			"status":       report.Status,
			"generated_at": report.GeneratedAt,
			"error":        report.Error,
		},
		"schedule": map[string]interface{}{
			"id":   schedule.ID,
			"name": schedule.Name,
			"type": schedule.Type,
		},
	}

	// Trigger webhooks asynchronously
	for _, webhook := range eventWebhooks {
		go s.triggerWebhook(webhook, payload, event)
	}

	utils.Logger.Infof("Triggered %d webhooks for event %s (report %d)", len(eventWebhooks), event, report.ID)
}

// triggerWebhook triggers a single webhook
func (s *ReportGenerationService) triggerWebhook(webhook models.Webhook, payload interface{}, event string) {
	// Create delivery record
	delivery := &models.WebhookDelivery{
		WebhookID: webhook.ID,
		Event:     event,
		Status:    "pending",
		Attempts:  0,
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.Logger.Errorf("Failed to marshal webhook payload: %v", err)
		delivery.Status = "failed"
		delivery.Response = "Failed to marshal payload"
		s.webhookRepo.CreateDelivery(delivery)
		return
	}
	delivery.Payload = string(payloadBytes)

	// Create delivery record
	if err := s.webhookRepo.CreateDelivery(delivery); err != nil {
		utils.Logger.Errorf("Failed to create webhook delivery record: %v", err)
		return
	}

	// Trigger webhook
	err = s.webhookTrigger.TriggerWebhook(webhook.URL, payload)

	// Update delivery record
	now := time.Now()
	delivery.SentAt = &now
	delivery.Attempts++

	if err != nil {
		delivery.Status = "failed"
		delivery.Response = err.Error()
		utils.Logger.Errorf("Webhook delivery failed for webhook %d: %v", webhook.ID, err)
	} else {
		delivery.Status = "success"
		delivery.Response = "Webhook delivered successfully"
		utils.Logger.Infof("Webhook delivered successfully for webhook %d", webhook.ID)
	}

	// Update delivery record
	if err := s.webhookRepo.UpdateDelivery(delivery); err != nil {
		utils.Logger.Errorf("Failed to update webhook delivery record: %v", err)
	}
}
