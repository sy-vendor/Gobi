package database

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"gobi/internal/models"
	"gobi/pkg/errors"
)

// IndexInfo represents database index information
type IndexInfo struct {
	TableName   string    `json:"table_name"`
	IndexName   string    `json:"index_name"`
	ColumnName  string    `json:"column_name"`
	IndexType   string    `json:"index_type"`
	IsUnique    bool      `json:"is_unique"`
	Cardinality int64     `json:"cardinality"`
	Size        int64     `json:"size"`
	LastUsed    time.Time `json:"last_used"`
}

// IndexSuggestion represents an index suggestion
type IndexSuggestion struct {
	TableName        string   `json:"table_name"`
	Columns          []string `json:"columns"`
	IndexType        string   `json:"index_type"`
	Priority         int      `json:"priority"`
	Reason           string   `json:"reason"`
	SQL              string   `json:"sql"`
	EstimatedBenefit string   `json:"estimated_benefit"`
}

// IndexManager manages database indexes for optimization
type IndexManager struct {
	mu          sync.RWMutex
	indexCache  map[string]*IndexInfo
	suggestions []*IndexSuggestion
	stats       *IndexStats
}

// IndexStats tracks index management statistics
type IndexStats struct {
	TotalIndexes    int64     `json:"total_indexes"`
	CreatedIndexes  int64     `json:"created_indexes"`
	DroppedIndexes  int64     `json:"dropped_indexes"`
	AnalyzedQueries int64     `json:"analyzed_queries"`
	LastAnalysis    time.Time `json:"last_analysis"`
	LastReset       time.Time `json:"last_reset"`
}

// NewIndexManager creates a new index manager instance
func NewIndexManager() *IndexManager {
	return &IndexManager{
		indexCache:  make(map[string]*IndexInfo),
		suggestions: []*IndexSuggestion{},
		stats: &IndexStats{
			LastAnalysis: time.Now(),
			LastReset:    time.Now(),
		},
	}
}

// AnalyzeIndexes analyzes existing indexes in the database
func (im *IndexManager) AnalyzeIndexes(ds models.DataSource) ([]*IndexInfo, error) {
	db, err := GetConnection(&ds)
	if err != nil {
		return nil, errors.WrapError(err, "could not get database connection")
	}

	var indexes []*IndexInfo

	switch ds.Type {
	case "mysql":
		indexes, err = im.analyzeMySQLIndexes(db)
	case "postgres":
		indexes, err = im.analyzePostgresIndexes(db)
	case "sqlite":
		indexes, err = im.analyzeSQLiteIndexes(db)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", ds.Type)
	}

	if err != nil {
		return nil, err
	}

	// Update cache
	im.mu.Lock()
	for _, index := range indexes {
		cacheKey := fmt.Sprintf("%s_%s_%s", ds.Type, index.TableName, index.IndexName)
		im.indexCache[cacheKey] = index
	}
	im.stats.TotalIndexes = int64(len(indexes))
	im.stats.LastAnalysis = time.Now()
	im.mu.Unlock()

	return indexes, nil
}

// SuggestIndexes suggests indexes based on query patterns
func (im *IndexManager) SuggestIndexes(queryPatterns []string, ds models.DataSource) ([]*IndexSuggestion, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	suggestions := []*IndexSuggestion{}
	tableUsage := make(map[string]map[string]int)

	// Analyze query patterns
	for _, query := range queryPatterns {
		tables, columns := im.extractTableColumnUsage(query)
		for _, table := range tables {
			if tableUsage[table] == nil {
				tableUsage[table] = make(map[string]int)
			}
			for _, column := range columns {
				tableUsage[table][column]++
			}
		}
	}

	// Generate suggestions
	for table, columns := range tableUsage {
		for column, usage := range columns {
			if usage >= 3 { // Suggest index if column is used in 3+ queries
				suggestion := &IndexSuggestion{
					TableName:        table,
					Columns:          []string{column},
					IndexType:        "BTREE",
					Priority:         usage,
					Reason:           fmt.Sprintf("Column used in %d queries", usage),
					EstimatedBenefit: "High",
				}

				// Generate SQL based on database type
				switch ds.Type {
				case "mysql":
					suggestion.SQL = fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);",
						table, column, table, column)
				case "postgres":
					suggestion.SQL = fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);",
						table, column, table, column)
				case "sqlite":
					suggestion.SQL = fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);",
						table, column, table, column)
				}

				suggestions = append(suggestions, suggestion)
			}
		}
	}

	im.suggestions = suggestions
	return suggestions, nil
}

// CreateIndex creates an index based on suggestion
func (im *IndexManager) CreateIndex(suggestion *IndexSuggestion, ds models.DataSource) error {
	db, err := GetConnection(&ds)
	if err != nil {
		return errors.WrapError(err, "could not get database connection")
	}

	// Execute the index creation SQL
	_, err = db.Exec(suggestion.SQL)
	if err != nil {
		return errors.WrapError(err, "failed to create index")
	}

	// Update statistics
	im.mu.Lock()
	im.stats.CreatedIndexes++
	im.mu.Unlock()

	return nil
}

// DropIndex drops an index
func (im *IndexManager) DropIndex(tableName, indexName string, ds models.DataSource) error {
	db, err := GetConnection(&ds)
	if err != nil {
		return errors.WrapError(err, "could not get database connection")
	}

	dropSQL := fmt.Sprintf("DROP INDEX %s ON %s", indexName, tableName)
	if ds.Type == "postgres" {
		dropSQL = fmt.Sprintf("DROP INDEX %s", indexName)
	}

	_, err = db.Exec(dropSQL)
	if err != nil {
		return errors.WrapError(err, "failed to drop index")
	}

	// Update statistics
	im.mu.Lock()
	im.stats.DroppedIndexes++
	im.mu.Unlock()

	return nil
}

