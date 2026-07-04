package controllers

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/gorm"
)

// --- INPUT SCHEMAS ---

type CreateAdminInput struct {
	Name            string `json:"name"`
	RollNo          string `json:"roll_no"`
	ContactNo       string `json:"contact_no"`
	EmailID         string `json:"email_id"`
	EnrollmentNo    string `json:"enrollment_no"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            string `json:"role"` // 'admin' or 'superadmin'
}

type UpdateUserRoleInput struct {
	Role string `json:"role"`
}

type VerifyUserInput struct {
	IsVerified bool `json:"is_verified"`
}

type ActivityInput struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Credits     int    `json:"credits"`
	Description string `json:"description"`
	Status      string `json:"status"` // 'active' or 'inactive'
}

type TrackInput struct {
	Name        string `json:"name"`
	MinCredits  int    `json:"min_credits"`
	Description string `json:"description"`
}

type AuditSubmissionInput struct {
	Status  string `json:"status"` // 'approved' or 'rejected'
	Credits int    `json:"credits"`
	Remarks string `json:"remarks"`
}

type AnnouncementInput struct {
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Category       string     `json:"category"`
	TargetAudience string     `json:"target_audience"` // 'all', 'student', 'admin'
	IsPinned       bool       `json:"is_pinned"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

type SettingInput struct {
	Value string `json:"value"`
}

// --- CONTROLLER HANDLERS ---

// GetAdminDashboard returns general system stats and summaries
func GetAdminDashboard(c *fiber.Ctx) error {
	var totalStudents int64
	var totalAdmins int64
	var pendingSubmissions int64
	var approvedSubmissions int64
	var totalActivities int64
	var totalTracks int64
	var totalCreditsApproved int64

	config.DB.Model(&models.Student{}).Where("role = ?", "student").Count(&totalStudents)
	config.DB.Model(&models.Student{}).Where("role IN ?", []string{"admin", "superadmin"}).Count(&totalAdmins)
	config.DB.Model(&models.ActivitySubmission{}).Where("status = ?", "pending").Count(&pendingSubmissions)
	config.DB.Model(&models.ActivitySubmission{}).Where("status = ?", "approved").Count(&approvedSubmissions)
	config.DB.Model(&models.Activity{}).Count(&totalActivities)
	config.DB.Model(&models.Track{}).Count(&totalTracks)

	// Calculate sum of approved credits
	row := config.DB.Model(&models.ActivitySubmission{}).
		Where("status = ?", "approved").
		Select("COALESCE(SUM(credits), 0)").
		Row()
	_ = row.Scan(&totalCreditsApproved)

	// Fetch 5 most recent submissions
	var submissions []models.ActivitySubmission
	if err := config.DB.Order("created_at desc").Limit(5).Find(&submissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch recent submissions"})
	}

	// Attach student names to submissions
	type RecentSubmissionResponse struct {
		models.ActivitySubmission
		StudentName string `json:"student_name"`
	}
	recent := make([]RecentSubmissionResponse, 0)
	for _, sub := range submissions {
		var stud models.Student
		name := "Unknown"
		if err := config.DB.Select("name").Where("roll_no = ?", sub.StudentRoll).First(&stud).Error; err == nil {
			name = stud.Name
		}
		recent = append(recent, RecentSubmissionResponse{
			ActivitySubmission: sub,
			StudentName:        name,
		})
	}

	return c.JSON(fiber.Map{
		"total_students":         totalStudents,
		"total_admins":           totalAdmins,
		"pending_submissions":    pendingSubmissions,
		"approved_submissions":   approvedSubmissions,
		"total_activities":       totalActivities,
		"total_tracks":           totalTracks,
		"total_credits_approved": totalCreditsApproved,
		"recent_submissions":     recent,
	})
}

// GetUsers lists users (students/admins) with filters
func GetUsers(c *fiber.Ctx) error {
	role := c.Query("role")
	status := c.Query("status")
	search := c.Query("search")

	var users []models.Student
	query := config.DB.Model(&models.Student{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		isVerified := status == "verified"
		query = query.Where("is_verified = ?", isVerified)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR roll_no LIKE ? OR email_id LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return c.JSON(users)
}

// GetUserDetail returns a specific user's profile and submissions
func GetUserDetail(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")

	var user models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var submissions []models.ActivitySubmission
	config.DB.Where("student_roll = ?", rollNo).Find(&submissions)

	return c.JSON(fiber.Map{
		"user":        user,
		"submissions": submissions,
	})
}

// UpdateUserRole changes a student's role (admin vs superadmin vs student)
func UpdateUserRole(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")
	var input UpdateUserRoleInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Role != "student" && input.Role != "admin" && input.Role != "superadmin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role value"})
	}

	var user models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Prevent updating self role
	currentAdminRoll := c.Locals("roll_no").(string)
	if rollNo == currentAdminRoll {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot update your own role"})
	}

	user.Role = input.Role
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	return c.JSON(fiber.Map{
		"message": "User role updated successfully",
		"user":    user,
	})
}

