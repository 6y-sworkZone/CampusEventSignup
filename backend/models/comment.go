package models

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	EventID  uint    `gorm:"not null;index" json:"event_id"`
	UserID   uint    `gorm:"not null;index" json:"user_id"`
	ParentID *uint   `gorm:"index" json:"parent_id"`
	Content  string  `gorm:"type:text;not null" json:"content"`
	
	User    User      `gorm:"foreignKey:UserID" json:"user"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}
