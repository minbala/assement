package services

import (
	"net/http"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/pkg/resources"
	"test_assessment/service/models"
)

// PublicService  represents public service.
type PublicService struct {
	logger          Logger
	passwordManager PasswordManager
	jwtManager      JwtManager
	sessionManager  repository.SessionRepositoryInterface
	userManager     repository.UserRepositoryInterface
}

func NewPublicService(
	logger Logger, passwordManager PasswordManager, jwtManager JwtManager,
	sessionManager repository.SessionRepositoryInterface, userManager repository.UserRepositoryInterface) *PublicService {
	return &PublicService{
		logger:          logger,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
		sessionManager:  sessionManager,
		userManager:     userManager,
	}
}

func (p PublicService) Login(input models.LoginInput) (string, uint, error) {
	user, err := p.userManager.Get(domainModel.GetUserInput{Email: input.Email, Limit: 1})
	if err != nil {
		p.logger.LogError(err)
		return "", 0, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}
	}
	if len(user.Users) == 0 {

		return "", 0, ErrorResponse{Code: http.StatusBadRequest, ResponseMessage: resources.RecordNotFound,
			ErrorString: resources.RecordNotFound}
	}
	err = p.CheckPassword(input.Password, user.Users[0].Password)
	if err != nil {
		p.logger.LogError(err)
		return "", 0, ErrorResponse{Code: http.StatusBadRequest, ResponseMessage: err.Error(),
			ErrorString: err.Error()}
	}
	accessToken, err := p.CreateUserSession(user.Users[0].Id)
	if err != nil {
		p.logger.LogError(err)
		return "", 0, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}
	}
	err = p.sessionManager.Create(domainModel.CreateSessionInput{UserId: user.Users[0].Id,
		AccessToken: accessToken})
	if err != nil {
		p.logger.LogError(err)
		return "", 0,
			ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
				ErrorString: err.Error()}
	}
	return accessToken, user.Users[0].Id, nil
}

func (p PublicService) CheckPassword(inputPassword, dbPassword string) error {
	return p.passwordManager.CheckPassword(inputPassword, dbPassword)
}

func (p PublicService) CreateUserSession(userID uint) (accessToken string, err error) {
	accessToken, err = p.jwtManager.CreateAuthTokens(userID)
	if err != nil {
		return
	}
	err = p.sessionManager.Create(domainModel.CreateSessionInput{UserId: userID, AccessToken: accessToken})
	return
}
