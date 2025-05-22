package models

import (
	"time"

	"gorm.io/gorm"
)

type RequestLog struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	IP           string         `gorm:"index" json:"ip"`
	Method       string         `json:"method"`
	Path         string         `json:"path"`
	Status       int            `json:"status"`
	UserAgent    string         `json:"user_agent"`
	UserID       *uint          `json:"user_id,omitempty"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category     string         `gorm:"index" json:"category"` // normal, warning, anomaly
	Details      string         `json:"details"`               // Additional request details
	ResponseTime float64        `json:"response_time"`         // Response time in milliseconds
}

// BannedIP represents a banned IP address
type BannedIP struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	IP          string    `gorm:"uniqueIndex" json:"ip"`
	BannedUntil time.Time `json:"banned_until"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `json:"created_at"`
}
