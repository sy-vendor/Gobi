package utils

import (
	"strings"
	"time"

	"fmt"
	"gobi/config"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
)

var (
	QueryCache *cache.Cache
	appConfig  *config.Config
)

// InitQueryCache initializes the query cache with configuration
func InitQueryCache(cfg *config.Config) {
	appConfig = cfg

	defaultExpiration := time.Duration(cfg.Cache.TTL) * time.Second
	cleanupInterval := defaultExpiration * 2

	QueryCache = cache.New(defaultExpiration, cleanupInterval)
}

// GetQueryCache retrieves a value from cache
func GetQueryCache(key string) (interface{}, bool) {
	return QueryCache.Get(key)
}

// SetQueryCache sets a value in cache with smart TTL
func SetQueryCache(key string, value interface{}, sql string) {
	ttl := getSmartTTL(sql)
	QueryCache.Set(key, value, ttl)
}

// DeleteQueryCache deletes a value from cache
func DeleteQueryCache(key string) {
	QueryCache.Delete(key)
}

// getSmartTTL determines cache TTL based on query complexity
func getSmartTTL(sql string) time.Duration {
	if appConfig == nil {
		return 5 * time.Minute // Default fallback
	}

	complexity := analyzeQueryComplexity(sql)

	if complexity == "simple" {
		return time.Duration(appConfig.Cache.Strategy.SimpleQueryTTL) * time.Second
	} else {
		return time.Duration(appConfig.Cache.Strategy.ComplexQueryTTL) * time.Second
	}
}

// analyzeQueryComplexity analyzes SQL query complexity
func analyzeQueryComplexity(sql string) string {
	upperSQL := strings.ToUpper(sql)

	simplePatterns := []string{
		"SELECT COUNT(*)",
		"SELECT * FROM",
		"SELECT ID,",
		"SELECT NAME,",
	}

	complexPatterns := []string{
		"JOIN",
		"UNION",
		"GROUP BY",
		"HAVING",
		"SUBQUERY",
		"EXISTS",
		"IN (",
		"WITH",
		"WINDOW",
		"OVER (",
	}

	for _, pattern := range complexPatterns {
		if strings.Contains(upperSQL, pattern) {
			return "complex"
		}
	}

	for _, pattern := range simplePatterns {
		if strings.Contains(upperSQL, pattern) {
			return "simple"
		}
	}

	return "complex"
}

// GetCacheStats returns cache statistics
func GetCacheStats() map[string]interface{} {
	if QueryCache == nil {
		return map[string]interface{}{
			"enabled": false,
			"error":   "Cache not initialized",
		}
	}

	return map[string]interface{}{
		"enabled":           true,
		"max_cache_size":    appConfig.Cache.Strategy.MaxCacheSize,
		"simple_query_ttl":  appConfig.Cache.Strategy.SimpleQueryTTL,
		"complex_query_ttl": appConfig.Cache.Strategy.ComplexQueryTTL,
		"note":              "Detailed statistics not available with go-cache",
	}
}

// ExecuteSQL connects to the given data source and executes the SQL, returning the result as []map[string]interface{} or error
func ExecuteSQL(ds models.DataSource, sqlStr string) ([]map[string]interface{}, error) {
	sqlStr = SanitizeSQL(sqlStr)

	db, err := database.GetConnection(&ds)
	if err != nil {
		return nil, errors.WrapError(err, "could not get database connection")
	}

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, errors.WrapError(err, "query execution failed")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.WrapError(err, "failed to get column information")
	}

	for _, col := range cols {
		if err := GetGlobalSQLValidator().ValidateColumnNameSmart(col); err != nil {
			return nil, errors.WrapError(err, "invalid column name detected")
		}
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range columns {
			scanArgs[i] = &columns[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, errors.WrapError(err, "failed to scan row")
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			val := columns[i]
			b, ok := val.([]byte)
			if ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WrapError(err, "error during result iteration")
	}

	return results, nil
}

// ExecuteSQLWithTimeout executes SQL with a timeout
func ExecuteSQLWithTimeout(ds models.DataSource, sqlStr string, timeout time.Duration) ([]map[string]interface{}, error) {
	resultChan := make(chan []map[string]interface{}, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := ExecuteSQL(ds, sqlStr)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(timeout):
		return nil, errors.NewError(408, "Query execution timeout", nil)
	}
}

// ExecuteSQLWithLimit executes SQL with a row limit
func ExecuteSQLWithLimit(ds models.DataSource, sqlStr string, limit int) ([]map[string]interface{}, error) {
	if !containsLimit(sqlStr) {
		sqlStr = addLimitClause(sqlStr, limit)
	}

	return ExecuteSQL(ds, sqlStr)
}

func containsLimit(sql string) bool {
	upperSQL := strings.ToUpper(sql)
	return strings.Contains(upperSQL, "LIMIT")
}

func addLimitClause(sql string, limit int) string {
	return fmt.Sprintf("%s LIMIT %d", sql, limit)
}
