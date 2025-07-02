package utils

import (
	"fmt"
	"gobi/pkg/security"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

// SQLValidator provides SQL injection protection with enhanced performance
type SQLValidator struct {
	config *security.SQLSecurityConfig

	// Performance optimizations
	compiledPatterns map[string]*regexp.Regexp
	patternCache     map[string]bool
	cacheMu          sync.RWMutex

	// Validation statistics
	validationCount int64
	cacheHitCount   int64
	lastValidation  time.Time
}

var (
	globalSQLValidator *SQLValidator
	validatorOnce      sync.Once
)

// GetGlobalSQLValidator returns a singleton instance of SQLValidator
func GetGlobalSQLValidator() *SQLValidator {
	validatorOnce.Do(func() {
		globalSQLValidator = NewSQLValidator()
	})
	return globalSQLValidator
}

// NewSQLValidator creates a new SQL validator with enhanced performance
func NewSQLValidator() *SQLValidator {
	config := security.GetGlobalSQLConfig()

	validator := &SQLValidator{
		config:           config,
		compiledPatterns: make(map[string]*regexp.Regexp),
		patternCache:     make(map[string]bool),
	}

	// Pre-compile common patterns
	validator.compileCommonPatterns()

	return validator
}

// compileCommonPatterns pre-compiles commonly used regex patterns
func (v *SQLValidator) compileCommonPatterns() {
	commonPatterns := map[string]string{
		"identifier":     `^[a-zA-Z_][a-zA-Z0-9_.]*$`,
		"quoted_string":  `'[^']*'`,
		"number":         `\b\d+(?:\.\d+)?\b`,
		"whitespace":     `\s+`,
		"comment_single": `--.*$`,
		"comment_multi":  `/\*.*?\*/`,
	}

	for name, pattern := range commonPatterns {
		v.compiledPatterns[name] = regexp.MustCompile(pattern)
	}
}

// ValidateSQL validates SQL query for security with caching
func (v *SQLValidator) ValidateSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	// Use cached validation if available
	valid, errorMsg := v.config.ValidateSQLWithCache(sql)
	if !valid {
		return fmt.Errorf(errorMsg)
	}

	// Additional validations
	if err := v.checkBalancedParentheses(sql); err != nil {
		return err
	}

	if err := v.checkComments(sql); err != nil {
		return err
	}

	v.validationCount++
	v.lastValidation = time.Now()

	return nil
}

// ValidateSQLSmart performs smart SQL validation with context awareness
func (v *SQLValidator) ValidateSQLSmart(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	// Use cached validation
	valid, errorMsg := v.config.ValidateSQLWithCache(sql)
	if !valid {
		// Check if this is a context-specific keyword issue
		if strings.Contains(errorMsg, "blocked keyword detected: END") {
			// Check if END is in a CASE WHEN context
			if v.isContextualKeyword("END", sql) {
				// END is allowed in CASE WHEN context, continue with other validations
			} else {
				return fmt.Errorf(errorMsg)
			}
		} else {
			return fmt.Errorf(errorMsg)
		}
	}

	// Check if read-only is required
	settings := v.config.GetSecuritySettings()
	if settings["require_readonly"] && !v.IsReadOnlyQuery(sql) {
		return fmt.Errorf("only SELECT queries are allowed")
	}

	// Additional smart validations
	if err := v.checkBalancedParentheses(sql); err != nil {
		return err
	}

	if settings["allow_comments"] {
		if err := v.checkCommentsLoose(sql); err != nil {
			return err
		}
	} else {
		if err := v.checkComments(sql); err != nil {
			return err
		}
	}

	v.validationCount++
	v.lastValidation = time.Now()

	return nil
}

