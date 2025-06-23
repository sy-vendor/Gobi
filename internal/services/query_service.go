package services

import (
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// QueryService handles query-related business logic
type QueryService struct {
	db *gorm.DB
}

// NewQueryService creates a new QueryService instance
func NewQueryService(db *gorm.DB) *QueryService {
	return &QueryService{db: db}
}

// CreateQuery creates a new query
func (s *QueryService) CreateQuery(query *models.Query, userID uint) error {
	query.UserID = userID

	if err := s.db.Create(query).Error; err != nil {
		return errors.WrapError(err, "Could not create query")
	}

	utils.QueryCache.Flush()
	return nil
}

// ListQueries retrieves queries based on user permissions
func (s *QueryService) ListQueries(userID uint, isAdmin bool) ([]models.Query, error) {
	var queries []models.Query

	query := s.db.Preload("DataSource").Preload("User").Model(&models.Query{})
	if !isAdmin {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&queries).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch queries")
	}

	return queries, nil
}

// GetQuery retrieves a specific query
func (s *QueryService) GetQuery(queryID uint, userID uint, isAdmin bool) (*models.Query, error) {
	var query models.Query
	if err := s.db.Preload("DataSource").Preload("User").First(&query, queryID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && query.UserID != userID && !query.IsPublic {
		return nil, errors.ErrForbidden
	}

	return &query, nil
}

// UpdateQuery updates a query
func (s *QueryService) UpdateQuery(queryID uint, updates *models.Query, userID uint, isAdmin bool) (*models.Query, error) {
	var query models.Query
	if err := s.db.First(&query, queryID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && query.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update allowed fields
	if updates.Name != "" {
		query.Name = updates.Name
	}
	if updates.DataSourceID != 0 {
		query.DataSourceID = updates.DataSourceID
	}
	if updates.SQL != "" {
		query.SQL = updates.SQL
	}
	if updates.Description != "" {
		query.Description = updates.Description
	}
	// Update IsPublic if provided
	query.IsPublic = updates.IsPublic

	if err := s.db.Save(&query).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update query")
	}

	utils.QueryCache.Flush()
	return &query, nil
}

// DeleteQuery deletes a query
func (s *QueryService) DeleteQuery(queryID uint, userID uint, isAdmin bool) error {
	var query models.Query
	if err := s.db.First(&query, queryID).Error; err != nil {
		return errors.ErrNotFound
	}

	if !isAdmin && query.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&query).Error; err != nil {
		return errors.WrapError(err, "Could not delete query")
	}

	utils.QueryCache.Flush()
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
	if result, found := utils.QueryCache.Get(cacheKey); found {
		return &ExecuteQueryResult{
			Data:          result.([]map[string]interface{}),
			Source:        "cache",
			ExecutionTime: "0ms",
		}, nil
	}

	// Get query with data source
	var query models.Query
	if err := s.db.Preload("DataSource").First(&query, queryID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	// Permission check
	if !isAdmin && query.UserID != userID && !query.IsPublic {
		return nil, errors.ErrForbidden
	}

	// Decrypt password
	if query.DataSource.Password != "" {
		decryptedPassword, err := utils.DecryptAES(query.DataSource.Password)
		if err != nil {
			return nil, errors.WrapError(err, "Could not decrypt password")
		}
		query.DataSource.Password = decryptedPassword
	}

	// Execute query
	startTime := time.Now()
	results, err := utils.ExecuteSQL(query.DataSource, query.SQL)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to execute query")
	}
	executionTime := time.Since(startTime)

	// Update execution count
	query.ExecCount++
	s.db.Save(&query)

	// Set cache
	utils.QueryCache.Set(cacheKey, results, 5*time.Minute)

	// Prepare columns info
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
