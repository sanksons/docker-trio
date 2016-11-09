package service

import (
	"fmt"
	"strings"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type UriInterpreter struct {
	id string
}

func (uriIdentifier UriInterpreter) Name() string {
	return "URL Interpreter"
}

func (n *UriInterpreter) SetID(id string) {
	n.id = id
}

func (n UriInterpreter) GetID() (id string, err error) {
	return n.id, nil
}

func (u UriInterpreter) Execute(data workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	rc, _ := data.ExecContext.Get(constants.REQUEST_CONTEXT)

	logger.Info(fmt.Sprintln("Entered ", u.Name()), rc)

	resource, version, action := u.getResource(data)
	data.IOData.Set(constants.RESOURCE, resource)
	data.IOData.Set(constants.VERSION, version)
	data.IOData.Set(constants.ACTION, action)
	data.IOData.Set(constants.RESPONSE_META_DATA, utilhttp.NewResponseMetaData())

	logger.Info(fmt.Sprintln("exiting ", u.Name()), rc)
	return data, nil
}

func (u UriInterpreter) getResource(data workflow.WorkFlowData) (resource string,
	version string,
	action string) {

	rc, _ := data.ExecContext.Get(constants.REQUEST_CONTEXT)
	uridata, _ := data.IOData.Get(constants.URI)
	actiondata, _ := data.IOData.Get(constants.HTTPVERB)

	var uri string
	if v, ok := uridata.(string); ok {
		uri = v
	}
	logger.Info(fmt.Sprintln("uri is ", uri), rc)

	uriArr := strings.Split(uri[1:], "/")

	if len(uriArr) >= 2 &&
		uriArr[0] == config.GlobalAppConfig.AppName &&
		strings.ToUpper(uriArr[1]) == constants.HEALTHCHECKAPI {
		resource = constants.HEALTHCHECKAPI
		version = ""
	} else if len(uriArr) >= 3 && uriArr[0] == config.GlobalAppConfig.AppName {
		resource = strings.ToUpper(uriArr[2])
		version = strings.ToUpper(uriArr[1])
	} else {
		//Badly formed URI
		resource = ""
		version = ""
	}

	if v, ok := actiondata.(utilhttp.HTTPMethod); ok {
		action = string(v)
	}

	return resource, version, action
}
