package config

import (
	"fmt"
	"log"
	"os"

	"github.com/heronhoga/bars-be/models/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	log.Println("Connected to the database successfully")

	err = db.AutoMigrate(entities.UserEntity{}, entities.BeatEntity{}, entities.LikedBeatEntity{})
	if err != nil {
		log.Println("Database entities failed to migrate")
	}
	DB = db
}
