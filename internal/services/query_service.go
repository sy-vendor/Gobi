package services

import (
	errs "errors"
	"fmt"
	"gobi/config"
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// QueryService handles query-related business logic
type QueryService struct {
	queryRepo           repositories.QueryRepository
	cacheService        CacheService
	validationService   ValidationService
	sqlExecutionService SQLExecutionService
	encryptionService   EncryptionService
}

// NewQueryService creates a new QueryService instance
func NewQueryService(
	queryRepo repositories.QueryRepository,
	cacheService CacheService,
	validationService ValidationService,
	sqlExecutionService SQLExecutionService,
	encryptionService EncryptionService,
) *QueryService {
	return &QueryService{
		queryRepo:           queryRepo,
		cacheService:        cacheService,
		validationService:   validationService,
		sqlExecutionService: sqlExecutionService,
		encryptionService:   encryptionService,
	}
}

// CreateQuery creates a new query
func (s *QueryService) CreateQuery(query *models.Query, userID uint) error {
	query.UserID = userID

	// Validate SQL
	if err := s.validationService.ValidateSQL(query.SQL); err != nil {
		return errors.WrapError(err, "Invalid SQL query")
	}

	// Create query
	if err := s.queryRepo.Create(query); err != nil {
		return errors.WrapError(err, "Could not create query")
	}

	// Flush cache
	s.cacheService.Flush()
	return nil
}

// ListQueries retrieves queries based on user permissions
func (s *QueryService) ListQueries(userID uint, isAdmin bool) ([]models.Query, error) {
	queries, err := s.queryRepo.FindByUser(userID, isAdmin)
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch queries")
	}
	return queries, nil
}

// GetQuery retrieves a specific query
func (s *QueryService) GetQuery(queryID uint, userID uint, isAdmin bool) (*models.Query, error) {
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch query")
	}

	// Check permissions
	if !isAdmin && query.UserID != userID && !query.IsPublic {
		return nil, errors.ErrForbidden
	}

	return query, nil
}

// UpdateQuery updates a query
func (s *QueryService) UpdateQuery(queryID uint, updates *models.Query, userID uint, isAdmin bool) (*models.Query, error) {
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not fetch query")
	}

	// Check permissions
	if !isAdmin && query.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update fields
	if updates.Name != "" {
		query.Name = updates.Name
	}
	if updates.DataSourceID != 0 {
		query.DataSourceID = updates.DataSourceID
	}
	if updates.SQL != "" {
		if err := s.validationService.ValidateSQL(updates.SQL); err != nil {
			return nil, errors.WrapError(err, "Invalid SQL query")
		}
		query.SQL = updates.SQL
	}
	if updates.Description != "" {
		query.Description = updates.Description
	}
	query.IsPublic = updates.IsPublic

	// Save changes
	if err := s.queryRepo.Update(query); err != nil {
		return nil, errors.WrapError(err, "Could not update query")
	}

	// Flush cache
	s.cacheService.Flush()
	return query, nil
}

// DeleteQuery deletes a query
func (s *QueryService) DeleteQuery(queryID uint, userID uint, isAdmin bool) error {
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, "Could not fetch query")
	}

	// Check permissions
	if !isAdmin && query.UserID != userID {
		return errors.ErrForbidden
	}

	// Delete query
	if err := s.queryRepo.Delete(queryID); err != nil {
		return errors.WrapError(err, "Could not delete query")
	}

	// Flush cache
	s.cacheService.Flush()
	return nil
}

// ExecuteQueryResult represents the result of query execution
type ExecuteQueryResult struct {
	Data          []map[string]interface{} `json:"data"`
	Columns       []map[string]string      `json:"columns"`
	RowCount      int                      `json:"rowCount"`
	ExecutionTime string                   `json:"executionTime"`
	Source        string                   `json:"source"`
}

// 判断是否为简单查询
func isSimpleQuery(sql string) bool {
	normalized := strings.ToUpper(sql)
	if strings.Contains(normalized, "JOIN") ||
		strings.Contains(normalized, "GROUP BY") ||
		strings.Contains(normalized, "ORDER BY") ||
		strings.Contains(normalized, "UNION") ||
		strings.Contains(normalized, "(") { // 粗略判定子查询
		return false
	}
	return true
}

