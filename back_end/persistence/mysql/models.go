package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"test_assessment/pkg/setting"
)

var psqlDB *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	psqlDB, err = gorm.Open(mysql.Open(setting.DatabaseSetting.DATABASEURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	sqlDB, err := psqlDB.DB()
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(100)
	if err != nil {
		log.Fatal(err)
	}

}
