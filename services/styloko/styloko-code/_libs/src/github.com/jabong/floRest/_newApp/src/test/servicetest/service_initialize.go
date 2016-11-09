package servicetest

import (
	"common/appconstant"
	"github.com/jabong/floRest/src/service"
	"hello"
)

func InitializeTestService() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
	service.RegisterApi(new(hello.HelloApi))
	initTestLogger()

	initTestConfig()

	service.InitVersionManager()

	initialiseTestWebServer()

	service.InitHealthCheck()

}

func PurgeTestService() {

}