// VerifyUser verifies or suspends a student's account
func VerifyUser(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")
	var input VerifyUserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var user models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	user.IsVerified = input.IsVerified
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update verification status"})
	}

	return c.JSON(fiber.Map{
		"message": "User verification status updated successfully",
		"user":    user,
	})
}

// DeleteUser deletes user account
func DeleteUser(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")

	var user models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Prevent deleting self
	currentAdminRoll := c.Locals("roll_no").(string)
	if rollNo == currentAdminRoll {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot delete your own account"})
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// CreateAdmin allows superadmin to directly create other admin accounts
func CreateAdmin(c *fiber.Ctx) error {
	var input CreateAdminInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	// Validation
	if input.Name == "" || input.RollNo == "" || input.EmailID == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, roll number, email, and password are required"})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	if input.Role != "admin" && input.Role != "superadmin" {
		input.Role = "admin"
	}

	var existingUser models.Student
	res := config.DB.Where("roll_no = ? OR email_id = ?", input.RollNo, input.EmailID).First(&existingUser)
	if res.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User with this Roll No or Email ID already exists"})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	newAdmin := models.Student{
		RollNo:       input.RollNo,
		Name:         input.Name,
		CourseName:   "Administration",
		Semester:     0,
		ContactNo:    input.ContactNo,
		EmailID:      input.EmailID,
		EnrollmentNo: input.EnrollmentNo,
		Password:     hashedPassword,
		IsVerified:   true, // Admins directly verified
		Role:         input.Role,
	}

	if err := config.DB.Create(&newAdmin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create administrator"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Administrator account created successfully",
		"user":    newAdmin,
	})
}

// --- MASTER ACTIVITIES HANDLERS ---

func GetActivities(c *fiber.Ctx) error {
	var activities []models.Activity
	if err := config.DB.Find(&activities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch activities"})
	}
	return c.JSON(activities)
}

func CreateActivity(c *fiber.Ctx) error {
	var input ActivityInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Name == "" || input.Category == "" || input.Credits <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, category, and credit points (>0) are required"})
	}

	if input.Status != "active" && input.Status != "inactive" {
		input.Status = "active"
	}

	activity := models.Activity{
		Name:        input.Name,
		Category:    input.Category,
		Credits:     input.Credits,
		Description: input.Description,
		Status:      input.Status,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create activity"})
	}

	return c.Status(fiber.StatusCreated).JSON(activity)
}

func UpdateActivity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var input ActivityInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var activity models.Activity
	if err := config.DB.First(&activity, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Activity not found"})
	}

	if input.Name != "" {
		activity.Name = input.Name
	}
	if input.Category != "" {
		activity.Category = input.Category
	}
	if input.Credits > 0 {
		activity.Credits = input.Credits
	}
	if input.Description != "" {
		activity.Description = input.Description
	}
	if input.Status == "active" || input.Status == "inactive" {
		activity.Status = input.Status
	}

	if err := config.DB.Save(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update activity"})
	}

	return c.JSON(activity)
}

func DeleteActivity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var activity models.Activity
	if err := config.DB.First(&activity, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Activity not found"})
	}

	if err := config.DB.Delete(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete activity"})
	}

	return c.JSON(fiber.Map{"message": "Activity deleted successfully"})
}

// --- SUBMISSION AUDITING HANDLERS ---

