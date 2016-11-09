package get

import (
	"amenities/brands/common"
	"common/utils"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//Struct for GetOne Decision Node
type GetOneDecision struct {
	id string
}

//Function to SetID for current node from orchestrator
func (d *GetOneDecision) SetID(id string) {
	d.id = id
}

//Function that returns current node ID to orchestrator
func (d GetOneDecision) GetID() (id string, err error) {
	return d.id, nil
}

//Function that returns node name to orchestrator
func (d GetOneDecision) Name() string {
	return "GetOneDecision node for GET"
}

//Function that returns bool to triger decision
func (d GetOneDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_GET)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_GET)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_BRAND_GET_ONE)

	//Reading path params
	pathParams, _ := utils.GetPathParams(io)
	if pathParams != nil {
		io.IOData.Set(common.BRAND_GET, pathParams[0])
		return true, nil
	}
	return false, nil
}
