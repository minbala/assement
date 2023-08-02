package model

type CreateSessionInput struct {
	UserId      uint
	AccessToken string
}

type DeleteSessionInput struct {
	UserId      uint
	AccessToken string
}

type GetSessionInput struct {
	AccessToken string
	UserId      uint
}

type GetSessionOutput struct {
	Id          uint
	UserId      uint   `json:"userId" gorm:"not null"`
	AccessToken string `json:"accessToken" gorm:"not null"`
	User        User
}