// ValidateSQLComplete performs comprehensive SQL validation
func (v *SQLValidator) ValidateSQLComplete(sql string) error {
	if err := v.ValidateSQLSmart(sql); err != nil {
		return err
	}

	// Get security settings to check if we should validate table/column names
	settings := v.config.GetSecuritySettings()

	// Extract and validate table names only if enabled
	if settings["validate_table_names"] {
		tableNames := v.extractTableNames(sql)
		for _, tableName := range tableNames {
			if !v.config.ValidateTableName(tableName) {
				return fmt.Errorf("invalid table name: %s", tableName)
			}
		}
	}

	// Extract and validate column names only if enabled
	if settings["validate_column_names"] {
		columnNames := v.extractColumnNames(sql)
		for _, columnName := range columnNames {
			// Skip validation for empty column names or wildcards
			if columnName == "" || columnName == "*" {
				continue
			}

			// Skip validation for SQL keywords that might be extracted as column names
			if v.isSQLKeyword(columnName) {
				continue
			}

			// Use smart validation for column names
			if err := v.ValidateColumnNameSmart(columnName); err != nil {
				return fmt.Errorf("invalid column name detected: %v", err)
			}
		}
	}

	return nil
}

// extractTableNames extracts table names from SQL query
func (v *SQLValidator) extractTableNames(sql string) []string {
	// Simple table name extraction (could be enhanced with proper SQL parsing)
	tableNames := []string{}

	// Look for FROM clause
	fromIndex := strings.Index(strings.ToUpper(sql), "FROM")
	if fromIndex == -1 {
		return tableNames
	}

	// Extract text after FROM
	afterFrom := sql[fromIndex+4:]

	// Simple regex to find table names
	tablePattern := v.compiledPatterns["identifier"]
	matches := tablePattern.FindAllString(afterFrom, -1)

	for _, match := range matches {
		// Skip common SQL keywords that might match the pattern
		if !v.isSQLKeyword(match) {
			tableNames = append(tableNames, match)
		}
	}

	return tableNames
}

// extractColumnNames extracts column names from SQL query
func (v *SQLValidator) extractColumnNames(sql string) []string {
	// Simple column name extraction with improved handling of functions and expressions
	columnNames := []string{}

	// Look for SELECT clause
	selectIndex := strings.Index(strings.ToUpper(sql), "SELECT")
	if selectIndex == -1 {
		return columnNames
	}

	// Extract text between SELECT and FROM
	fromIndex := strings.Index(strings.ToUpper(sql), "FROM")
	if fromIndex == -1 {
		return columnNames
	}

	selectClause := sql[selectIndex+6 : fromIndex]

	// Split by comma and extract column names
	parts := strings.Split(selectClause, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "*" {
			continue
		}

		// Remove AS aliases
		if asIndex := strings.Index(strings.ToUpper(part), " AS "); asIndex != -1 {
			part = part[:asIndex]
		}

		// Handle function calls like SUM(amount), COUNT(*), etc.
		if strings.Contains(part, "(") && strings.Contains(part, ")") {
			// This is likely a function call, extract the function name and arguments
			funcName := v.extractFunctionName(part)
			if funcName != "" {
				// For now, skip function validation as they're handled by allowed_functions config
				continue
			}
		}

		// Handle simple column names
		columnName := strings.TrimSpace(part)
		if columnName != "" && !v.isSQLKeyword(columnName) {
			// Remove table prefixes like "table.column"
			if dotIndex := strings.LastIndex(columnName, "."); dotIndex != -1 {
				columnName = columnName[dotIndex+1:]
			}

			// Only add if it's a valid identifier and not a SQL keyword
			if v.isValidIdentifier(columnName) && !v.isSQLKeyword(columnName) {
				columnNames = append(columnNames, columnName)
			}
		}
	}

	return columnNames
}

// extractFunctionName extracts function name from function call
func (v *SQLValidator) extractFunctionName(part string) string {
	part = strings.TrimSpace(part)

	// Find the opening parenthesis
	openParen := strings.Index(part, "(")
	if openParen == -1 {
		return ""
	}

	// Extract function name
	funcName := strings.TrimSpace(part[:openParen])

	// Check if it's a valid function name
	if v.isValidIdentifier(funcName) {
		return funcName
	}

	return ""
}

// isValidIdentifier checks if a string is a valid SQL identifier
func (v *SQLValidator) isValidIdentifier(name string) bool {
	if name == "" {
		return false
	}

	// Check if it starts with a letter or underscore
	firstChar := rune(name[0])
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		return false
	}

	// Check for valid characters
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '.' {
			return false
		}
	}

	return true
}

