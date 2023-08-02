package models

import "time"

type CreateUserInput struct {
	Email    string
	Name     string
	Password string
	UserRole string
}

type CreateUserOutput struct {
	Id        uint
	Email     string
	Name      string
	UserRole  string
	CreatedAt time.Time
}

type User struct {
	Id        uint
	Name      string
	Email     string
	Password  string
	UserRole  string
	CreatedAt time.Time
}

type GetUserOutput struct {
	Users []User
	Total int64
}

type UpdateUserInput struct {
	Email    string
	Name     string
	Password string
	UserRole string
	UserId   uint
}
