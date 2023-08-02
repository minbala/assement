package schema

type LoginInput struct {
	Email    string `validate:"required" json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
