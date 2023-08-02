package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	domainModel "test_assessment/domain/model"
	pg "test_assessment/persistence/mysql"
	"test_assessment/pkg/jwt"
	"test_assessment/service/models"
	"testing"
	"time"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Get(input domainModel.GetUserInput) (domainModel.GetUserOutput, error) {
	args := m.Called(input)
	var output domainModel.GetUserOutput
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(domainModel.GetUserOutput); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) GetTotalUserCount(input domainModel.GetTotalUserCountInput) (int64, error) {
	args := m.Called(input)
	var output int64
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(int64); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) Create(input domainModel.CreateUserInput) (domainModel.CreateUserOutput, error) {
	args := m.Called(input)
	var output domainModel.CreateUserOutput
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(domainModel.CreateUserOutput); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) Delete(input domainModel.DeleteUserInput) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserRepo) Update(input domainModel.UpdateUserInput) error {
	args := m.Called(input)
	return args.Error(0)
}

type MockPasswordManager struct {
	mock.Mock
}

func (m *MockPasswordManager) CheckPassword(password string, hashPassword string) error {
	args := m.Called(password, hashPassword)
	return args.Error(0)
}
func (m *MockPasswordManager) Hash(plain string) (string, error) {
	args := m.Called(plain)
	var output string
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(string); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogError(err error) {
	log.Println(err)
}

func (m *MockLogger) LogInfo(message string) {
	log.Println(message)
}

func (m *MockLogger) LogDebugMessage(message string) {
	log.Println(message)
}

// need to add mock jwt manager, session manager, userLog manager
func TestCreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockPasswordManager := new(MockPasswordManager)
	mockLogger := new(MockLogger)

	adminService := NewAdminService(mockLogger, mockPasswordManager, jwt.Manager,
		pg.NewSessionRepository(), mockUserRepo, pg.NewUserLogRepository())
	_, err := adminService.CreateUser(models.CreateUserInput{Name: "bala", Email: "minbala33@gmail.com", Password: "",
		UserRole: "admin"})
	assert.Error(t, err)

	_, err = adminService.CreateUser(models.CreateUserInput{Name: "bala", Email: "minbala33", Password: "minbala33",
		UserRole: "admin"})
	assert.Error(t, err)
	createdAt := time.Now()
	mockPasswordManager.On("Hash", "minbala33").Return("minbala3333", nil)
	mockUserRepo.On("Create", mock.Anything).Return(domainModel.CreateUserOutput{
		Name: "bala", Email: "minbala33@gmail.com", Id: 1, CreatedAt: createdAt,
		UserRole: "admin",
	}, nil)
	output, err := adminService.CreateUser(models.CreateUserInput{Name: "bala", Email: "minbala33@gmail.com", Password: "minbala33",
		UserRole: "admin"})

	mockUserCall := mockUserRepo.On("Get", mock.Anything).Return(domainModel.GetUserOutput{
		Users: []domainModel.User{{Name: "bala", Email: "minbala33@gmail.com", Id: 1, CreatedAt: createdAt,
			UserRole: "admin", Password: "minbala3333"}}}, nil)
	assert.NoError(t, err)
	user, err := adminService.GetUserById(output.Id)
	assert.NoError(t, err)
	hashPassword, err := adminService.passwordManager.Hash("minbala33")
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		Id:        output.Id,
		Name:      output.Name,
		Email:     output.Email,
		Password:  hashPassword,
		UserRole:  output.UserRole,
		CreatedAt: output.CreatedAt,
	}, user)
	mockUserRepo.On("Update", mock.Anything).Return(nil)
	err = adminService.UpdateUser(models.UpdateUserInput{UserId: output.Id, Name: "minbala", Password: "",
		UserRole: "user", Email: "minbala22@gmail.com"})
	assert.NoError(t, err)
	mockUserCall.Unset()
	mockUserCall = mockUserRepo.On("Get", mock.Anything).Return(domainModel.GetUserOutput{
		Users: []domainModel.User{{Name: "minbala", Email: "minbala22@gmail.com", Id: 1, CreatedAt: createdAt,
			UserRole: "user", Password: "minbala3333"}}}, nil)
	user, err = adminService.GetUserById(output.Id)
	assert.NoError(t, err)
	hashPassword, err = adminService.passwordManager.Hash("minbala33")
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		Id:        output.Id,
		Name:      "minbala",
		Email:     "minbala22@gmail.com",
		Password:  hashPassword,
		UserRole:  "user",
		CreatedAt: output.CreatedAt,
	}, user)
	mockUserRepo.On("Delete", mock.Anything).Return(nil)
	err = adminService.DeleteUser(output.Id)
	assert.NoError(t, err)
	mockUserCall.Unset()
	mockUserRepo.On("Get", mock.Anything).Return(domainModel.GetUserOutput{}, errors.New("no record found"))
	user, err = adminService.GetUserById(output.Id)
	assert.Error(t, err)
}
