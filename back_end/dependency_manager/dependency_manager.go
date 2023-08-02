package dependency_manager

import (
	"github.com/gin-gonic/gin"
	"test_assessment/controller/app"
	pg "test_assessment/persistence/mysql"
	"test_assessment/pkg/jwt"
	"test_assessment/pkg/logging"
	"test_assessment/pkg/password"
	services "test_assessment/service"
)

var (
	publicService *services.PublicService
	adminService  *services.AdminService
)

func GetPublicService() *services.PublicService {

	if publicService == nil {
		publicService = services.NewPublicService(logging.LogManager, password.Manager, jwt.Manager,
			pg.NewSessionRepository(), pg.NewUserRepository())
	}

	return publicService
}

func GetAdminService() *services.AdminService {

	if adminService == nil {
		adminService = services.NewAdminService(logging.LogManager, password.Manager, jwt.Manager,
			pg.NewSessionRepository(), pg.NewUserRepository(), pg.NewUserLogRepository())
	}

	return adminService
}

func GetAppManager(c *gin.Context) app.Gin {
	return app.Gin{C: c, Logger: logging.LogManager, UserLogManager: pg.NewUserLogRepository()}
}
