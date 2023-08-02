package repository

import "test_assessment/domain/model"

type SessionRepositoryInterface interface {
	Create(model.CreateSessionInput) error
	Delete(input model.DeleteSessionInput) error
	Get(input model.GetSessionInput) (model.GetSessionOutput, error)
}
