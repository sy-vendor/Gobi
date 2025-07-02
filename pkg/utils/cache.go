package utils

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"gobi/config"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
)

// CacheEntry represents a cache entry with metadata
type CacheEntry struct {
	Data        interface{}   `json:"data"`
	CreatedAt   time.Time     `json:"created_at"`
	AccessedAt  time.Time     `json:"accessed_at"`
	AccessCount int64         `json:"access_count"`
	Size        int64         `json:"size"`
	TTL         time.Duration `json:"ttl"`
	Priority    int           `json:"priority"`
	Tags        []string      `json:"tags"`
}

// CacheStats represents detailed cache statistics
type CacheStats struct {
	Enabled           bool      `json:"enabled"`
	TotalItems        int       `json:"total_items"`
	HitCount          int64     `json:"hit_count"`
	MissCount         int64     `json:"miss_count"`
	HitRate           float64   `json:"hit_rate"`
	MemoryUsage       int64     `json:"memory_usage"`
	EvictionCount     int64     `json:"eviction_count"`
	AverageAccessTime float64   `json:"average_access_time"`
	LastReset         time.Time `json:"last_reset"`
}

// CacheManager manages the intelligent cache system
type CacheManager struct {
	primaryCache *cache.Cache
	hotCache     *cache.Cache
	stats        *CacheStats
	config       *config.Config
	mutex        sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

var (
	CacheManagerInstance *CacheManager
	appConfig            *config.Config
)

// InitQueryCache initializes the intelligent cache system
func InitQueryCache(cfg *config.Config) {
	appConfig = cfg
	ctx, cancel := context.WithCancel(context.Background())

	CacheManagerInstance = &CacheManager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
		stats: &CacheStats{
			Enabled:   true,
			LastReset: time.Now(),
		},
	}

	// Initialize primary cache
	defaultExpiration := time.Duration(cfg.Cache.TTL) * time.Second
	cleanupInterval := defaultExpiration * 2
	CacheManagerInstance.primaryCache = cache.New(defaultExpiration, cleanupInterval)

	// Initialize hot cache for frequently accessed items
	hotExpiration := time.Duration(cfg.Cache.TTL/2) * time.Second
	CacheManagerInstance.hotCache = cache.New(hotExpiration, hotExpiration)

	// Start background tasks
	go CacheManagerInstance.startBackgroundTasks()
}

// GetQueryCache retrieves a value from cache with intelligent routing
func GetQueryCache(key string) (interface{}, bool) {
	if CacheManagerInstance == nil {
		return nil, false
	}
	return CacheManagerInstance.Get(key)
}

// SetQueryCache sets a value in cache with intelligent TTL and priority
func SetQueryCache(key string, value interface{}, sql string) {
	if CacheManagerInstance == nil {
		return
	}
	CacheManagerInstance.Set(key, value, sql)
}

// DeleteQueryCache deletes a value from cache
func DeleteQueryCache(key string) {
	if CacheManagerInstance == nil {
		return
	}
	CacheManagerInstance.Delete(key)
}

// Get retrieves a value from cache with intelligent routing
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	start := time.Now()
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Try hot cache first
	if value, found := cm.hotCache.Get(key); found {
		cm.updateStats(true, time.Since(start))
		cm.updateAccessStats(key, true)
		return value, true
	}

	// Try primary cache
	if value, found := cm.primaryCache.Get(key); found {
		cm.updateStats(true, time.Since(start))
		cm.updateAccessStats(key, true)

		// Promote to hot cache if accessed frequently
		if cm.shouldPromoteToHot(key) {
			cm.promoteToHot(key, value)
		}

		return value, true
	}

	cm.updateStats(false, time.Since(start))
	return nil, false
}

