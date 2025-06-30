package infrastructure

import (
	"encoding/json"
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
)

// ValidationServiceImpl implements ValidationService
type ValidationServiceImpl struct{}

// NewValidationService creates a new ValidationService instance
func NewValidationService() *ValidationServiceImpl {
	return &ValidationServiceImpl{}
}

// ValidateSQL validates SQL query
func (s *ValidationServiceImpl) ValidateSQL(sql string) error {
	return utils.ValidateSQLComplete(sql)
}

// ValidateChartType validates chart type
func (s *ValidationServiceImpl) ValidateChartType(chartType string) error {
	validChartTypes := map[string]bool{
		"bar":               true,
		"line":              true,
		"pie":               true,
		"scatter":           true,
		"radar":             true,
		"heatmap":           true,
		"gauge":             true,
		"funnel":            true,
		"area":              true,
		"3d-bar":            true,
		"3d-scatter":        true,
		"3d-surface":        true,
		"3d-bubble":         true,
		"treemap":           true,
		"sunburst":          true,
		"tree":              true,
		"boxplot":           true,
		"candlestick":       true,
		"wordcloud":         true,
		"graph":             true,
		"waterfall":         true,
		"polar":             true,
		"gantt":             true,
		"rose":              true,
		"geo":               true,
		"map":               true,
		"choropleth":        true,
		"progress":          true,
		"circular-progress": true,
	}

	if !validChartTypes[chartType] {
		return errors.ErrInvalidChartType
	}

	return nil
}

// ValidateDataSource validates datasource configuration
func (s *ValidationServiceImpl) ValidateDataSource(ds *models.DataSource) error {
	if ds.Name == "" {
		return errors.ErrDataSourceNameRequired
	}

	if ds.Type == "" {
		return errors.ErrDataSourceTypeRequired
	}

	if ds.Host == "" {
		return errors.ErrDataSourceHostRequired
	}

	if ds.Port == 0 {
		return errors.ErrDataSourcePortRequired
	}

	if ds.Database == "" {
		return errors.ErrDataSourceDatabaseRequired
	}

	return nil
}

// ValidateChartConfig validates chart configuration JSON
func (s *ValidationServiceImpl) ValidateChartConfig(config string) error {
	if config == "" {
		return nil // Empty config is valid
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return errors.ErrInvalidChartConfig
	}

	return nil
}

// ValidateChartData validates chart data JSON
func (s *ValidationServiceImpl) ValidateChartData(data string) error {
	if data == "" {
		return nil // Empty data is valid
	}

	var dataArray []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dataArray); err != nil {
		return errors.ErrInvalidChartData
	}

	return nil
}
