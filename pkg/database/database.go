package database

import (
	"gobi/config"
	"gobi/internal/models"
	"gobi/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error
	switch cfg.Database.Type {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	case "mysql":
		DB, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	case "postgres":
		DB, err = gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	default:
		return errors.NewError(errors.ErrCodeInvalidRequest, "unsupported database type: "+cfg.Database.Type, nil)
	}
	if err != nil {
		return errors.WrapError(err, "Failed to open database connection")
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.DataSource{},
		&models.Query{},
		&models.Chart{},
		&models.ExcelTemplate{},
		&models.Report{},
		&models.ReportSchedule{},
		&models.APIKey{},
		&models.Webhook{},
		&models.WebhookDelivery{},
	)
	if err != nil {
		return errors.WrapError(err, "Failed to auto-migrate database schema")
	}

	return nil
}

// GetDB returns the global database instance.
func GetDB() *gorm.DB {
	return DB
}
