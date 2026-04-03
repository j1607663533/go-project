package models

import (
	"time"
	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Path      string         `gorm:"size:255;not null" json:"path"`
	Size      int64          `json:"size"`
	Ext       string         `gorm:"size:20" json:"ext"`
	Type      string         `gorm:"size:50" json:"type"` // image, document, etc.
	UserID    uint           `json:"user_id"`
	Username  string         `gorm:"size:100" json:"username"`
}
