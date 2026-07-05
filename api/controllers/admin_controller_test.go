package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/routes"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/gorm"
)

// SetupAdminTestDB overrides GORM instance to use memory sqlite with all migrated/seeded schemas
func SetupAdminTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory SQLite database: %v", err)
	}

	// Auto Migrate all models
	err = db.AutoMigrate(
		&models.Student{},
		&models.OTP{},
		&models.Activity{},
		&models.ActivitySubmission{},
		&models.Track{},
		&models.Announcement{},
		&models.SystemSetting{},
	)
	if err != nil {
		t.Fatalf("Failed to run database migrations: %v", err)
	}

	config.DB = db

	// Run standard seeding
	// 1. Seed Superadmin
	hashedPassword, _ := utils.HashPassword("admin123")
	db.Create(&models.Student{
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
	})

	// 2. Seed Tracks
	db.Create(&[]models.Track{
		{Name: "Literary", MinCredits: 10, Description: "Debates, poetry, essay writing, and public speaking events."},
		{Name: "Cultural", MinCredits: 10, Description: "Dance, drama, music, art, and festival participations."},
	})

	// 3. Seed Activities
	db.Create(&[]models.Activity{
		{Name: "Inter-College Debate Championship", Category: "Literary", Credits: 15, Description: "Annual inter-college debate tournament.", Status: "active"},
		{Name: "Robotics Workshop", Category: "Technical", Credits: 12, Description: "Hands-on robotics building workshop.", Status: "active"},
	})

	// 4. Seed Settings
	db.Create(&[]models.SystemSetting{
		{Key: "REGISTRATION_OPEN", Value: "true", Description: "Controls if new student registrations are open (true/false)"},
	})
}

func TestSuperAdminModule(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_jwt_secret_key_1234567890")
	SetupAdminTestDB(t)

	app := fiber.New()
	routes.SetupRoutes(app)

	// Create test tokens
	adminToken, err := utils.GenerateAccessToken("admin", "admin@admin.com", "superadmin")
	if err != nil {
		t.Fatalf("Failed to generate admin token: %v", err)
	}

	// Create a regular student in DB and generate token
	studentEmail := "student@student.com"
	studentRoll := "student123"
	config.DB.Create(&models.Student{
		RollNo:       studentRoll,
		Name:         "Student Tester",
		CourseName:   "MCA",
		Semester:     4,
		ContactNo:    "1234567890",
		EmailID:      studentEmail,
		EnrollmentNo: "STUDENT123",
		Password:     "password",
		IsVerified:   true,
		Role:         "student",
	})
	studentToken, err := utils.GenerateAccessToken(studentRoll, studentEmail, "student")
	if err != nil {
		t.Fatalf("Failed to generate student token: %v", err)
	}

	// 1. TEST AUTHORIZATION RULES
	t.Run("AccessDashboard_Admin_Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("AccessDashboard_Student_Forbidden", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden, got %d", resp.StatusCode)
		}
	})

	// 2. TEST SYSTEM SETTINGS API
	t.Run("GetSettings_Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/admin/settings", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var settings []models.SystemSetting
		json.NewDecoder(resp.Body).Decode(&settings)
		if len(settings) != 1 {
			t.Errorf("Expected 1 setting, got %d", len(settings))
		}
	})

	t.Run("UpdateSetting_Success", func(t *testing.T) {
		payload := map[string]string{"value": "false"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/admin/settings/REGISTRATION_OPEN", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var setting models.SystemSetting
		config.DB.Where("key = ?", "REGISTRATION_OPEN").First(&setting)
		if setting.Value != "false" {
			t.Errorf("Expected REGISTRATION_OPEN to be false, got %s", setting.Value)
		}
	})

	// 3. TEST TRACKS CRUD API
	t.Run("CreateTrack_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Technical",
			"min_credits": 15,
			"description": "Coding and tech challenges",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/admin/tracks", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
		}

		var track models.Track
		err := config.DB.Where("name = ?", "Technical").First(&track).Error
		if err != nil {
			t.Errorf("Failed to find created track: %v", err)
		}
	})

	// 4. TEST MASTER ACTIVITIES CRUD API
	var activityID uint
	t.Run("CreateActivity_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Annual Codeathon 2026",
			"category":    "Technical",
			"credits":     25,
			"description": "24 hour coding contest",
			"status":      "active",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/admin/activities", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
		}

		var activity models.Activity
		config.DB.Where("name = ?", "Annual Codeathon 2026").First(&activity)
		activityID = activity.ID
	})

	// 5. TEST STUDENT SUBMISSION FLOW
	var submissionID uint
	t.Run("StudentSubmitCertificate_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"activity_name":   "Annual Codeathon 2026",
			"category":        "Technical",
			"description":     "Won 2nd prize in the Codeathon",
			"credits":         25,
			"activity_id":     activityID,
			"certificate_url": "http://example.com/certificate.pdf",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/student/submissions", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+studentToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
		}

		var sub models.ActivitySubmission
		config.DB.Where("student_roll = ? AND activity_name = ?", studentRoll, "Annual Codeathon 2026").First(&sub)
		submissionID = sub.ID
		if sub.Status != "pending" {
			t.Errorf("Expected status to be pending, got %s", sub.Status)
		}
	})

	// 6. TEST ADMIN AUDITING API
	t.Run("AuditSubmission_Approve_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"status":  "approved",
			"credits": 20, // Approved with 20 instead of 25 requested
			"remarks": "Verified, awarded 20 credits.",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/admin/submissions/"+strconv.Itoa(int(submissionID))+"/audit", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var sub models.ActivitySubmission
		config.DB.First(&sub, submissionID)
		if sub.Status != "approved" || sub.Credits != 20 || sub.Remarks != "Verified, awarded 20 credits." {
			t.Errorf("Submission audit details mismatch: status=%s, credits=%d", sub.Status, sub.Credits)
		}
	})

	// 7. TEST REPORTS CENTER API
	t.Run("GetReportsSummary_Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/admin/reports/summary", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("ExportStudentsReport_Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/admin/reports/export/students", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var report []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&report)
		if len(report) != 1 {
			t.Errorf("Expected 1 student in report, got %d", len(report))
		}

		totalCredits := report[0]["total_credits"].(float64)
		if totalCredits != 20 {
			t.Errorf("Expected student's total credits to be 20, got %f", totalCredits)
		}
	})

	// 8. TEST USER MANAGEMENT API
	t.Run("CreateAdminDirectly_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":             "New Admin Person",
			"roll_no":          "admin2",
			"contact_no":       "0000000001",
			"email_id":         "admin2@admin.com",
			"enrollment_no":    "ADMIN002",
			"password":         "adminpassword1",
			"confirm_password": "adminpassword1",
			"role":             "admin",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/admin/users/create-admin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
		}

		var newAdmin models.Student
		err := config.DB.Where("roll_no = ?", "admin2").First(&newAdmin).Error
		if err != nil {
			t.Errorf("Failed to find created admin account: %v", err)
		}
		if newAdmin.Role != "admin" {
			t.Errorf("Expected role 'admin', got '%s'", newAdmin.Role)
		}
	})

	t.Run("UpdateStudentRole_Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"role": "admin",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/admin/users/"+studentRoll+"/role", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var stud models.Student
		config.DB.Where("roll_no = ?", studentRoll).First(&stud)
		if stud.Role != "admin" {
			t.Errorf("Expected student role to become 'admin', got '%s'", stud.Role)
		}
	})
}
