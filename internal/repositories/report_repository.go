package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/gorm"
)

// ReportRepositoryImpl implements ReportRepository interface
type ReportRepositoryImpl struct {
	db *gorm.DB
}

// NewReportRepository creates a new ReportRepository instance
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{db: db}
}

// Create creates a new report
func (r *ReportRepositoryImpl) Create(report *models.Report) error {
	if err := r.db.Create(report).Error; err != nil {
		return errors.WrapError(err, "Could not create report")
	}
	return nil
}

// FindByID finds a report by ID
func (r *ReportRepositoryImpl) FindByID(id uint) (*models.Report, error) {
	var report models.Report
	if err := r.db.Preload("User").First(&report, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "Could not find report")
	}
	return &report, nil
}

// FindByUser finds reports by user ID
func (r *ReportRepositoryImpl) FindByUser(userID uint, isAdmin bool) ([]models.Report, error) {
	var reports []models.Report
	query := r.db.Preload("User").Model(&models.Report{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&reports).Error; err != nil {
		return nil, errors.WrapError(err, "Could not find reports")
	}
	return reports, nil
}

// Update updates a report
func (r *ReportRepositoryImpl) Update(report *models.Report) error {
	if err := r.db.Save(report).Error; err != nil {
		return errors.WrapError(err, "Could not update report")
	}
	return nil
}

// Delete deletes a report
func (r *ReportRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.Report{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete report")
	}
	return nil
}

// UpdateStatus updates report status and error message
func (r *ReportRepositoryImpl) UpdateStatus(reportID uint, status string, error string) error {
	if err := r.db.Model(&models.Report{}).Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"status": status,
			"error":  error,
		}).Error; err != nil {
		return errors.WrapError(err, "Could not update report status")
	}
	return nil
}