// Set sets a value in cache with intelligent TTL and priority
func (cm *CacheManager) Set(key string, value interface{}, sql string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	entry := &CacheEntry{
		Data:        value,
		CreatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 0,
		Size:        cm.calculateSize(value),
		TTL:         cm.calculateSmartTTL(sql),
		Priority:    cm.calculatePriority(sql),
		Tags:        cm.extractTags(sql),
	}

	// Store in primary cache
	cm.primaryCache.Set(key, entry, entry.TTL)

	// If high priority, also store in hot cache
	if entry.Priority >= 8 {
		cm.hotCache.Set(key, entry, entry.TTL/2)
	}
}

// Delete removes a value from both caches
func (cm *CacheManager) Delete(key string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.primaryCache.Delete(key)
	cm.hotCache.Delete(key)
}

// calculateSmartTTL determines cache TTL based on multiple factors
func (cm *CacheManager) calculateSmartTTL(sql string) time.Duration {
	if cm.config == nil {
		return 5 * time.Minute
	}

	complexity := cm.analyzeQueryComplexity(sql)
	baseTTL := cm.getBaseTTL(complexity)

	// Adjust TTL based on query characteristics
	multiplier := 1.0

	// Time-based adjustments
	hour := time.Now().Hour()
	if hour >= 9 && hour <= 17 {
		// Business hours - shorter TTL for real-time data
		multiplier *= 0.8
	} else {
		// Off-hours - longer TTL
		multiplier *= 1.2
	}

	// Query type adjustments
	if strings.Contains(strings.ToUpper(sql), "COUNT") {
		multiplier *= 1.5 // Aggregation queries can be cached longer
	}

	if strings.Contains(strings.ToUpper(sql), "WHERE") {
		multiplier *= 0.9 // Filtered queries might change more frequently
	}

	// Ensure minimum and maximum bounds
	finalTTL := time.Duration(float64(baseTTL) * multiplier)
	if finalTTL < 30*time.Second {
		finalTTL = 30 * time.Second
	}
	if finalTTL > 24*time.Hour {
		finalTTL = 24 * time.Hour
	}

	return finalTTL
}

// getBaseTTL returns base TTL based on query complexity
func (cm *CacheManager) getBaseTTL(complexity string) time.Duration {
	if complexity == "simple" {
		return time.Duration(cm.config.Cache.Strategy.SimpleQueryTTL) * time.Second
	}
	return time.Duration(cm.config.Cache.Strategy.ComplexQueryTTL) * time.Second
}

// analyzeQueryComplexity analyzes SQL query complexity with enhanced logic
func (cm *CacheManager) analyzeQueryComplexity(sql string) string {
	upperSQL := strings.ToUpper(sql)

	// Count complexity indicators
	complexityScore := 0

	complexityPatterns := map[string]int{
		"JOIN":     2,
		"UNION":    3,
		"GROUP BY": 2,
		"HAVING":   2,
		"SUBQUERY": 3,
		"EXISTS":   2,
		"IN (":     1,
		"WITH":     3,
		"WINDOW":   3,
		"OVER (":   3,
		"ORDER BY": 1,
		"LIMIT":    0,
		"OFFSET":   0,
	}

	for pattern, score := range complexityPatterns {
		if strings.Contains(upperSQL, pattern) {
			complexityScore += score
		}
	}

	// Table count estimation
	tableCount := strings.Count(upperSQL, "FROM") + strings.Count(upperSQL, "JOIN")
	if tableCount > 3 {
		complexityScore += tableCount - 3
	}

	if complexityScore >= 5 {
		return "complex"
	} else if complexityScore <= 1 {
		return "simple"
	}
	return "medium"
}

// calculatePriority calculates cache priority based on query characteristics
func (cm *CacheManager) calculatePriority(sql string) int {
	priority := 5 // Base priority

	upperSQL := strings.ToUpper(sql)

	// High priority patterns
	if strings.Contains(upperSQL, "COUNT") {
		priority += 2 // Aggregation queries are often expensive
	}
	if strings.Contains(upperSQL, "GROUP BY") {
		priority += 2
	}
	if strings.Contains(upperSQL, "JOIN") {
		priority += 1
	}

	// Low priority patterns
	if strings.Contains(upperSQL, "LIMIT 1") {
		priority -= 1 // Simple lookups
	}
	if strings.Contains(upperSQL, "ORDER BY") {
		priority += 1 // Sorting can be expensive
	}

	// Ensure priority is within bounds
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}

	return priority
}

