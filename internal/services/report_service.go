package services

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"time"

	errs "errors"

	"gorm.io/gorm"
)

// ReportService handles report-related business logic
type ReportService struct {
	db *gorm.DB
}

// NewReportService creates a new ReportService instance
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// CreateReport creates a new report
func (s *ReportService) CreateReport(report *models.Report, userID uint) error {
	report.UserID = userID

	if err := s.db.Create(report).Error; err != nil {
		return errors.WrapError(err, "Could not create report")
	}

	return nil
}

// ListReports retrieves reports based on user permissions
func (s *ReportService) ListReports(userID uint, isAdmin bool) ([]models.Report, error) {
	var reports []models.Report

	query := s.db.Preload("User").Model(&models.Report{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&reports).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch reports")
	}

	return reports, nil
}

// GetReport retrieves a specific report
func (s *ReportService) GetReport(reportID uint, userID uint, isAdmin bool) (*models.Report, error) {
	var report models.Report
	if err := s.db.Preload("User").First(&report, reportID).Error; err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch report")
	}

	if !isAdmin && report.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return &report, nil
}

// UpdateReport updates a report
func (s *ReportService) UpdateReport(reportID uint, updates *models.Report, userID uint, isAdmin bool) (*models.Report, error) {
	var report models.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch report")
	}

	if !isAdmin && report.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Only update allowed fields
	if updates.Name != "" {
		report.Name = updates.Name
	}
	if updates.Type != "" {
		report.Type = updates.Type
	}
	// Note: We don't update Content, GeneratedAt, Status, Error as these are managed by the system
	// during report generation

	if err := s.db.Save(&report).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update report")
	}

	return &report, nil
}

// DeleteReport deletes a report
func (s *ReportService) DeleteReport(reportID uint, userID uint, isAdmin bool) error {
	var report models.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "Could not fetch report")
	}

	if !isAdmin && report.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&report).Error; err != nil {
		return errors.WrapError(err, "Could not delete report")
	}

	return nil
}

// GenerateReportResult represents the result of report generation
type GenerateReportResult struct {
	ReportID     uint      `json:"reportId"`
	FileName     string    `json:"fileName"`
	FileSize     int64     `json:"fileSize"`
	GeneratedAt  time.Time `json:"generatedAt"`
	DownloadURL  string    `json:"downloadUrl"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

// GenerateReport generates a report based on the report configuration
func (s *ReportService) GenerateReport(reportID uint, userID uint, isAdmin bool) (*GenerateReportResult, error) {
	// Get report with all related data
	var report models.Report
	if err := s.db.Preload("User").First(&report, reportID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check
	if !isAdmin && report.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update report status
	report.Status = "generating"
	report.GeneratedAt = time.Now()
	s.db.Save(&report)

	// For now, we'll create a simple Excel report
	// In a real implementation, you would use the actual template and data
	content, err := utils.GenerateExcelFromTemplate("[]", []byte{}, "report")
	if err != nil {
		// Update report status to failed
		report.Status = "failed"
		report.Error = err.Error()
		s.db.Save(&report)

		return &GenerateReportResult{
			ReportID:     reportID,
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, nil
	}

	// Update report status to success
	report.Status = "success"
	report.Content = content
	report.Error = ""
	s.db.Save(&report)

	return &GenerateReportResult{
		ReportID:    reportID,
		FileName:    "report.xlsx",
		FileSize:    int64(len(content)),
		GeneratedAt: time.Now(),
		DownloadURL: "/api/reports/" + string(rune(reportID)) + "/download",
		Status:      "success",
	}, nil
}

// GetReportStatus gets the current status of a report
func (s *ReportService) GetReportStatus(reportID uint, userID uint, isAdmin bool) (*GenerateReportResult, error) {
	var report models.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && report.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return &GenerateReportResult{
		ReportID:     reportID,
		FileName:     "report.xlsx",
		FileSize:     int64(len(report.Content)),
		GeneratedAt:  report.GeneratedAt,
		DownloadURL:  "/api/reports/" + string(rune(reportID)) + "/download",
		Status:       report.Status,
		ErrorMessage: report.Error,
	}, nil
}
