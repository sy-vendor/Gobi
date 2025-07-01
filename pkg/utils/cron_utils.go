package utils

import (
	"time"

	"github.com/robfig/cron/v3"
)

// CronUtils provides common cron-related utilities
type CronUtils struct{}

// NewCronUtils creates a new CronUtils instance
func NewCronUtils() *CronUtils {
	return &CronUtils{}
}

// CalculateNextRunFromCron calculates the next run time based on cron pattern
// with customizable default time for error cases
func (c *CronUtils) CalculateNextRunFromCron(cronPattern string, defaultTime time.Time) time.Time {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronPattern)
	if err != nil {
		return defaultTime
	}
	return schedule.Next(time.Now())
}

// ValidateCronPattern validates if a cron pattern is valid
func (c *CronUtils) ValidateCronPattern(cronPattern string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronPattern)
	return err
}

// Global convenience functions
var cronUtils = NewCronUtils()

// CalculateNextRunFromCron is a convenience function using the global utils
func CalculateNextRunFromCron(cronPattern string, defaultTime time.Time) time.Time {
	return cronUtils.CalculateNextRunFromCron(cronPattern, defaultTime)
}

// ValidateCronPattern is a convenience function using the global utils
func ValidateCronPattern(cronPattern string) error {
	return cronUtils.ValidateCronPattern(cronPattern)
}
