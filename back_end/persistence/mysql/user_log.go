package mysql

import (
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/persistence/mysql/models"
)

type UserLogRepository struct {
}

func (UserLogRepository) Create(
	userLog domainModel.CreateUserLogInput) error {
	return psqlDB.Create(&models.UserLogs{UserId: userLog.UserId, Method: userLog.Method, RequestUrl: userLog.RequestUrl,
		ServiceType: userLog.ServiceType, Status: userLog.Status, ErrorMessage: userLog.ErrorMessage}).Error
}

func (UserLogRepository) Get(
	input domainModel.GetUserLogInput) (response domainModel.GetUserLogOutput, err error) {
	query := psqlDB.Model(&models.UserLogs{})
	if input.UserId != 0 {
		query = query.Where("user_id = ?", input.UserId)
	}
	var userLogs []models.UserLogs
	err = query.Order("id desc").Limit(int(input.Limit)).Offset(int(input.Offset)).Find(&userLogs).Error
	if err != nil {
		return
	}
	for _, user := range userLogs {
		response.UserLogs = append(response.UserLogs, domainModel.UserLog{
			Id:           user.ID,
			UserId:       user.UserId,
			Method:       user.Method,
			RequestUrl:   user.RequestUrl,
			ServiceType:  user.ServiceType,
			Status:       user.Status,
			ErrorMessage: user.ErrorMessage,
			CreatedAt:    user.CreatedAt,
		})
	}
	return
}

func (UserLogRepository) GetTotalUserLogCount(
	input domainModel.GetTotalUserLogCountInput) (count int64, err error) {
	query := psqlDB.Model(&models.UserLogs{})
	if input.UserId != 0 {
		query = query.Where("user_id = ?", input.UserId)
	}
	err = query.Count(&count).Error
	if err != nil {
		return
	}
	return
}

// NewUserLogRepository  gets a newUserLogRepository instance.
func NewUserLogRepository() repository.UserLogRepositoryInterface {
	return UserLogRepository{}
}
