package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := GetDSN()

	config := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		),
	}

	//connting the newly created database during setup
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the database schema
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
