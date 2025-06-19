package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(64);uniqueIndex"`
	Email     string `gorm:"type:varchar(128);uniqueIndex"`
	Password  string
	Role      string    // admin or user
	LastLogin time.Time `json:"last_login"`
}

type DataSource struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Type        string // mysql, postgres, sqlite, etc.
	Host        string
	Port        int
	Database    string
	Username    string
	Password    string
	Description string
	IsPublic    bool
}

type Query struct {
	gorm.Model
	UserID       uint
	User         User
	DataSourceID uint
	DataSource   DataSource
	Name         string
	SQL          string
	Description  string
	IsPublic     bool
	ExecCount    int64 // 新增：执行次数
}

type Chart struct {
	gorm.Model
	QueryID     uint
	Query       Query
	UserID      uint
	User        User
	Name        string
	Type        string // bar, line, pie, scatter, radar, heatmap, gauge, funnel
	Config      string // JSON configuration
	Data        string // JSON data
	Description string `json:"description"`
}

type ExcelTemplate struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Template    []byte
	Description string `json:"description"`
}

type Report struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Type        string    // daily, weekly, monthly
	Content     []byte    // report content in PDF or Excel format
	GeneratedAt time.Time // when the report was generated
	Status      string    // pending, success, failed
	Error       string    // error message if generation failed
}

type ReportSchedule struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Type        string    // daily, weekly, monthly
	Queries     string    // JSON array of query IDs to include
	Charts      string    // JSON array of chart IDs to include
	Templates   string    // JSON array of template IDs to use
	LastRun     time.Time // last time the report was generated
	NextRun     time.Time // next scheduled run time
	Active      bool      // whether the schedule is active
	CronPattern string    // cron pattern for scheduling
}
