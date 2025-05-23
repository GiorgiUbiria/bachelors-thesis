package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;uniqueIndex" json:"name"`
	Description string         `json:"description"`
	ParentID    *uint          `json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	ImageURL    string         `json:"image_url"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Products    []Product      `gorm:"foreignKey:Category;references:Name" json:"products,omitempty"`
}
