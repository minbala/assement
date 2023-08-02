package services

type Logger interface {
	LogError(error)
	LogInfo(string)
	LogDebugMessage(string)
}

type PasswordManager interface {
	CheckPassword(password string, hashPassword string) error
	Hash(plain string) (string, error)
}

type ErrorResponse struct {
	ResponseMessage string
	ErrorString     string
	Code            int
}

func (e ErrorResponse) Error() string {
	return e.ErrorString
}

type JwtManager interface {
	CreateAuthTokens(userId uint) (accessToken string, err error)
	ParseAuthJwtToken(token string, secretKey string) (uint, error)
}

type Authenticator interface {
	Validate(userId uint) (accessToken string, err error)
}
