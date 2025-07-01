package handlers

import (
	"encoding/json"
	"gobi/internal/models"
	"gobi/internal/services"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// ReportHandler holds dependencies for report-related handlers.
type ReportHandler struct {
	DB                      *gorm.DB
	ReportService           *services.ReportService
	ReportScheduleService   *services.ReportScheduleService
	ReportGenerationService *services.ReportGenerationService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(db *gorm.DB, serviceFactory *services.ServiceFactory) *ReportHandler {
	return &ReportHandler{
		DB:                      db,
		ReportService:           serviceFactory.CreateReportService(),
		ReportScheduleService:   services.NewReportScheduleService(db),
		ReportGenerationService: serviceFactory.CreateReportGenerationService(),
	}
}

// CreateReport creates a new report
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var report models.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "create_report",
			"error":  err.Error(),
		}).Warn("Invalid report data")
		c.Error(errors.NewBadRequestError("Invalid report data", err))
		return
	}

	userID, _ := c.Get("userID")
	if err := h.ReportService.CreateReport(&report, userID.(uint)); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "create_report",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to create report")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":   "create_report",
		"userID":   userID,
		"reportID": report.ID,
	}).Info("Report created successfully")

	c.JSON(http.StatusCreated, report)
}

// ListReports lists all reports for the user
func (h *ReportHandler) ListReports(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	reports, err := h.ReportService.ListReports(userID.(uint), isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "list_reports",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to list reports")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":  "list_reports",
		"userID":  userID,
		"count":   len(reports),
		"isAdmin": isAdmin,
	}).Info("Reports listed successfully")

	c.JSON(http.StatusOK, reports)
}

// GetReport gets a specific report
func (h *ReportHandler) GetReport(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	report, err := h.ReportService.GetReport(uint(reportID), userID.(uint), isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "get_report",
			"userID":   userID,
			"reportID": reportID,
			"error":    err.Error(),
		}).Error("Failed to get report")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":   "get_report",
		"userID":   userID,
		"reportID": reportID,
	}).Info("Report retrieved successfully")

	c.JSON(http.StatusOK, report)
}

// UpdateReport updates a report
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	var req models.Report
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "update_report",
			"reportID": reportID,
			"error":    err.Error(),
		}).Warn("Invalid report data")
		c.Error(errors.NewBadRequestError("Invalid report data", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	report, err := h.ReportService.UpdateReport(uint(reportID), &req, userID.(uint), isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "update_report",
			"userID":   userID,
			"reportID": reportID,
			"error":    err.Error(),
		}).Error("Failed to update report")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":   "update_report",
		"userID":   userID,
		"reportID": reportID,
	}).Info("Report updated successfully")

	c.JSON(http.StatusOK, report)
}

// DeleteReport deletes a report
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	if err := h.ReportService.DeleteReport(uint(reportID), userID.(uint), isAdmin); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "delete_report",
			"userID":   userID,
			"reportID": reportID,
			"error":    err.Error(),
		}).Error("Failed to delete report")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":   "delete_report",
		"userID":   userID,
		"reportID": reportID,
	}).Info("Report deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}

// GenerateReport generates a report based on the report configuration
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	result, err := h.ReportService.GenerateReport(uint(reportID), userID.(uint), isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "generate_report",
			"userID":   userID,
			"reportID": reportID,
			"error":    err.Error(),
		}).Error("Failed to generate report")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":   "generate_report",
		"userID":   userID,
		"reportID": reportID,
		"status":   result.Status,
	}).Info("Report generation completed")

	c.JSON(http.StatusOK, result)
}

