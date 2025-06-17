package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex"`
	Email     string `gorm:"uniqueIndex"`
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
}

type Chart struct {
	gorm.Model
	QueryID uint
	Query   Query
	UserID  uint
	User    User
	Name    string
	Type    string // bar, line, pie, scatter, radar, heatmap, gauge, funnel
	Config  string // JSON configuration
	Data    string // JSON data
}

type ExcelTemplate struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Template    []byte
	Description string
}
