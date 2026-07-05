package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivitySubmission struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	StudentRoll    string         `gorm:"type:varchar(50);not null;index" json:"student_roll"`
	ActivityID     *uint          `json:"activity_id"` // null for custom external activity
	ActivityName   string         `gorm:"type:varchar(150);not null" json:"activity_name"`
	Category       string         `gorm:"type:varchar(50);not null" json:"category"`
	Description    string         `gorm:"type:text" json:"description"`
	CertificateURL string         `gorm:"type:varchar(255)" json:"certificate_url"`
	Credits        int            `gorm:"not null" json:"credits"`                          // credits requested initially / approved eventually
	Status         string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // 'pending', 'approved', 'rejected'
	Remarks        string         `gorm:"type:text" json:"remarks"`
	VerifiedBy     string         `gorm:"type:varchar(100)" json:"verified_by"`
	VerifiedAt     *time.Time     `json:"verified_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
