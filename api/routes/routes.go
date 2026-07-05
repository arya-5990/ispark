package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/controllers"
	"github.com/iips-oss/ispark/api/middleware"
)

// SetupRoutes configures the endpoints for the API
func SetupRoutes(app *fiber.App) {
	// Base API group
	api := app.Group("/api")

	// Auth group
	auth := api.Group("/auth")

	// Public Auth routes
	auth.Get("/captcha", controllers.GetCaptcha)
	auth.Post("/register", controllers.Register)
	auth.Post("/verify-otp", controllers.VerifyOTP)
	auth.Post("/login", controllers.Login)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/refresh", controllers.RefreshToken)

	// Protected routes (Require login)
	auth.Use(middleware.AuthRequired())
	auth.Post("/logout", controllers.Logout)
	auth.Get("/profile", controllers.GetProfile)

	// Student route group (Require student role)
	student := api.Group("/student", middleware.AuthRequired(), middleware.RoleRequired("student"))
	student.Post("/submissions", controllers.SubmitCertificate)
	student.Get("/submissions", controllers.GetStudentSubmissions)
	student.Get("/tracks", controllers.GetStudentTracks)
	student.Get("/announcements", controllers.GetStudentAnnouncements)
	student.Get("/activities", controllers.GetStudentActivities)

	// Admin route group (Require admin/superadmin roles)
	admin := api.Group("/admin", middleware.AuthRequired(), middleware.RoleRequired("admin", "superadmin"))

	// Dashboard
	admin.Get("/dashboard", controllers.GetAdminDashboard)

	// User Management
	admin.Get("/users", controllers.GetUsers)
	admin.Get("/users/:roll_no", controllers.GetUserDetail)
	admin.Put("/users/:roll_no/verify", controllers.VerifyUser)

	// User Management (Superadmin Only)
	admin.Put("/users/:roll_no/role", middleware.RoleRequired("superadmin"), controllers.UpdateUserRole)
	admin.Delete("/users/:roll_no", middleware.RoleRequired("superadmin"), controllers.DeleteUser)
	admin.Post("/users/create-admin", middleware.RoleRequired("superadmin"), controllers.CreateAdmin)

	// Activity templates (Master list)
	admin.Get("/activities", controllers.GetActivities)
	admin.Post("/activities", controllers.CreateActivity)
	admin.Put("/activities/:id", controllers.UpdateActivity)
	admin.Delete("/activities/:id", controllers.DeleteActivity)

	// Activity Submissions / Audits
	admin.Get("/submissions", controllers.GetSubmissions)
	admin.Get("/submissions/:id", controllers.GetSubmissionDetail)
	admin.Put("/submissions/:id/audit", controllers.AuditSubmission)

	// Tracks
	admin.Get("/tracks", controllers.GetTracks)
	admin.Post("/tracks", controllers.CreateTrack)
	admin.Put("/tracks/:id", controllers.UpdateTrack)
	admin.Delete("/tracks/:id", controllers.DeleteTrack)

	// Announcements
	admin.Get("/announcements", controllers.GetAnnouncements)
	admin.Post("/announcements", controllers.CreateAnnouncement)
	admin.Put("/announcements/:id", controllers.UpdateAnnouncement)
	admin.Delete("/announcements/:id", controllers.DeleteAnnouncement)

	// Reports Center
	admin.Get("/reports/summary", controllers.GetReportsSummary)
	admin.Get("/reports/export/students", controllers.ExportStudentsReport)

	// System Settings
	admin.Get("/settings", controllers.GetSystemSettings)
	admin.Put("/settings/:key", controllers.UpdateSystemSetting)
}
