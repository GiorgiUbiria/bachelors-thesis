package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Try to connect to the database
	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			// Wait for 2 seconds before retrying
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		log.Fatal("Failed to connect to database after multiple attempts:", err)
	}

	DB = db

	// Drop all tables if they exist (for development)
	if os.Getenv("ENV") == "development" {
		log.Println("Dropping all tables...")
		DB.Migrator().DropTable(
			&models.User{},
			&models.Product{},
			&models.Category{},
			&models.Order{},
			&models.OrderItem{},
			&models.Cart{},
			&models.CartItem{},
			&models.Favorite{},
			&models.UserActivity{},
			&models.RequestLog{},
			&models.BannedIP{},
			&models.Review{},
			&models.ReviewHelpful{},
			&models.ReviewReport{},
			&models.Payment{},
			&models.PaymentAttempt{},
		)
	}

	// Run migrations
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Favorite{},
		&models.UserActivity{},
		&models.RequestLog{},
		&models.BannedIP{},
		&models.Review{},
		&models.ReviewHelpful{},
		&models.ReviewReport{},
		&models.Payment{},
		&models.PaymentAttempt{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed successfully")

	// Seed the database if it's empty
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		log.Println("Seeding database...")
		SeedDatabase()
		log.Println("Database seeding completed")
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
