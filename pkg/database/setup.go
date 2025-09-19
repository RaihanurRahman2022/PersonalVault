package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() error {
	dsn := GetDSNWithoutDBName()

	// connecting default database of postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("faild to connect to postgres %v", err)
	}

	// create our own database
	dbName := os.Getenv("DB_NAME")
	result := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

	if result.Error != nil {
		if !isDatabaseExistsError(result.Error) {
			return fmt.Errorf("failed to create database: %v", result.Error)
		}
	}

	// getting the newly created database instance
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("faile to get database instance: %v", err)
	}

	sqlDB.Close()

	// connecting newly created database.
	dsn = GetDSN()
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to new database: %v", err)
	}

	// Enable uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to enable uuid-ossp extension: %v", err)
	}

	return nil
}

func isDatabaseExistsError(err error) bool {
	return err != nil && err.Error() == fmt.Sprintf("ERROR: database \"%s\" already exists (SQLSTATE 42P04)", os.Getenv("DB_NAME"))
}

func GetDSNWithoutDBName() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
}

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}