// GetIndexStats returns index management statistics
func (im *IndexManager) GetIndexStats() *IndexStats {
	im.mu.RLock()
	defer im.mu.RUnlock()

	stats := *im.stats
	return &stats
}

// GetUnusedIndexes returns indexes that haven't been used recently
func (im *IndexManager) GetUnusedIndexes(threshold time.Duration) []*IndexInfo {
	im.mu.RLock()
	defer im.mu.RUnlock()

	var unusedIndexes []*IndexInfo
	for _, index := range im.indexCache {
		if time.Since(index.LastUsed) > threshold {
			unusedIndexes = append(unusedIndexes, index)
		}
	}
	return unusedIndexes
}

// analyzeMySQLIndexes analyzes MySQL indexes
func (im *IndexManager) analyzeMySQLIndexes(db *sql.DB) ([]*IndexInfo, error) {
	query := `
		SELECT 
			TABLE_NAME,
			INDEX_NAME,
			COLUMN_NAME,
			INDEX_TYPE,
			NON_UNIQUE = 0 as IS_UNIQUE,
			CARDINALITY,
			INDEX_LENGTH as SIZE
		FROM INFORMATION_SCHEMA.STATISTICS 
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []*IndexInfo
	for rows.Next() {
		var index IndexInfo
		var isUnique int
		err := rows.Scan(
			&index.TableName,
			&index.IndexName,
			&index.ColumnName,
			&index.IndexType,
			&isUnique,
			&index.Cardinality,
			&index.Size,
		)
		if err != nil {
			return nil, err
		}
		index.IsUnique = isUnique == 1
		indexes = append(indexes, &index)
	}

	return indexes, nil
}

// analyzePostgresIndexes analyzes PostgreSQL indexes
func (im *IndexManager) analyzePostgresIndexes(db *sql.DB) ([]*IndexInfo, error) {
	query := `
		SELECT 
			t.schemaname,
			t.tablename as table_name,
			i.indexname as index_name,
			a.attname as column_name,
			am.amname as index_type,
			i.indisunique as is_unique,
			pg_stat_get_live_tuples(c.oid) as cardinality,
			pg_relation_size(i.indexrelid) as size
		FROM pg_index i
		JOIN pg_class c ON i.indrelid = c.oid
		JOIN pg_class ic ON i.indexrelid = ic.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		JOIN pg_tables t ON t.tablename = c.relname AND t.schemaname = n.nspname
		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = ANY(i.indkey)
		JOIN pg_am am ON ic.relam = am.oid
		WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY t.tablename, i.indexname
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []*IndexInfo
	for rows.Next() {
		var index IndexInfo
		var schemaName string
		var isUnique bool
		err := rows.Scan(
			&schemaName,
			&index.TableName,
			&index.IndexName,
			&index.ColumnName,
			&index.IndexType,
			&isUnique,
			&index.Cardinality,
			&index.Size,
		)
		if err != nil {
			return nil, err
		}
		index.IsUnique = isUnique
		indexes = append(indexes, &index)
	}

	return indexes, nil
}

// analyzeSQLiteIndexes analyzes SQLite indexes
func (im *IndexManager) analyzeSQLiteIndexes(db *sql.DB) ([]*IndexInfo, error) {
	query := `
		SELECT 
			tbl_name as table_name,
			name as index_name,
			sql as index_sql
		FROM sqlite_master 
		WHERE type = 'index' AND tbl_name NOT LIKE 'sqlite_%'
		ORDER BY tbl_name, name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []*IndexInfo
	for rows.Next() {
		var index IndexInfo
		var indexSQL string
		err := rows.Scan(
			&index.TableName,
			&index.IndexName,
			&indexSQL,
		)
		if err != nil {
			return nil, err
		}

		// Parse index SQL to extract column and type information
		index.ColumnName = im.extractColumnFromSQLiteIndex(indexSQL)
		index.IndexType = "BTREE" // SQLite default
		index.IsUnique = strings.Contains(strings.ToUpper(indexSQL), "UNIQUE")

		indexes = append(indexes, &index)
	}

	return indexes, nil
}

// extractColumnFromSQLiteIndex extracts column name from SQLite index SQL
func (im *IndexManager) extractColumnFromSQLiteIndex(indexSQL string) string {
	// Simple extraction - look for column name in parentheses
	start := strings.Index(indexSQL, "(")
	end := strings.Index(indexSQL, ")")
	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(indexSQL[start+1 : end])
	}
	return ""
}

// extractTableColumnUsage extracts table and column usage from SQL
func (im *IndexManager) extractTableColumnUsage(sql string) ([]string, []string) {
	upperSQL := strings.ToUpper(sql)
	tables := []string{}
	columns := []string{}

	// Extract tables
	parts := strings.Split(upperSQL, "FROM")
	if len(parts) > 1 {
		tablePart := strings.Fields(parts[1])[0]
		tables = append(tables, strings.ToLower(tablePart))
	}

	// Extract columns from WHERE clause
	if strings.Contains(upperSQL, "WHERE") {
		whereParts := strings.Split(upperSQL, "WHERE")
		if len(whereParts) > 1 {
			whereClause := whereParts[1]
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

	return tables, columns
}

// ResetStats resets index management statistics
func (im *IndexManager) ResetStats() {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.stats = &IndexStats{
		LastAnalysis: time.Now(),
		LastReset:    time.Now(),
	}
	im.indexCache = make(map[string]*IndexInfo)
	im.suggestions = []*IndexSuggestion{}
}
