package utils

import (
	"fmt"
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
	return &SQLValidator{
		allowedKeywords: map[string]bool{
			"SELECT": true, "FROM": true, "WHERE": true, "AND": true, "OR": true,
			"ORDER": true, "BY": true, "GROUP": true, "HAVING": true, "LIMIT": true,
			"OFFSET": true, "JOIN": true, "LEFT": true, "RIGHT": true, "INNER": true,
			"OUTER": true, "ON": true, "AS": true, "DISTINCT": true, "COUNT": true,
			"SUM": true, "AVG": true, "MIN": true, "MAX": true, "CASE": true,
			"WHEN": true, "THEN": true, "ELSE": true, "END": true, "IN": true,
			"NOT": true, "LIKE": true, "IS": true, "NULL": true, "BETWEEN": true,
			"ASC": true, "DESC": true, "TOP": true, "FIRST": true, "LAST": true,
		},
		allowedFunctions: map[string]bool{
			"COUNT": true, "SUM": true, "AVG": true, "MIN": true, "MAX": true,
			"UPPER": true, "LOWER": true, "TRIM": true, "LENGTH": true, "SUBSTR": true,
			"CONCAT": true, "COALESCE": true, "NULLIF": true, "ROUND": true,
			"DATE": true, "DATETIME": true, "STRFTIME": true, "JULIANDAY": true,
			"YEAR": true, "MONTH": true, "DAY": true, "HOUR": true, "MINUTE": true,
		},
		blockedKeywords: map[string]bool{
			"DROP": true, "DELETE": true, "UPDATE": true, "INSERT": true, "CREATE": true,
			"ALTER": true, "TRUNCATE": true, "EXEC": true, "EXECUTE": true, "EXECUTE_IMMEDIATE": true,
			"UNION": true, "UNION_ALL": true, "INTERSECT": true, "EXCEPT": true,
			"GRANT": true, "REVOKE": true, "COMMIT": true, "ROLLBACK": true, "SAVEPOINT": true,
			"BEGIN": true, "TRANSACTION": true, "LOCK": true, "UNLOCK": true,
			"SHUTDOWN": true, "KILL": true, "PROCESS": true, "SHOW": true, "DESCRIBE": true,
			"EXPLAIN": true, "ANALYZE": true, "VACUUM": true, "REINDEX": true,
		},
		blockedFunctions: map[string]bool{
			"LOAD_FILE": true, "SLEEP": true, "BENCHMARK": true, "UPDATEXML": true,
			"EXTRACTVALUE": true, "USER": true, "DATABASE": true, "VERSION": true,
			"CONNECTION_ID": true, "LAST_INSERT_ID": true, "ROW_COUNT": true,
		},
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
