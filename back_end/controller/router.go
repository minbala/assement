package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"log"
	"os"
	controllers "test_assessment/controller/api/v1"
	docs "test_assessment/docs"
	"test_assessment/middleware"
	"test_assessment/pkg/setting"
)

func InitRouter() *gin.Engine {
	f, err := os.OpenFile(setting.AppSetting.LogFilePath+"/"+setting.AppSetting.LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {

		c.String(200, "pong")
	})

	publicAPI := controllers.NewPublicAPI()
	adminAPI := controllers.NewAdminV1API()

	{
		r.POST("/v1/login", publicAPI.UserLoginV1)
	}
	v1AdminRoutes := r.Group("/v1")
	v1AdminRoutes.Use(middleware.AdminAuthenticationManager.ValidateAdmin())
	{
		//!Version 1 Routes
		v1AdminRoutes.DELETE("/logout", adminAPI.UserLogoutV1)
		v1AdminRoutes.PUT("/user", adminAPI.UpdateUser)
		v1AdminRoutes.DELETE("/user/:userId", adminAPI.DeleteUserById)
		v1AdminRoutes.POST("/user", adminAPI.CreateUser)
		v1AdminRoutes.GET("/user", adminAPI.GetUsers)
		v1AdminRoutes.GET("/user/:userId", adminAPI.GetUserById)
		v1AdminRoutes.GET("/user-logs", adminAPI.GetUserLogs)

	}

	return r
}