func GetSubmissions(c *fiber.Ctx) error {
	status := c.Query("status")
	rollNo := c.Query("roll_no")
	category := c.Query("category")

	var submissions []models.ActivitySubmission
	query := config.DB.Model(&models.ActivitySubmission{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if rollNo != "" {
		query = query.Where("student_roll = ?", rollNo)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&submissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch submissions"})
	}

	return c.JSON(submissions)
}

func GetSubmissionDetail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var submission models.ActivitySubmission
	if err := config.DB.First(&submission, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Submission not found"})
	}

	var student models.Student
	_ = config.DB.Where("roll_no = ?", submission.StudentRoll).First(&student)

	return c.JSON(fiber.Map{
		"submission": submission,
		"student":    student,
	})
}

func AuditSubmission(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var input AuditSubmissionInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Status != "approved" && input.Status != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Audit status must be either 'approved' or 'rejected'"})
	}

	var submission models.ActivitySubmission
	if err := config.DB.First(&submission, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Submission not found"})
	}

	if submission.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Submission is already audited"})
	}

	adminEmail := c.Locals("email").(string)
	now := time.Now()

	submission.Status = input.Status
	submission.Remarks = input.Remarks
	submission.VerifiedBy = adminEmail
	submission.VerifiedAt = &now

	if input.Status == "approved" {
		if input.Credits >= 0 {
			submission.Credits = input.Credits
		}
	} else {
		submission.Credits = 0 // Rejections result in 0 points
	}

	if err := config.DB.Save(&submission).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to audit submission"})
	}

	return c.JSON(submission)
}

// --- CURRICULUM TRACK HANDLERS ---

func GetTracks(c *fiber.Ctx) error {
	var tracks []models.Track
	if err := config.DB.Find(&tracks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tracks"})
	}
	return c.JSON(tracks)
}

func CreateTrack(c *fiber.Ctx) error {
	var input TrackInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Name == "" || input.MinCredits < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and min credits are required"})
	}

	track := models.Track{
		Name:        input.Name,
		MinCredits:  input.MinCredits,
		Description: input.Description,
	}

	if err := config.DB.Create(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create track"})
	}

	return c.Status(fiber.StatusCreated).JSON(track)
}

func UpdateTrack(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var input TrackInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var track models.Track
	if err := config.DB.First(&track, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
	}

	if input.Name != "" {
		track.Name = input.Name
	}
	if input.MinCredits >= 0 {
		track.MinCredits = input.MinCredits
	}
	if input.Description != "" {
		track.Description = input.Description
	}

	if err := config.DB.Save(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update track"})
	}

	return c.JSON(track)
}

func DeleteTrack(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var track models.Track
	if err := config.DB.First(&track, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
	}

	if err := config.DB.Delete(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete track"})
	}

	return c.JSON(fiber.Map{"message": "Track deleted successfully"})
}

// --- ANNOUNCEMENT HANDLERS ---

func GetAnnouncements(c *fiber.Ctx) error {
	var announcements []models.Announcement
	if err := config.DB.Order("created_at desc").Find(&announcements).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch announcements"})
	}
	return c.JSON(announcements)
}

func CreateAnnouncement(c *fiber.Ctx) error {
	var input AnnouncementInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Title == "" || input.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title and content are required"})
	}

	adminEmail := c.Locals("email").(string)

	ann := models.Announcement{
		Title:          input.Title,
		Content:        input.Content,
		Category:       input.Category,
		TargetAudience: input.TargetAudience,
		IsPinned:       input.IsPinned,
		CreatedBy:      adminEmail,
		ExpiresAt:      input.ExpiresAt,
	}

	if err := config.DB.Create(&ann).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create announcement"})
	}

	return c.Status(fiber.StatusCreated).JSON(ann)
}

func UpdateAnnouncement(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var input AnnouncementInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var ann models.Announcement
	if err := config.DB.First(&ann, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Announcement not found"})
	}

	if input.Title != "" {
		ann.Title = input.Title
	}
	if input.Content != "" {
		ann.Content = input.Content
	}
	if input.Category != "" {
		ann.Category = input.Category
	}
	if input.TargetAudience != "" {
		ann.TargetAudience = input.TargetAudience
	}
	ann.IsPinned = input.IsPinned
	if input.ExpiresAt != nil {
		ann.ExpiresAt = input.ExpiresAt
	}

	if err := config.DB.Save(&ann).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update announcement"})
	}

	return c.JSON(ann)
}

