package simplifier

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type ReturnError struct {
}

func (a *ReturnError) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: ERRORS,
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (a *ReturnError) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("CSV Error Pipeline Creation begin")

	csvErrorOrchestrator := new(orchestrator.Orchestrator)
	csvErrorWorkflow := new(orchestrator.WorkFlowDefinition)
	csvErrorWorkflow.Create()

	//Creation of the nodes in the workflow definition
	errorNode := new(csvError)
	errorNode.SetID("csv error node 1")
	eerr := csvErrorWorkflow.AddExecutionNode(errorNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	csvErrorWorkflow.SetStartNode(errorNode)

	//Assign the workflow definition to the Orchestrator
	csvErrorOrchestrator.Create(csvErrorWorkflow)

	logger.Info(csvErrorOrchestrator.String())
	return *csvErrorOrchestrator
}

func (a *ReturnError) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CsvErrorHeathCheck)
}

func (a *ReturnError) Init() {
	//api initialization should come here
}
