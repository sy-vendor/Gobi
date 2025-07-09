package database

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gobi/internal/models"
	"gobi/pkg/errors"
)

// QueryPlan represents a database query execution plan
type QueryPlan struct {
	QueryID       string                 `json:"query_id"`
	SQL           string                 `json:"sql"`
	ExecutionTime time.Duration          `json:"execution_time"`
	RowCount      int64                  `json:"row_count"`
	IndexUsed     []string               `json:"index_used"`
	TableScans    []string               `json:"table_scans"`
	Joins         []string               `json:"joins"`
	Complexity    string                 `json:"complexity"`
	Suggestions   []string               `json:"suggestions"`
	Metrics       map[string]interface{} `json:"metrics"`
}

// QueryOptimizer provides database query optimization features
type QueryOptimizer struct {
	mu           sync.RWMutex
	queryHistory map[string]*QueryPlan
	stats        *OptimizationStats
}

// OptimizationStats tracks optimization statistics
type OptimizationStats struct {
	TotalQueries     int64         `json:"total_queries"`
	OptimizedQueries int64         `json:"optimized_queries"`
	AverageTime      time.Duration `json:"average_time"`
	SlowQueries      int64         `json:"slow_queries"`
	LastReset        time.Time     `json:"last_reset"`
}

// NewQueryOptimizer creates a new query optimizer instance
func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		queryHistory: make(map[string]*QueryPlan),
		stats: &OptimizationStats{
			LastReset: time.Now(),
		},
	}
}

// AnalyzeQuery analyzes a SQL query and provides optimization suggestions
func (qo *QueryOptimizer) AnalyzeQuery(sql string, ds models.DataSource) (*QueryPlan, error) {
	plan := &QueryPlan{
		QueryID:     generateQueryID(sql),
		SQL:         sql,
		Suggestions: []string{},
		Metrics:     make(map[string]interface{}),
	}

	// Analyze query complexity
	plan.Complexity = qo.analyzeComplexity(sql)

	// Extract table and join information
	plan.Joins = qo.extractJoins(sql)
	plan.TableScans = qo.extractTables(sql)

	// Generate optimization suggestions
	plan.Suggestions = qo.generateSuggestions(sql, plan)

	// Store in history
	qo.mu.Lock()
	qo.queryHistory[plan.QueryID] = plan
	qo.stats.TotalQueries++
	qo.mu.Unlock()

	return plan, nil
}

// ExecuteWithOptimization executes a query with optimization analysis
func (qo *QueryOptimizer) ExecuteWithOptimization(ctx context.Context, ds models.DataSource, sql string) ([]map[string]interface{}, *QueryPlan, error) {
	// Analyze query first
	plan, err := qo.AnalyzeQuery(sql, ds)
	if err != nil {
		return nil, nil, err
	}

	// Execute query with timing
	startTime := time.Now()
	results, err := qo.executeQuery(ctx, ds, sql)
	plan.ExecutionTime = time.Since(startTime)
	plan.RowCount = int64(len(results))

	// Update statistics
	qo.mu.Lock()
	if plan.ExecutionTime > 1*time.Second {
		qo.stats.SlowQueries++
	}
	qo.stats.AverageTime = qo.calculateAverageTime()
	qo.mu.Unlock()

	// Add execution metrics
	plan.Metrics["memory_usage"] = qo.estimateMemoryUsage(results)
	plan.Metrics["network_transfer"] = qo.estimateNetworkTransfer(results)

	return results, plan, err
}

// GetOptimizationStats returns optimization statistics
func (qo *QueryOptimizer) GetOptimizationStats() *OptimizationStats {
	qo.mu.RLock()
	defer qo.mu.RUnlock()

	stats := *qo.stats
	return &stats
}

// GetSlowQueries returns queries that took longer than threshold
func (qo *QueryOptimizer) GetSlowQueries(threshold time.Duration) []*QueryPlan {
	qo.mu.RLock()
	defer qo.mu.RUnlock()

	var slowQueries []*QueryPlan
	for _, plan := range qo.queryHistory {
		if plan.ExecutionTime > threshold {
			slowQueries = append(slowQueries, plan)
		}
	}
	return slowQueries
}

