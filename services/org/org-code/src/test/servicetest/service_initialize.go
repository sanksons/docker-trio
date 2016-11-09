package servicetest

import (
	"apiservice"
	"common/appconstant"
	"github.com/jabong/floRest/src/common/env"
	"github.com/jabong/floRest/src/service"
	"hello"
)

func InitializeTestService() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
	service.RegisterApi(new(hello.HelloApi))
	apiservice.Register()

	env.GetOsEnviron()

	initTestConfig()

	initTestLogger()

	service.InitDBAdapterManager()

	service.InitVersionManager()

	service.InitCustomApiInit()

	service.InitApis()

	service.InitHealthCheck()

	initialiseTestWebServer()

}

func PurgeTestService() {

}
