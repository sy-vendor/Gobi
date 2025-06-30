package services

import (
	"fmt"
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/errors"
	"strconv"
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
		return err
	}

	// Flush cache
	s.cacheService.Flush()
	return nil
}

// ListQueries retrieves queries based on user permissions
func (s *QueryService) ListQueries(userID uint, isAdmin bool) ([]models.Query, error) {
	return s.queryRepo.FindByUser(userID, isAdmin)
}

// GetQuery retrieves a specific query
func (s *QueryService) GetQuery(queryID uint, userID uint, isAdmin bool) (*models.Query, error) {
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	// Flush cache
	s.cacheService.Flush()
	return query, nil
}

// DeleteQuery deletes a query
func (s *QueryService) DeleteQuery(queryID uint, userID uint, isAdmin bool) error {
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		return err
	}

	// Check permissions
	if !isAdmin && query.UserID != userID {
		return errors.ErrForbidden
	}

	// Delete query
	if err := s.queryRepo.Delete(queryID); err != nil {
		return err
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

// ExecuteQuery executes a query and returns the results
func (s *QueryService) ExecuteQuery(queryID uint, userID uint, isAdmin bool) (*ExecuteQueryResult, error) {
	// Check cache first
	cacheKey := "query_result_" + strconv.FormatUint(uint64(queryID), 10)
	if result, found := s.cacheService.Get(cacheKey); found {
		return &ExecuteQueryResult{
			Data:          result.([]map[string]interface{}),
			Source:        "cache",
			ExecutionTime: "0ms",
		}, nil
	}

	// Get query
	query, err := s.queryRepo.FindByID(queryID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !isAdmin && query.UserID != userID && !query.IsPublic {
		return nil, errors.ErrForbidden
	}

	// Validate SQL
	if err := s.validationService.ValidateSQL(query.SQL); err != nil {
		return nil, errors.WrapError(err, "Failed to execute query")
	}

	// Decrypt password if needed
	if query.DataSource.Password != "" {
		decryptedPassword, err := s.encryptionService.Decrypt(query.DataSource.Password)
		if err != nil {
			return nil, errors.WrapError(err, "Could not decrypt password")
		}
		query.DataSource.Password = decryptedPassword
	}

	// Execute query
	startTime := time.Now()
	results, err := s.sqlExecutionService.ExecuteSQL(query.DataSource, query.SQL)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to execute query")
	}
	executionTime := time.Since(startTime)

	// Update execution count
	s.queryRepo.IncrementExecCount(queryID)

	// Cache results
	s.cacheService.Set(cacheKey, results, 5*time.Minute)

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
