package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Username      string    `gorm:"size:50;not null" json:"username"`
	StudentID     *string   `gorm:"size:20;uniqueIndex:idx_student_id" json:"student_id"`
	EmployeeID    *string   `gorm:"size:20;uniqueIndex:idx_employee_id" json:"employee_id"`
	Phone         string    `gorm:"size:11;not null" json:"phone"`
	College       string    `gorm:"size:100;not null" json:"college"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	Avatar        string    `gorm:"size:255" json:"avatar"`
	Role          string    `gorm:"size:20;default:'student'" json:"role"`
	Status        string    `gorm:"size:20;default:'active'" json:"status"`
	LoginAttempts int       `gorm:"default:0" json:"-"`
	LockedUntil   time.Time `json:"-"`
}

type LoginLog struct {
	ID         uint      `gorm:"primarykey"`
	UserID     uint      `gorm:"index"`
	IP         string    `gorm:"size:50"`
	Success    bool      `gorm:"default:true"`
	Locked     bool      `gorm:"default:false"`
	CreatedAt  time.Time
}

const (
	RoleStudent  = "student"
	RoleAdmin    = "admin"
	StatusActive   = "active"
	StatusDisabled = "disabled"
)
