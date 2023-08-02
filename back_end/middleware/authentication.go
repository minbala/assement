package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/pkg/resources"
	"test_assessment/pkg/setting"
)

type JwtManager interface {
	ParseAuthJwtToken(token string, secretKey string) (uint, error)
}

type authenticator struct {
	jwtManager     JwtManager
	sessionManager repository.SessionRepositoryInterface
}

var AdminAuthenticationManager *authenticator

func SetUp(jwtManager JwtManager, sessionManager repository.SessionRepositoryInterface) {
	AdminAuthenticationManager = &authenticator{
		jwtManager:     jwtManager,
		sessionManager: sessionManager,
	}
}

func (a authenticator) CheckTokenIsValid(token string, userId uint) (bool, string) {
	session, err := a.sessionManager.Get(domainModel.GetSessionInput{AccessToken: token, UserId: userId})
	if err != nil {
		return false, ""
	}
	if session.Id == 0 {
		return false, ""
	}
	return true, session.User.UserRole
}

type UserData struct {
	UserId   uint
	UserRole string
}

func isStringValid(testString string) bool {
	if m := strings.TrimSpace(testString); len(m) == 0 || len(testString) == 0 {
		return false
	}
	return true
}

func parseToken(token string) (string, bool) {
	if !isStringValid(token) {
		return "", false
	}
	validationTokens := strings.Split(token, " ")
	if len(validationTokens) != 2 {
		return "", false
	}
	return validationTokens[1], true
}

func (a authenticator) ValidateToken(token string) (data UserData, message string, status int) {
	if !isStringValid(token) {
		return UserData{}, resources.ClientError, http.StatusUnauthorized
	}
	token, ok := parseToken(token)
	if !ok {
		return UserData{}, resources.ClientError, http.StatusUnauthorized
	}
	userId, err := a.jwtManager.ParseAuthJwtToken(token, setting.AppSetting.JWTSECRET)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			return UserData{}, resources.JwtTokenExpired, http.StatusUnauthorized
		default:

			return UserData{}, resources.UnAuthorized, http.StatusUnauthorized
		}
	}
	ok, userRole := a.CheckTokenIsValid(token, userId)
	if !ok {
		return UserData{}, resources.UnAuthorized, http.StatusUnauthorized
	}
	return UserData{UserId: userId, UserRole: userRole}, resources.Success, http.StatusOK
}

func (a authenticator) ValidateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, message, statusCode := a.ValidateToken(c.GetHeader("Authorization"))
		if statusCode != 200 {
			c.JSON(statusCode, gin.H{
				"message": message,
			})
			c.Abort()
			return
		}
		if userInfo.UserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": resources.UnAuthorized,
			})
			c.Abort()
			return
		}
		c.Set("UserData", userInfo)
		c.Next()
	}
}

func (a authenticator) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, message, statusCode := a.ValidateToken(c.GetHeader("Authorization"))
		if statusCode != 200 {
			c.JSON(statusCode, gin.H{
				"message": message,
			})
			c.Abort()
			return
		}
		c.Set("UserData", userInfo)
		c.Next()
	}
}
