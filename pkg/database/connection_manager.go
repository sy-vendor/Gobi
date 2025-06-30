package database

import (
	"database/sql"
	"fmt"
	"gobi/config"
	"gobi/internal/models"
	"sync"
	"time"
)

var (
	connectionPools = make(map[uint]*sql.DB)
	mu              sync.Mutex
	appConfig       *config.Config
)

// InitConnectionManager initializes the connection manager with configuration
func InitConnectionManager(cfg *config.Config) {
	appConfig = cfg
}

// GetConnection retrieves a cached database connection pool for a given data source.
// If a pool does not exist, it creates a new one and caches it.
func GetConnection(ds *models.DataSource) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if db, ok := connectionPools[ds.ID]; ok {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		db.Close()
		delete(connectionPools, ds.ID)
	}

	var dsn, driver string
	switch ds.Type {
	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", ds.Username, ds.Password, ds.Host, ds.Port, ds.Database)
	case "postgres":
		driver = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", ds.Host, ds.Port, ds.Username, ds.Password, ds.Database)
	case "sqlite":
		driver = "sqlite3"
		dsn = ds.Database
	default:
		return nil, fmt.Errorf("unsupported data source type: %s", ds.Type)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if appConfig != nil {
		pool := appConfig.Database.ConnectionPool
		db.SetMaxOpenConns(pool.MaxOpenConns)
		db.SetMaxIdleConns(pool.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(pool.ConnMaxLifetime) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(pool.ConnMaxIdleTime) * time.Second)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)
	}

	connectionPools[ds.ID] = db

	return db, nil
}

// CloseAllConnections closes all cached database connection pools.
// This should be called on application shutdown.
func CloseAllConnections() {
	mu.Lock()
	defer mu.Unlock()

	for _, db := range connectionPools {
		db.Close()
	}
	connectionPools = make(map[uint]*sql.DB) // Clear the map
}

// GetConnectionStats returns statistics about connection pools
func GetConnectionStats() map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()

	stats := make(map[string]interface{})
	stats["total_pools"] = len(connectionPools)

	poolDetails := make(map[uint]map[string]interface{})
	for id, db := range connectionPools {
		poolDetails[id] = map[string]interface{}{
			"max_open_connections": db.Stats().MaxOpenConnections,
			"open_connections":     db.Stats().OpenConnections,
			"in_use":               db.Stats().InUse,
			"idle":                 db.Stats().Idle,
		}
	}
	stats["pool_details"] = poolDetails

	return stats
}
