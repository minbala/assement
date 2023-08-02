package models

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name" gorm:"not null"`
	Email     string `json:"email" gorm:"not null"`
	Password  string `json:"-" gorm:"not null"`
	UserRole  string `json:"userRole" gorm:"not null"`
}
