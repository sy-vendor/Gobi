package security

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// SQLPatterns defines SQL security patterns
type SQLPatterns struct {
	AllowedKeywords    []string `yaml:"allowed_keywords"`
	BlockedKeywords    []string `yaml:"blocked_keywords"`
	AllowedFunctions   []string `yaml:"allowed_functions"`
	BlockedFunctions   []string `yaml:"blocked_functions"`
	SuspiciousPatterns []string `yaml:"suspicious_patterns"`
}

// SQLSecurityConfig provides centralized SQL security configuration
type SQLSecurityConfig struct {
	patterns *SQLPatterns
	mu       sync.RWMutex
}

var (
	globalSQLConfig *SQLSecurityConfig
	configOnce      sync.Once
)

// GetGlobalSQLConfig returns a singleton instance of SQLSecurityConfig
func GetGlobalSQLConfig() *SQLSecurityConfig {
	configOnce.Do(func() {
		globalSQLConfig = &SQLSecurityConfig{}
		globalSQLConfig.LoadFromFile("config/sql_whitelist.yaml")
	})
	return globalSQLConfig
}

// LoadFromFile loads SQL patterns from a YAML file
func (c *SQLSecurityConfig) LoadFromFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		// Fallback to default patterns if file not found
		c.patterns = getDefaultPatterns()
		return nil
	}

	var patterns SQLPatterns
	if err := yaml.Unmarshal(data, &patterns); err != nil {
		c.patterns = getDefaultPatterns()
		return err
	}

	c.patterns = &patterns
	return nil
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
	}
}
