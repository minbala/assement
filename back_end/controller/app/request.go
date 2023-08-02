package app

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/mail"
	"strconv"
	"test_assessment/pkg/resources"
	services "test_assessment/service"
)

var userRole map[string]bool = map[string]bool{"admin": true, "user": true}
var validate *validator.Validate

func SetUp() {
	validate = validator.New()
	//intentionally ignore errors
	validate.RegisterValidation("user_role", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(string)
		return userRole[value]
	})

}

// BindAndValid binds and validates data
func (g *Gin) BindAndValid(form interface{}) error {
	err := g.C.Bind(form)
	if err != nil {
		return services.ErrorResponse{Code: http.StatusBadRequest, ResponseMessage: resources.ErrUnmarshalData}
	}
	err = validate.Struct(form)

	if err != nil {
		if valid, ok := err.(validator.ValidationErrors); ok {
			for _, err := range valid {
				g.Logger.LogError(err)
			}
			return services.ErrorResponse{Code: http.StatusBadRequest,
				ResponseMessage: resources.ClientError}
		}
		g.Logger.LogError(err)
		return services.ErrorResponse{Code: http.StatusBadRequest, ResponseMessage: resources.InternalServerError}
	}
	return nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// GetUint returns input as an uint or the default value while it's present and input is blank
func (g *Gin) GetUint(key string, def ...uint) (uint, error) {
	strv := g.C.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 32)
	return uint(u64), err
}

func (g *Gin) DefaultQueryUint(key string, defaultValue uint) uint {
	value, _ := g.GetUint(key)
	if value == 0 {
		value = defaultValue
	}
	return value
}

func (g *Gin) GetID(key string) (uint64, error) {
	idStr := g.C.Param(key)
	if idStr == "" || idStr == "0" {
		return 0, errors.New(resources.ClientError)
	}
	return strconv.ParseUint(idStr, 64, 10)
}
