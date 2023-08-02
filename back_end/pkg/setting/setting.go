package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	APIHOST                string
	APIPORT                string
	JWTSECRET              string
	LogFilePath            string
	LogFile                string
	AccessTokenExpiredTime int
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	DATABASEURL string
}

var DatabaseSetting = &Database{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup(source string) {
	var err error
	cfg, err = ini.Load(source)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("database", DatabaseSetting)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
