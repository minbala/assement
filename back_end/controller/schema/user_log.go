package schema

import "time"

type GetUserLogsResponse struct {
	UserLogs []UserLog `json:"userLogs"`
	Total    int64     `json:"total"`
}

type UserLog struct {
	Id           uint      `json:"id"`
	UserId       uint      `json:"userId"`
	Method       string    `json:"method"`
	RequestUrl   string    `json:"requestUrl"`
	ServiceType  string    `json:"serviceType"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}
