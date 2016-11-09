package get

import (
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ValidateNode struct {
	id string
}

func (cs *ValidateNode) SetID(id string) {
	cs.id = id
}

func (cs ValidateNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ValidateNode) Name() string {
	return "ValidateNode"
}

func (cs ValidateNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter Validate node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_TAXCLASS_VALIDATE_NODE)
	defer logger.EndProfile(profiler, GET_TAXCLASS_VALIDATE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(DEBUG_KEY_NODE, "Validate")

	logger.Info("Exit Validate node")
	return io, nil
}