// isSQLKeyword checks if a string is a SQL keyword
func (v *SQLValidator) isSQLKeyword(s string) bool {
	upperS := strings.ToUpper(s)

	// Check allowed keywords first
	allowedKeywords := v.config.GetAllowedKeywords()
	for _, keyword := range allowedKeywords {
		if strings.ToUpper(keyword) == upperS {
			return true
		}
	}

	// Check blocked keywords
	blockedKeywords := v.config.GetBlockedKeywords()
	for _, keyword := range blockedKeywords {
		if strings.ToUpper(keyword) == upperS {
			return true
		}
	}

	return false
}

// isContextualKeyword checks if a keyword is allowed in its current context
func (v *SQLValidator) isContextualKeyword(keyword, sql string) bool {
	upperKeyword := strings.ToUpper(keyword)
	upperSQL := strings.ToUpper(sql)

	// Handle CASE WHEN END context
	if upperKeyword == "END" {
		// Check if END is part of a CASE WHEN statement
		caseIndex := strings.Index(upperSQL, "CASE")
		endIndex := strings.Index(upperSQL, "END")

		if caseIndex != -1 && endIndex != -1 && endIndex > caseIndex {
			// Check if there's a WHEN between CASE and END
			whenIndex := strings.Index(upperSQL[caseIndex:endIndex], "WHEN")
			if whenIndex != -1 {
				return true // END is part of CASE WHEN statement
			}
		}
	}

	// Handle BEGIN END context for transactions
	if upperKeyword == "BEGIN" || upperKeyword == "END" {
		// Check if BEGIN/END is part of a transaction
		if strings.Contains(upperSQL, "TRANSACTION") {
			return false // BEGIN/END in transaction context is blocked
		}
	}

	// Default: check if keyword is in allowed list
	allowedKeywords := v.config.GetAllowedKeywords()
	for _, allowed := range allowedKeywords {
		if strings.ToUpper(allowed) == upperKeyword {
			return true
		}
	}

	return false
}

// checkCommentsLoose checks for SQL comments with more lenient rules
func (v *SQLValidator) checkCommentsLoose(sql string) error {
	if strings.Contains(sql, "--") {
		if !v.isInString(sql, "--") {
			return fmt.Errorf("SQL comments not allowed")
		}
	}

	if strings.Contains(sql, "/*") || strings.Contains(sql, "*/") {
		if !v.isInString(sql, "/*") && !v.isInString(sql, "*/") {
			return fmt.Errorf("SQL comments not allowed")
		}
	}

	if strings.Contains(sql, "#") {
		if !v.isInString(sql, "#") {
			return fmt.Errorf("SQL comments not allowed")
		}
	}

	return nil
}

// isInString checks if a pattern is inside a string literal
func (v *SQLValidator) isInString(sql, pattern string) bool {
	index := strings.Index(sql, pattern)
	if index == -1 {
		return false
	}

	before := sql[:index]
	beforeQuotes := strings.Count(before, "'")
	return beforeQuotes%2 == 1
}

// checkBalancedParentheses checks if parentheses are balanced
func (v *SQLValidator) checkBalancedParentheses(sql string) error {
	stack := 0
	for _, char := range sql {
		if char == '(' {
			stack++
		} else if char == ')' {
			stack--
			if stack < 0 {
				return fmt.Errorf("unbalanced parentheses")
			}
		}
	}

	if stack != 0 {
		return fmt.Errorf("unbalanced parentheses")
	}

	return nil
}

// checkComments checks for SQL comments
func (v *SQLValidator) checkComments(sql string) error {
	if strings.Contains(sql, "--") {
		return fmt.Errorf("SQL comments not allowed")
	}

	if strings.Contains(sql, "/*") || strings.Contains(sql, "*/") {
		return fmt.Errorf("SQL comments not allowed")
	}

	if strings.Contains(sql, "#") {
		return fmt.Errorf("SQL comments not allowed")
	}

	return nil
}

