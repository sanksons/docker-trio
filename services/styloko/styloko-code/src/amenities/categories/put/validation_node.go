package put

import (
	"amenities/categories/common"
	"common/appconstant"
	"common/utils"
	"strconv"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// ValidationNode -> struct for node based data
type ValidationNode struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *ValidationNode) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs ValidationNode) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs ValidationNode) Name() string {
	return "ValidationNode"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (cs ValidationNode) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.PUT_VALIDATION_NODE)
	defer func() {
		logger.EndProfile(profiler, common.PUT_VALIDATION_NODE)
	}()
	pathParams, err := utils.GetPathParams(io)
	logger.Debug("Path Parameters")
	logger.Debug(pathParams)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.InconsistantDataStateErrorCode, Message: "Could not fetch path params", DeveloperMessage: err.Error()}
	}

	if len(pathParams) == 1 {
		id, err := strconv.Atoi(pathParams[0])
		if err != nil {
			return false, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode, Message: "ID is not an integer", DeveloperMessage: err.Error()}
		}
		io.IOData.Set(common.CATEGORY_ID, id)
		return true, nil
	}
	return false, &florest_constants.AppError{Code: appconstant.FunctionalityNotImplementedErrorCode, Message: "Invalid path provided.", DeveloperMessage: "Path length can be only 1."}
}
