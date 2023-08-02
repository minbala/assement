package main

import (
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"test_assessment/persistence/mysql/models"
	"test_assessment/pkg/password"
	"test_assessment/pkg/setting"
)

func init() {
	setting.Setup("../../../conf/app.ini")
	password.SetUp()
}

func main() {
	db, err := gorm.Open(mysql.Open(setting.DatabaseSetting.DATABASEURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to open Database")
	}
	driver, err := mysqlMigrate.WithInstance(sqlDB, &mysqlMigrate.Config{})
	if err != nil {
		log.Fatal("Failed to open Database")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:////home/user/go/src/assement/back_end/persistence/mysql/migration/migration",
		"mysql", driver)
	fakePassword, _ := password.Manager.Hash("minbala33")
	err = m.Steps(3) // or m.Step(2) if you want to explicitly set the number of migrations to run
	log.Println(err)
	user := &models.User{Password: fakePassword, Email: "minbala33@gmail.com", ID: 1, Name: "minbala", UserRole: "admin"}
	db.Create(user)

}
