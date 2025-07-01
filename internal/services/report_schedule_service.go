package services

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// ReportScheduleService handles report schedule-related business logic
type ReportScheduleService struct {
	db *gorm.DB
}

// NewReportScheduleService creates a new ReportScheduleService instance
func NewReportScheduleService(db *gorm.DB) *ReportScheduleService {
	return &ReportScheduleService{db: db}
}

// CreateReportSchedule creates a new report schedule
func (s *ReportScheduleService) CreateReportSchedule(schedule *models.ReportSchedule, userID uint) error {
	// Validate cron pattern
	if err := utils.ValidateCronPattern(schedule.CronPattern); err != nil {
		return errors.NewBadRequestError("Invalid cron pattern", err)
	}

	schedule.UserID = userID
	schedule.Active = true
	schedule.NextRun = s.calculateNextRunFromCron(schedule.CronPattern)

	if err := s.db.Create(schedule).Error; err != nil {
		return errors.WrapError(err, "Could not create report schedule")
	}

	return nil
}

// ListReportSchedules retrieves report schedules based on user permissions
func (s *ReportScheduleService) ListReportSchedules(userID uint, isAdmin bool) ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule

	query := s.db.Model(&models.ReportSchedule{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch report schedules")
	}

	return schedules, nil
}

// GetReportSchedule retrieves a specific report schedule
func (s *ReportScheduleService) GetReportSchedule(scheduleID uint, userID uint, isAdmin bool) (*models.ReportSchedule, error) {
	var schedule models.ReportSchedule
	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && schedule.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return &schedule, nil
}

// UpdateReportSchedule updates a report schedule
func (s *ReportScheduleService) UpdateReportSchedule(scheduleID uint, updates *models.ReportSchedule, userID uint, isAdmin bool) (*models.ReportSchedule, error) {
	var schedule models.ReportSchedule
	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if !isAdmin && schedule.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if updates.Active != schedule.Active {
		schedule.Active = updates.Active
	}

	if updates.Name != "" {
		schedule.Name = updates.Name
	}
	if updates.Type != "" {
		schedule.Type = updates.Type
	}
	if updates.Queries != "" {
		schedule.Queries = updates.Queries
	}
	if updates.Charts != "" {
		schedule.Charts = updates.Charts
	}
	if updates.Templates != "" {
		schedule.Templates = updates.Templates
	}
	if updates.CronPattern != "" {
		if err := utils.ValidateCronPattern(updates.CronPattern); err != nil {
			return nil, errors.NewBadRequestError("Invalid cron pattern", err)
		}
		schedule.CronPattern = updates.CronPattern
		schedule.NextRun = s.calculateNextRunFromCron(updates.CronPattern)
	}

	if err := s.db.Save(&schedule).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update report schedule")
	}

	return &schedule, nil
}

// DeleteReportSchedule deletes a report schedule
func (s *ReportScheduleService) DeleteReportSchedule(scheduleID uint, userID uint, isAdmin bool) error {
	var schedule models.ReportSchedule
	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		return errors.ErrNotFound
	}

	if !isAdmin && schedule.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.db.Delete(&schedule).Error; err != nil {
		return errors.WrapError(err, "Could not delete report schedule")
	}

	return nil
}

// calculateNextRunFromCron calculates the next run time based on cron pattern
func (s *ReportScheduleService) calculateNextRunFromCron(cronPattern string) time.Time {
	return utils.CalculateNextRunFromCron(cronPattern, time.Now().Add(24*time.Hour))
}
