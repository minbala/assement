package app

import (
	"github.com/gin-gonic/gin"
	domainModel "test_assessment/domain/model"
	"test_assessment/domain/repository"
	services "test_assessment/service"
)

type Gin struct {
	C              *gin.Context
	Logger         services.Logger
	UserLogManager repository.UserLogRepositoryInterface
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, data interface{}) {
	g.C.JSON(httpCode, data)
	return
}

func (g *Gin) RecordUserAction(userLog *domainModel.CreateUserLogInput, err error) {
	if err != nil {
		userLog.ErrorMessage = err.Error()
		userLog.Status = "fail"
	}
	err = g.UserLogManager.Create(*userLog)
	if err != nil {
		g.Logger.LogError(err)
	}
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func (g *Gin) ResponseErrorToClient(err error) {
	holder, ok := err.(services.ErrorResponse)
	if ok {
		g.Response(holder.Code, ResponseMessage{Message: holder.ResponseMessage})
		return
	}
}
