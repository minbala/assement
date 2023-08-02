package mysql

import (
	"strings"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	"test_assessment/persistence/mysql/models"
)

type UserRepository struct {
}

func isStringValid(testString string) bool {
	if m := strings.TrimSpace(testString); len(m) == 0 || len(testString) == 0 {
		return false
	}
	return true
}

func (UserRepository) Get(
	input domainModel.GetUserInput) (response domainModel.GetUserOutput, err error) {
	query := psqlDB.Model(&models.User{})
	if input.UserId != 0 {
		query = query.Where("id = ?", input.UserId)
	}
	if isStringValid(input.Email) {
		query = query.Where("email = ?", input.Email)
	}
	if isStringValid(input.UserRole) {
		query = query.Where("user_role = ?", input.UserRole)
	}
	if isStringValid(input.Name) {
		query = query.Where("name like ?", "%"+input.Name+"%")
	}
	var users []models.User
	err = query.Order("id desc").Limit(int(input.Limit)).Offset(int(input.Offset)).Scan(&users).Error

	if err != nil {
		return
	}
	for _, user := range users {
		response.Users = append(response.Users, domainModel.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			UserRole:  user.UserRole,
			CreatedAt: user.CreatedAt,
		})
	}
	return
}

func (UserRepository) GetTotalUserCount(
	input domainModel.GetTotalUserCountInput) (count int64, err error) {
	query := psqlDB.Model(&models.User{})
	if isStringValid(input.UserRole) {
		query = query.Where("user_role = ?", input.UserRole)
	}
	if isStringValid(input.Name) {
		query = query.Where("name like ?", "%"+input.Name+"%")
	}
	err = query.Count(&count).Error
	if err != nil {
		return
	}
	return
}

func (UserRepository) Create(
	input domainModel.CreateUserInput) (response domainModel.CreateUserOutput, err error) {
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		UserRole: input.UserRole,
	}
	err = psqlDB.Create(&user).Error
	if err != nil {
		return
	}
	return domainModel.CreateUserOutput{Id: user.ID, Name: user.Name,
		Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt}, nil
}

func (UserRepository) Delete(
	input domainModel.DeleteUserInput) (err error) {
	return psqlDB.Unscoped().Delete(models.User{ID: input.UserId}).Error

}

func (UserRepository) Update(
	input domainModel.UpdateUserInput) (err error) {
	userData := map[string]interface{}{"name": input.Name,
		"email": input.Email, "user_role": input.UserRole}
	if isStringValid(input.Password) {
		userData["password"] = input.Password
	}
	return psqlDB.Model(&models.User{ID: input.UserId}).Updates(userData).Error

}

// NewUserRepository  gets a new UserRepository instance.
func NewUserRepository() repository.UserRepositoryInterface {
	return UserRepository{}
}
