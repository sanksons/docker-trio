package set

import (
	"common/appconstant"
	"common/utils"
	"fmt"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type PathParamDecision struct {
	id string
}

func (d *PathParamDecision) SetID(id string) {
	d.id = id
}

func (d PathParamDecision) GetID() (id string, err error) {
	return d.id, nil
}

func (d PathParamDecision) Name() string {
	return "Decision node for path params"
}

func (d PathParamDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PATH_DECISION)
	defer func() {
		logger.EndProfile(profiler, PATH_DECISION)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, CUSTOM_ATTRIBUTESETS_GET_ONE)
	pathParams, err := utils.GetPathParams(io)
	if len(pathParams) == 0 {
		return false, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting path parameters: %s", err.Error()))
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Path Parameters",
			DeveloperMessage: "Error while getting path parameters"}
	}
	params := d.GetPathParmas(pathParams)
	io.IOData.Set(PATH_PARAMETERS, params)
	if params.SetId != nil {
		return true, nil
	}
	return false, nil
}

func (d PathParamDecision) GetPathParmas(pathParams []string) Parameters {
	var params Parameters
	l := len(pathParams)
	switch l {
	case 1:
		params.SetId = &pathParams[0]
		break
	case 2:
		params.SetId = &pathParams[0]
		params.EndPoint = &pathParams[1]
		break
	case 3:
		params.SetId = &pathParams[0]
		params.EndPoint = &pathParams[1]
		params.AttrIdName = &pathParams[2]
		break
	}
	params.Count = l
	return params
}
