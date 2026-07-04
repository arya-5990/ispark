package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(200);not null" json:"title"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Category       string         `gorm:"type:varchar(50);default:'general'" json:"category"`
	TargetAudience string         `gorm:"type:varchar(50);default:'all'" json:"target_audience"` // 'all', 'student', 'admin'
	IsPinned       bool           `gorm:"default:false" json:"is_pinned"`
	CreatedBy      string         `gorm:"type:varchar(100);not null" json:"created_by"`
	ExpiresAt      *time.Time     `json:"expires_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
