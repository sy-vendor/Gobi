package utils

import (
	"encoding/json"
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/database"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xuri/excelize/v2"
)

var reportCron *cron.Cron

// InitReportGenerator initializes the report generator cron jobs
func InitReportGenerator() {
	reportCron = cron.New()
	reportCron.Start()

	// Schedule report generation check every minute
	reportCron.AddFunc("* * * * *", checkAndGenerateReports)
}

// StopReportGenerator stops the report generator cron jobs
func StopReportGenerator() {
	if reportCron != nil {
		reportCron.Stop()
	}
}

// checkAndGenerateReports checks for reports that need to be generated
func checkAndGenerateReports() {
	now := time.Now()
	var schedules []models.ReportSchedule

	// Find all active schedules that are due
	if err := database.DB.Where("active = ? AND next_run <= ?", true, now).Find(&schedules).Error; err != nil {
		Logger.WithFields(map[string]interface{}{
			"action": "check_reports",
			"error":  err.Error(),
		}).Error("Failed to fetch report schedules")
		return
	}

	for _, schedule := range schedules {
		go generateReport(&schedule)
	}
}

// generateReport generates a report based on the schedule
func generateReport(schedule *models.ReportSchedule) {
	// Create a new report record
	report := models.Report{
		UserID:      schedule.UserID,
		Name:        schedule.Name,
		Type:        schedule.Type,
		Status:      "pending",
		GeneratedAt: time.Now(),
	}

	if err := database.DB.Create(&report).Error; err != nil {
		Logger.WithFields(map[string]interface{}{
			"action":     "generate_report",
			"scheduleID": schedule.ID,
			"error":      err.Error(),
		}).Error("Failed to create report record")
		return
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Process queries
	var queryIDs []uint
	if err := json.Unmarshal([]byte(schedule.Queries), &queryIDs); err == nil {
		for i, queryID := range queryIDs {
			var query models.Query
			if err := database.DB.First(&query, queryID).Error; err != nil {
				continue
			}

			// Execute query
			var ds models.DataSource
			if err := database.DB.First(&ds, query.DataSourceID).Error; err != nil {
				continue
			}

			results, err := ExecuteSQL(ds, query.SQL)
			if err != nil {
				continue
			}

			// Create sheet for query results
			sheetName := fmt.Sprintf("Query_%d", i+1)
			f.NewSheet(sheetName)

			// Write headers
			if len(results) > 0 {
				col := 1
				for key := range results[0] {
					cell, _ := excelize.CoordinatesToCellName(col, 1)
					f.SetCellValue(sheetName, cell, key)
					col++
				}
			}

			// Write data
			for row, result := range results {
				col := 1
				for _, value := range result {
					cell, _ := excelize.CoordinatesToCellName(col, row+2)
					f.SetCellValue(sheetName, cell, value)
					col++
				}
			}
		}
	}

	// Save the report
	content, err := f.WriteToBuffer()
	if err != nil {
		report.Status = "failed"
		report.Error = "Failed to generate Excel file: " + err.Error()
	} else {
		report.Status = "success"
		report.Content = content.Bytes()
	}

	// Update report status
	if err := database.DB.Save(&report).Error; err != nil {
		Logger.WithFields(map[string]interface{}{
			"action":     "generate_report",
			"scheduleID": schedule.ID,
			"reportID":   report.ID,
			"error":      err.Error(),
		}).Error("Failed to update report status")
	}

	// Update schedule next run time using cron pattern
	schedule.LastRun = time.Now()
	schedule.NextRun = calculateNextRunFromCron(schedule.CronPattern)
	if err := database.DB.Save(&schedule).Error; err != nil {
		Logger.WithFields(map[string]interface{}{
			"action":     "generate_report",
			"scheduleID": schedule.ID,
			"error":      err.Error(),
		}).Error("Failed to update schedule next run time")
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
