package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Category    string         `gorm:"type:varchar(50);not null" json:"category"` // e.g. Literary, Cultural, Technical, Academic, Social Service
	Credits     int            `gorm:"not null" json:"credits"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // 'active', 'inactive'
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
