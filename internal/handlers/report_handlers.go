package handlers

import (
	"encoding/json"
	"gobi/internal/models"
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
	DB *gorm.DB
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{DB: db}
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

	if err := h.DB.Create(&schedule).Error; err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "create_report_schedule",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to create report schedule")
		c.Error(errors.WrapError(err, "Could not create report schedule"))
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

	var schedules []models.ReportSchedule
	query := h.DB.Model(&models.ReportSchedule{})

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&schedules).Error; err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "list_report_schedules",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to list report schedules")
		c.Error(errors.WrapError(err, "Could not fetch report schedules"))
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// GetReportSchedule gets a specific report schedule
func (h *ReportHandler) GetReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := h.DB.First(&schedule, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	if role != "admin" && schedule.UserID != userID {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateReportSchedule updates a report schedule
func (h *ReportHandler) UpdateReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := h.DB.First(&schedule, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	if role != "admin" && schedule.UserID != userID {
		c.Error(errors.ErrForbidden)
		return
	}

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

	if req.Name != "" {
		schedule.Name = req.Name
	}
	if req.Type != "" {
		schedule.Type = req.Type
	}
	if req.QueryIDs != nil {
		queryIDs, _ := json.Marshal(req.QueryIDs)
		schedule.Queries = string(queryIDs)
	}
	if req.ChartIDs != nil {
		chartIDs, _ := json.Marshal(req.ChartIDs)
		schedule.Charts = string(chartIDs)
	}
	if req.TemplateIDs != nil {
		templateIDs, _ := json.Marshal(req.TemplateIDs)
		schedule.Templates = string(templateIDs)
	}
	if req.CronPattern != "" {
		// 验证新的cron表达式
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		_, err := parser.Parse(req.CronPattern)
		if err != nil {
			utils.Logger.WithFields(map[string]interface{}{
				"action":      "update_report_schedule",
				"cronPattern": req.CronPattern,
				"error":       err.Error(),
			}).Warn("Invalid cron pattern")
			c.Error(errors.NewBadRequestError("Invalid cron pattern", err))
			return
		}
		schedule.CronPattern = req.CronPattern
		// 重新计算下次运行时间
		schedule.NextRun = calculateNextRunFromCron(req.CronPattern)
	}
	if req.Active != nil {
		schedule.Active = *req.Active
	}

	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "update_report_schedule",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to update report schedule")
		c.Error(errors.WrapError(err, "Could not update report schedule"))
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":     "update_report_schedule",
		"userID":     userID,
		"scheduleID": schedule.ID,
		"nextRun":    schedule.NextRun,
	}).Info("Report schedule updated successfully")

	c.JSON(http.StatusOK, schedule)
}

// DeleteReportSchedule deletes a report schedule
func (h *ReportHandler) DeleteReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := h.DB.First(&schedule, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	if role != "admin" && schedule.UserID != userID {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := h.DB.Delete(&schedule).Error; err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "delete_report_schedule",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to delete report schedule")
		c.Error(errors.WrapError(err, "Could not delete report schedule"))
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"action":     "delete_report_schedule",
		"userID":     userID,
		"scheduleID": schedule.ID,
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
	var chart models.Chart
	if err := h.DB.First(&chart, req.ChartID).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}
	// TODO: Replace with actual PDF generation logic
	pdfBytes := []byte("Fake PDF content for chart: " + chart.Name)
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
	var chart models.Chart
	if err := h.DB.First(&chart, req.ChartID).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}
	var template models.ExcelTemplate
	if err := h.DB.First(&template, req.TemplateID).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}
	excelBytes, err := utils.GenerateExcelFromTemplate(chart.Data, template.Template, strconv.Itoa(int(chart.ID)))
	if err != nil {
		c.Error(errors.WrapError(err, "Could not generate excel report"))
		return
	}
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}

// ListReports lists all generated reports for the user
func (h *ReportHandler) ListReports(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var reports []models.Report
	query := h.DB.Model(&models.Report{})

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&reports).Error; err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"action": "list_reports",
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to list reports")
		c.Error(errors.WrapError(err, "Could not fetch reports"))
		return
	}

	c.JSON(http.StatusOK, reports)
}

// DownloadReport downloads a generated report
func (h *ReportHandler) DownloadReport(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var report models.Report
	if err := h.DB.First(&report, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	if role != "admin" && report.UserID != userID {
		c.Error(errors.ErrForbidden)
		return
	}

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

// calculateNextRun calculates the next run time based on report type
func calculateNextRun(reportType string) time.Time {
	now := time.Now()
	switch reportType {
	case "daily":
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	case "weekly":
		daysUntilNextWeek := 7 - int(now.Weekday())
		if daysUntilNextWeek == 0 {
			daysUntilNextWeek = 7
		}
		return time.Date(now.Year(), now.Month(), now.Day()+daysUntilNextWeek, 0, 0, 0, 0, now.Location())
	case "monthly":
		return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	default:
		return now
	}
}

// calculateNextRunFromCron calculates the next run time based on cron pattern
func calculateNextRunFromCron(cronPattern string) time.Time {
	now := time.Now()
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronPattern)
	if err != nil {
		// 返回一个零时时间或者一个表示错误的将来时间
		return time.Time{}
	}
	return schedule.Next(now)
}
