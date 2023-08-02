package jwt

import (
	"test_assessment/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthJwtTokenInput struct {
	UserId uint  `json:"userId"`
	Iat    int64 `json:"iat"`
}

type AuthJwtClaims struct {
	AuthJwtTokenInput
	jwt.StandardClaims
}

// GenerateAuthJwtToken generate tokens
func GenerateAuthJwtToken(data AuthJwtTokenInput, secret string, expiryTime int64) (string, error) {
	claims := AuthJwtClaims{
		data,
		jwt.StandardClaims{
			ExpiresAt: expiryTime,
			Issuer:    "authentication-service",
		},
	}
	token, err := GenerateJwtToken(claims, secret)
	return token, err
}

func GenerateJwtToken(claims jwt.Claims, secretKey string) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

type jwtManager struct {
}

var Manager *jwtManager

func SetUp() {
	Manager = &jwtManager{}
}

func (jwtManager) CreateAuthTokens(userId uint) (accessToken string, err error) {
	return CreateAccessToken(userId, time.Now().AddDate(0, 0, setting.AppSetting.AccessTokenExpiredTime))
}

func CreateAccessToken(userId uint, expiredTime time.Time) (accessToken string, err error) {
	return GenerateAuthJwtToken(AuthJwtTokenInput{UserId: userId, Iat: time.Now().Unix()},
		setting.AppSetting.JWTSECRET, expiredTime.Unix())
}

// ParseAuthJwtToken parsing token
func (jwtManager) ParseAuthJwtToken(token string, secretKey string) (uint, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &AuthJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*AuthJwtClaims); ok && tokenClaims.Valid {
			return claims.UserId, nil
		}
	}
	return 0, err
}
