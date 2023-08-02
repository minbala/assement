package models

import (
	"time"
)

// UserLogs ! Database Model
type UserLogs struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserId       uint
	Method       string
	RequestUrl   string
	ServiceType  string
	Status       string `json:"status" example:"fail,success"`
	ErrorMessage string
}
