package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     string // admin or user
}

type Query struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	SQL         string
	Description string
	IsPublic    bool
}

type Chart struct {
	gorm.Model
	QueryID uint
	Query   Query
	UserID  uint
	User    User
	Name    string
	Type    string // bar, line, pie
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
