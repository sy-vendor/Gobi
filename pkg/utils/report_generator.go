package utils

import (
	"bytes"
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

// GenerateExcelFromTemplate populates an Excel template with chart data.
func GenerateExcelFromTemplate(chartData string, templateData []byte, chartID string) ([]byte, error) {
	// Unmarshal chart data
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(chartData), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chart data: %w", err)
	}

	// Create a reader from the template data
	r := bytes.NewReader(templateData)

	// Open the template
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel template: %w", err)
	}
	defer f.Close()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in the template")
	}

	// Write headers if there's data
	if len(data) > 0 {
		col := 1
		for key := range data[0] {
			cell, _ := excelize.CoordinatesToCellName(col, 1)
			f.SetCellValue(sheetName, cell, key)
			col++
		}
	}

	// Write data rows
	for row, rowData := range data {
		col := 1
		// This assumes consistent key order, which is not guaranteed for maps.
		// For production, it's better to have a fixed order of columns.
		for _, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(col, row+2)
			f.SetCellValue(sheetName, cell, value)
			col++
		}
	}

	// Write to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel to buffer: %w", err)
	}

	return buf.Bytes(), nil
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
	return CalculateNextRunFromCron(cronPattern, time.Now())
}
