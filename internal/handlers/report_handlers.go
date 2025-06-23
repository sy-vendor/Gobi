package handlers

import (
	"encoding/json"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// CreateReportSchedule creates a new report schedule
func CreateReportSchedule(c *gin.Context) {
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

	if err := database.DB.Create(&schedule).Error; err != nil {
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
func ListReportSchedules(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedules []models.ReportSchedule
	query := database.DB.Model(&models.ReportSchedule{})

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
func GetReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := database.DB.First(&schedule, id).Error; err != nil {
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
func UpdateReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := database.DB.First(&schedule, id).Error; err != nil {
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

	if err := database.DB.Save(&schedule).Error; err != nil {
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
func DeleteReportSchedule(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var schedule models.ReportSchedule
	if err := database.DB.First(&schedule, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	if role != "admin" && schedule.UserID != userID {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := database.DB.Delete(&schedule).Error; err != nil {
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

// ListReports lists all reports for the user
func ListReports(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var reports []models.Report
	query := database.DB.Model(&models.Report{})

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

// DownloadReport downloads a specific report
func DownloadReport(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var report models.Report
	if err := database.DB.First(&report, id).Error; err != nil {
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
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronPattern)
	if err != nil {
		return time.Now()
	}
	return schedule.Next(time.Now())
}
