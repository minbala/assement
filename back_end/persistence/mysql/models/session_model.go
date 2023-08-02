package models

import (
	"time"
)

type Session struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      uint   `json:"userId" gorm:"not null"`
	User        *User  `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" swaggerignore:"true"`
	AccessToken string `json:"accessToken" gorm:"not null"`
}
