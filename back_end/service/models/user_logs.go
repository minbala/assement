package models

import "time"

type GetUserLogsOutput struct {
	UserLogs []UserLog
	Total    int64
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
