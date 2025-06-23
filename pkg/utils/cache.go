package utils

import (
	"time"

	"fmt"
	"gobi/internal/models"
	"gobi/pkg/database"

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
	db, err := database.GetConnection(&ds)
	if err != nil {
		return nil, fmt.Errorf("could not get database connection: %w", err)
	}

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range columns {
			scanArgs[i] = &columns[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
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
	return results, nil
}
