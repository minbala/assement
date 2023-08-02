package mysql

import (
	"errors"
	"gorm.io/gorm/clause"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/persistence/mysql/models"
	"test_assessment/pkg/resources"
)

type SessionRepository struct {
}

func (SessionRepository) Create(
	session domainModel.CreateSessionInput) error {
	return psqlDB.Create(&models.Session{UserId: session.UserId, AccessToken: session.AccessToken}).Error
}

func (SessionRepository) Delete(
	token domainModel.DeleteSessionInput) error {
	return psqlDB.Where("user_id = ? and access_token = ?", token.UserId, token.AccessToken).Delete(&models.Session{}).Error
}

func (SessionRepository) Get(data domainModel.GetSessionInput) (response domainModel.GetSessionOutput, err error) {
	if data.UserId != 0 && !isStringValid(data.AccessToken) {
		return domainModel.GetSessionOutput{}, errors.New(resources.ClientError)
	}
	var session models.Session
	err = psqlDB.Preload(clause.Associations).Where("user_id = ? and access_token = ? ", data.UserId, data.AccessToken).Find(&session).Error
	if err != nil {
		return domainModel.GetSessionOutput{}, err
	}
	if session.ID == 0 {
		return domainModel.GetSessionOutput{}, errors.New(resources.RecordNotFound)
	}
	return domainModel.GetSessionOutput{UserId: session.UserId, AccessToken: session.AccessToken, Id: session.ID,
		User: domainModel.User{Id: session.User.ID, Name: session.User.Name, Email: session.User.Email,
			UserRole: session.User.UserRole}}, nil
}

// NewSessionRepository  gets a new SessionRepository instance.
func NewSessionRepository() repository.SessionRepositoryInterface {
	return SessionRepository{}
}
