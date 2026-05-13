package models

import (
	"time"
	"gorm.io/gorm"
)

type Event struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	
	Title           string    `gorm:"size:100;not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	CoverImage      string    `gorm:"size:255" json:"cover_image"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	RegistrationEnd time.Time `json:"registration_end"`
	MaxParticipants int       `gorm:"not null" json:"max_participants"`
	CurrentCount    int       `gorm:"default:0" json:"current_count"`
	WaitlistCount   int       `gorm:"default:0" json:"waitlist_count"`
	Status          string    `gorm:"size:20;default:'draft'" json:"status"`
	SignCode        string    `gorm:"size:4" json:"-"`
	CreatorID       uint      `gorm:"not null" json:"creator_id"`
	Tags            []EventTag `gorm:"foreignKey:EventID" json:"tags"`
}

type EventTag struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	EventID uint   `gorm:"index" json:"event_id"`
	Name    string `gorm:"size:50;not null" json:"name"`
}

const (
	EventStatusDraft     = "draft"
	EventStatusOpen      = "open"
	EventStatusOngoing   = "ongoing"
	EventStatusEnded     = "ended"
	EventStatusCancelled = "cancelled"
)