// ExecuteQuery executes a query and returns the results
func (s *QueryService) ExecuteQuery(queryID uint, userID uint, isAdmin bool) (*ExecuteQueryResult, error) {
	// Check cache first
	cacheKey := "query_result_" + strconv.FormatUint(uint64(queryID), 10)
	if result, found := s.cacheService.Get(cacheKey); found {
		// Handle different cache result types
		var data []map[string]interface{}

		// Check if it's a CacheEntry
		if cacheEntry, ok := result.(*utils.CacheEntry); ok {
			// Extract data from CacheEntry
			if entryData, ok := cacheEntry.Data.([]map[string]interface{}); ok {
				data = entryData
			} else {
				// If data is not the expected type, treat as cache miss
				found = false
			}
		} else if directData, ok := result.([]map[string]interface{}); ok {
			// Direct data (backward compatibility)
			data = directData
		} else {
			// Unknown type, treat as cache miss
			found = false
		}

		if found && data != nil {
			return &ExecuteQueryResult{
				Data:          data,
				Source:        "cache",
				ExecutionTime: "0ms",
			}, nil
		}
	}

	// Get query
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		if errs.Is(err, errors.ErrNotFound) {
			return nil, errors.NewErrorWithSeverity(
				errors.ErrCodeNotFound,
				"Query not found",
				err,
				errors.SeverityLow,
				errors.CategoryBusiness,
			)
		}
		return nil, errors.NewDatabaseError("Could not fetch query", err)
	}

	// Check permissions
	if !isAdmin && query.UserID != userID && !query.IsPublic {
		return nil, errors.NewErrorWithSeverity(
			errors.ErrCodeForbidden,
			"Access denied to this query",
			nil,
			errors.SeverityMedium,
			errors.CategoryAuthz,
		)
	}

	// Validate SQL
	if err := s.validationService.ValidateSQL(query.SQL); err != nil {
		return nil, errors.NewErrorWithSeverity(
			errors.ErrCodeInvalidSQL,
			"Invalid SQL query",
			err,
			errors.SeverityMedium,
			errors.CategoryValidation,
		)
	}

	// Decrypt password if needed
	if query.DataSource.Password != "" {
		decryptedPassword, err := s.encryptionService.Decrypt(query.DataSource.Password)
		if err != nil {
			return nil, errors.NewErrorWithSeverity(
				errors.ErrCodeInternalServer,
				"Could not decrypt password",
				err,
				errors.SeverityHigh,
				errors.CategorySecurity,
			)
		}
		query.DataSource.Password = decryptedPassword
	}

	// Execute query with retry mechanism
	startTime := time.Now()
	var results []map[string]interface{}

	// 使用重试机制执行SQL查询
	err = errors.Retry(func() error {
		var execErr error
		results, execErr = s.sqlExecutionService.ExecuteSQL(query.DataSource, query.SQL)
		if execErr != nil {
			// 记录重试
			errors.RecordRetry()
			return execErr
		}
		return nil
	}, errors.RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableErrors: []errors.ErrorCode{
			errors.ErrCodeDatabaseConnection,
			errors.ErrCodeDatabaseTimeout,
			errors.ErrCodeQueryTimeout,
		},
	})

	if err != nil {
		return nil, errors.NewErrorWithSeverity(
			errors.ErrCodeDatabaseQuery,
			"Failed to execute query",
			err,
			errors.SeverityHigh,
			errors.CategoryDatabase,
		)
	}

	executionTime := time.Since(startTime)

	// Update execution count
	if err := s.queryRepo.IncrementExecCount(queryID); err != nil {
		// 记录更新失败但不影响查询结果
		errors.RecordError(errors.NewDatabaseError("Failed to increment execution count", err))
	}

	// Cache results
	var ttl time.Duration
	if isSimpleQuery(query.SQL) {
		ttl = time.Duration(config.AppConfig.Cache.Strategy.SimpleQueryTTL) * time.Second
	} else {
		ttl = time.Duration(config.AppConfig.Cache.Strategy.ComplexQueryTTL) * time.Second
	}
	s.cacheService.Set(cacheKey, results, ttl)

	// Build columns info
	var columns []map[string]string
	if len(results) > 0 {
		for key := range results[0] {
			columns = append(columns, map[string]string{"name": key, "type": "unknown"})
		}
	}

	return &ExecuteQueryResult{
		Data:          results,
		Columns:       columns,
		RowCount:      len(results),
		ExecutionTime: fmt.Sprintf("%.2fms", float64(executionTime.Nanoseconds())/1e6),
		Source:        "database",
	}, nil
}
