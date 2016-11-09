package main

import (
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/examples"
	"github.com/jabong/floRest/src/service"
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
	service.RegisterApi(new(examples.HelloApi))
}

func registerConfig() {
	service.RegisterConfig(new(config.AppConfig))
}

func registerErrors() {
}

func registerBuckets() {
}
