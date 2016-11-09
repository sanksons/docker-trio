package main

import (
	"apiservice"
	"fmt"
	"github.com/jabong/floRest/src/service"
	_ "net/http/pprof"
)

//main is the entry point of the florest web service
func main() {
	fmt.Println("APPLICATION BEGIN")
	webserver := new(service.Webserver)
	apiservice.Register()
	webserver.Start()
}
