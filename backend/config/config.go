package config

import (
	"campus-events/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("campus_events.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.EventTag{},
		&models.Registration{},
		&models.Comment{},
	)
	
	DB = db
}

const JWTSecret = "campus-events-secret-key-2024"
