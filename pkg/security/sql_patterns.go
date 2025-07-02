package security

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// SQLPatterns defines SQL security patterns
type SQLPatterns struct {
	AllowedKeywords       []string `yaml:"allowed_keywords"`
	BlockedKeywords       []string `yaml:"blocked_keywords"`
	AllowedFunctions      []string `yaml:"allowed_functions"`
	BlockedFunctions      []string `yaml:"blocked_functions"`
	SuspiciousPatterns    []string `yaml:"suspicious_patterns"`
	AllowedTablePatterns  []string `yaml:"allowed_table_patterns"`
	AllowedColumnPatterns []string `yaml:"allowed_column_patterns"`
	QueryLimits           struct {
		MaxExecutionTime string `yaml:"max_execution_time"`
		MaxRows          int    `yaml:"max_rows"`
		MaxQueryLength   int    `yaml:"max_query_length"`
	} `yaml:"query_limits"`
	Security struct {
		AllowComments           bool `yaml:"allow_comments"`
		AllowMultipleStatements bool `yaml:"allow_multiple_statements"`
		RequireReadonly         bool `yaml:"require_readonly"`
		ValidateTableNames      bool `yaml:"validate_table_names"`
		ValidateColumnNames     bool `yaml:"validate_column_names"`
		SanitizeInput           bool `yaml:"sanitize_input"`
	} `yaml:"security"`
}

// SQLSecurityConfig provides centralized SQL security configuration with caching
type SQLSecurityConfig struct {
	patterns *SQLPatterns
	mu       sync.RWMutex

	// Cached compiled patterns for performance
	allowedKeywordsMap    map[string]bool
	blockedKeywordsMap    map[string]bool
	allowedFunctionsMap   map[string]bool
	blockedFunctionsMap   map[string]bool
	suspiciousPatternsMap map[string]bool

	// Compiled regex patterns
	tableNameRegex  *regexp.Regexp
	columnNameRegex *regexp.Regexp

	// Cache for validation results
	validationCache map[string]validationResult
	cacheMu         sync.RWMutex
	cacheExpiry     time.Duration

	// Performance metrics
	validationCount int64
	cacheHitCount   int64
	lastReload      time.Time
}

// validationResult represents a cached validation result
type validationResult struct {
	Valid   bool
	Error   string
	Expires time.Time
}

var (
	globalSQLConfig *SQLSecurityConfig
	configOnce      sync.Once
)

// GetGlobalSQLConfig returns a singleton instance of SQLSecurityConfig
func GetGlobalSQLConfig() *SQLSecurityConfig {
	configOnce.Do(func() {
		globalSQLConfig = &SQLSecurityConfig{
			validationCache: make(map[string]validationResult),
			cacheExpiry:     5 * time.Minute, // Cache validation results for 5 minutes
		}
		globalSQLConfig.LoadFromFile("config/sql_whitelist.yaml")

		// Start cache cleanup goroutine
		go globalSQLConfig.cleanupCache()
	})
	return globalSQLConfig
}

// LoadFromFile loads SQL patterns from a YAML file with caching
func (c *SQLSecurityConfig) LoadFromFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		// Fallback to default patterns if file not found
		c.patterns = getDefaultPatterns()
		c.compilePatterns()
		return nil
	}

	var patterns SQLPatterns
	if err := yaml.Unmarshal(data, &patterns); err != nil {
		c.patterns = getDefaultPatterns()
		c.compilePatterns()
		return err
	}

	c.patterns = &patterns
	c.compilePatterns()
	c.lastReload = time.Now()

	// Clear validation cache when config is reloaded
	c.clearValidationCache()

	return nil
}

// compilePatterns compiles patterns for faster lookup
func (c *SQLSecurityConfig) compilePatterns() {
	// Compile keyword maps
	c.allowedKeywordsMap = make(map[string]bool)
	for _, keyword := range c.patterns.AllowedKeywords {
		c.allowedKeywordsMap[strings.ToUpper(keyword)] = true
	}

	c.blockedKeywordsMap = make(map[string]bool)
	for _, keyword := range c.patterns.BlockedKeywords {
		c.blockedKeywordsMap[strings.ToUpper(keyword)] = true
	}

	c.allowedFunctionsMap = make(map[string]bool)
	for _, function := range c.patterns.AllowedFunctions {
		c.allowedFunctionsMap[strings.ToUpper(function)] = true
	}

	c.blockedFunctionsMap = make(map[string]bool)
	for _, function := range c.patterns.BlockedFunctions {
		c.blockedFunctionsMap[strings.ToUpper(function)] = true
	}

	c.suspiciousPatternsMap = make(map[string]bool)
	for _, pattern := range c.patterns.SuspiciousPatterns {
		c.suspiciousPatternsMap[strings.ToUpper(pattern)] = true
	}

	// Compile regex patterns
	if len(c.patterns.AllowedTablePatterns) > 0 {
		c.tableNameRegex = regexp.MustCompile(c.patterns.AllowedTablePatterns[0])
	}
	if len(c.patterns.AllowedColumnPatterns) > 0 {
		c.columnNameRegex = regexp.MustCompile(c.patterns.AllowedColumnPatterns[0])
	}
}

