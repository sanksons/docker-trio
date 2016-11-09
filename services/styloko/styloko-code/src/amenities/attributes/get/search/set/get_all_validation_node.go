package set

import (
	"common/utils"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type QueryParamDecision struct {
	id string
}

func (d *QueryParamDecision) SetID(id string) {
	d.id = id
}

func (d QueryParamDecision) GetID() (id string, err error) {
	return d.id, nil
}

func (d QueryParamDecision) Name() string {
	return "QueryParamDecision node for GET"
}

func (d QueryParamDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PATH_DECISION)
	defer func() {
		logger.EndProfile(profiler, PATH_DECISION)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, CUSTOM_ATTRIBUTESETS_GET_ALL)
	global, ok := utils.GetQueryParams(io, QUERY_PARAM)
	if !ok {
		return false, nil
	}
	io.IOData.Set(GET_SEARCH, global)
	return true, nil
}
