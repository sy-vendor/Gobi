package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gobi/pkg/security"
	"gobi/pkg/utils"
)

// SQLValidationService provides comprehensive SQL validation functionality
type SQLValidationService struct {
	config    *security.SQLSecurityConfig
	validator *utils.SQLValidator
}

// NewSQLValidationService creates a new SQL validation service
func NewSQLValidationService() *SQLValidationService {
	return &SQLValidationService{
		config:    security.GetGlobalSQLConfig(),
		validator: utils.GetGlobalSQLValidator(),
	}
}

// ValidateSQL validates SQL query with caching and performance optimization
func (s *SQLValidationService) ValidateSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	// Use cached validation for performance
	valid, errorMsg := s.config.ValidateSQLWithCache(sql)
	if !valid {
		return fmt.Errorf(errorMsg)
	}

	// Additional validations
	return s.validator.ValidateSQL(sql)
}

// ValidateSQLWithContext validates SQL with context and timeout
func (s *SQLValidationService) ValidateSQLWithContext(ctx context.Context, sql string) error {
	// Check context timeout
	select {
	case <-ctx.Done():
		return fmt.Errorf("validation timeout: %v", ctx.Err())
	default:
	}

	return s.ValidateSQL(sql)
}

// ValidateSQLBatch validates multiple SQL queries efficiently
func (s *SQLValidationService) ValidateSQLBatch(sqls []string) ([]error, []bool) {
	errors := make([]error, len(sqls))
	valid := make([]bool, len(sqls))

	for i, sql := range sqls {
		if err := s.ValidateSQL(sql); err != nil {
			errors[i] = err
			valid[i] = false
		} else {
			valid[i] = true
		}
	}

	return errors, valid
}

// ValidateSQLWithLimits validates SQL with query limits
func (s *SQLValidationService) ValidateSQLWithLimits(sql string) error {
	// Basic validation
	if err := s.ValidateSQL(sql); err != nil {
		return err
	}

	// Check query limits
	_, maxRows, maxQueryLength := s.config.GetQueryLimits()

	// Check query length
	if maxQueryLength > 0 && len(sql) > maxQueryLength {
		return fmt.Errorf("query too long: %d characters (max: %d)", len(sql), maxQueryLength)
	}

	// Check for LIMIT clause if maxRows is set
	if maxRows > 0 {
		if !strings.Contains(strings.ToUpper(sql), "LIMIT") {
			return fmt.Errorf("query must include LIMIT clause (max: %d rows)", maxRows)
		}
	}

	// Note: maxExecutionTime is checked at execution time, not validation time
	return nil
}

// ValidateTableName validates table name with enhanced security
func (s *SQLValidationService) ValidateTableName(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	// Use config-based validation
	if !s.config.ValidateTableName(tableName) {
		return fmt.Errorf("invalid table name format: %s", tableName)
	}

	// Additional security checks
	return s.validator.ValidateTableName(tableName)
}

// ValidateColumnName validates column name with enhanced security
func (s *SQLValidationService) ValidateColumnName(columnName string) error {
	if columnName == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	// Use config-based validation
	if !s.config.ValidateColumnName(columnName) {
		return fmt.Errorf("invalid column name format: %s", columnName)
	}

	// Additional security checks
	return s.validator.ValidateColumnNameSmart(columnName)
}

// ValidateSQLComplete performs comprehensive SQL validation
func (s *SQLValidationService) ValidateSQLComplete(sql string) error {
	// Basic validation
	if err := s.ValidateSQL(sql); err != nil {
		return err
	}

	// Smart validation
	if err := s.validator.ValidateSQLSmart(sql); err != nil {
		return err
	}

	// Complete validation (includes table and column name validation)
	return s.validator.ValidateSQLComplete(sql)
}

// SanitizeSQL sanitizes SQL input
func (s *SQLValidationService) SanitizeSQL(sql string) string {
	return s.validator.SanitizeSQL(sql)
}

// IsReadOnlyQuery checks if SQL query is read-only
func (s *SQLValidationService) IsReadOnlyQuery(sql string) bool {
	return s.validator.IsReadOnlyQuery(sql)
}

