package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_assessment/controller/app"
	"test_assessment/controller/schema"
	"test_assessment/dependency_manager"
	domainModel "test_assessment/domain/model"
	"test_assessment/pkg/resources"
	services "test_assessment/service"
	serviceModel "test_assessment/service/models"
)

// PublicAPIV1   represents the public routes.
type PublicAPIV1 struct{}

// NewPublicAPI  gets a new PublicAPI instance.
func NewPublicAPI() *PublicAPIV1 {
	return &PublicAPIV1{}
}

// UserLoginV1   godoc
//
//	@Summary		User Login
//	@Description	User Login Version 1
//	@Tags			Authentication
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body		schema.LoginInput	true	"Request body in JSON format"
//	@Success		200			{object}	schema.LoginResponse
//	@Failure		400			{object}	app.ResponseMessage
//	@Failure		500			{object}	app.ResponseMessage
//	@Router			/v1/login [post]
func (PublicAPIV1) UserLoginV1(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	var requestData schema.LoginInput
	err := appG.BindAndValid(&requestData)
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	if !app.IsValidEmail(requestData.Email) {
		appG.ResponseErrorToClient(services.ErrorResponse{Code: 400, ResponseMessage: resources.ClientError})
		return
	}
	publicService := dependency_manager.GetPublicService()
	accessToken, userId, err := publicService.Login(serviceModel.LoginInput{Email: requestData.Email, Password: requestData.Password})
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}

	go appG.RecordUserAction(
		&domainModel.CreateUserLogInput{UserId: userId, Method: "post",
			RequestUrl: "/v1/login", ServiceType: "admin panel", Status: "success"}, err)

	appG.Response(http.StatusCreated, schema.LoginResponse{
		AccessToken: accessToken,
	})
}
