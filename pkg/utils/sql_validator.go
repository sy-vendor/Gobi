package utils

import (
	"fmt"
	"gobi/pkg/security"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// SQLValidator provides SQL injection protection
type SQLValidator struct {
	allowedKeywords        map[string]bool
	allowedFunctions       map[string]bool
	blockedKeywords        map[string]bool
	blockedFunctions       map[string]bool
	suspiciousPatterns     []string
	strictColumnValidation bool
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

// NewSQLValidator creates a new SQL validator with default security rules
func NewSQLValidator() *SQLValidator {
	config := security.GetGlobalSQLConfig()

	// Convert slices to maps for faster lookup
	allowedKeywordsMap := make(map[string]bool)
	for _, keyword := range config.GetAllowedKeywords() {
		allowedKeywordsMap[keyword] = true
	}

	allowedFunctionsMap := make(map[string]bool)
	for _, function := range config.GetAllowedFunctions() {
		allowedFunctionsMap[function] = true
	}

	blockedKeywordsMap := make(map[string]bool)
	for _, keyword := range config.GetBlockedKeywords() {
		blockedKeywordsMap[keyword] = true
	}

	blockedFunctionsMap := make(map[string]bool)
	for _, function := range config.GetBlockedFunctions() {
		blockedFunctionsMap[function] = true
	}

	return &SQLValidator{
		allowedKeywords:        allowedKeywordsMap,
		allowedFunctions:       allowedFunctionsMap,
		blockedKeywords:        blockedKeywordsMap,
		blockedFunctions:       blockedFunctionsMap,
		suspiciousPatterns:     config.GetSuspiciousPatterns(),
		strictColumnValidation: false, // 默认关闭严格列名验证
	}
}

// SetStrictColumnValidation sets whether to use strict column name validation
func (v *SQLValidator) SetStrictColumnValidation(strict bool) {
	v.strictColumnValidation = strict
}

// ValidateSQL validates SQL query for security
func (v *SQLValidator) ValidateSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	normalizedSQL := v.normalizeSQL(sql)

	if err := v.checkBlockedKeywords(normalizedSQL); err != nil {
		return err
	}

	if err := v.checkBlockedFunctions(normalizedSQL); err != nil {
		return err
	}

	if err := v.checkSuspiciousPatterns(normalizedSQL); err != nil {
		return err
	}

	if err := v.checkBalancedParentheses(sql); err != nil {
		return err
	}

	if err := v.checkComments(sql); err != nil {
		return err
	}

	return nil
}

// ValidateSQLSmart validates SQL query with more intelligence
func (v *SQLValidator) ValidateSQLSmart(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	if !v.IsReadOnlyQuery(sql) {
		return fmt.Errorf("only SELECT queries are allowed")
	}

	normalizedSQL := v.normalizeSQL(sql)

	writeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE"}
	for _, keyword := range writeKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return fmt.Errorf("write operation not allowed: %s", keyword)
		}
	}

	dangerousKeywords := []string{"EXEC", "EXECUTE", "EXECUTE_IMMEDIATE", "UNION", "UNION_ALL"}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return fmt.Errorf("dangerous operation not allowed: %s", keyword)
		}
	}

	suspiciousPatterns := []string{
		"1=1", "TRUE", "FALSE",
		"OR 1", "OR TRUE", "OR FALSE",
		"AND 1", "AND TRUE", "AND FALSE",
		"';--", "';/*", "';#",
		"UNION SELECT", "UNION ALL SELECT",
		"INFORMATION_SCHEMA",
		"SYSTEM_TABLES",
		"DUAL",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(normalizedSQL, pattern) {
			return fmt.Errorf("suspicious SQL pattern detected: %s", pattern)
		}
	}

	if err := v.checkBalancedParentheses(sql); err != nil {
		return err
	}

	if err := v.checkCommentsLoose(sql); err != nil {
		return err
	}

	return nil
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

// normalizeSQL normalizes SQL for easier parsing
func (v *SQLValidator) normalizeSQL(sql string) string {
	upperSQL := strings.ToUpper(sql)

	upperSQL = regexp.MustCompile(`\s+`).ReplaceAllString(upperSQL, " ")

	return strings.TrimSpace(upperSQL)
}

// checkBlockedKeywords checks for blocked SQL keywords
func (v *SQLValidator) checkBlockedKeywords(sql string) error {
	for keyword := range v.blockedKeywords {
		if strings.Contains(sql, keyword) {
			return fmt.Errorf("blocked SQL keyword detected: %s", keyword)
		}
	}
	return nil
}

// checkBlockedFunctions checks for blocked SQL functions
func (v *SQLValidator) checkBlockedFunctions(sql string) error {
	for function := range v.blockedFunctions {
		if strings.Contains(sql, function+"(") {
			return fmt.Errorf("blocked SQL function detected: %s", function)
		}
	}
	return nil
}

// checkSuspiciousPatterns checks for suspicious SQL patterns
func (v *SQLValidator) checkSuspiciousPatterns(sql string) error {
	for _, pattern := range v.suspiciousPatterns {
		if strings.Contains(sql, pattern) {
			return fmt.Errorf("suspicious SQL pattern detected: %s", pattern)
		}
	}

	return nil
}

