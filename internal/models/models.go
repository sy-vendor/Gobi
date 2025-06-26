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
	DataSourceID uint `json:"data_source_id"`
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
	Type        string // bar, line, pie, scatter, radar, heatmap, gauge, funnel, area, 3d-bar, 3d-scatter, 3d-surface, 3d-bubble, treemap, sunburst, tree, boxplot, candlestick, wordcloud, graph, waterfall, polar, gantt, rose, geo, map, choropleth
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

// APIKey represents an API key for service-to-service authentication
// Each key is associated with a user (owner), has a name, key value, and status
// Key is stored as a hash for security
// ExpiresAt is optional (nil means never expires)
// Add this comment for migration reference
// NOTE: Remember to add &APIKey{} to your auto-migration in main.go or database.go
type APIKey struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	User      User       `gorm:"constraint:OnDelete:CASCADE"`
	Name      string     `gorm:"type:varchar(64)" json:"name"`
	KeyHash   string     `gorm:"type:varchar(128);uniqueIndex" json:"-"` // hashed key
	Prefix    string     `gorm:"type:varchar(16);index" json:"prefix"`   // for fast lookup
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	Revoked   bool       `gorm:"default:false" json:"revoked"`
}

// Webhook represents a webhook configuration for event notifications
// Supports multiple event types and custom headers
type Webhook struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	Name      string    `gorm:"type:varchar(64)" json:"name"`
	URL       string    `gorm:"type:varchar(512)" json:"url"`
	Events    string    `gorm:"type:text" json:"events"`    // JSON array of event types
	Headers   string    `gorm:"type:text" json:"headers"`   // JSON object of custom headers
	Secret    string    `gorm:"type:varchar(128)" json:"-"` // for signature verification
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	WebhookID uint       `gorm:"index" json:"webhook_id"`
	Webhook   Webhook    `gorm:"constraint:OnDelete:CASCADE"`
	Event     string     `gorm:"type:varchar(64)" json:"event"`
	Payload   string     `gorm:"type:text" json:"payload"`       // JSON payload sent
	Status    string     `gorm:"type:varchar(32)" json:"status"` // success, failed, pending
	Response  string     `gorm:"type:text" json:"response"`      // response from webhook URL
	Attempts  int        `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at"`
}
