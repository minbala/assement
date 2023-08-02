package main

import (
	"fmt"
	"log"
	"test_assessment/middleware"
	"test_assessment/pkg/jwt"
	"test_assessment/pkg/password"

	"github.com/gin-gonic/gin"
	routers "test_assessment/controller"
	"test_assessment/controller/app"
	mysql "test_assessment/persistence/mysql"
	pg "test_assessment/persistence/mysql"
	"test_assessment/pkg/logging"
	"test_assessment/pkg/setting"
)

func init() {
	setting.Setup("./conf/app.ini")
	logging.Setup()
	mysql.Setup()
	password.SetUp()
	jwt.SetUp()
	middleware.SetUp(jwt.Manager, pg.NewSessionRepository())
	app.SetUp()
}

//	@title			AdminPanel Service API
//	@version		1.0
//	@description	AdminPanel Service API in Go using Gin framework
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	minbala33@gmail.com
//	@contact.email	minbala33@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:8080

//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	err := routersInit.Run(":" + setting.AppSetting.APIPORT)
	log.Println(" start http server listening " + setting.AppSetting.APIPORT)
	if err != nil {
		log.Fatal(fmt.Errorf("Server err: " + err.Error()))
	}
}
