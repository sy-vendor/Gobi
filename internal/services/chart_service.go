package services

import (
	"encoding/json"
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// ChartService handles chart-related business logic
type ChartService struct {
	db *gorm.DB
}

// NewChartService creates a new ChartService instance
func NewChartService(db *gorm.DB) *ChartService {
	return &ChartService{db: db}
}

// CreateChart creates a new chart
func (s *ChartService) CreateChart(chart *models.Chart, userID uint) error {
	// Validate chart type
	if err := s.validateChartType(chart.Type); err != nil {
		return err
	}

	chart.UserID = userID

	if err := s.db.Create(chart).Error; err != nil {
		return errors.WrapError(err, "Could not create chart")
	}

	return nil
}

// ListCharts retrieves charts based on user permissions
func (s *ChartService) ListCharts(userID uint, isAdmin bool) ([]models.Chart, error) {
	var charts []models.Chart

	query := s.db.Preload("Query").Preload("User").Model(&models.Chart{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&charts).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch charts")
	}

	return charts, nil
}

// GetChart retrieves a specific chart
func (s *ChartService) GetChart(chartID uint, userID uint, isAdmin bool) (*models.Chart, error) {
	var chart models.Chart
	if err := s.db.Preload("Query").Preload("User").First(&chart, chartID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && chart.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return &chart, nil
}

// UpdateChart updates a chart
func (s *ChartService) UpdateChart(chartID uint, updates *models.Chart, userID uint, isAdmin bool) (*models.Chart, error) {
	var chart models.Chart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && chart.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Validate chart type if it's being updated
	if updates.Type != "" {
		if err := s.validateChartType(updates.Type); err != nil {
			return nil, err
		}
		chart.Type = updates.Type
	}

	// Update allowed fields
	if updates.Name != "" {
		chart.Name = updates.Name
	}
	if updates.QueryID != 0 {
		chart.QueryID = updates.QueryID
	}
	if updates.Config != "" {
		chart.Config = updates.Config
	}
	if updates.Data != "" {
		chart.Data = updates.Data
	}
	if updates.Description != "" {
		chart.Description = updates.Description
	}

	if err := s.db.Save(&chart).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update chart")
	}

	return &chart, nil
}

// DeleteChart deletes a chart
func (s *ChartService) DeleteChart(chartID uint, userID uint, isAdmin bool) error {
	var chart models.Chart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		return errors.ErrNotFound
	}

	if !isAdmin && chart.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&chart).Error; err != nil {
		return errors.WrapError(err, "Could not delete chart")
	}

	return nil
}

// validateChartType validates if the chart type is supported
func (s *ChartService) validateChartType(chartType string) error {
	validChartTypes := map[string]bool{
		"bar":         true,
		"line":        true,
		"pie":         true,
		"scatter":     true,
		"radar":       true,
		"heatmap":     true,
		"gauge":       true,
		"funnel":      true,
		"area":        true,
		"3d-bar":      true,
		"3d-scatter":  true,
		"3d-surface":  true,
		"3d-bubble":   true,
		"treemap":     true,
		"sunburst":    true,
		"tree":        true,
		"boxplot":     true,
		"candlestick": true,
	}

	if !validChartTypes[chartType] {
		return errors.NewBadRequestError("Invalid chart type", nil)
	}

	return nil
}

// ValidateChartConfig validates chart configuration JSON
func (s *ChartService) ValidateChartConfig(config string) error {
	if config == "" {
		return nil // Empty config is valid
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return errors.NewBadRequestError("Invalid chart configuration JSON", err)
	}

	return nil
}

// ValidateChartData validates chart data JSON
func (s *ChartService) ValidateChartData(data string) error {
	if data == "" {
		return nil // Empty data is valid
	}

	var dataArray []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dataArray); err != nil {
		return errors.NewBadRequestError("Invalid chart data JSON", err)
	}

	return nil
}