// GetValidationStats returns comprehensive validation statistics
func (s *SQLValidationService) GetValidationStats() map[string]interface{} {
	validatorStats := s.validator.GetValidationStats()
	configStats := s.config.GetStats()

	return map[string]interface{}{
		"service_stats": map[string]interface{}{
			"service_type": "SQLValidationService",
			"created_at":   time.Now(),
		},
		"validator_stats": validatorStats,
		"config_stats":    configStats,
		"performance": map[string]interface{}{
			"cache_enabled":     true,
			"batch_validation":  true,
			"context_support":   true,
			"limits_validation": true,
		},
	}
}

// ReloadConfig reloads SQL security configuration
func (s *SQLValidationService) ReloadConfig(filename string) error {
	return s.config.ReloadConfig(filename)
}

// GetSecuritySettings returns current security settings
func (s *SQLValidationService) GetSecuritySettings() map[string]bool {
	return s.config.GetSecuritySettings()
}

// GetQueryLimits returns current query limits
func (s *SQLValidationService) GetQueryLimits() (maxExecutionTime string, maxRows int, maxQueryLength int) {
	return s.config.GetQueryLimits()
}

// ValidateSQLWithCustomRules validates SQL with custom validation rules
func (s *SQLValidationService) ValidateSQLWithCustomRules(sql string, customRules map[string]bool) error {
	// Apply custom rules
	if customRules["require_select"] && !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
		return fmt.Errorf("query must start with SELECT")
	}

	if customRules["require_limit"] && !strings.Contains(strings.ToUpper(sql), "LIMIT") {
		return fmt.Errorf("query must include LIMIT clause")
	}

	if customRules["require_where"] && !strings.Contains(strings.ToUpper(sql), "WHERE") {
		return fmt.Errorf("query must include WHERE clause")
	}

	// Apply standard validation
	return s.ValidateSQL(sql)
}

// ValidateSQLWithWhitelist validates SQL against a custom whitelist
func (s *SQLValidationService) ValidateSQLWithWhitelist(sql string, allowedKeywords []string, blockedKeywords []string) error {
	normalizedSQL := strings.ToUpper(strings.TrimSpace(sql))

	// Check blocked keywords
	for _, keyword := range blockedKeywords {
		if strings.Contains(normalizedSQL, strings.ToUpper(keyword)) {
			return fmt.Errorf("blocked keyword detected: %s", keyword)
		}
	}

	// Check if all keywords are allowed (if whitelist is provided)
	if len(allowedKeywords) > 0 {
		// Extract keywords from SQL (simplified)
		words := strings.Fields(normalizedSQL)
		for _, word := range words {
			// Skip common SQL elements
			if word == "SELECT" || word == "FROM" || word == "WHERE" || word == "AND" || word == "OR" {
				continue
			}

			// Check if word is in allowed keywords
			found := false
			for _, allowed := range allowedKeywords {
				if strings.ToUpper(allowed) == word {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("keyword not in whitelist: %s", word)
			}
		}
	}

	return s.ValidateSQL(sql)
}

// GetValidationReport generates a detailed validation report
func (s *SQLValidationService) GetValidationReport(sql string) map[string]interface{} {
	report := map[string]interface{}{
		"sql":             sql,
		"timestamp":       time.Now(),
		"validations":     map[string]interface{}{},
		"recommendations": []string{},
	}

	// Basic validation
	if err := s.ValidateSQL(sql); err != nil {
		report["validations"].(map[string]interface{})["basic"] = map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		}
	} else {
		report["validations"].(map[string]interface{})["basic"] = map[string]interface{}{
			"valid": true,
		}
	}

	// Read-only check
	isReadOnly := s.IsReadOnlyQuery(sql)
	report["validations"].(map[string]interface{})["readonly"] = map[string]interface{}{
		"valid": isReadOnly,
	}

	// Query length
	report["validations"].(map[string]interface{})["length"] = map[string]interface{}{
		"length": len(sql),
		"valid":  len(sql) <= 10000, // Default limit
	}

	// Generate recommendations
	if !isReadOnly {
		report["recommendations"] = append(report["recommendations"].([]string), "Consider using SELECT queries only for security")
	}

	if len(sql) > 5000 {
		report["recommendations"] = append(report["recommendations"].([]string), "Consider optimizing query length for better performance")
	}

	if !strings.Contains(strings.ToUpper(sql), "LIMIT") {
		report["recommendations"] = append(report["recommendations"].([]string), "Add LIMIT clause to prevent large result sets")
	}

	return report
}