func DeleteAnnouncement(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var ann models.Announcement
	if err := config.DB.First(&ann, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Announcement not found"})
	}

	if err := config.DB.Delete(&ann).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete announcement"})
	}

	return c.JSON(fiber.Map{"message": "Announcement deleted successfully"})
}

// --- REPORTS CENTER HANDLERS ---

// GetReportsSummary yields general summaries of credits earned by courses
func GetReportsSummary(c *fiber.Ctx) error {
	type CourseSummary struct {
		CourseName    string  `json:"course_name"`
		TotalStudents int64   `json:"total_students"`
		AvgCredits    float64 `json:"avg_credits"`
	}

	var courses []string
	config.DB.Model(&models.Student{}).Where("role = ?", "student").Distinct("course_name").Pluck("course_name", &courses)

	summaries := make([]CourseSummary, 0)

	for _, course := range courses {
		var studCount int64
		config.DB.Model(&models.Student{}).Where("role = ? AND course_name = ?", "student", course).Count(&studCount)

		// Calculate average credits for this course
		var totalCredits int64
		var studentRolls []string
		config.DB.Model(&models.Student{}).Where("role = ? AND course_name = ?", "student", course).Pluck("roll_no", &studentRolls)

		if len(studentRolls) > 0 {
			config.DB.Model(&models.ActivitySubmission{}).
				Where("status = ? AND student_roll IN ?", "approved", studentRolls).
				Select("COALESCE(SUM(credits), 0)").
				Row().
				Scan(&totalCredits)
		}

		avg := 0.0
		if studCount > 0 {
			avg = float64(totalCredits) / float64(studCount)
			avg = math.Round(avg*100) / 100 // Round to 2 decimal places
		}

		summaries = append(summaries, CourseSummary{
			CourseName:    course,
			TotalStudents: studCount,
			AvgCredits:    avg,
		})
	}

	return c.JSON(summaries)
}

// ExportStudentsReport returns a structured data table list of students and their track credits
func ExportStudentsReport(c *fiber.Ctx) error {
	type StudentReportItem struct {
		RollNo       string         `json:"roll_no"`
		Name         string         `json:"name"`
		CourseName   string         `json:"course_name"`
		Semester     int            `json:"semester"`
		EmailID      string         `json:"email_id"`
		EnrollmentNo string         `json:"enrollment_no"`
		TrackCredits map[string]int `json:"track_credits"`
		TotalCredits int            `json:"total_credits"`
	}

	var students []models.Student
	if err := config.DB.Where("role = ?", "student").Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch student list"})
	}

	var tracks []models.Track
	config.DB.Find(&tracks)

	report := make([]StudentReportItem, 0)

	for _, s := range students {
		trackCreds := make(map[string]int)
		for _, t := range tracks {
			trackCreds[t.Name] = 0
		}

		// Query approved credits by category for this student
		type CreditsByCategory struct {
			Category string
			Sum      int
		}
		var categorySums []CreditsByCategory
		config.DB.Model(&models.ActivitySubmission{}).
			Select("category, SUM(credits) as sum").
			Where("student_roll = ? AND status = ?", s.RollNo, "approved").
			Group("category").
			Scan(&categorySums)

		total := 0
		for _, cs := range categorySums {
			trackCreds[cs.Category] = cs.Sum
			total += cs.Sum
		}

		report = append(report, StudentReportItem{
			RollNo:       s.RollNo,
			Name:         s.Name,
			CourseName:   s.CourseName,
			Semester:     s.Semester,
			EmailID:      s.EmailID,
			EnrollmentNo: s.EnrollmentNo,
			TrackCredits: trackCreds,
			TotalCredits: total,
		})
	}

	return c.JSON(report)
}

// --- SYSTEM SETTINGS HANDLERS ---

func GetSystemSettings(c *fiber.Ctx) error {
	var settings []models.SystemSetting
	if err := config.DB.Find(&settings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch settings"})
	}
	return c.JSON(settings)
}

func UpdateSystemSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	var input SettingInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var setting models.SystemSetting
	if err := config.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Setting not found"})
	}

	setting.Value = input.Value
	setting.UpdatedAt = time.Now()

	if err := config.DB.Save(&setting).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update setting"})
	}

	return c.JSON(setting)
}
