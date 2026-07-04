package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
)

// SubmitCertificate handles dynamic certificate upload and submission
func SubmitCertificate(c *fiber.Ctx) error {
	studentRoll := c.Locals("roll_no").(string)

	activityName := c.FormValue("activity_name")
	category := c.FormValue("category")
	description := c.FormValue("description")
	creditsStr := c.FormValue("credits")
	activityIDStr := c.FormValue("activity_id")
	certificateURL := c.FormValue("certificate_url") // fallback text URL parameter for testing

	if activityName == "" || category == "" || creditsStr == "" {
		// Fallback check: try JSON binding if form-data parsing is not used (e.g., in unit tests)
		type JSONSubmissionInput struct {
			ActivityName   string `json:"activity_name"`
			Category       string `json:"category"`
			Description    string `json:"description"`
			Credits        int    `json:"credits"`
			ActivityID     *uint  `json:"activity_id"`
			CertificateURL string `json:"certificate_url"`
		}
		var jsonInput JSONSubmissionInput
		if err := c.BodyParser(&jsonInput); err == nil && jsonInput.ActivityName != "" {
			activityName = jsonInput.ActivityName
			category = jsonInput.Category
			description = jsonInput.Description
			creditsStr = strconv.Itoa(jsonInput.Credits)
			certificateURL = jsonInput.CertificateURL
			if jsonInput.ActivityID != nil {
				activityIDStr = strconv.Itoa(int(*jsonInput.ActivityID))
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Activity name, category, and credits are required",
			})
		}
	}

	credits, err := strconv.Atoi(creditsStr)
	if err != nil || credits <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Credits must be a valid positive integer",
		})
	}

	var activityID *uint
	if activityIDStr != "" {
		if idVal, err := strconv.Atoi(activityIDStr); err == nil {
			val := uint(idVal)
			activityID = &val
		}
	}

	// Handle file upload
	fileHeader, err := c.FormFile("certificate")
	if err == nil && fileHeader != nil {
		uploadDir := "./uploads/certificates"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create uploads directory",
			})
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		filePath := fmt.Sprintf("%s/%s", uploadDir, filename)
		if err := c.SaveFile(fileHeader, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save certificate file",
			})
		}
		certificateURL = "/uploads/certificates/" + filename
	}

	submission := models.ActivitySubmission{
		StudentRoll:    studentRoll,
		ActivityID:     activityID,
		ActivityName:   activityName,
		Category:       category,
		Description:    description,
		CertificateURL: certificateURL,
		Credits:        credits,
		Status:         "pending", // default pending
		Remarks:        "",
	}

	if err := config.DB.Create(&submission).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record certificate submission",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(submission)
}

// GetStudentSubmissions lists current student's uploads
func GetStudentSubmissions(c *fiber.Ctx) error {
	studentRoll := c.Locals("roll_no").(string)

	var submissions []models.ActivitySubmission
	if err := config.DB.Where("student_roll = ?", studentRoll).Order("created_at desc").Find(&submissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve submissions",
		})
	}

	return c.JSON(submissions)
}

// GetStudentTracks retrieves master tracks list mapped with student's earned credits
func GetStudentTracks(c *fiber.Ctx) error {
	studentRoll := c.Locals("roll_no").(string)

	var tracks []models.Track
	if err := config.DB.Find(&tracks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tracks",
		})
	}

	type TrackProgressResponse struct {
		models.Track
		EarnedCredits int `json:"earned_credits"`
	}

	progress := make([]TrackProgressResponse, 0)

	for _, track := range tracks {
		var earnedCredits int64
		config.DB.Model(&models.ActivitySubmission{}).
			Where("student_roll = ? AND category = ? AND status = ?", studentRoll, track.Name, "approved").
			Select("COALESCE(SUM(credits), 0)").
			Row().
			Scan(&earnedCredits)

		progress = append(progress, TrackProgressResponse{
			Track:         track,
			EarnedCredits: int(earnedCredits),
		})
	}

	return c.JSON(progress)
}

// GetStudentAnnouncements fetches announcements targeted to students
func GetStudentAnnouncements(c *fiber.Ctx) error {
	now := time.Now()
	var announcements []models.Announcement

	err := config.DB.Where(
		"target_audience IN ? AND (expires_at IS NULL OR expires_at > ?)",
		[]string{"all", "student"}, now,
	).Order("is_pinned desc, created_at desc").Find(&announcements).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve announcements",
		})
	}

	return c.JSON(announcements)
}

// GetStudentActivities fetches active master activities
func GetStudentActivities(c *fiber.Ctx) error {
	var activities []models.Activity
	err := config.DB.Where("status = ?", "active").Order("name asc").Find(&activities).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve master activities",
		})
	}

	return c.JSON(activities)
}
