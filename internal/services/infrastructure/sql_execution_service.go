package infrastructure

import (
	"gobi/internal/models"
	"gobi/pkg/utils"
	"time"
)

// SQLExecutionServiceImpl implements SQLExecutionService
type SQLExecutionServiceImpl struct{}

// NewSQLExecutionService creates a new SQLExecutionService instance
func NewSQLExecutionService() *SQLExecutionServiceImpl {
	return &SQLExecutionServiceImpl{}
}

// ExecuteSQL executes SQL query
func (s *SQLExecutionServiceImpl) ExecuteSQL(ds models.DataSource, sql string) ([]map[string]interface{}, error) {
	return utils.ExecuteSQL(ds, sql)
}

// ExecuteSQLWithTimeout executes SQL query with timeout
func (s *SQLExecutionServiceImpl) ExecuteSQLWithTimeout(ds models.DataSource, sql string, timeout time.Duration) ([]map[string]interface{}, error) {
	return utils.ExecuteSQLWithTimeout(ds, sql, timeout)
}

// ExecuteSQLWithLimit executes SQL query with row limit
func (s *SQLExecutionServiceImpl) ExecuteSQLWithLimit(ds models.DataSource, sql string, limit int) ([]map[string]interface{}, error) {
	return utils.ExecuteSQLWithLimit(ds, sql, limit)
}
