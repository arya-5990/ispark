package config

import (
	"fmt"
	"log"
	"os"

	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	log.Printf("Connecting to database at %s:%s...", dbHost, dbPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto Migration
	log.Println("Running AutoMigration...")
	err = DB.AutoMigrate(
		&models.Student{},
		&models.OTP{},
		&models.Activity{},
		&models.ActivitySubmission{},
		&models.Track{},
		&models.Announcement{},
		&models.SystemSetting{},
	)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed.")

	// Run Seeding
	log.Println("Seeding database...")
	seedDatabase(DB)
}

func seedDatabase(db *gorm.DB) {
	// 1. Seed Superadmin
	var adminCount int64
	db.Model(&models.Student{}).Where("roll_no = ? OR email_id = ?", "admin", "admin@admin.com").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := utils.HashPassword("admin123")
		if err == nil {
			adminUser := models.Student{
				RollNo:       "admin",
				Name:         "System Super Admin",
				CourseName:   "Administration",
				Semester:     0,
				ContactNo:    "0000000000",
				EmailID:      "admin@admin.com",
				EnrollmentNo: "ADMIN001",
				Password:     hashedPassword,
				IsVerified:   true,
				Role:         "superadmin",
			}
			if err := db.Create(&adminUser).Error; err != nil {
				log.Printf("Failed to seed superadmin user: %v", err)
			} else {
				log.Println("Seeded default superadmin: admin@admin.com / admin123")
			}
		}
	}

	// 2. Seed System Settings
	var settingsCount int64
	db.Model(&models.SystemSetting{}).Count(&settingsCount)
	if settingsCount == 0 {
		settings := []models.SystemSetting{
			{Key: "REGISTRATION_OPEN", Value: "true", Description: "Controls if new student registrations are open (true/false)"},
			{Key: "MIN_CREDITS_GRADUATION", Value: "60", Description: "Minimum extracurricular credits required for graduation"},
			{Key: "SEMESTER_DURATION_MONTHS", Value: "6", Description: "Standard academic semester duration in months"},
		}
		if err := db.Create(&settings).Error; err != nil {
			log.Printf("Failed to seed system settings: %v", err)
		} else {
			log.Println("Seeded system settings")
		}
	}

	// 3. Seed Tracks
	var tracksCount int64
	db.Model(&models.Track{}).Count(&tracksCount)
	if tracksCount == 0 {
		tracks := []models.Track{
			{Name: "Literary", MinCredits: 10, Description: "Debates, poetry, essay writing, and public speaking events."},
			{Name: "Cultural", MinCredits: 10, Description: "Dance, drama, music, art, and festival participations."},
			{Name: "Technical", MinCredits: 15, Description: "Coding contests, hackathons, robotics, and tech presentations."},
			{Name: "Academic", MinCredits: 15, Description: "Research papers, science olympiads, and subject-specific quizzes."},
			{Name: "Social Service", MinCredits: 10, Description: "Blood donation, NGO support, cleanliness drives, and social awareness."},
		}
		if err := db.Create(&tracks).Error; err != nil {
			log.Printf("Failed to seed curriculum tracks: %v", err)
		} else {
			log.Println("Seeded curriculum tracks")
		}
	}

	// 4. Seed Activities
	var activitiesCount int64
	db.Model(&models.Activity{}).Count(&activitiesCount)
	if activitiesCount == 0 {
		activities := []models.Activity{
			{Name: "Inter-College Debate Championship", Category: "Literary", Credits: 15, Description: "Annual inter-college debate tournament.", Status: "active"},
			{Name: "National Science Olympiad", Category: "Academic", Credits: 20, Description: "National level science exam and quiz.", Status: "active"},
			{Name: "Annual Cultural Fest - Dance", Category: "Cultural", Credits: 10, Description: "Solo or group dance performance at cultural fest.", Status: "active"},
			{Name: "Blood Donation Camp", Category: "Social Service", Credits: 8, Description: "Participating or volunteering in campus blood donation camp.", Status: "active"},
			{Name: "Robotics Workshop", Category: "Technical", Credits: 12, Description: "Hands-on robotics building workshop.", Status: "active"},
		}
		if err := db.Create(&activities).Error; err != nil {
			log.Printf("Failed to seed master activities: %v", err)
		} else {
			log.Println("Seeded master activities")
		}
	}
}
