package repository

import "test_assessment/domain/model"

type UserLogRepositoryInterface interface {
	Create(model.CreateUserLogInput) error
	Get(model.GetUserLogInput) (model.GetUserLogOutput, error)
	GetTotalUserLogCount(model.GetTotalUserLogCountInput) (int64, error)
}