// SuggestIndexes suggests indexes for frequently used queries
func (qo *QueryOptimizer) SuggestIndexes() []string {
	qo.mu.RLock()
	defer qo.mu.RUnlock()

	indexSuggestions := []string{}
	tableUsage := make(map[string]map[string]int)

	// Analyze table and column usage patterns
	for _, plan := range qo.queryHistory {
		for _, table := range plan.TableScans {
			if tableUsage[table] == nil {
				tableUsage[table] = make(map[string]int)
			}

			// Extract WHERE conditions for index suggestions
			columns := qo.extractWhereColumns(plan.SQL)
			for _, col := range columns {
				tableUsage[table][col]++
			}
		}
	}

	// Generate index suggestions
	for table, columns := range tableUsage {
		for column, usage := range columns {
			if usage >= 3 { // Suggest index if column is used in 3+ queries
				indexSuggestions = append(indexSuggestions,
					fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);",
						table, column, table, column))
			}
		}
	}

	return indexSuggestions
}

// analyzeComplexity analyzes query complexity
func (qo *QueryOptimizer) analyzeComplexity(sql string) string {
	upperSQL := strings.ToUpper(sql)
	complexityScore := 0

	complexityFactors := map[string]int{
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
		"DISTINCT": 1,
	}

	for factor, score := range complexityFactors {
		if strings.Contains(upperSQL, factor) {
			complexityScore += score
		}
	}

	// Count tables
	tableCount := strings.Count(upperSQL, "FROM") + strings.Count(upperSQL, "JOIN")
	if tableCount > 3 {
		complexityScore += tableCount - 3
	}

	if complexityScore >= 5 {
		return "high"
	} else if complexityScore >= 2 {
		return "medium"
	}
	return "low"
}

