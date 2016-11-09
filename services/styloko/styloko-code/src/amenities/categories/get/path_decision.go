package get

import (
	"amenities/categories/common"
	"common/appconstant"
	"common/utils"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// PathDecision -> struct for node based data
type PathDecision struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *PathDecision) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs PathDecision) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs PathDecision) Name() string {
	return "PathDecision"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (cs PathDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_PATH_DECISION)
	defer func() {
		logger.EndProfile(profiler, common.GET_PATH_DECISION)
	}()
	pathParams, err := utils.GetPathParams(io)
	if err != nil {
		return false, nil
	}
	if len(pathParams) == 1 {
		io.IOData.Set(common.CATEGORY_PATH_PARAMS, pathParams)
		return true, nil
	}
	logger.Debug(pathParams)
	return false, &florest_constants.AppError{Code: appconstant.FunctionalityNotImplementedErrorCode, Message: "Invalid Path.", DeveloperMessage: "Path length can be 0 or 1 only."}
}
