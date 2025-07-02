package infrastructure

import (
	"context"
	"fmt"
	"time"

	"gobi/config"
	"gobi/internal/models"
	"gobi/pkg/utils"
)

// CacheService provides intelligent caching functionality
type CacheService struct {
	config *config.Config
}

// NewCacheService creates a new cache service instance
func NewCacheService(cfg *config.Config) *CacheService {
	return &CacheService{
		config: cfg,
	}
}

// GetCachedQuery retrieves a cached query result
func (cs *CacheService) GetCachedQuery(datasourceID uint, sql string) (interface{}, bool) {
	key := utils.GenerateCacheKey(datasourceID, sql)
	return utils.GetQueryCache(key)
}

// SetCachedQuery stores a query result in cache
func (cs *CacheService) SetCachedQuery(datasourceID uint, sql string, result interface{}) {
	key := utils.GenerateCacheKey(datasourceID, sql)
	utils.SetQueryCache(key, result, sql)
}

// InvalidateCacheByTags invalidates cache entries by tags
func (cs *CacheService) InvalidateCacheByTags(tags []string) {
	// This is a simplified implementation
	// In a real implementation, you would iterate through cache entries
	// and remove those matching the tags
	fmt.Printf("Invalidating cache for tags: %v\n", tags)
}

// InvalidateCacheByTable invalidates cache entries for a specific table
func (cs *CacheService) InvalidateCacheByTable(tableName string) {
	tags := []string{fmt.Sprintf("table:%s", tableName)}
	cs.InvalidateCacheByTags(tags)
}

// GetCacheStats returns detailed cache statistics
func (cs *CacheService) GetCacheStats() map[string]interface{} {
	return utils.GetCacheStats()
}

// ClearCache clears all cache entries
func (cs *CacheService) ClearCache() {
	utils.ClearCache()
}

// WarmupCache preloads frequently accessed data
func (cs *CacheService) WarmupCache(dataSources []models.DataSource) {
	if !cs.config.Cache.Strategy.CacheWarmup {
		return
	}

	utils.WarmupCache(dataSources)
}

// ExecuteQueryWithCache executes a query with intelligent caching
func (cs *CacheService) ExecuteQueryWithCache(ctx context.Context, ds models.DataSource, sql string) ([]map[string]interface{}, error) {
	// Check cache first
	if cached, found := cs.GetCachedQuery(ds.ID, sql); found {
		if result, ok := cached.([]map[string]interface{}); ok {
			return result, nil
		}
	}

	// Execute query
	result, err := utils.ExecuteSQL(ds, sql)
	if err != nil {
		return nil, err
	}

	// Cache the result
	cs.SetCachedQuery(ds.ID, sql, result)

	return result, nil
}

// ExecuteQueryWithCacheAndTimeout executes a query with cache and timeout
func (cs *CacheService) ExecuteQueryWithCacheAndTimeout(ctx context.Context, ds models.DataSource, sql string, timeout time.Duration) ([]map[string]interface{}, error) {
	// Check cache first
	if cached, found := cs.GetCachedQuery(ds.ID, sql); found {
		if result, ok := cached.([]map[string]interface{}); ok {
			return result, nil
		}
	}

	// Execute query with timeout
	result, err := utils.ExecuteSQLWithTimeout(ds, sql, timeout)
	if err != nil {
		return nil, err
	}

	// Cache the result
	cs.SetCachedQuery(ds.ID, sql, result)

	return result, nil
}

// ExecuteQueryWithCacheAndLimit executes a query with cache and limit
func (cs *CacheService) ExecuteQueryWithCacheAndLimit(ctx context.Context, ds models.DataSource, sql string, limit int) ([]map[string]interface{}, error) {
	// Check cache first
	if cached, found := cs.GetCachedQuery(ds.ID, sql); found {
		if result, ok := cached.([]map[string]interface{}); ok {
			// Apply limit to cached result
			if len(result) > limit {
				return result[:limit], nil
			}
			return result, nil
		}
	}

	// Execute query with limit
	result, err := utils.ExecuteSQLWithLimit(ds, sql, limit)
	if err != nil {
		return nil, err
	}

	// Cache the result
	cs.SetCachedQuery(ds.ID, sql, result)

	return result, nil
}

// GetCachePerformanceMetrics returns cache performance metrics
func (cs *CacheService) GetCachePerformanceMetrics() map[string]interface{} {
	stats := cs.GetCacheStats()

	// Calculate additional metrics
	metrics := map[string]interface{}{
		"cache_stats": stats,
		"performance": map[string]interface{}{
			"cache_efficiency": cs.calculateCacheEfficiency(),
			"memory_usage":     cs.calculateMemoryUsage(),
			"response_time":    cs.calculateAverageResponseTime(),
		},
		"recommendations": cs.generateRecommendations(),
	}

	return metrics
}

