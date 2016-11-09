package attributes

import (
	"common/appconstant"
	"common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
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
	return "Decision node Attribute updattion for path params"
}

func (d PathParamDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PATH_DECISION)
	defer func() {
		logger.EndProfile(profiler, PATH_DECISION)
	}()
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
	if params.EndPoint != nil && *params.EndPoint == "options" {
		io.IOData.Set(PATH_PARAMETERS, params)
		return true, nil
	}
	return false, nil
}

func (d PathParamDecision) GetPathParmas(pathParams []string) Parameters {
	var params Parameters
	l := len(pathParams)
	switch l {
	case 1:
		params.AttrId, _ = strconv.Atoi(pathParams[0])
	case 3:
		params.AttrId, _ = strconv.Atoi(pathParams[0])
		params.OptionId, _ = strconv.Atoi(pathParams[2])
		params.EndPoint = &pathParams[1]
	}
	params.Count = l
	return params
}
