package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"test_assessment/controller/app"
	"test_assessment/controller/schema"
	"test_assessment/dependency_manager"
	domainModel "test_assessment/domain/model"
	"test_assessment/middleware"
	"test_assessment/pkg/resources"
	services "test_assessment/service"
	serviceModel "test_assessment/service/models"
)

type AdminAPIV1 struct{}

func NewAdminV1API() *AdminAPIV1 {
	return &AdminAPIV1{}
}

// UserLogoutV1   godoc
//
//	@Summary		User Logout
//	@Description	User Logout Version 1
//	@Tags			Authentication
//	@Accept			application/json
//	@Produce		application/json
//	@Success		204	{object}	nil
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/logout [delete]
//
//	@Security		Bearer
func (AdminAPIV1) UserLogoutV1(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	adminService := dependency_manager.GetAdminService()
	err := adminService.Logout(appG.C.GetHeader("Authorization"), userInformation.UserId)
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "delete",
				RequestUrl: "/v1/logout", ServiceType: "admin panel", Status: "success"}, err)
	}()
	appG.Response(http.StatusNoContent, nil)
}

// CreateUser go doc
//
//	@Summary		create user account
//	@Description	create user account
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body	schema.CreateUserInput	true	"Request body in JSON format"
//	@Security		Bearer
//	@Success		201	{object}	schema.CreateUserResponse
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user [post]
func (AdminAPIV1) CreateUser(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	var requestData schema.CreateUserInput
	err := appG.BindAndValid(&requestData)
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "post",
				RequestUrl: "/v1/user", ServiceType: "admin panel", Status: "success"}, err)
	}()
	if !app.IsValidEmail(requestData.Email) {
		appG.ResponseErrorToClient(services.ErrorResponse{Code: 400, ResponseMessage: resources.ClientError,
			ErrorString: "invalid email"})
		return
	}
	adminService := dependency_manager.GetAdminService()
	user, err := adminService.CreateUser(serviceModel.CreateUserInput{Name: requestData.Name, Email: requestData.Email,
		Password: requestData.Password, UserRole: requestData.UserRole})
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}

	appG.Response(http.StatusCreated, schema.CreateUserResponse{Id: user.Id, Name: user.Name, Email: user.Email,
		UserRole: user.UserRole, CreatedAt: user.CreatedAt})
}

// GetUsers go doc
//
//	@Summary		get users
//	@Description	get users
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			name		query	string	false	"bala"
//	@Param			user_role	query	string	false	"user"
//	@Param			limit		query	uint	false	"20"
//	@Param			offset		query	uint	false	"0"
//	@Security		Bearer
//	@Success		200	{object}	schema.GetUserResponse
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user [get]
func (AdminAPIV1) GetUsers(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	name, _ := c.GetQuery("name")
	userRole, _ := c.GetQuery("user_role")
	adminService := dependency_manager.GetAdminService()
	user, err := adminService.GetUsers(name, userRole, appG.DefaultQueryUint("limit", 20),
		appG.DefaultQueryUint("offset", 0))
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "get",
				RequestUrl: "/v1/user", ServiceType: "admin panel", Status: "success"}, err)
	}()
	response := schema.GetUserResponse{Total: user.Total}
	for _, user := range user.Users {
		response.Users = append(response.Users, schema.User{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			UserRole:  user.UserRole,
			CreatedAt: user.CreatedAt,
		})
	}
	appG.Response(http.StatusOK, response)
}

// GetUserById go doc
//
//	@Summary		get  user
//	@Description	get  user
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			userId	path	uint	true	"2"
//	@Security		Bearer
//	@Success		200	{object}	schema.User
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user/{userId} [get]
func (AdminAPIV1) GetUserById(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	userId, err := strconv.ParseUint(appG.C.Param("userId"), 10, 64)
	if err != nil {
		appG.ResponseErrorToClient(services.ErrorResponse{Code: 400, ResponseMessage: resources.ClientError,
			ErrorString: err.Error()})
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "get",
				RequestUrl: "/v1/user/{userId}", ServiceType: "admin panel", Status: "success"}, err)
	}()
	adminService := dependency_manager.GetAdminService()
	user, err := adminService.GetUserById(uint(userId))
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	appG.Response(http.StatusOK, schema.User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		UserRole:  user.UserRole,
		CreatedAt: user.CreatedAt,
	})
}

// DeleteUserById go doc
//
//	@Summary		delete user
//	@Description	delete user
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			userId	path	uint	true	"2"
//	@Security		Bearer
//	@Success		204	{object}	nil
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user/{userId} [delete]
func (AdminAPIV1) DeleteUserById(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	userId, err := strconv.ParseUint(appG.C.Param("userId"), 10, 64)
	if err != nil {
		appG.ResponseErrorToClient(services.ErrorResponse{Code: 400, ResponseMessage: resources.ClientError,
			ErrorString: err.Error()})
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "delete",
				RequestUrl: "/v1/user/{userId}", ServiceType: "admin panel", Status: "success"}, err)
	}()
	adminService := dependency_manager.GetAdminService()
	err = adminService.DeleteUser(uint(userId))
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	appG.Response(http.StatusNoContent, nil)
}

// UpdateUser go doc
//
//	@Summary		update user
//	@Description	update user
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body	schema.UpdateUserInput	true	"Request body in JSON format"
//	@Security		Bearer
//	@Success		200	{object}	nil
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user [put]
func (AdminAPIV1) UpdateUser(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	var requestData schema.UpdateUserInput
	err := appG.BindAndValid(&requestData)
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "put",
				RequestUrl: "/v1/user", ServiceType: "admin panel", Status: "success"}, err)
	}()
	if !app.IsValidEmail(requestData.Email) {
		appG.ResponseErrorToClient(services.ErrorResponse{Code: 400, ResponseMessage: resources.ClientError,
			ErrorString: "invalid email"})
		return
	}
	adminService := dependency_manager.GetAdminService()
	err = adminService.UpdateUser(serviceModel.UpdateUserInput{Name: requestData.Name, Email: requestData.Email,
		Password: requestData.Password, UserRole: requestData.UserRole, UserId: requestData.UserId})
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}

	appG.Response(http.StatusOK, nil)
}

// GetUserLogs go doc
//
//	@Summary		get  user logs
//	@Description	get  user logs
//	@Tags			Admin
//	@Accept			application/json
//	@Produce		application/json
//	@Param			user_id	query	uint	false	"2"
//	@Param			limit	query	uint	false	"20"
//	@Param			offset	query	uint	false	"0"
//	@Security		Bearer
//	@Success		200	{object}	schema.GetUserLogsResponse
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user-logs [get]
func (AdminAPIV1) GetUserLogs(c *gin.Context) {
	appG := dependency_manager.GetAppManager(c)
	userInformation := appG.C.MustGet("UserData").(middleware.UserData)
	userId, _ := appG.GetUint("user_id")
	adminService := dependency_manager.GetAdminService()
	user, err := adminService.GetUserLog(userId, appG.DefaultQueryUint("limit", 20),
		appG.DefaultQueryUint("offset", 0))
	if err != nil {
		appG.ResponseErrorToClient(err)
		return
	}
	defer func() {
		go appG.RecordUserAction(
			&domainModel.CreateUserLogInput{UserId: userInformation.UserId, Method: "get",
				RequestUrl: "/v1/user-logs", ServiceType: "admin panel", Status: "success"}, err)
	}()
	response := schema.GetUserLogsResponse{Total: user.Total}
	for _, userLog := range user.UserLogs {
		response.UserLogs = append(response.UserLogs, schema.UserLog{
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
	appG.Response(http.StatusOK, response)
}
