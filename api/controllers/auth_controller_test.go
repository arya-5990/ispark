package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/routes"
	"gorm.io/gorm"
)

// SetupTestDB initializes an in-memory SQLite database for testing and overrides config.DB
func SetupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory SQLite database: %v", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(&models.Student{}, &models.OTP{})
	if err != nil {
		t.Fatalf("Failed to run database migrations: %v", err)
	}

	config.DB = db
}

// Helper to solve the simple math captcha challenge
func solveCaptcha(question string) (string, error) {
	// Question format is: "X + Y = ?"
	re := regexp.MustCompile(`(\d+)\s*\+\s*(\d+)`)
	matches := re.FindStringSubmatch(question)
	if len(matches) < 3 {
		return "", fmt.Errorf("unable to parse captcha question: %s", question)
	}

	val1, err1 := strconv.Atoi(matches[1])
	val2, err2 := strconv.Atoi(matches[2])
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("failed to convert numbers in captcha: %v, %v", err1, err2)
	}

	return strconv.Itoa(val1 + val2), nil
}

func TestAuthFlow(t *testing.T) {
	// Set test environment variables
	os.Setenv("JWT_SECRET", "test_jwt_secret_key_1234567890")
	os.Setenv("JWT_REFRESH_SECRET", "test_jwt_refresh_secret_key_1234567890")

	// Initialize test DB
	SetupTestDB(t)

	// Create and setup Fiber App
	app := fiber.New()
	routes.SetupRoutes(app)

	// Test registration details
	testStudentEmail := "test.student@example.com"
	testPassword := "SecurePassword123!"

	// 1. TEST REGISTRATION FLOW
	t.Run("RegisterStudent_Success", func(t *testing.T) {
		regPayload := map[string]interface{}{
			"name":             "John Doe",
			"roll_no":          "12345",
			"course_name":      "MCA",
			"semester":         2,
			"contact_no":       "9876543210",
			"email_id":         testStudentEmail,
			"enrollment_no":    "EN12345",
			"password":         testPassword,
			"confirm_password": testPassword,
		}
		body, _ := json.Marshal(regPayload)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute register request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
		}

		var respBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		if respBody["email"] != testStudentEmail {
			t.Errorf("Expected registered email to be %s, got %v", testStudentEmail, respBody["email"])
		}
	})

	t.Run("RegisterStudent_DuplicateConflict", func(t *testing.T) {
		// Attempting to register again with same details
		regPayload := map[string]interface{}{
			"name":             "John Doe Duplicate",
			"roll_no":          "12345",
			"course_name":      "MCA",
			"semester":         2,
			"contact_no":       "9876543210",
			"email_id":         testStudentEmail,
			"enrollment_no":    "EN12345",
			"password":         testPassword,
			"confirm_password": testPassword,
		}
		body, _ := json.Marshal(regPayload)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute register request: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
		}

		// Since they are not verified, the old registration is cleared if it's the same student.
		// However, let's verify if registering a different field raises conflict or returns StatusCreated.
		// Wait, according to controllers/auth_controller.go, if they are not verified yet,
		// the unverified record is deleted and it re-registers (status 201 Created).
		// But let's check with a verified student.
		// Let's mark the student as verified in the DB to test the conflict response.
		var dbStudent models.Student
		config.DB.Where("email_id = ?", testStudentEmail).First(&dbStudent)
		dbStudent.IsVerified = true
		config.DB.Save(&dbStudent)

		// Now attempt duplicate registration
		resp2, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute duplicate register request: %v", err)
		}

		if resp2.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 Conflict for duplicate registration, got %d", resp2.StatusCode)
		}
	})

	// 2. TEST OTP VALIDATION
	t.Run("VerifyOTP_Success", func(t *testing.T) {
		// Reset database and re-register unverified student
		SetupTestDB(t)

		regPayload := map[string]interface{}{
			"name":             "Alice Smith",
			"roll_no":          "67890",
			"course_name":      "BTech",
			"semester":         4,
			"contact_no":       "9876543211",
			"email_id":         "alice@example.com",
			"enrollment_no":    "EN67890",
			"password":         testPassword,
			"confirm_password": testPassword,
		}
		body, _ := json.Marshal(regPayload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if _, err := app.Test(req); err != nil {
			t.Fatalf("Failed to execute register request: %v", err)
		}

		// Get OTP code from SQLite
		var otp models.OTP
		err := config.DB.Where("email = ? AND purpose = ?", "alice@example.com", "register").First(&otp).Error
		if err != nil {
			t.Fatalf("Could not find generated OTP in database: %v", err)
		}

		// Verify OTP with correct code
		verifyPayload := map[string]string{
			"email": "alice@example.com",
			"code":  otp.Code,
		}
		vBody, _ := json.Marshal(verifyPayload)

		vReq := httptest.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(vBody))
		vReq.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(vReq)
		if err != nil {
			t.Fatalf("Failed to execute verify OTP request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var respBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		if respBody["access_token"] == nil {
			t.Error("Expected access token in response, got nil")
		}
	})

	t.Run("VerifyOTP_Failure_InvalidCode", func(t *testing.T) {
		verifyPayload := map[string]string{
			"email": "alice@example.com",
			"code":  "999999", // Wrong code
		}
		vBody, _ := json.Marshal(verifyPayload)

		vReq := httptest.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(vBody))
		vReq.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(vReq)
		if err != nil {
			t.Fatalf("Failed to execute verify OTP request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", resp.StatusCode)
		}
	})

	// 3. TEST LOGIN FAILURE/SUCCESS & PROFILE ACCESS
	t.Run("LoginAndProfileFlow", func(t *testing.T) {
		// Fetch a captcha challenge first
		capReq := httptest.NewRequest("GET", "/api/auth/captcha", nil)
		capResp, err := app.Test(capReq)
		if err != nil {
			t.Fatalf("Failed to get captcha: %v", err)
		}
		if capResp.StatusCode != http.StatusOK {
			t.Fatalf("Expected captcha status 200 OK, got %d", capResp.StatusCode)
		}

		var capBody map[string]interface{}
		if err := json.NewDecoder(capResp.Body).Decode(&capBody); err != nil {
			t.Fatalf("Failed to decode captcha response: %v", err)
		}
		captchaID := capBody["captcha_id"].(string)
		question := capBody["question"].(string)

		correctAnswer, err := solveCaptcha(question)
		if err != nil {
			t.Fatalf("Error solving captcha: %v", err)
		}

		// A. Login Failure - Wrong Captcha
		loginWrongCaptcha := map[string]string{
			"email_id":       "alice@example.com",
			"password":       testPassword,
			"captcha_id":     captchaID,
			"captcha_answer": "999", // Wrong answer
		}
		wBody, _ := json.Marshal(loginWrongCaptcha)
		wReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(wBody))
		wReq.Header.Set("Content-Type", "application/json")
		wResp, err := app.Test(wReq)
		if err != nil {
			t.Fatalf("Failed to execute login request: %v", err)
		}
		if wResp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request for wrong captcha, got %d", wResp.StatusCode)
		}

		// B. Login Failure - Wrong Password
		loginWrongPass := map[string]string{
			"email_id":       "alice@example.com",
			"password":       "WrongPassword!",
			"captcha_id":     captchaID,
			"captcha_answer": correctAnswer,
		}
		wpBody, _ := json.Marshal(loginWrongPass)
		wpReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(wpBody))
		wpReq.Header.Set("Content-Type", "application/json")
		wpResp, err := app.Test(wpReq)
		if err != nil {
			t.Fatalf("Failed to execute login request: %v", err)
		}
		if wpResp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for wrong credentials, got %d", wpResp.StatusCode)
		}

		// C. Login Success
		loginSuccessPayload := map[string]string{
			"email_id":       "alice@example.com",
			"password":       testPassword,
			"captcha_id":     captchaID,
			"captcha_answer": correctAnswer,
		}
		sBody, _ := json.Marshal(loginSuccessPayload)
		sReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(sBody))
		sReq.Header.Set("Content-Type", "application/json")
		sResp, err := app.Test(sReq)
		if err != nil {
			t.Fatalf("Failed to execute login request: %v", err)
		}
		if sResp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK for successful login, got %d", sResp.StatusCode)
		}

		var sBodyMap map[string]interface{}
		if err := json.NewDecoder(sResp.Body).Decode(&sBodyMap); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}
		accessToken := sBodyMap["access_token"].(string)

		// D. Profile Access - Unauthenticated
		profReqUnauth := httptest.NewRequest("GET", "/api/auth/profile", nil)
		profRespUnauth, err := app.Test(profReqUnauth)
		if err != nil {
			t.Fatalf("Failed to execute profile request: %v", err)
		}
		if profRespUnauth.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized for missing token, got %d", profRespUnauth.StatusCode)
		}

		// E. Profile Access - Authenticated
		profReqAuth := httptest.NewRequest("GET", "/api/auth/profile", nil)
		profReqAuth.Header.Set("Authorization", "Bearer "+accessToken)
		profRespAuth, err := app.Test(profReqAuth)
		if err != nil {
			t.Fatalf("Failed to execute profile request: %v", err)
		}
		if profRespAuth.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK for authenticated profile access, got %d", profRespAuth.StatusCode)
		}

		var profBody map[string]interface{}
		if err := json.NewDecoder(profRespAuth.Body).Decode(&profBody); err != nil {
			t.Fatalf("Failed to decode profile response: %v", err)
		}
		studentData := profBody["student"].(map[string]interface{})
		if studentData["email_id"] != "alice@example.com" {
			t.Errorf("Expected profile email to be alice@example.com, got %v", studentData["email_id"])
		}
	})
}