// IsReadOnlyQuery checks if SQL query is read-only
func (v *SQLValidator) IsReadOnlyQuery(sql string) bool {
	normalizedSQL := strings.ToUpper(strings.TrimSpace(sql))

	// Must start with SELECT or WITH (for CTEs)
	if !strings.HasPrefix(normalizedSQL, "SELECT") && !strings.HasPrefix(normalizedSQL, "WITH") {
		return false
	}

	// Check for write operations
	writeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE", "EXEC", "EXECUTE"}
	for _, keyword := range writeKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return false
		}
	}

	// Additional checks for dangerous operations
	dangerousKeywords := []string{"GRANT", "REVOKE", "COMMIT", "ROLLBACK", "SAVEPOINT", "LOCK", "UNLOCK"}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return false
		}
	}

	return true
}

// SanitizeSQL sanitizes SQL input
func (v *SQLValidator) SanitizeSQL(sql string) string {
	// Remove comments
	sql = v.compiledPatterns["comment_single"].ReplaceAllString(sql, "")
	sql = v.compiledPatterns["comment_multi"].ReplaceAllString(sql, "")

	// Normalize whitespace
	sql = v.compiledPatterns["whitespace"].ReplaceAllString(sql, " ")

	// Trim spaces
	sql = strings.TrimSpace(sql)

	return sql
}

// ValidateTableName validates table name
func (v *SQLValidator) ValidateTableName(tableName string) error {
	if !v.config.ValidateTableName(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}
	return nil
}

// ValidateColumnName validates column name
func (v *SQLValidator) ValidateColumnName(columnName string) error {
	if !v.config.ValidateColumnName(columnName) {
		return fmt.Errorf("invalid column name: %s", columnName)
	}
	return nil
}

// ValidateColumnNameSmart validates column name with enhanced rules
func (v *SQLValidator) ValidateColumnNameSmart(columnName string) error {
	if columnName == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	// Check for SQL injection patterns (only truly dangerous ones)
	// These patterns are only dangerous if they appear as standalone keywords
	dangerousPatterns := []string{
		"1=1", "TRUE", "FALSE", "UNION", "SELECT", "INSERT", "UPDATE", "DELETE",
		"DROP", "CREATE", "ALTER", "EXEC", "EXECUTE",
	}

	upperColumnName := strings.ToUpper(columnName)
	for _, pattern := range dangerousPatterns {
		// Only flag if the pattern appears as a standalone word, not as part of a larger identifier
		if upperColumnName == pattern {
			return fmt.Errorf("suspicious pattern in column name: %s", pattern)
		}
	}

	// Check for special characters
	for _, char := range columnName {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '.' {
			return fmt.Errorf("invalid character in column name: %c", char)
		}
	}

	// Check if it starts with a letter or underscore
	if len(columnName) > 0 {
		firstChar := rune(columnName[0])
		if !unicode.IsLetter(firstChar) && firstChar != '_' {
			return fmt.Errorf("column name must start with a letter or underscore")
		}
	}

	return nil
}

// GetValidationStats returns validation statistics
func (v *SQLValidator) GetValidationStats() map[string]interface{} {
	configStats := v.config.GetStats()

	return map[string]interface{}{
		"validator_stats": map[string]interface{}{
			"validation_count":  v.validationCount,
			"cache_hit_count":   v.cacheHitCount,
			"last_validation":   v.lastValidation,
			"compiled_patterns": len(v.compiledPatterns),
		},
		"config_stats": configStats,
	}
}

// ReloadConfig reloads the SQL security configuration
func (v *SQLValidator) ReloadConfig(filename string) error {
	return v.config.ReloadConfig(filename)
}

// Convenience functions for backward compatibility
func ValidateSQL(sql string) error {
	return GetGlobalSQLValidator().ValidateSQL(sql)
}

func IsReadOnlyQuery(sql string) bool {
	return GetGlobalSQLValidator().IsReadOnlyQuery(sql)
}

func SanitizeSQL(sql string) string {
	return GetGlobalSQLValidator().SanitizeSQL(sql)
}

func ValidateSQLComplete(sql string) error {
	return GetGlobalSQLValidator().ValidateSQLComplete(sql)
}
