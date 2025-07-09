package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/utils"
)

// OptimizedSQLExecutionService provides optimized SQL execution with performance monitoring
type OptimizedSQLExecutionService struct {
	optimizer    *database.QueryOptimizer
	cacheService *CacheService
	mu           sync.RWMutex
	stats        *ExecutionStats
}

// ExecutionStats tracks execution statistics
type ExecutionStats struct {
	TotalExecutions int64         `json:"total_executions"`
	CacheHits       int64         `json:"cache_hits"`
	CacheMisses     int64         `json:"cache_misses"`
	AverageTime     time.Duration `json:"average_time"`
	SlowQueries     int64         `json:"slow_queries"`
	FailedQueries   int64         `json:"failed_queries"`
	LastReset       time.Time     `json:"last_reset"`
}

// ExecutionResult represents the result of an optimized query execution
type ExecutionResult struct {
	Data          []map[string]interface{} `json:"data"`
	ExecutionTime time.Duration            `json:"execution_time"`
	CacheHit      bool                     `json:"cache_hit"`
	QueryPlan     *database.QueryPlan      `json:"query_plan,omitempty"`
	Error         error                    `json:"error,omitempty"`
}

// NewOptimizedSQLExecutionService creates a new optimized SQL execution service
func NewOptimizedSQLExecutionService(cacheService *CacheService) *OptimizedSQLExecutionService {
	return &OptimizedSQLExecutionService{
		optimizer:    database.NewQueryOptimizer(),
		cacheService: cacheService,
		stats: &ExecutionStats{
			LastReset: time.Now(),
		},
	}
}

// ExecuteWithOptimization executes a query with full optimization analysis
func (s *OptimizedSQLExecutionService) ExecuteWithOptimization(ctx context.Context, ds models.DataSource, sql string) (*ExecutionResult, error) {
	startTime := time.Now()

	// Check cache first
	cacheKey := s.generateCacheKey(ds.ID, sql)
	if cached, found := s.cacheService.Get(cacheKey); found {
		s.updateStats(true, time.Since(startTime), false)
		if result, ok := cached.([]map[string]interface{}); ok {
			return &ExecutionResult{
				Data:          result,
				ExecutionTime: time.Since(startTime),
				CacheHit:      true,
			}, nil
		}
	}

	// Execute with optimization analysis
	results, plan, err := s.optimizer.ExecuteWithOptimization(ctx, ds, sql)

	executionTime := time.Since(startTime)
	s.updateStats(false, executionTime, err != nil)

	if err != nil {
		return &ExecutionResult{
			ExecutionTime: executionTime,
			Error:         err,
		}, err
	}

	// Cache results if successful
	s.cacheResults(cacheKey, results, sql)

	return &ExecutionResult{
		Data:          results,
		ExecutionTime: executionTime,
		CacheHit:      false,
		QueryPlan:     plan,
	}, nil
}

// ExecuteWithTimeout executes a query with timeout and optimization
func (s *OptimizedSQLExecutionService) ExecuteWithTimeout(ctx context.Context, ds models.DataSource, sql string, timeout time.Duration) (*ExecutionResult, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Check cache first
	cacheKey := s.generateCacheKey(ds.ID, sql)
	if cached, found := s.cacheService.Get(cacheKey); found {
		if result, ok := cached.([]map[string]interface{}); ok {
			return &ExecutionResult{
				Data:          result,
				ExecutionTime: 0,
				CacheHit:      true,
			}, nil
		}
	}

	// Execute with timeout
	startTime := time.Now()
	results, plan, err := s.optimizer.ExecuteWithOptimization(ctx, ds, sql)
	executionTime := time.Since(startTime)

	s.updateStats(false, executionTime, err != nil)

	if err != nil {
		return &ExecutionResult{
			ExecutionTime: executionTime,
			Error:         err,
		}, err
	}

	// Cache results
	s.cacheResults(cacheKey, results, sql)

	return &ExecutionResult{
		Data:          results,
		ExecutionTime: executionTime,
		CacheHit:      false,
		QueryPlan:     plan,
	}, nil
}

// ExecuteWithLimit executes a query with row limit and optimization
func (s *OptimizedSQLExecutionService) ExecuteWithLimit(ctx context.Context, ds models.DataSource, sql string, limit int) (*ExecutionResult, error) {
	// Check if SQL already has LIMIT
	if !s.containsLimit(sql) {
		sql = s.addLimitClause(sql, limit)
	}

	return s.ExecuteWithOptimization(ctx, ds, sql)
}

// ExecuteBatch executes multiple queries in batch with optimization
func (s *OptimizedSQLExecutionService) ExecuteBatch(ctx context.Context, ds models.DataSource, queries []string) ([]*ExecutionResult, error) {
	results := make([]*ExecutionResult, len(queries))

	// Execute queries concurrently with semaphore to limit concurrency
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent queries
	var wg sync.WaitGroup

	for i, query := range queries {
		wg.Add(1)
		go func(index int, sql string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := s.ExecuteWithOptimization(ctx, ds, sql)
			if err != nil {
				result = &ExecutionResult{Error: err}
			}
			results[index] = result
		}(i, query)
	}

	wg.Wait()

	return results, nil
}

// GetOptimizationStats returns optimization statistics
func (s *OptimizedSQLExecutionService) GetOptimizationStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	optStats := s.optimizer.GetOptimizationStats()

	return map[string]interface{}{
		"execution_stats":    s.stats,
		"optimization_stats": optStats,
		"cache_stats":        s.cacheService.GetStats(),
	}
}

