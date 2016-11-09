package service

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/monitor"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/utils/orchestratorhelper"
)

type BusinessLogicExecutor struct {
	id string
}

func (n BusinessLogicExecutor) Name() string {
	return "Business Logic Executor"
}

func (n *BusinessLogicExecutor) SetID(id string) {
	n.id = id
}

func (n BusinessLogicExecutor) GetID() (id string, err error) {
	return n.id, nil
}

func (n BusinessLogicExecutor) Execute(data workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	rc, _ := data.ExecContext.Get(constants.REQUEST_CONTEXT)
	logger.Info(fmt.Sprintln("entered ", n.Name()), rc)

	resource, version, action, orchBucket := getServiceVersion(data)

	logger.Info(fmt.Sprintf("Resource: %s, Version: %s, Action: %s, BucketId: %s", resource,
		version, action, orchBucket), rc)

	orchestrator, oerr := orchestratorhelper.GetOrchestrator(resource, version,
		action, orchBucket)
	if oerr != nil {
		data.IOData.Set(constants.APPERROR, oerr)
		return data, nil
	}

	dderr := monitor.GetInstance().Count(
		fmt.Sprintf("%v_%v_%v_%v_%vrequest_count", action, version, resource, orchBucket, getCustomMetricPrefix(data)), 1, nil, 1)
	if dderr != nil {
		logger.Error(fmt.Sprintln("Monitoring Error ", dderr.Error()), rc)
	}

	profiler := logger.NewProfiler()
	nameOforchestratorExecuted := fmt.Sprintf("%v_%v_%v_%v_execution", action, version,
		resource, orchBucket)

	logger.StartProfile(profiler, nameOforchestratorExecuted)
	res, err := orchestratorhelper.ExecuteOrchestrator(&data, orchestrator)

	customProfilerMetric := fmt.Sprintf("%v_%v_%v_%v_%vexecution", action, version,
		resource, orchBucket, getCustomMetricPrefix(data))

	t := logger.EndProfileCustomMetric(profiler, nameOforchestratorExecuted, customProfilerMetric)

	threshold := ResourceToThreshold[resource]
	if threshold != 0 && t != 0 && t >= threshold {
		logger.Error(fmt.Sprintf("%s_THRESHOLD_REACHED : Response time is more than threshold : time taken(MS): %v threshold value(MS): %v", resource, t, threshold), rc)
	}

	data.IOData.Set(constants.RESPONSE_DATA, res)

	if err != nil {
		data.IOData.Set(constants.APPERROR, err)
		return data, nil
	}

	logger.Info(fmt.Sprintln("exiting ", n.Name()), rc)

	return data, nil
}
