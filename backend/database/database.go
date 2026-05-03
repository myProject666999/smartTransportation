package database

import (
	"fmt"
	"log"
	"smart_transportation/config"
	"smart_transportation/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbConfig := config.AppConfig.Database

	switch dbConfig.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(dbConfig.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		DB, err = gorm.Open(sqlite.Open("./transportation.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected successfully")

	AutoMigrate()
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Vehicle{},
		&models.VehicleMonitor{},
		&models.RoadCondition{},
		&models.Violation{},
		&models.Ticket{},
		&models.Booking{},
		&models.News{},
		&models.Announcement{},
		&models.Carousel{},
		&models.Feedback{},
		&models.TrafficFlow{},
		&models.SignalLight{},
		&models.Route{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully")
	SeedData()
}

func SeedData() {
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		admin := &models.Admin{
			Username: "admin",
			Password: "admin123",
			Role:     "super_admin",
			Status:   1,
		}
		DB.Create(admin)
		fmt.Println("Default admin created: username=admin, password=admin123")
	}
}
