package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

// Input schemas
type RegisterInput struct {
	Name            string `json:"name"`
	RollNo          string `json:"roll_no"`
	CourseName      string `json:"course_name"`
	Semester        int    `json:"semester"`
	ContactNo       string `json:"contact_no"`
	EmailID         string `json:"email_id"`
	EnrollmentNo    string `json:"enrollment_no"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type VerifyOTPInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginInput struct {
	EmailID       string `json:"email_id"`
	Password      string `json:"password"`
	CaptchaID     string `json:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer"`
}

type ForgotPasswordInput struct {
	EmailID string `json:"email_id"`
}

type ResetPasswordInput struct {
	EmailID         string `json:"email_id"`
	OTPCode         string `json:"otp_code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// GetCaptcha returns a new math-based captcha challenge
func GetCaptcha(c *fiber.Ctx) error {
	challenge, err := utils.GenerateCaptchaChallenge()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate captcha",
		})
	}
	return c.JSON(challenge)
}

// Register registers a new student and sends a verification OTP
func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// Basic Validation
	if input.Name == "" || input.RollNo == "" || input.CourseName == "" || input.Semester <= 0 ||
		input.ContactNo == "" || input.EmailID == "" || input.EnrollmentNo == "" ||
		input.Password == "" || input.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required and semester must be greater than 0",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password and Confirm Password do not match",
		})
	}

	// Check if student already exists (by RollNo, EmailID, or EnrollmentNo)
	var existingStudent models.Student
	result := config.DB.Where("roll_no = ? OR email_id = ? OR enrollment_no = ?", input.RollNo, input.EmailID, input.EnrollmentNo).First(&existingStudent)
	if result.Error == nil {
		// If the existing student is not verified yet, we can delete the record to allow re-registration
		if !existingStudent.IsVerified {
			config.DB.Unscoped().Delete(&existingStudent)
			// Also clean up any pending OTPs for this email to prevent conflict
			config.DB.Where("email = ? AND purpose = ?", existingStudent.EmailID, "register").Delete(&models.OTP{})
		} else {
			if existingStudent.RollNo == input.RollNo {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Student with this Roll Number already exists",
				})
			}
			if existingStudent.EmailID == input.EmailID {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Student with this Email ID already exists",
				})
			}
			if existingStudent.EnrollmentNo == input.EnrollmentNo {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Student with this Enrollment Number already exists",
				})
			}
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create student (unverified)
	student := models.Student{
		RollNo:       input.RollNo,
		Name:         input.Name,
		CourseName:   input.CourseName,
		Semester:     input.Semester,
		ContactNo:    input.ContactNo,
		EmailID:      input.EmailID,
		EnrollmentNo: input.EnrollmentNo,
		Password:     hashedPassword,
		IsVerified:   false,
	}

	if err := config.DB.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register student",
		})
	}

	// Generate OTP
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		// Clean up the created student since we failed to generate OTP
		config.DB.Unscoped().Delete(&student)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate verification code",
		})
	}

	// Save OTP to database
	otp := models.OTP{
		Email:     input.EmailID,
		Code:      otpCode,
		Purpose:   "register",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := config.DB.Create(&otp).Error; err != nil {
		log.Printf("Error saving OTP to DB: %v", err)
		config.DB.Unscoped().Delete(&student)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save verification code",
		})
	}

	// Send OTP
	err = utils.SendOTP(input.EmailID, otpCode, "student registration")
	if err != nil {
		log.Printf("Error sending OTP email: %v", err)
		// We still proceed since the email sender falls back to printing in local environment console
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Student registered successfully. Please verify the OTP sent to your email.",
		"email":   input.EmailID,
	})
}

// VerifyOTP verifies the registration OTP and activates the account, logging the user in
func VerifyOTP(c *fiber.Ctx) error {
	var input VerifyOTPInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if input.Email == "" || input.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Code are required",
		})
	}

	// Find the OTP code in database
	var otp models.OTP
	err := config.DB.Where("email = ? AND code = ? AND purpose = ? AND expires_at > ?",
		input.Email, input.Code, "register", time.Now()).First(&otp).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired OTP verification code",
		})
	}

	// Mark student as verified
	var student models.Student
	if err := config.DB.Where("email_id = ?", input.Email).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student account not found",
		})
	}

	student.IsVerified = true
	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify student account",
		})
	}

	// Delete the verified OTP
	config.DB.Delete(&otp)

	// Automatically log the user in
	accessToken, err := utils.GenerateAccessToken(student.RollNo, student.EmailID, student.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(student.RollNo, student.EmailID, student.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Set refresh token in HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"message":      "Account verified successfully and logged in",
		"access_token": accessToken,
		"student": fiber.Map{
			"roll_no":       student.RollNo,
			"name":          student.Name,
			"email_id":      student.EmailID,
			"course_name":   student.CourseName,
			"semester":      student.Semester,
			"contact_no":    student.ContactNo,
			"enrollment_no": student.EnrollmentNo,
		},
	})
}

// Login authenticates a student and generates tokens
func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if input.EmailID == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Password are required",
		})
	}

	// Captcha verification (Required)
	if input.CaptchaID == "" || input.CaptchaAnswer == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Captcha identification and answer are required",
		})
	}

	captchaValid, err := utils.VerifyCaptcha(input.CaptchaID, input.CaptchaAnswer)
	if err != nil || !captchaValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired captcha answer. Please try again.",
		})
	}

	// Retrieve student
	var student models.Student
	err = config.DB.Where("email_id = ?", input.EmailID).First(&student).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Verify Password
	if !utils.CheckPasswordHash(input.Password, student.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Verify if account is activated via OTP
	if !student.IsVerified {
		// Send a new verification OTP
		otpCode, err := utils.GenerateOTP()
		if err == nil {
			otp := models.OTP{
				Email:     student.EmailID,
				Code:      otpCode,
				Purpose:   "register",
				ExpiresAt: time.Now().Add(15 * time.Minute),
			}
			config.DB.Create(&otp)
			if err := utils.SendOTP(student.EmailID, otpCode, "student registration reactivation"); err != nil {
				log.Printf("Error sending registration reactivation OTP: %v", err)
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Student email is not verified yet. A new OTP has been sent to your email.",
			"email": student.EmailID,
		})
	}

	// Generate Tokens
	accessToken, err := utils.GenerateAccessToken(student.RollNo, student.EmailID, student.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(student.RollNo, student.EmailID, student.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Set refresh token in HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"message":      "Logged in successfully",
		"access_token": accessToken,
		"student": fiber.Map{
			"roll_no":       student.RollNo,
			"name":          student.Name,
			"email_id":      student.EmailID,
			"course_name":   student.CourseName,
			"semester":      student.Semester,
			"contact_no":    student.ContactNo,
			"enrollment_no": student.EnrollmentNo,
		},
	})
}

// RefreshToken issues a new access token based on a valid refresh token
func RefreshToken(c *fiber.Ctx) error {
	// Try to get refresh token from cookie, then body
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		type BodyStruct struct {
			RefreshToken string `json:"refresh_token"`
		}
		var body BodyStruct
		if err := c.BodyParser(&body); err == nil {
			refreshToken = body.RefreshToken
		}
	}

	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is missing",
		})
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateAccessToken(claims.RollNo, claims.Email, claims.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

// Logout clears the user tokens
func Logout(c *fiber.Ctx) error {
	// Clear the refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expired
		HTTPOnly: true,
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// ForgotPassword initiates password reset by sending an OTP
func ForgotPassword(c *fiber.Ctx) error {
	var input ForgotPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if input.EmailID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// Check if student exists
	var student models.Student
	err := config.DB.Where("email_id = ?", input.EmailID).First(&student).Error
	if err != nil {
		// To prevent email harvesting, return a success message even if the email is not registered
		return c.JSON(fiber.Map{
			"message": "If the email is registered, a password reset code has been sent.",
		})
	}

	// Generate OTP
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate verification code",
		})
	}

	// Invalidate any older forgot_password OTPs for this email
	config.DB.Where("email = ? AND purpose = ?", input.EmailID, "forgot_password").Delete(&models.OTP{})

	// Save new OTP
	otp := models.OTP{
		Email:     input.EmailID,
		Code:      otpCode,
		Purpose:   "forgot_password",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := config.DB.Create(&otp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save verification code",
		})
	}

	// Send OTP email
	if err := utils.SendOTP(input.EmailID, otpCode, "password reset"); err != nil {
		log.Printf("Error sending password reset OTP: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "A verification code has been sent to your email.",
		"email":   input.EmailID,
	})
}

// ResetPassword resets password using the OTP received
func ResetPassword(c *fiber.Ctx) error {
	var input ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if input.EmailID == "" || input.OTPCode == "" || input.Password == "" || input.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password and Confirm Password do not match",
		})
	}

	// Check OTP Validity
	var otp models.OTP
	err := config.DB.Where("email = ? AND code = ? AND purpose = ? AND expires_at > ?",
		input.EmailID, input.OTPCode, "forgot_password", time.Now()).First(&otp).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired verification code",
		})
	}

	// Update password
	var student models.Student
	if err := config.DB.Where("email_id = ?", input.EmailID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student account not found",
		})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	student.Password = hashedPassword
	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reset password",
		})
	}

	// Delete used OTP
	config.DB.Delete(&otp)

	return c.JSON(fiber.Map{
		"message": "Password reset successfully. You can now login with your new password.",
	})
}

// GetProfile returns the authenticated student's profile
func GetProfile(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	var student models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student profile not found",
		})
	}

	return c.JSON(fiber.Map{
		"student": fiber.Map{
			"roll_no":       student.RollNo,
			"name":          student.Name,
			"email_id":      student.EmailID,
			"course_name":   student.CourseName,
			"semester":      student.Semester,
			"contact_no":    student.ContactNo,
			"enrollment_no": student.EnrollmentNo,
			"is_verified":   student.IsVerified,
			"role":          student.Role,
			"created_at":    student.CreatedAt,
			"updated_at":    student.UpdatedAt,
		},
	})
}
