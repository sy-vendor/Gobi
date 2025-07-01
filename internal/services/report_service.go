package services

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"time"

	errs "errors"
)

// ReportService handles report-related business logic
type ReportService struct {
	reportRepo        ReportRepository
	reportGenerator   ReportGeneratorService
	permissionService PermissionService
}

// NewReportService creates a new ReportService instance
func NewReportService(
	reportRepo ReportRepository,
	reportGenerator ReportGeneratorService,
	permissionService PermissionService,
) *ReportService {
	return &ReportService{
		reportRepo:        reportRepo,
		reportGenerator:   reportGenerator,
		permissionService: permissionService,
	}
}

// CreateReport creates a new report
func (s *ReportService) CreateReport(report *models.Report, userID uint) error {
	report.UserID = userID

	if err := s.reportRepo.Create(report); err != nil {
		return errors.WrapError(err, "Could not create report")
	}

	return nil
}

// ListReports retrieves reports based on user permissions
func (s *ReportService) ListReports(userID uint, isAdmin bool) ([]models.Report, error) {
	reports, err := s.reportRepo.FindByUser(userID, isAdmin)
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch reports")
	}

	return reports, nil
}

// GetReport retrieves a specific report
func (s *ReportService) GetReport(reportID uint, userID uint, isAdmin bool) (*models.Report, error) {
	report, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch report")
	}

	if !s.permissionService.CanAccess(userID, reportID, "report", isAdmin) {
		return nil, errors.ErrForbidden
	}

	return report, nil
}

// UpdateReport updates a report
func (s *ReportService) UpdateReport(reportID uint, updates *models.Report, userID uint, isAdmin bool) (*models.Report, error) {
	report, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch report")
	}

	if !s.permissionService.CanAccess(userID, reportID, "report", isAdmin) {
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

	if err := s.reportRepo.Update(report); err != nil {
		return nil, errors.WrapError(err, "Could not update report")
	}

	return report, nil
}

// DeleteReport deletes a report
func (s *ReportService) DeleteReport(reportID uint, userID uint, isAdmin bool) error {
	_, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "Could not fetch report")
	}

	if !s.permissionService.CanAccess(userID, reportID, "report", isAdmin) {
		return errors.ErrForbidden
	}

	if err := s.reportRepo.Delete(reportID); err != nil {
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
	report, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check
	if !s.permissionService.CanAccess(userID, reportID, "report", isAdmin) {
		return nil, errors.ErrForbidden
	}

	// Update report status
	report.Status = "generating"
	report.GeneratedAt = time.Now()
	if err := s.reportRepo.UpdateStatus(reportID, "generating", ""); err != nil {
		return nil, errors.WrapError(err, "Could not update report status")
	}

	// For now, we'll create a simple Excel report
	// In a real implementation, you would use the actual template and data
	content, err := s.reportGenerator.GenerateExcelFromTemplate("[]", []byte{}, "report")
	if err != nil {
		// Update report status to failed
		if updateErr := s.reportRepo.UpdateStatus(reportID, "failed", err.Error()); updateErr != nil {
			// Log the error but don't fail the operation
		}

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
	if err := s.reportRepo.Update(report); err != nil {
		return nil, errors.WrapError(err, "Could not update report")
	}

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
	report, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if !s.permissionService.CanAccess(userID, reportID, "report", isAdmin) {
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
