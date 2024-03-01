package database

import (
	"task-5-pbi-btpns/models"

	"gorm.io/gorm"
)

// AutoMigrate
func AutoMigrate(db *gorm.DB) error {
	// Migrate model PhotoModel
	if err := db.AutoMigrate(&models.PhotoModel{}); err != nil {
		return err
	}

	// Migrate model UserModel
	if err := db.AutoMigrate(&models.UserModel{}); err != nil {
		return err
	}

	return nil
}
