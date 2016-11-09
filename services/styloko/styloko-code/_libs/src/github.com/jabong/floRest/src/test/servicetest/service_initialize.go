package servicetest

import (
	"github.com/jabong/floRest/src/common/env"
	"github.com/jabong/floRest/src/service"
)

func InitializeTestService() {

	//apiservice.Register()

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
