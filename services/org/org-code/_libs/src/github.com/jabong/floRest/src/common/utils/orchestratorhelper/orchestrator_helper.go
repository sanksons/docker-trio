package orchestratorhelper

import (
	"fmt"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

func GetOrchestrator(resource string,
	version string,
	action string,
	bucketId string) (*workflow.Orchestrator, error) {

	logFormatter := "GetOrchestrator ==== Resource: %v === Version: %v === Action: %v === BucketId: %v"
	logger.Info(fmt.Sprintf(logFormatter, resource, version, action, bucketId))
	orchestratorVersion, gerr := versionmanager.Get(resource, version, action, bucketId)
	if gerr != nil {
		return nil, &constants.AppError{Code: constants.InvalidRequestUri, Message: gerr.Error()}
	}

	orchestrator, ok := orchestratorVersion.(workflow.Orchestrator)
	if !ok {
		return nil, &constants.AppError{Code: constants.ResourceErrorCode,
			Message: "Error retrieving orchestrator"}
	}

	return &orchestrator, nil

}

func ExecuteOrchestrator(input *workflow.WorkFlowData,
	orchestrator *workflow.Orchestrator) (interface{}, error) {

	output := orchestrator.Start(input)
	res, _ := output.IOData.Get(constants.RESULT)

	orchestratorStates := output.GetWorkflowState()
	var orchestratorError error
	for _, err := range orchestratorStates {
		if v, ok := err.(error); ok {
			orchestratorError = v
		}
	}

	rc, _ := input.ExecContext.Get(constants.REQUEST_CONTEXT)
	logger.Info(fmt.Sprintf("%v", orchestratorError), rc)
	return res, orchestratorError
}
