package models

import (
	"time"

	"gorm.io/gorm"
)

type UserActivity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Type      string         `gorm:"not null" json:"type"` // view, click, search, etc.
	ProductID *uint          `json:"product_id,omitempty"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Details   string         `json:"details"` // Additional activity details
	SessionID string         `gorm:"index" json:"session_id"`
}
