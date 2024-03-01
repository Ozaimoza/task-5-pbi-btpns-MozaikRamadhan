package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB membuka koneksi ke database PostgreSQL
func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	fmt.Println("Connected successfully to the database")

	//add connection to global variabel
	DB = db

	//return
	return db, nil
}
