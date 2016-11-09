package main

import (
	"apiservices"
	"log"

	"github.com/jabong/floRest/src/service"
)

//main is the entry point of the florest web service
func main() {
	log.Println("APPLICATION BEGIN")
	webServer := new(service.Webserver)
	apiservices.Register()
	webServer.Start()
}
