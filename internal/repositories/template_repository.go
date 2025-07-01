package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"time"

	"gorm.io/gorm"
)

// TemplateRepositoryImpl implements TemplateRepository interface
type TemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewTemplateRepository creates a new TemplateRepository instance
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &TemplateRepositoryImpl{db: db}
}

// Create creates a new template
func (r *TemplateRepositoryImpl) Create(template *models.ExcelTemplate) error {
	if err := r.db.Create(template).Error; err != nil {
		return errors.WrapError(err, "Could not create template")
	}
	return nil
}

// FindByID finds a template by ID
func (r *TemplateRepositoryImpl) FindByID(id uint) (*models.ExcelTemplate, error) {
	var template models.ExcelTemplate
	if err := r.db.Preload("User").First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not find template")
	}
	return &template, nil
}

// FindByUser finds templates by user ID
func (r *TemplateRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.ExcelTemplate, error) {
	var templates []models.ExcelTemplate
	query := r.db.Preload("User").Model(&models.ExcelTemplate{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&templates).Error; err != nil {
		return nil, errors.WrapError(err, "Could not find templates")
	}
	return templates, nil
}

// Update updates a template
func (r *TemplateRepositoryImpl) Update(template *models.ExcelTemplate) error {
	if err := r.db.Save(template).Error; err != nil {
		return errors.WrapError(err, "Could not update template")
	}
	return nil
}

// Delete deletes a template
func (r *TemplateRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.ExcelTemplate{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete template")
	}
	return nil
}

// GetStats retrieves dashboard statistics
func (r *TemplateRepositoryImpl) GetStats() (map[string]interface{}, error) {
	var totalQueries int64
	var totalCharts int64
	var totalUsers int64
	var todayQueries int64

	today := time.Now().Format("2006-01-02")
	r.db.Model(&models.Query{}).Count(&totalQueries)
	r.db.Model(&models.Chart{}).Count(&totalCharts)
	r.db.Model(&models.User{}).Count(&totalUsers)
	r.db.Model(&models.Query{}).Where("DATE(created_at) = ?", today).Count(&todayQueries)

	// 查询趋势（最近7天每天的查询数）
	queryTrends := []map[string]interface{}{}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		r.db.Model(&models.Query{}).Where("DATE(created_at) = ?", date).Count(&count)
		queryTrends = append(queryTrends, map[string]interface{}{"date": date, "count": count})
	}

	// 热门查询（执行次数最多的前5个查询）
	type HotQuery struct {
		Name  string
		Count int64
	}
	hotQueries := []HotQuery{}
	r.db.Table("queries").Select("name, exec_count as count").Order("exec_count desc").Limit(5).Scan(&hotQueries)

	return map[string]interface{}{
		"totalQueries": totalQueries,
		"totalCharts":  totalCharts,
		"totalUsers":   totalUsers,
		"todayQueries": todayQueries,
		"queryTrends":  queryTrends,
		"hotQueries":   hotQueries,
	}, nil
}
