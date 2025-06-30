package utils

import (
	"strings"
	"time"

	"fmt"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
)

var QueryCache *cache.Cache

func InitQueryCache(defaultExpiration, cleanupInterval time.Duration) {
	QueryCache = cache.New(defaultExpiration, cleanupInterval)
}

func GetQueryCache(key string) (interface{}, bool) {
	return QueryCache.Get(key)
}

func SetQueryCache(key string, value interface{}, ttl time.Duration) {
	QueryCache.Set(key, value, ttl)
}

func DeleteQueryCache(key string) {
	QueryCache.Delete(key)
}

// ExecuteSQL connects to the given data source and executes the SQL, returning the result as []map[string]interface{} or error
func ExecuteSQL(ds models.DataSource, sqlStr string) ([]map[string]interface{}, error) {
	sqlStr = SanitizeSQL(sqlStr)

	if err := ValidateSQL(sqlStr); err != nil {
		return nil, errors.WrapError(err, "SQL validation failed")
	}

	if !IsReadOnlyQuery(sqlStr) {
		return nil, errors.NewError(403, "Only SELECT queries are allowed", nil)
	}

	db, err := database.GetConnection(&ds)
	if err != nil {
		return nil, fmt.Errorf("could not get database connection: %w", err)
	}

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column information
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column information: %w", err)
	}

	// Validate column names for security (using smart validation with relaxed mode)
	for _, col := range cols {
		if err := GlobalSQLValidator.ValidateColumnNameSmart(col); err != nil {
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
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
		return nil, fmt.Errorf("error during result iteration: %w", err)
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
