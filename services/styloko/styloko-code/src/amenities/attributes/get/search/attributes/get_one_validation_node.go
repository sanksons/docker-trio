package attributes

import (
	"common/utils"

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
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, CUSTOM_ATTRIBUTE_GET_ONE)
	pathParams, err := utils.GetPathParams(io)
	if err != nil {
		return false, nil
	}
	if len(pathParams) == 1 {
		io.IOData.Set(PATH_PARAMETERS, pathParams)
		return true, nil
	}
	return false, nil
}
