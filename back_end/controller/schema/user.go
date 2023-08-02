package schema

import "time"

type CreateUserInput struct {
	Email    string `validate:"required" json:"email"`
	Name     string `validate:"required" json:"name"`
	Password string `validate:"required" json:"password"`
	UserRole string `validate:"required,user_role" json:"userRole"`
}

type CreateUserResponse struct {
	Id        uint      `json:"id"`
	Email     string    `validate:"required" json:"email"`
	Name      string    `validate:"required" json:"name"`
	UserRole  string    `validate:"required" json:"userRole,user_role"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	Id        uint      `json:"id"`
	Email     string    `validate:"required" json:"email"`
	Name      string    `validate:"required" json:"name"`
	UserRole  string    `validate:"required" json:"userRole,user_role"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetUserResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}

type UpdateUserInput struct {
	UserId   uint   `validate:"required" json:"userId"`
	Email    string `validate:"required" json:"email"`
	Name     string
	Password string ` json:"password"`
	UserRole string `validate:"required" json:"userRole,user_role"`
}