// checkBalancedParentheses checks if parentheses are balanced
func (v *SQLValidator) checkBalancedParentheses(sql string) error {
	count := 0
	for _, char := range sql {
		if char == '(' {
			count++
		} else if char == ')' {
			count--
			if count < 0 {
				return fmt.Errorf("unbalanced parentheses in SQL")
			}
		}
	}
	if count != 0 {
		return fmt.Errorf("unbalanced parentheses in SQL")
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

// IsReadOnlyQuery checks if the query is read-only
func (v *SQLValidator) IsReadOnlyQuery(sql string) bool {
	cleanedSQL := v.SanitizeSQL(sql)
	normalizedSQL := v.normalizeSQL(cleanedSQL)

	if !strings.Contains(normalizedSQL, "SELECT") {
		return false
	}

	writeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE"}
	for _, keyword := range writeKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return false
		}
	}

	dangerousKeywords := []string{"EXEC", "EXECUTE", "EXECUTE_IMMEDIATE", "UNION", "UNION_ALL"}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(normalizedSQL, keyword) {
			return false
		}
	}

	return true
}

// SanitizeSQL performs basic SQL sanitization
func (v *SQLValidator) SanitizeSQL(sql string) string {
	sql = strings.ReplaceAll(sql, "\x00", "")

	var sanitized strings.Builder
	for _, char := range sql {
		if unicode.IsPrint(char) || unicode.IsSpace(char) {
			sanitized.WriteRune(char)
		}
	}

	return sanitized.String()
}

// ValidateTableName validates table name for security
func (v *SQLValidator) ValidateTableName(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	sqlInjectionPatterns := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "JAVASCRIPT", "VBSCRIPT",
		"<", ">", "\"", "'", ";", "--", "/*", "*/", "#",
	}

	upperTableName := strings.ToUpper(tableName)
	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(upperTableName, pattern) {
			return fmt.Errorf("invalid table name: contains forbidden pattern '%s'", pattern)
		}
	}

	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "SCRIPT", "JAVASCRIPT", "VBSCRIPT",
	}

	for _, keyword := range sqlKeywords {
		if strings.EqualFold(tableName, keyword) {
			return fmt.Errorf("invalid table name: cannot be SQL keyword '%s'", keyword)
		}
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)
	if !validPattern.MatchString(tableName) {
		return fmt.Errorf("invalid table name format: must start with letter or underscore and contain only alphanumeric characters, underscores, and dots")
	}

	return nil
}

// ValidateColumnNameSmart validates column name with more intelligence
func (v *SQLValidator) ValidateColumnNameSmart(columnName string) error {
	if columnName == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "SCRIPT", "JAVASCRIPT", "VBSCRIPT",
		"FROM", "WHERE", "AND", "OR", "ORDER", "BY", "GROUP", "HAVING",
		"LIMIT", "OFFSET", "JOIN", "LEFT", "RIGHT", "INNER", "OUTER",
		"ON", "AS", "DISTINCT", "COUNT", "SUM", "AVG", "MIN", "MAX",
	}

	for _, keyword := range sqlKeywords {
		if strings.EqualFold(columnName, keyword) {
			return fmt.Errorf("invalid column name: cannot be SQL keyword '%s'", keyword)
		}
	}

	sqlInjectionPatterns := []string{
		"<", ">", "\"", "'", ";", "--", "/*", "*/", "#",
	}

	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(columnName, pattern) {
			return fmt.Errorf("invalid column name: contains forbidden character '%s'", pattern)
		}
	}

	if v.strictColumnValidation {
		validPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.\s\-\(\)\[\]]*$`)
		if !validPattern.MatchString(columnName) {
			if len(columnName) == 0 {
				return fmt.Errorf("column name cannot be empty")
			}

			firstChar := rune(columnName[0])
			if !unicode.IsLetter(firstChar) && firstChar != '_' {
				return fmt.Errorf("column name must start with a letter or underscore, got: %c", firstChar)
			}

			dangerousChars := []rune{'<', '>', '"', '\'', ';', '#', '\\', '\x00'}
			for _, char := range dangerousChars {
				if strings.ContainsRune(columnName, char) {
					return fmt.Errorf("column name contains dangerous character: %c", char)
				}
			}

			return nil
		}
	}

	return nil
}

// ValidateColumnName validates column name for security
func (v *SQLValidator) ValidateColumnName(columnName string) error {
	return v.ValidateColumnNameSmart(columnName)
}

// ValidateSQL is a convenience function using the global validator
func ValidateSQL(sql string) error {
	return GetGlobalSQLValidator().ValidateSQL(sql)
}

// IsReadOnlyQuery is a convenience function using the global validator
func IsReadOnlyQuery(sql string) bool {
	return GetGlobalSQLValidator().IsReadOnlyQuery(sql)
}

// SanitizeSQL is a convenience function using the global validator
func SanitizeSQL(sql string) string {
	return GetGlobalSQLValidator().SanitizeSQL(sql)
}

// ValidateSQLComplete performs complete SQL validation including sanitization, security checks, and read-only validation
func (v *SQLValidator) ValidateSQLComplete(sql string) error {
	if sql == "" {
		return fmt.Errorf("SQL query cannot be empty")
	}

	sanitizedSQL := v.SanitizeSQL(sql)

	if !v.IsReadOnlyQuery(sanitizedSQL) {
		return fmt.Errorf("only SELECT queries are allowed")
	}

	return v.ValidateSQLSmart(sanitizedSQL)
}

// ValidateSQLComplete is a convenience function using the global validator
func ValidateSQLComplete(sql string) error {
	return GetGlobalSQLValidator().ValidateSQLComplete(sql)
}
