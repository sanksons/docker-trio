package get

import (
	"amenities/categories/common"
	"common/utils"

	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// QueryDecision -> struct for node based data
type QueryDecision struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *QueryDecision) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs QueryDecision) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs QueryDecision) Name() string {
	return "QueryDecision"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (cs QueryDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_QUERY_DECISION)
	defer func() {
		logger.EndProfile(profiler, common.GET_QUERY_DECISION)
	}()
	qParams, ok := utils.GetQueryParams(io, "status")
	if !ok {
		return false, nil
	}
	logger.Debug(qParams)
	io.IOData.Set(common.CATEGORY_QUERY_PARAMS, qParams)
	return true, nil
}
