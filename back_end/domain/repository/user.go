package repository

import "test_assessment/domain/model"

type UserRepositoryInterface interface {
	Get(input model.GetUserInput) (model.GetUserOutput, error)
	GetTotalUserCount(input model.GetTotalUserCountInput) (int64, error)
	Create(input model.CreateUserInput) (model.CreateUserOutput, error)
	Delete(input model.DeleteUserInput) error
	Update(input model.UpdateUserInput) error
}