// clearValidationCache clears the validation cache
func (c *SQLSecurityConfig) clearValidationCache() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.validationCache = make(map[string]validationResult)
}

// cleanupCache periodically cleans up expired cache entries
func (c *SQLSecurityConfig) cleanupCache() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cacheMu.Lock()
		now := time.Now()
		for key, result := range c.validationCache {
			if now.After(result.Expires) {
				delete(c.validationCache, key)
			}
		}
		c.cacheMu.Unlock()
	}
}

// ValidateSQLWithCache validates SQL with caching for performance
func (c *SQLSecurityConfig) ValidateSQLWithCache(sql string) (bool, string) {
	// Check cache first
	c.cacheMu.RLock()
	if result, exists := c.validationCache[sql]; exists && time.Now().Before(result.Expires) {
		c.cacheMu.RUnlock()
		c.cacheHitCount++
		return result.Valid, result.Error
	}
	c.cacheMu.RUnlock()

	// Perform validation
	valid, errorMsg := c.validateSQLInternal(sql)

	// Cache result
	c.cacheMu.Lock()
	c.validationCache[sql] = validationResult{
		Valid:   valid,
		Error:   errorMsg,
		Expires: time.Now().Add(c.cacheExpiry),
	}
	c.cacheMu.Unlock()

	c.validationCount++
	return valid, errorMsg
}

// validateSQLInternal performs the actual SQL validation
func (c *SQLSecurityConfig) validateSQLInternal(sql string) (bool, string) {
	if sql == "" {
		return false, "SQL query cannot be empty"
	}

	normalizedSQL := strings.ToUpper(strings.TrimSpace(sql))

	// Check blocked keywords
	for keyword := range c.blockedKeywordsMap {
		if strings.Contains(normalizedSQL, keyword) {
			return false, "blocked keyword detected: " + keyword
		}
	}

	// Check blocked functions
	for function := range c.blockedFunctionsMap {
		if strings.Contains(normalizedSQL, function) {
			return false, "blocked function detected: " + function
		}
	}

	// Check suspicious patterns
	for pattern := range c.suspiciousPatternsMap {
		if strings.Contains(normalizedSQL, pattern) {
			return false, "suspicious pattern detected: " + pattern
		}
	}

	// Check query length
	if c.patterns.QueryLimits.MaxQueryLength > 0 && len(sql) > c.patterns.QueryLimits.MaxQueryLength {
		return false, "query too long"
	}

	// Check if read-only is required
	if c.patterns.Security.RequireReadonly {
		writeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE"}
		for _, keyword := range writeKeywords {
			if strings.Contains(normalizedSQL, keyword) {
				return false, "write operations not allowed"
			}
		}
	}

	return true, ""
}

// GetAllowedKeywords returns allowed SQL keywords
func (c *SQLSecurityConfig) GetAllowedKeywords() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.patterns.AllowedKeywords
}

// GetBlockedKeywords returns blocked SQL keywords
func (c *SQLSecurityConfig) GetBlockedKeywords() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.patterns.BlockedKeywords
}

// GetAllowedFunctions returns allowed SQL functions
func (c *SQLSecurityConfig) GetAllowedFunctions() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.patterns.AllowedFunctions
}

// GetBlockedFunctions returns blocked SQL functions
func (c *SQLSecurityConfig) GetBlockedFunctions() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.patterns.BlockedFunctions
}

// GetSuspiciousPatterns returns suspicious SQL patterns
func (c *SQLSecurityConfig) GetSuspiciousPatterns() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.patterns.SuspiciousPatterns
}

// ValidateTableName validates table name using regex pattern
func (c *SQLSecurityConfig) ValidateTableName(tableName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.tableNameRegex == nil {
		return true // If no pattern defined, allow all
	}

	return c.tableNameRegex.MatchString(tableName)
}

// ValidateColumnName validates column name using regex pattern
func (c *SQLSecurityConfig) ValidateColumnName(columnName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.columnNameRegex == nil {
		return true // If no pattern defined, allow all
	}

	return c.columnNameRegex.MatchString(columnName)
}