// extractTags extracts cache tags from SQL for invalidation
func (cm *CacheManager) extractTags(sql string) []string {
	tags := []string{}
	upperSQL := strings.ToUpper(sql)

	// Extract table names
	if strings.Contains(upperSQL, "FROM") {
		// Simple table extraction (could be enhanced with proper SQL parsing)
		parts := strings.Split(upperSQL, "FROM")
		if len(parts) > 1 {
			tablePart := strings.Fields(parts[1])[0]
			tags = append(tags, "table:"+strings.ToLower(tablePart))
		}
	}

	// Add query type tags
	if strings.Contains(upperSQL, "SELECT") {
		tags = append(tags, "type:select")
	}
	if strings.Contains(upperSQL, "COUNT") {
		tags = append(tags, "type:aggregation")
	}

	return tags
}

// shouldPromoteToHot determines if an item should be promoted to hot cache
func (cm *CacheManager) shouldPromoteToHot(key string) bool {
	// Check access count in primary cache
	if entry, found := cm.primaryCache.Get(key); found {
		if cacheEntry, ok := entry.(*CacheEntry); ok {
			return cacheEntry.AccessCount >= 3 // Promote after 3 accesses
		}
	}
	return false
}

// promoteToHot promotes an item to hot cache
func (cm *CacheManager) promoteToHot(key string, value interface{}) {
	if entry, ok := value.(*CacheEntry); ok {
		entry.AccessedAt = time.Now()
		entry.AccessCount++
		cm.hotCache.Set(key, entry, entry.TTL/2)
	}
}

// updateAccessStats updates access statistics for an item
func (cm *CacheManager) updateAccessStats(key string, hit bool) {
	if entry, found := cm.primaryCache.Get(key); found {
		if cacheEntry, ok := entry.(*CacheEntry); ok {
			cacheEntry.AccessedAt = time.Now()
			cacheEntry.AccessCount++
			cm.primaryCache.Set(key, cacheEntry, cacheEntry.TTL)
		}
	}
}

// updateStats updates cache statistics
func (cm *CacheManager) updateStats(hit bool, accessTime time.Duration) {
	if hit {
		cm.stats.HitCount++
	} else {
		cm.stats.MissCount++
	}

	total := cm.stats.HitCount + cm.stats.MissCount
	if total > 0 {
		cm.stats.HitRate = float64(cm.stats.HitCount) / float64(total) * 100
	}

	// Update average access time
	if cm.stats.HitCount > 0 {
		cm.stats.AverageAccessTime = (cm.stats.AverageAccessTime*float64(cm.stats.HitCount-1) + float64(accessTime.Microseconds())) / float64(cm.stats.HitCount)
	}
}

// calculateSize estimates the size of a value in bytes
func (cm *CacheManager) calculateSize(value interface{}) int64 {
	data, err := json.Marshal(value)
	if err != nil {
		return 0
	}
	return int64(len(data))
}

// startBackgroundTasks starts background cache maintenance tasks
func (cm *CacheManager) startBackgroundTasks() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.performMaintenance()
		case <-cm.ctx.Done():
			return
		}
	}
}

// performMaintenance performs cache maintenance tasks
func (cm *CacheManager) performMaintenance() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Update statistics
	cm.stats.TotalItems = cm.primaryCache.ItemCount() + cm.hotCache.ItemCount()

	// Clean up expired items
	cm.primaryCache.DeleteExpired()
	cm.hotCache.DeleteExpired()

	// Adjust cache sizes based on hit rates
	cm.adjustCacheSizes()
}

// adjustCacheSizes adjusts cache sizes based on performance
func (cm *CacheManager) adjustCacheSizes() {
	// If hit rate is low, reduce cache size to free memory
	if cm.stats.HitRate < 50 && cm.stats.TotalItems > 100 {
		// This is a simplified adjustment - in a real implementation,
		// you might want to implement LRU eviction or other strategies
	}
}