// GetReportStatus gets the current status of a report
func (h *ReportHandler) GetReportStatus(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role.(string) == "admin"

	result, err := h.ReportService.GetReportStatus(uint(reportID), userID.(uint), isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":   "get_report_status",
			"userID":   userID,
			"reportID": reportID,
			"error":    err.Error(),
		}).Error("Failed to get report status")
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateReportSchedule creates a new report schedule
func (h *ReportHandler) CreateReportSchedule(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required,oneof=daily weekly monthly"`
		QueryIDs    []uint `json:"query_ids"`
		ChartIDs    []uint `json:"chart_ids"`
		TemplateIDs []uint `json:"template_ids"`
		CronPattern string `json:"cron_pattern" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "create_report_schedule",
			"error":  err.Error(),
		}).Warn("Invalid report schedule data")
		c.Error(errors.NewBadRequestError("Invalid report schedule data", err))
		return
	}

	// 验证cron表达式
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(req.CronPattern)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action":      "create_report_schedule",
			"cronPattern": req.CronPattern,
			"error":       err.Error(),
		}).Warn("Invalid cron pattern")
		c.Error(errors.NewBadRequestError("Invalid cron pattern", err))
		return
	}

	userID := c.GetUint("userID")

	// Convert arrays to JSON strings
	queryIDs, _ := json.Marshal(req.QueryIDs)
	chartIDs, _ := json.Marshal(req.ChartIDs)
	templateIDs, _ := json.Marshal(req.TemplateIDs)

	// 使用cron表达式计算下次运行时间
	nextRun := calculateNextRunFromCron(req.CronPattern)

	schedule := models.ReportSchedule{
		UserID:      userID,
		Name:        req.Name,
		Type:        req.Type,
		Queries:     string(queryIDs),
		Charts:      string(chartIDs),
		Templates:   string(templateIDs),
		CronPattern: req.CronPattern,
		Active:      true,
		NextRun:     nextRun,
	}

	if err := h.ReportScheduleService.CreateReportSchedule(&schedule, userID); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "create_report_schedule",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to create report schedule")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":     "create_report_schedule",
		"userID":     userID,
		"scheduleID": schedule.ID,
		"nextRun":    nextRun,
	}).Info("Report schedule created successfully")

	c.JSON(http.StatusCreated, schedule)
}

// ListReportSchedules lists all report schedules for the user
func (h *ReportHandler) ListReportSchedules(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	schedules, err := h.ReportScheduleService.ListReportSchedules(userID, isAdmin)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "list_report_schedules",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to list report schedules")
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// GetReportSchedule gets a specific report schedule
func (h *ReportHandler) GetReportSchedule(c *gin.Context) {
	id := c.Param("id")
	scheduleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid schedule ID", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	schedule, err := h.ReportScheduleService.GetReportSchedule(uint(scheduleID), userID, isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateReportSchedule updates a report schedule
func (h *ReportHandler) UpdateReportSchedule(c *gin.Context) {
	id := c.Param("id")
	scheduleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid schedule ID", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type" binding:"omitempty,oneof=daily weekly monthly"`
		QueryIDs    []uint `json:"query_ids"`
		ChartIDs    []uint `json:"chart_ids"`
		TemplateIDs []uint `json:"template_ids"`
		CronPattern string `json:"cron_pattern"`
		Active      *bool  `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid report schedule data", err))
		return
	}

	// Convert arrays to JSON strings if provided
	var queryIDs, chartIDs, templateIDs string
	if req.QueryIDs != nil {
		queryIDsBytes, _ := json.Marshal(req.QueryIDs)
		queryIDs = string(queryIDsBytes)
	}
	if req.ChartIDs != nil {
		chartIDsBytes, _ := json.Marshal(req.ChartIDs)
		chartIDs = string(chartIDsBytes)
	}
	if req.TemplateIDs != nil {
		templateIDsBytes, _ := json.Marshal(req.TemplateIDs)
		templateIDs = string(templateIDsBytes)
	}

	updates := &models.ReportSchedule{
		Name:        req.Name,
		Type:        req.Type,
		Queries:     queryIDs,
		Charts:      chartIDs,
		Templates:   templateIDs,
		CronPattern: req.CronPattern,
	}
	if req.Active != nil {
		updates.Active = *req.Active
	}

	schedule, err := h.ReportScheduleService.UpdateReportSchedule(uint(scheduleID), updates, userID, isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// DeleteReportSchedule deletes a report schedule
func (h *ReportHandler) DeleteReportSchedule(c *gin.Context) {
	id := c.Param("id")
	scheduleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid schedule ID", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	if err := h.ReportScheduleService.DeleteReportSchedule(uint(scheduleID), userID, isAdmin); err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "delete_report_schedule",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to delete report schedule")
		c.Error(err)
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":     "delete_report_schedule",
		"userID":     userID,
		"scheduleID": scheduleID,
	}).Info("Report schedule deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Report schedule deleted successfully"})
}

// GeneratePDFReport generates a PDF report from a chart
func (h *ReportHandler) GeneratePDFReport(c *gin.Context) {
	var req struct {
		ChartID uint `json:"chart_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid request", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	pdfBytes, err := h.ReportGenerationService.GeneratePDFReport(req.ChartID, userID, isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateExcelReport generates an Excel report from a chart and template
func (h *ReportHandler) GenerateExcelReport(c *gin.Context) {
	var req struct {
		ChartID    uint `json:"chart_id"`
		TemplateID uint `json:"template_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid request", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	excelBytes, err := h.ReportGenerationService.GenerateExcelReport(req.ChartID, req.TemplateID, userID, isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}

// DownloadReport downloads a generated report
func (h *ReportHandler) DownloadReport(c *gin.Context) {
	id := c.Param("id")
	reportID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid report ID", err))
		return
	}

	userID := c.GetUint("userID")
	role := c.GetString("role")
	isAdmin := role == "admin"

	report, err := h.ReportGenerationService.DownloadReport(uint(reportID), userID, isAdmin)
	if err != nil {
		c.Error(err)
		return
	}

	// Generate filename
	fileName := report.Name
	if report.Type == "daily" {
		fileName += "_" + report.GeneratedAt.Format("2006-01-02")
	} else if report.Type == "weekly" {
		fileName += "_week_" + report.GeneratedAt.Format("2006-01-02")
	} else if report.Type == "monthly" {
		fileName += "_" + report.GeneratedAt.Format("2006-01")
	}
	fileName += ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(report.Content)))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", report.Content)
}

// calculateNextRunFromCron calculates the next run time based on cron pattern
func calculateNextRunFromCron(cronPattern string) time.Time {
	return utils.CalculateNextRunFromCron(cronPattern, time.Time{})
}