// GetQueryLimits returns query execution limits
func (c *SQLSecurityConfig) GetQueryLimits() (maxExecutionTime string, maxRows int, maxQueryLength int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.patterns.QueryLimits.MaxExecutionTime,
		c.patterns.QueryLimits.MaxRows,
		c.patterns.QueryLimits.MaxQueryLength
}

// GetSecuritySettings returns security settings
func (c *SQLSecurityConfig) GetSecuritySettings() map[string]bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]bool{
		"allow_comments":            c.patterns.Security.AllowComments,
		"allow_multiple_statements": c.patterns.Security.AllowMultipleStatements,
		"require_readonly":          c.patterns.Security.RequireReadonly,
		"validate_table_names":      c.patterns.Security.ValidateTableNames,
		"validate_column_names":     c.patterns.Security.ValidateColumnNames,
		"sanitize_input":            c.patterns.Security.SanitizeInput,
	}
}

// GetStats returns validation statistics
func (c *SQLSecurityConfig) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.cacheMu.RLock()
	cacheSize := len(c.validationCache)
	c.cacheMu.RUnlock()

	hitRate := 0.0
	if c.validationCount > 0 {
		hitRate = float64(c.cacheHitCount) / float64(c.validationCount) * 100
	}

	return map[string]interface{}{
		"validation_count": c.validationCount,
		"cache_hit_count":  c.cacheHitCount,
		"cache_hit_rate":   hitRate,
		"cache_size":       cacheSize,
		"last_reload":      c.lastReload,
		"cache_expiry":     c.cacheExpiry.String(),
	}
}

// ReloadConfig reloads configuration from file
func (c *SQLSecurityConfig) ReloadConfig(filename string) error {
	return c.LoadFromFile(filename)
}

// getDefaultPatterns returns default SQL patterns if config file is not available
func getDefaultPatterns() *SQLPatterns {
	return &SQLPatterns{
		AllowedKeywords: []string{
			"SELECT", "FROM", "WHERE", "AND", "OR", "ORDER", "BY", "GROUP", "HAVING",
			"LIMIT", "OFFSET", "JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "ON", "AS",
			"DISTINCT", "COUNT", "SUM", "AVG", "MIN", "MAX", "CASE", "WHEN", "THEN",
			"ELSE", "END", "IN", "NOT", "LIKE", "IS", "NULL", "BETWEEN", "ASC", "DESC",
		},
		BlockedKeywords: []string{
			"DROP", "DELETE", "UPDATE", "INSERT", "CREATE", "ALTER", "TRUNCATE",
			"EXEC", "EXECUTE", "UNION", "GRANT", "REVOKE", "COMMIT", "ROLLBACK",
		},
		AllowedFunctions: []string{
			"COUNT", "SUM", "AVG", "MIN", "MAX", "UPPER", "LOWER", "TRIM", "LENGTH",
			"SUBSTR", "CONCAT", "COALESCE", "NULLIF", "ROUND", "DATE", "DATETIME",
		},
		BlockedFunctions: []string{
			"LOAD_FILE", "SLEEP", "BENCHMARK", "UPDATEXML", "EXTRACTVALUE",
		},
		SuspiciousPatterns: []string{
			"1=1", "TRUE", "FALSE", "OR 1", "AND 1", "';--", "';/*", "';#",
			"UNION SELECT", "INFORMATION_SCHEMA", "SYSTEM_TABLES", "DUAL",
		},
		AllowedTablePatterns:  []string{"^[a-zA-Z_][a-zA-Z0-9_.]*$"},
		AllowedColumnPatterns: []string{"^[a-zA-Z_][a-zA-Z0-9_.]*$"},
		QueryLimits: struct {
			MaxExecutionTime string `yaml:"max_execution_time"`
			MaxRows          int    `yaml:"max_rows"`
			MaxQueryLength   int    `yaml:"max_query_length"`
		}{
			MaxExecutionTime: "30s",
			MaxRows:          10000,
			MaxQueryLength:   10000,
		},
		Security: struct {
			AllowComments           bool `yaml:"allow_comments"`
			AllowMultipleStatements bool `yaml:"allow_multiple_statements"`
			RequireReadonly         bool `yaml:"require_readonly"`
			ValidateTableNames      bool `yaml:"validate_table_names"`
			ValidateColumnNames     bool `yaml:"validate_column_names"`
			SanitizeInput           bool `yaml:"sanitize_input"`
		}{
			AllowComments:           false,
			AllowMultipleStatements: false,
			RequireReadonly:         true,
			ValidateTableNames:      true,
			ValidateColumnNames:     true,
			SanitizeInput:           true,
		},
	}
}