// GetCacheStats returns detailed cache statistics
func GetCacheStats() map[string]interface{} {
	if CacheManagerInstance == nil {
		return map[string]interface{}{
			"enabled": false,
			"error":   "Cache not initialized",
		}
	}

	CacheManagerInstance.mutex.RLock()
	defer CacheManagerInstance.mutex.RUnlock()

	return map[string]interface{}{
		"enabled":             CacheManagerInstance.stats.Enabled,
		"total_items":         CacheManagerInstance.stats.TotalItems,
		"hit_count":           CacheManagerInstance.stats.HitCount,
		"miss_count":          CacheManagerInstance.stats.MissCount,
		"hit_rate":            fmt.Sprintf("%.2f%%", CacheManagerInstance.stats.HitRate),
		"average_access_time": fmt.Sprintf("%.2fμs", CacheManagerInstance.stats.AverageAccessTime),
		"last_reset":          CacheManagerInstance.stats.LastReset,
		"primary_cache_items": CacheManagerInstance.primaryCache.ItemCount(),
		"hot_cache_items":     CacheManagerInstance.hotCache.ItemCount(),
		"config": map[string]interface{}{
			"simple_query_ttl":  CacheManagerInstance.config.Cache.Strategy.SimpleQueryTTL,
			"complex_query_ttl": CacheManagerInstance.config.Cache.Strategy.ComplexQueryTTL,
			"max_cache_size":    CacheManagerInstance.config.Cache.Strategy.MaxCacheSize,
		},
	}
}

// ClearCache clears all caches
func ClearCache() {
	if CacheManagerInstance != nil {
		CacheManagerInstance.mutex.Lock()
		defer CacheManagerInstance.mutex.Unlock()

		CacheManagerInstance.primaryCache.Flush()
		CacheManagerInstance.hotCache.Flush()

		// Reset statistics
		CacheManagerInstance.stats.HitCount = 0
		CacheManagerInstance.stats.MissCount = 0
		CacheManagerInstance.stats.HitRate = 0
		CacheManagerInstance.stats.TotalItems = 0
		CacheManagerInstance.stats.LastReset = time.Now()
	}
}

// WarmupCache preloads frequently accessed data
func WarmupCache(dataSources []models.DataSource) {
	if CacheManagerInstance == nil {
		return
	}

	// Common queries to preload
	commonQueries := []string{
		"SELECT COUNT(*) FROM users",
		"SELECT COUNT(*) FROM reports",
		"SELECT COUNT(*) FROM datasources",
	}

	for _, ds := range dataSources {
		for _, query := range commonQueries {
			key := GenerateCacheKey(ds.ID, query)
			if _, found := CacheManagerInstance.Get(key); !found {
				// Execute query and cache result
				go func(ds models.DataSource, query string) {
					if result, err := ExecuteSQL(ds, query); err == nil {
						CacheManagerInstance.Set(GenerateCacheKey(ds.ID, query), result, query)
					}
				}(ds, query)
			}
		}
	}
}

// generateCacheKey generates a unique cache key
func GenerateCacheKey(datasourceID uint, sql string) string {
	data := fmt.Sprintf("%d:%s", datasourceID, sql)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("query:%x", hash)
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

// ExecuteSQLWithLimit executes SQL with a limit clause
func ExecuteSQLWithLimit(ds models.DataSource, sqlStr string, limit int) ([]map[string]interface{}, error) {
	if !containsLimit(sqlStr) {
		sqlStr = addLimitClause(sqlStr, limit)
	}
	return ExecuteSQL(ds, sqlStr)
}

// containsLimit checks if SQL already contains a LIMIT clause
func containsLimit(sql string) bool {
	return strings.Contains(strings.ToUpper(sql), "LIMIT")
}

// addLimitClause adds a LIMIT clause to SQL
func addLimitClause(sql string, limit int) string {
	return fmt.Sprintf("%s LIMIT %d", sql, limit)
}
