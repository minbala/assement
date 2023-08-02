package services

import (
	"net/http"
	"net/mail"
	"strings"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/pkg/resources"
	"test_assessment/service/models"
)

// AdminService  represents admin service.
type AdminService struct {
	logger          Logger
	passwordManager PasswordManager
	jwtManager      JwtManager
	sessionManager  repository.SessionRepositoryInterface
	userManager     repository.UserRepositoryInterface
	userLogsManager repository.UserLogRepositoryInterface
}

func NewAdminService(
	logger Logger, passwordManager PasswordManager, jwtManager JwtManager, sessionManager repository.SessionRepositoryInterface,
	userManager repository.UserRepositoryInterface, userLogsManager repository.UserLogRepositoryInterface) *AdminService {
	return &AdminService{
		logger:          logger,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
		sessionManager:  sessionManager,
		userManager:     userManager,
		userLogsManager: userLogsManager,
	}
}

func (a AdminService) Logout(token string, userId uint) error {
	token = strings.Split(token, " ")[1]
	err := a.sessionManager.Delete(domainModel.DeleteSessionInput{AccessToken: token, UserId: userId})
	if err != nil {
		a.logger.LogError(err)
		return ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}
	}
	return nil
}

func (a AdminService) CreateUser(input models.CreateUserInput) (models.CreateUserOutput, error) {
	if !isStringValid(input.Password) {
		return models.CreateUserOutput{}, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: resources.ClientError}
	}
	if _, ok := mail.ParseAddress(input.Email); ok != nil {
		return models.CreateUserOutput{}, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: resources.ClientError}
	}
	hashPassword, err := a.passwordManager.Hash(input.Password)
	if err != nil {
		a.logger.LogError(err)
		return models.CreateUserOutput{}, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}

	}
	user, err := a.userManager.Create(domainModel.CreateUserInput{Name: input.Name, Email: input.Email, Password: hashPassword,
		UserRole: input.UserRole})
	if err != nil {
		a.logger.LogError(err)
		return models.CreateUserOutput{},
			ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
				ErrorString: err.Error()}
	}

	return models.CreateUserOutput{Id: user.Id, Name: user.Name, Email: user.Email, UserRole: user.UserRole,
		CreatedAt: user.CreatedAt}, nil
}

func (a AdminService) GetUsers(name, role string, limit, offset uint) (response models.GetUserOutput, err error) {
	out := make(chan error)
	go func() {
		var err error
		response.Total, err = a.userManager.GetTotalUserCount(domainModel.GetTotalUserCountInput{Name: name, UserRole: role})
		out <- err
	}()
	go func() {
		user, err := a.userManager.Get(domainModel.GetUserInput{Name: name, UserRole: role, Limit: limit, Offset: offset})
		if err != nil {
			out <- err
			return
		}
		for _, user := range user.Users {
			response.Users = append(response.Users, models.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Password:  user.Password,
				UserRole:  user.UserRole,
				CreatedAt: user.CreatedAt,
			})

		}
		out <- nil
	}()
	for i := 0; i < 2; i++ {
		err = <-out
		if err != nil {
			a.logger.LogError(err)
			return models.GetUserOutput{}, ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
				ErrorString: err.Error()}
		}
	}
	return response, nil
}

func (a AdminService) GetUserById(id uint) (response models.User, err error) {
	user, err := a.userManager.Get(domainModel.GetUserInput{UserId: id, Limit: 1})
	if err != nil {
		a.logger.LogError(err)
		return models.User{},
			ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
				ErrorString: err.Error()}

	}
	for _, user := range user.Users {
		response = models.User{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			UserRole:  user.UserRole,
			CreatedAt: user.CreatedAt,
		}
	}
	return response, nil
}

func (a AdminService) DeleteUser(id uint) (err error) {
	err = a.userManager.Delete(domainModel.DeleteUserInput{UserId: id})
	if err != nil {
		a.logger.LogError(err)
		return ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}

	}
	return nil
}

func isStringValid(testString string) bool {
	if m := strings.TrimSpace(testString); len(m) == 0 || len(testString) == 0 {
		return false
	}
	return true
}

func (a AdminService) UpdateUser(input models.UpdateUserInput) (err error) {
	var hashPassword string
	if isStringValid(input.Password) {
		hashPassword, err = a.passwordManager.Hash(input.Password)
		if err != nil {
			a.logger.LogError(err)
			return ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
				ErrorString: err.Error()}
		}
	}

	err = a.userManager.Update(domainModel.UpdateUserInput{Name: input.Name, Email: input.Email, Password: hashPassword,
		UserRole: input.UserRole, UserId: input.UserId})
	if err != nil {
		a.logger.LogError(err)
		return ErrorResponse{Code: http.StatusInternalServerError, ResponseMessage: resources.InternalServerError,
			ErrorString: err.Error()}
	}
	return nil
}

func (a AdminService) GetUserLog(userId, limit, offset uint) (response models.GetUserLogsOutput, err error) {
	out := make(chan error)
	go func() {
		var err error
		response.Total, err = a.userLogsManager.GetTotalUserLogCount(domainModel.GetTotalUserLogCountInput{UserId: userId})
		out <- err
	}()
	go func() {
		userLogs, err := a.userLogsManager.Get(domainModel.GetUserLogInput{UserId: userId, Limit: limit, Offset: offset})
		if err != nil {
			out <- err
			return
		}
		for _, userLog := range userLogs.UserLogs {
			response.UserLogs = append(response.UserLogs, models.UserLog{
				Id:           userLog.Id,
				UserId:       userLog.UserId,
				Method:       userLog.Method,
				RequestUrl:   userLog.RequestUrl,
				ServiceType:  userLog.ServiceType,
				Status:       userLog.Status,
				ErrorMessage: userLog.ErrorMessage,
				CreatedAt:    userLog.CreatedAt,
			})
		}
		out <- nil
	}()
	for i := 0; i < 2; i++ {
		err = <-out
		if err != nil {
			a.logger.LogError(err)
			return models.GetUserLogsOutput{}, ErrorResponse{Code: http.StatusInternalServerError,
				ResponseMessage: resources.InternalServerError, ErrorString: err.Error()}
		}
	}
	return response, nil
}
