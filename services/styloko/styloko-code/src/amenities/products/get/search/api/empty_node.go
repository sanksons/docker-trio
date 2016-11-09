package api

import (
	proUtil "amenities/products/common"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type EmptyNode struct {
	id string
}

func (cs *EmptyNode) SetID(id string) {
	cs.id = id
}

func (cs EmptyNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs EmptyNode) Name() string {
	return "EmptyNode"
}

func (cs EmptyNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Info("Enter " + cs.id)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, cs.id)

	logger.Info("Exit " + cs.id)
	return io, nil
}
