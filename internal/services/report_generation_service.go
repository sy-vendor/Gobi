package services

import (
	"encoding/json"
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ReportGenerationService handles report generation business logic
type ReportGenerationService struct {
	db *gorm.DB
}

// NewReportGenerationService creates a new ReportGenerationService instance
func NewReportGenerationService(db *gorm.DB) *ReportGenerationService {
	return &ReportGenerationService{db: db}
}

// GenerateExcelReport generates an Excel report from chart data and template
func (s *ReportGenerationService) GenerateExcelReport(chartID uint, templateID uint, userID uint, isAdmin bool) ([]byte, error) {
	// Get chart
	var chart models.Chart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check
	if !isAdmin && chart.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Get template
	var template models.ExcelTemplate
	if err := s.db.First(&template, templateID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check for template
	if !isAdmin && template.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Generate Excel report
	excelBytes, err := utils.GenerateExcelFromTemplate(chart.Data, template.Template, strconv.Itoa(int(chart.ID)))
	if err != nil {
		return nil, errors.WrapError(err, "Could not generate Excel report")
	}

	return excelBytes, nil
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
			return errors.WrapError(err, "Invalid queries configuration")
		}
	}

	if schedule.Charts != "" {
		if err := json.Unmarshal([]byte(schedule.Charts), &chartIDs); err != nil {
			report.Status = "failed"
			report.Error = "Invalid charts configuration"
			s.db.Save(&report)
			return errors.WrapError(err, "Invalid charts configuration")
		}
	}

	if schedule.Templates != "" {
		if err := json.Unmarshal([]byte(schedule.Templates), &templateIDs); err != nil {
			report.Status = "failed"
			report.Error = "Invalid templates configuration"
			s.db.Save(&report)
			return errors.WrapError(err, "Invalid templates configuration")
		}
	}

	// Generate report content based on configuration
	content, err := s.generateReportContent(queryIDs, chartIDs, templateIDs)
	if err != nil {
		report.Status = "failed"
		report.Error = err.Error()
		s.db.Save(&report)
		return err
	}

	// Update report with success status
	report.Status = "success"
	report.Content = content
	report.Error = ""
	if err := s.db.Save(&report).Error; err != nil {
		return errors.WrapError(err, "Could not update report status")
	}

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
