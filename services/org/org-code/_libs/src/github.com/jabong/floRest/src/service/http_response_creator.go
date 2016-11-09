package service

import (
	"fmt"

	"encoding/json"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/monitor"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type HttpResponseCreator struct {
	id string
}

func (n HttpResponseCreator) Name() string {
	return "Http Response Creator"
}

func (n *HttpResponseCreator) SetID(id string) {
	n.id = id
}

func (n HttpResponseCreator) GetID() (id string, err error) {
	return n.id, nil
}

func (n *HttpResponseCreator) Execute(data workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	rc, _ := data.ExecContext.Get(constants.REQUEST_CONTEXT)
	logger.Info(fmt.Sprintln("entered ", n.Name()), rc)

	resStatus, _ := data.IOData.Get(constants.APPERROR)
	resData, _ := data.IOData.Get(constants.RESPONSE_DATA)

	appError := new(constants.AppErrors)

	if resStatus != nil {
		if v, ok := resStatus.(*constants.AppError); ok {
			if v != nil { //if v is of type *AppError and is not nil
				appError.Errors = []constants.AppError{*v}
			}
		} else if v, ok := resStatus.(*constants.AppErrors); ok {
			if v != nil { //v is of type *AppErrors and is not nil
				appError = v
			}
		} else {
			appError.Errors = []constants.AppError{constants.AppError{Code: constants.InvalidErrorCode,
				Message: "Invalid App error"}}
		}
	}
	status := constants.GetAppHttpError(*appError)
	debugData, _ := data.ExecContext.GetDebugMsg()

	resource, version, action, orchBucket := getServiceVersion(data)

	serviceStatusKey := fmt.Sprintf("%v_%v_%v_%v_%vHttp_%v", action,
		version, resource, orchBucket, getCustomMetricPrefix(data), status.HttpStatusCode)

	if status.HttpStatusCode != constants.HttpStatusSuccessCode {
		logger.Error(fmt.Sprintf("%s_%v Application Errors : %v", resource, status.HttpStatusCode, appError), rc)
	}

	dderr := monitor.GetInstance().Count(serviceStatusKey, 1, nil, 1)
	if dderr != nil {
		logger.Error(fmt.Sprintln("Monitoring Error ", dderr.Error()), rc)
	}

	var appDebugData []utilhttp.Debug
	for _, d := range debugData {
		if v, ok := d.(workflow.WorkflowDebugDataInMemory); ok {
			appDebugData = append(appDebugData, utilhttp.Debug{Key: v.Key, Value: v.Value})
		}
	}

	m, _ := data.IOData.Get(constants.RESPONSE_META_DATA)
	md, _ := m.(*utilhttp.ResponseMetaData)
	appResponse := utilhttp.Response{Status: *status, Data: resData, DebugData: appDebugData, MetaData: md}
	data.IOData.Set(constants.RESPONSE, appResponse)
	jsonBody, err := json.Marshal(appResponse)
	if err != nil {
		return data, err
	}
	r, _ := data.IOData.Get(constants.API_RESPONSE)
	apiResponse, _ := r.(utilhttp.APIResponse)
	apiResponse.HttpStatus = appResponse.Status.HttpStatusCode
	apiResponse.Body = jsonBody
	data.IOData.Set(constants.API_RESPONSE, apiResponse)

	logger.Info(fmt.Sprintln("exiting ", n.Name()), rc)

	return data, nil
}
