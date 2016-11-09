package get

import (
	"common/utils"

	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// PathDecision -> struct for node based data
type PathDecision struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (pd *PathDecision) SetID(id string) {
	pd.id = id
}

// GetID -> returns current node ID to orchestrator
func (pd PathDecision) GetID() (id string, err error) {
	return pd.id, nil
}

// Name -> Returns node name to orchestrator
func (pd PathDecision) Name() string {
	return "PathDecision"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (pd PathDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	pathParams, _ := utils.GetPathParams(io)
	if len(pathParams) == 1 {
		io.IOData.Set(PATH, pathParams[0])
		return false, nil
	}
	logger.Debug(pathParams)
	return true, nil
}
