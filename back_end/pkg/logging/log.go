package logging

import (
	"github.com/beego/beego/v2/core/logs"
	"log"
)

type logger struct {
	record *logs.BeeLogger
}

var LogManager logger

func Setup() {
	LogManager.record = logs.NewLogger(10000)
	err := LogManager.record.SetLogger(logs.AdapterFile, `{"filename":"./logs/log.log","color":true,"perm":"0644"}`)
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}
	err = LogManager.record.SetLogger(logs.AdapterConsole)
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}
	LogManager.record.EnableFuncCallDepth(true)
	LogManager.record.Async(10000)

}

func (l logger) LogError(err error) {
	LogManager.record.Error(err.Error())

}

func (l logger) LogInfo(message string) {
	LogManager.record.Info(message)
}

func (l logger) LogDebugMessage(message string) {
	LogManager.record.Info(message)
}