// extractJoins extracts JOIN information from SQL
func (qo *QueryOptimizer) extractJoins(sql string) []string {
	upperSQL := strings.ToUpper(sql)
	joins := []string{}

	joinTypes := []string{"JOIN", "LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "OUTER JOIN"}
	for _, joinType := range joinTypes {
		if strings.Contains(upperSQL, joinType) {
			joins = append(joins, joinType)
		}
	}

	return joins
}

// extractTables extracts table names from SQL
func (qo *QueryOptimizer) extractTables(sql string) []string {
	upperSQL := strings.ToUpper(sql)
	tables := []string{}

	// Simple table extraction (could be enhanced with proper SQL parsing)
	parts := strings.Split(upperSQL, "FROM")
	if len(parts) > 1 {
		tablePart := strings.Fields(parts[1])[0]
		tables = append(tables, strings.ToLower(tablePart))
	}

	// Extract tables from JOIN clauses
	joinParts := strings.Split(upperSQL, "JOIN")
	for i := 1; i < len(joinParts); i++ {
		fields := strings.Fields(joinParts[i])
		if len(fields) > 0 {
			tables = append(tables, strings.ToLower(fields[0]))
		}
	}

	return tables
}

// extractWhereColumns extracts columns used in WHERE clauses
func (qo *QueryOptimizer) extractWhereColumns(sql string) []string {
	upperSQL := strings.ToUpper(sql)
	columns := []string{}

	// Simple column extraction from WHERE clause
	if strings.Contains(upperSQL, "WHERE") {
		whereParts := strings.Split(upperSQL, "WHERE")
		if len(whereParts) > 1 {
			whereClause := whereParts[1]
			// Extract column names (simplified)
			words := strings.Fields(whereClause)
			for i, word := range words {
				if i > 0 && (words[i-1] == "AND" || words[i-1] == "OR" || i == 0) {
					if !strings.Contains(word, "=") && !strings.Contains(word, ">") &&
						!strings.Contains(word, "<") && !strings.Contains(word, "LIKE") {
						columns = append(columns, strings.ToLower(word))
					}
				}
			}
		}
	}

	return columns
}

// generateSuggestions generates optimization suggestions
func (qo *QueryOptimizer) generateSuggestions(sql string, plan *QueryPlan) []string {
	suggestions := []string{}
	upperSQL := strings.ToUpper(sql)

	// Check for SELECT *
	if strings.Contains(upperSQL, "SELECT *") {
		suggestions = append(suggestions, "Consider specifying columns instead of SELECT *")
	}

	// Check for missing LIMIT
	if !strings.Contains(upperSQL, "LIMIT") && plan.Complexity == "low" {
		suggestions = append(suggestions, "Consider adding LIMIT clause for large result sets")
	}

	// Check for missing WHERE clause
	if !strings.Contains(upperSQL, "WHERE") && len(plan.TableScans) > 0 {
		suggestions = append(suggestions, "Consider adding WHERE clause to filter results")
	}

	// Check for ORDER BY without LIMIT
	if strings.Contains(upperSQL, "ORDER BY") && !strings.Contains(upperSQL, "LIMIT") {
		suggestions = append(suggestions, "Consider adding LIMIT when using ORDER BY")
	}

	// Check for complex queries without proper indexing
	if plan.Complexity == "high" && len(plan.Joins) > 2 {
		suggestions = append(suggestions, "Consider adding indexes on join columns")
	}

	// Check for GROUP BY without proper columns
	if strings.Contains(upperSQL, "GROUP BY") {
		suggestions = append(suggestions, "Ensure GROUP BY columns are indexed")
	}

	return suggestions
}

// executeQuery executes the actual query
func (qo *QueryOptimizer) executeQuery(ctx context.Context, ds models.DataSource, sql string) ([]map[string]interface{}, error) {
	db, err := GetConnection(&ds)
	if err != nil {
		return nil, errors.WrapError(err, "could not get database connection")
	}

	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.WrapError(err, "query execution failed")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.WrapError(err, "failed to get column information")
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
			if b, ok := val.([]byte); ok {
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

// estimateMemoryUsage estimates memory usage of query results
func (qo *QueryOptimizer) estimateMemoryUsage(results []map[string]interface{}) int64 {
	var totalSize int64
	for _, row := range results {
		for _, value := range row {
			switch v := value.(type) {
			case string:
				totalSize += int64(len(v))
			case int64:
				totalSize += 8
			case float64:
				totalSize += 8
			default:
				totalSize += 16 // Default estimate
			}
		}
	}
	return totalSize
}

// estimateNetworkTransfer estimates network transfer size
func (qo *QueryOptimizer) estimateNetworkTransfer(results []map[string]interface{}) int64 {
	var totalSize int64
	for _, row := range results {
		for _, value := range row {
			switch v := value.(type) {
			case string:
				totalSize += int64(len(v))
			case int64:
				totalSize += 8
			case float64:
				totalSize += 8
			default:
				totalSize += 16
			}
		}
	}
	return totalSize
}

// calculateAverageTime calculates average query execution time
func (qo *QueryOptimizer) calculateAverageTime() time.Duration {
	var totalTime time.Duration
	var count int64

	for _, plan := range qo.queryHistory {
		totalTime += plan.ExecutionTime
		count++
	}

	if count == 0 {
		return 0
	}
	return totalTime / time.Duration(count)
}

// generateQueryID generates a unique ID for a query
func generateQueryID(sql string) string {
	// Simple hash-based ID generation
	hash := 0
	for _, char := range sql {
		hash = ((hash << 5) - hash) + int(char)
		hash = hash & hash // Convert to 32-bit integer
	}
	return fmt.Sprintf("q_%d", hash)
}

// ResetStats resets optimization statistics
func (qo *QueryOptimizer) ResetStats() {
	qo.mu.Lock()
	defer qo.mu.Unlock()

	qo.stats = &OptimizationStats{
		LastReset: time.Now(),
	}
	qo.queryHistory = make(map[string]*QueryPlan)
}
