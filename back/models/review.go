package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	ProductID    uint           `gorm:"not null" json:"product_id"`
	Product      Product        `gorm:"foreignKey:ProductID" json:"product"`
	Rating       int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Title        string         `gorm:"not null" json:"title"`
	Comment      string         `json:"comment"`
	IsVerified   bool           `gorm:"default:false" json:"is_verified"` // Verified purchase
	HelpfulCount int            `gorm:"default:0" json:"helpful_count"`
	ReportCount  int            `gorm:"default:0" json:"report_count"`
	Status       string         `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
}

// ReviewHelpful tracks which users found a review helpful
type ReviewHelpful struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ReviewID  uint      `gorm:"not null" json:"review_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Review    Review    `gorm:"foreignKey:ReviewID" json:"review"`
	CreatedAt time.Time `json:"created_at"`
}

// ReviewReport tracks review reports
type ReviewReport struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ReviewID  uint      `gorm:"not null" json:"review_id"`
	Reason    string    `gorm:"not null" json:"reason"`
	Details   string    `json:"details"`
	Status    string    `gorm:"default:'pending'" json:"status"` // pending, resolved, dismissed
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Review    Review    `gorm:"foreignKey:ReviewID" json:"review"`
}
