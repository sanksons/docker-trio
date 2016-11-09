package commissions

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"sellers/common"
)

type CommissionsApi struct {
}

func (ca *CommissionsApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.COMMISSSIONS,
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (ca *CommissionsApi) GetOrchestrator() orchestrator.Orchestrator {

	logger.Info("Commissions API pipeline begins")

	getCommissionsWorkflow := new(orchestrator.WorkFlowDefinition)
	getCommissionsWorkflow.Create()

	getCommissionsNode := new(GetCommissions)
	getCommissionsNode.SetID("Seller Get Commissions Node")
	err := getCommissionsWorkflow.AddExecutionNode(getCommissionsNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Seller Get Commission Failure Node")
	err = getCommissionsWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	decisionNode := new(Validate)
	decisionNode.SetID("Seller Get Commission Validation Node")
	err = getCommissionsWorkflow.AddDecisionNode(decisionNode, getCommissionsNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	getCommissionsWorkflow.SetStartNode(decisionNode)

	getCommissionsOrchestrator := new(orchestrator.Orchestrator)
	getCommissionsOrchestrator.Create(getCommissionsWorkflow)
	logger.Info(getCommissionsOrchestrator.String())
	return *getCommissionsOrchestrator
}

func (ca *CommissionsApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(getCommissionsHealthCheck)
}

func (ca *CommissionsApi) Init() {
	//api initialization should come here
}
