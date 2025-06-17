package database

import (
	"gobi/config"
	"gobi/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.DataSource{},
		&models.Query{},
		&models.Chart{},
		&models.ExcelTemplate{},
	)
	if err != nil {
		return err
	}

	return nil
}
