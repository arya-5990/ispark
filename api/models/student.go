package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	RollNo       string         `gorm:"primaryKey;type:varchar(50)" json:"roll_no"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	CourseName   string         `gorm:"type:varchar(100);not null" json:"course_name"`
	Semester     int            `gorm:"not null" json:"semester"`
	ContactNo    string         `gorm:"type:varchar(20);not null" json:"contact_no"`
	EmailID      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email_id"`
	EnrollmentNo string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"enrollment_no"`
	Password     string         `gorm:"type:varchar(255);not null" json:"-"` // "-" hides password from json marshalling
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	Role         string         `gorm:"type:varchar(20);default:'student'" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