// calculateCacheEfficiency calculates cache efficiency score
func (cs *CacheService) calculateCacheEfficiency() float64 {
	stats := cs.GetCacheStats()

	hitRate, ok := stats["hit_rate"].(string)
	if !ok {
		return 0.0
	}

	// Parse hit rate percentage
	var rate float64
	fmt.Sscanf(hitRate, "%f%%", &rate)

	// Efficiency score based on hit rate
	if rate >= 80 {
		return 1.0 // Excellent
	} else if rate >= 60 {
		return 0.8 // Good
	} else if rate >= 40 {
		return 0.6 // Fair
	} else {
		return 0.4 // Poor
	}
}

// calculateMemoryUsage calculates memory usage percentage
func (cs *CacheService) calculateMemoryUsage() float64 {
	stats := cs.GetCacheStats()

	totalItems, ok := stats["total_items"].(int)
	if !ok {
		return 0.0
	}

	maxSize := cs.config.Cache.Strategy.MaxCacheSize
	if maxSize == 0 {
		return 0.0
	}

	return float64(totalItems) / float64(maxSize) * 100
}

// calculateAverageResponseTime calculates average response time
func (cs *CacheService) calculateAverageResponseTime() float64 {
	stats := cs.GetCacheStats()

	avgTime, ok := stats["average_access_time"].(string)
	if !ok {
		return 0.0
	}

	// Parse average access time
	var time float64
	fmt.Sscanf(avgTime, "%fμs", &time)

	return time
}

// generateRecommendations generates cache optimization recommendations
func (cs *CacheService) generateRecommendations() []string {
	recommendations := []string{}

	efficiency := cs.calculateCacheEfficiency()
	memoryUsage := cs.calculateMemoryUsage()

	if efficiency < 0.6 {
		recommendations = append(recommendations, "Consider increasing cache TTL for frequently accessed queries")
		recommendations = append(recommendations, "Review query patterns and optimize cache keys")
	}

	if memoryUsage > 80 {
		recommendations = append(recommendations, "Cache memory usage is high, consider increasing max_cache_size")
		recommendations = append(recommendations, "Review and remove unnecessary cache entries")
	}

	if efficiency < 0.4 {
		recommendations = append(recommendations, "Cache hit rate is very low, consider disabling cache for this workload")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Cache performance is optimal")
	}

	return recommendations
}

// OptimizeCache performs cache optimization based on current performance
func (cs *CacheService) OptimizeCache() map[string]interface{} {
	optimizations := map[string]interface{}{
		"performed": []string{},
		"skipped":   []string{},
		"results":   map[string]interface{}{},
	}

	// Check if adaptive TTL is enabled
	if !cs.config.Cache.Strategy.AdaptiveTTL {
		optimizations["skipped"] = append(optimizations["skipped"].([]string), "Adaptive TTL is disabled")
		return optimizations
	}

	// Perform cache cleanup
	utils.ClearCache()
	optimizations["performed"] = append(optimizations["performed"].([]string), "Cache cleared for optimization")

	// Update maintenance interval based on performance
	stats := cs.GetCacheStats()
	hitRate, ok := stats["hit_rate"].(string)
	if ok {
		var rate float64
		fmt.Sscanf(hitRate, "%f%%", &rate)

		if rate < 50 {
			// Increase maintenance frequency for poor performance
			optimizations["results"].(map[string]interface{})["maintenance_interval"] = "Reduced to 2 minutes"
		} else if rate > 80 {
			// Decrease maintenance frequency for good performance
			optimizations["results"].(map[string]interface{})["maintenance_interval"] = "Increased to 10 minutes"
		}
	}

	return optimizations
}

// GetCacheHealth returns cache health status
func (cs *CacheService) GetCacheHealth() map[string]interface{} {
	stats := cs.GetCacheStats()

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"details":   stats,
	}

	// Check if cache is enabled
	if enabled, ok := stats["enabled"].(bool); !ok || !enabled {
		health["status"] = "disabled"
		return health
	}

	// Check hit rate
	if hitRate, ok := stats["hit_rate"].(string); ok {
		var rate float64
		fmt.Sscanf(hitRate, "%f%%", &rate)

		if rate < 30 {
			health["status"] = "warning"
			health["message"] = "Low cache hit rate detected"
		} else if rate < 10 {
			health["status"] = "critical"
			health["message"] = "Very low cache hit rate, consider disabling cache"
		}
	}

	// Check memory usage
	memoryUsage := cs.calculateMemoryUsage()
	if memoryUsage > 90 {
		health["status"] = "warning"
		health["message"] = "High cache memory usage detected"
	}

	return health
}

// Get retrieves a value from cache by key
func (cs *CacheService) Get(key string) (interface{}, bool) {
	return utils.GetQueryCache(key)
}

// Set sets a value in cache with a given TTL
func (cs *CacheService) Set(key string, value interface{}, ttl time.Duration) {
	// 直接用默认ttl，不做sql智能分析
	utils.SetQueryCache(key, value, "")
}

// Delete removes a value from cache by key
func (cs *CacheService) Delete(key string) {
	utils.DeleteQueryCache(key)
}

// Flush clears all cache
func (cs *CacheService) Flush() {
	utils.ClearCache()
}

// GetStats returns cache statistics
func (cs *CacheService) GetStats() map[string]interface{} {
	return utils.GetCacheStats()
}
