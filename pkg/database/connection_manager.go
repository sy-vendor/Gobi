package database

import (
	"database/sql"
	"fmt"
	"gobi/internal/models"
	"sync"
	"time"
)

var (
	connectionPools = make(map[uint]*sql.DB)
	mu              sync.Mutex
)

// GetConnection retrieves a cached database connection pool for a given data source.
// If a pool does not exist, it creates a new one and caches it.
func GetConnection(ds *models.DataSource) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// If a connection pool already exists, return it.
	if db, ok := connectionPools[ds.ID]; ok {
		// Check if the connection is still alive.
		if err := db.Ping(); err == nil {
			return db, nil
		}
		// If not alive, close it and remove from the pool.
		db.Close()
		delete(connectionPools, ds.ID)
	}

	// Create a new connection if not found in the cache.
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

	// Configure the connection pool.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Cache the new connection pool.
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