// GetSlowQueries returns slow queries for analysis
func (s *OptimizedSQLExecutionService) GetSlowQueries(threshold time.Duration) []*database.QueryPlan {
	return s.optimizer.GetSlowQueries(threshold)
}

// SuggestIndexes returns index suggestions based on query patterns
func (s *OptimizedSQLExecutionService) SuggestIndexes() []string {
	return s.optimizer.SuggestIndexes()
}

// AnalyzeQuery analyzes a query without executing it
func (s *OptimizedSQLExecutionService) AnalyzeQuery(sql string, ds models.DataSource) (*database.QueryPlan, error) {
	return s.optimizer.AnalyzeQuery(sql, ds)
}

// ExecuteSQL implements SQLExecutionService interface
func (s *OptimizedSQLExecutionService) ExecuteSQL(ds models.DataSource, sql string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	result, err := s.ExecuteWithOptimization(ctx, ds, sql)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ExecuteSQLWithTimeout implements SQLExecutionService interface
func (s *OptimizedSQLExecutionService) ExecuteSQLWithTimeout(ds models.DataSource, sql string, timeout time.Duration) ([]map[string]interface{}, error) {
	ctx := context.Background()
	result, err := s.ExecuteWithTimeout(ctx, ds, sql, timeout)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ExecuteSQLWithLimit implements SQLExecutionService interface
func (s *OptimizedSQLExecutionService) ExecuteSQLWithLimit(ds models.DataSource, sql string, limit int) ([]map[string]interface{}, error) {
	ctx := context.Background()
	result, err := s.ExecuteWithLimit(ctx, ds, sql, limit)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ResetStats resets all statistics
func (s *OptimizedSQLExecutionService) ResetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats = &ExecutionStats{
		LastReset: time.Now(),
	}
	s.optimizer.ResetStats()
}

// updateStats updates execution statistics
func (s *OptimizedSQLExecutionService) updateStats(cacheHit bool, executionTime time.Duration, failed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalExecutions++
	if cacheHit {
		s.stats.CacheHits++
	} else {
		s.stats.CacheMisses++
	}

	if failed {
		s.stats.FailedQueries++
	}

	if executionTime > 1*time.Second {
		s.stats.SlowQueries++
	}

	// Update average time
	if s.stats.TotalExecutions > 0 {
		totalTime := s.stats.AverageTime * time.Duration(s.stats.TotalExecutions-1)
		s.stats.AverageTime = (totalTime + executionTime) / time.Duration(s.stats.TotalExecutions)
	}
}

// cacheResults caches query results with intelligent TTL
func (s *OptimizedSQLExecutionService) cacheResults(cacheKey string, results []map[string]interface{}, sql string) {
	// Calculate TTL based on query complexity
	ttl := s.calculateTTL(sql)
	s.cacheService.Set(cacheKey, results, ttl)
}

// calculateTTL calculates cache TTL based on query characteristics
func (s *OptimizedSQLExecutionService) calculateTTL(sql string) time.Duration {
	// Base TTL
	baseTTL := 5 * time.Minute

	// Adjust based on query type
	if s.isAggregationQuery(sql) {
		baseTTL = 10 * time.Minute // Aggregation queries can be cached longer
	}

	if s.isSimpleQuery(sql) {
		baseTTL = 3 * time.Minute // Simple queries shorter TTL
	}

	if s.isComplexQuery(sql) {
		baseTTL = 15 * time.Minute // Complex queries longer TTL
	}

	// Time-based adjustments
	hour := time.Now().Hour()
	if hour >= 9 && hour <= 17 {
		baseTTL = baseTTL * 8 / 10 // Business hours - shorter TTL
	} else {
		baseTTL = baseTTL * 12 / 10 // Off-hours - longer TTL
	}

	return baseTTL
}

// isAggregationQuery checks if query is an aggregation query
func (s *OptimizedSQLExecutionService) isAggregationQuery(sql string) bool {
	upperSQL := utils.SanitizeSQL(sql)
	return utils.ContainsAny(upperSQL, []string{"COUNT", "SUM", "AVG", "MIN", "MAX", "GROUP BY"})
}

// isSimpleQuery checks if query is a simple query
func (s *OptimizedSQLExecutionService) isSimpleQuery(sql string) bool {
	upperSQL := utils.SanitizeSQL(sql)
	return !utils.ContainsAny(upperSQL, []string{"JOIN", "GROUP BY", "ORDER BY", "UNION", "("})
}

// isComplexQuery checks if query is a complex query
func (s *OptimizedSQLExecutionService) isComplexQuery(sql string) bool {
	upperSQL := utils.SanitizeSQL(sql)
	complexityScore := 0

	complexityFactors := []string{"JOIN", "UNION", "GROUP BY", "HAVING", "SUBQUERY", "WITH", "WINDOW"}
	for _, factor := range complexityFactors {
		if utils.Contains(upperSQL, factor) {
			complexityScore++
		}
	}

	return complexityScore >= 3
}

// generateCacheKey generates a unique cache key for a query
func (s *OptimizedSQLExecutionService) generateCacheKey(datasourceID uint, sql string) string {
	return utils.GenerateCacheKey(datasourceID, sql)
}

// containsLimit checks if SQL already contains a LIMIT clause
func (s *OptimizedSQLExecutionService) containsLimit(sql string) bool {
	return utils.Contains(strings.ToUpper(sql), "LIMIT")
}

// addLimitClause adds a LIMIT clause to SQL
func (s *OptimizedSQLExecutionService) addLimitClause(sql string, limit int) string {
	return fmt.Sprintf("%s LIMIT %d", sql, limit)
}
