package model

import "time"

type GetUserInput struct {
	Email    string
	UserRole string
	Name     string
	Limit    uint
	Offset   uint
	UserId   uint
}

type GetTotalUserCountInput struct {
	Name     string
	UserRole string
}

type User struct {
	Id        uint
	Name      string `json:"name" gorm:"not null"`
	Email     string `json:"email" gorm:"not null"`
	Password  string `json:"-" gorm:"not null"`
	UserRole  string `json:"userRole" gorm:"not null"`
	CreatedAt time.Time
}

type GetUserOutput struct {
	Users []User
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	UserRole string
}

type CreateUserOutput struct {
	Id        uint
	Name      string
	Email     string
	UserRole  string
	CreatedAt time.Time
}

type DeleteUserInput struct {
	UserId uint
}

type UpdateUserInput struct {
	UserId   uint
	Name     string
	Email    string
	Password string
	UserRole string
}
