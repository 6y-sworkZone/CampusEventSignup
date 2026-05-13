package models

import (
	"time"
	"gorm.io/gorm"
)

type Registration struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID      uint      `gorm:"not null;index:idx_user_event,unique" json:"user_id"`
	EventID     uint      `gorm:"not null;index:idx_user_event,unique" json:"event_id"`
	Status      string    `gorm:"size:20;default:'registered'" json:"status"`
	IsWaitlist  bool      `gorm:"default:false" json:"is_waitlist"`
	WaitlistNum int       `json:"waitlist_num"`
	Signed      bool      `gorm:"default:false" json:"signed"`
	SignedAt    time.Time `json:"signed_at"`
	
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event       Event     `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

const (
	RegStatusRegistered = "registered"
	RegStatusWaitlist   = "waitlist"
	RegStatusCancelled  = "cancelled"
	RegStatusSigned     = "signed"
)
