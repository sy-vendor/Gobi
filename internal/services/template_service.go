package services

import (
	"fmt"
	"gobi/internal/models"
	"gobi/pkg/errors"
	"time"

	"gorm.io/gorm"
)

// TemplateService handles template-related business logic
type TemplateService struct {
	db *gorm.DB
}

// NewTemplateService creates a new TemplateService instance
func NewTemplateService(db *gorm.DB) *TemplateService {
	return &TemplateService{db: db}
}

// CreateTemplate creates a new template
func (s *TemplateService) CreateTemplate(template *models.ExcelTemplate, userID uint) error {
	template.UserID = userID

	if err := s.db.Create(template).Error; err != nil {
		return errors.WrapError(err, "Could not create template")
	}

	return nil
}

// ListTemplates retrieves templates based on user permissions
func (s *TemplateService) ListTemplates(userID uint, isAdmin bool) ([]models.ExcelTemplate, error) {
	var templates []models.ExcelTemplate

	query := s.db.Preload("User").Model(&models.ExcelTemplate{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&templates).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch templates")
	}

	// Do not return template content in list view
	for i := range templates {
		templates[i].Template = nil
	}

	return templates, nil
}

// GetTemplate retrieves a specific template
func (s *TemplateService) GetTemplate(templateID uint, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	var template models.ExcelTemplate
	if err := s.db.Preload("User").First(&template, templateID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && template.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Do not return template content in detail view
	template.Template = nil

	return &template, nil
}

// UpdateTemplate updates a template
func (s *TemplateService) UpdateTemplate(templateID uint, updates *models.ExcelTemplate, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	var template models.ExcelTemplate
	if err := s.db.First(&template, templateID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && template.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Update allowed fields
	if updates.Name != "" {
		template.Name = updates.Name
	}
	if updates.Description != "" {
		template.Description = updates.Description
	}

	if err := s.db.Save(&template).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update template")
	}

	// Do not return template content
	template.Template = nil

	return &template, nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(templateID uint, userID uint, isAdmin bool) error {
	var template models.ExcelTemplate
	if err := s.db.First(&template, templateID).Error; err != nil {
		return errors.ErrNotFound
	}

	if !isAdmin && template.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&template).Error; err != nil {
		return errors.WrapError(err, "Could not delete template")
	}

	return nil
}

// DownloadTemplate retrieves a template for download
func (s *TemplateService) DownloadTemplate(templateID uint, userID uint, isAdmin bool) (*models.ExcelTemplate, error) {
	var template models.ExcelTemplate
	if err := s.db.First(&template, templateID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && template.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return &template, nil
}

// GetDashboardStats retrieves dashboard statistics
func (s *TemplateService) GetDashboardStats() (map[string]interface{}, error) {
	var totalQueries int64
	var totalCharts int64
	var totalUsers int64
	var todayQueries int64

	today := time.Now().Format("2006-01-02")
	fmt.Println("[DEBUG] Go-side today:", today)

	var latestCreatedAt time.Time
	s.db.Model(&models.Query{}).Select("created_at").Order("created_at desc").Limit(1).Scan(&latestCreatedAt)
	fmt.Println("[DEBUG] Latest created_at in queries:", latestCreatedAt)

	s.db.Model(&models.Query{}).Count(&totalQueries)
	s.db.Model(&models.Chart{}).Count(&totalCharts)
	s.db.Model(&models.User{}).Count(&totalUsers)

	// Use database's current date for todayQueries
	s.db.Model(&models.Query{}).Where("DATE(created_at) = CURRENT_DATE").Count(&todayQueries)

	// 查询趋势（最近7天每天的查询数）
	queryTrends := []map[string]interface{}{}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		s.db.Model(&models.Query{}).Where("DATE(created_at) = ?", date).Count(&count)
		queryTrends = append(queryTrends, map[string]interface{}{"date": date, "count": count})
	}

	// 热门查询（执行次数最多的前5个查询）
	type HotQuery struct {
		Name  string
		Count int64
	}
	hotQueries := []HotQuery{}
	s.db.Table("queries").Select("name, exec_count as count").Order("exec_count desc").Limit(5).Scan(&hotQueries)

	return map[string]interface{}{
		"totalQueries": totalQueries,
		"totalCharts":  totalCharts,
		"totalUsers":   totalUsers,
		"todayQueries": todayQueries,
		"queryTrends":  queryTrends,
		"hotQueries":   hotQueries,
	}, nil
}
