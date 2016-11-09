package main

import (
	"common/appconfig"
	"common/appconstant"
	"fmt"
	"github.com/jabong/floRest/src/service"
	"hello"
)

//main is the entry point of the florest web service
func main() {
	fmt.Println("APPLICATION BEGIN")
	webserver := new(service.Webserver)
	registerConfig()
	registerErrors()
	registerAllApis()
	webserver.Start()
}

func registerAllApis() {
	service.RegisterApi(new(hello.HelloApi))
}

func registerConfig() {
	service.RegisterConfig(new(appconfig.AppConfig))
}

func registerErrors() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
}
