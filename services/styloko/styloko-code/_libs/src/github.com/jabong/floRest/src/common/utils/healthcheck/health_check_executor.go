package healthcheck

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type HealthCheckExecutor struct {
	id string
}

func (n HealthCheckExecutor) Name() string {
	return "Health Check Executor"
}

func (n *HealthCheckExecutor) SetID(id string) {
	n.id = id
}

func (n HealthCheckExecutor) GetID() (id string, err error) {
	return n.id, nil
}

func (n HealthCheckExecutor) Execute(data workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	rc, _ := data.ExecContext.Get(constants.REQUEST_CONTEXT)
	logger.Info(fmt.Sprintln("entered ", n.Name()), rc)

	if healthCheckApiList == nil {
		return data, &constants.AppError{Code: constants.ResourceErrorCode, Message: "Health Chech Api not Initialized"}
	}

	var res = make(map[string]interface{})

	for _, apiResource := range healthCheckApiList {
		res[apiResource.GetName()] = apiResource.GetHealth()
	}

	data.IOData.Set(constants.RESULT, res)

	logger.Info(fmt.Sprintln("exiting ", n.Name()), rc)
	return data, nil

}
