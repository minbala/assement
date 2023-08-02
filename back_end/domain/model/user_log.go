package model

import "time"

type CreateUserLogInput struct {
	UserId       uint
	Method       string
	RequestUrl   string
	ServiceType  string
	Status       string
	ErrorMessage string
}

type UserLog struct {
	Id           uint
	UserId       uint
	Method       string
	RequestUrl   string
	ServiceType  string
	Status       string
	ErrorMessage string
	CreatedAt    time.Time
}
type GetUserLogOutput struct {
	UserLogs []UserLog
}

type GetUserLogInput struct {
	UserId uint
	Limit  uint
	Offset uint
}

type GetTotalUserLogCountInput struct {
	UserId uint
}
