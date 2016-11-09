package get

import (
	"amenities/standardsize/common"
	"fmt"

	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

// GetStandardSize base struct
type GetStandardSizeApi struct {
}

// GetVersion returns version number
func (a *GetStandardSizeApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.STANDARDSIZE_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator returns orchestrator
func (a *GetStandardSizeApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Standard Size Get Creation begin")

	ssGetOrchestrator := new(orchestrator.Orchestrator)
	ssGetWorkflow := new(orchestrator.WorkFlowDefinition)
	ssGetWorkflow.Create()

	//Creation of the nodes in the workflow definition
	ssGetNode := new(StandardSizeGet)
	ssGetNode.SetID("StandardSize Get node")
	err := ssGetWorkflow.AddExecutionNode(ssGetNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	ssErrorNode := new(StandardSizeErrorResponse)
	ssErrorNode.SetID("Standard Size Get Error Response Node")
	err = ssGetWorkflow.AddExecutionNode(ssErrorNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	validationNode := new(StandardSizeGetValidation)
	validationNode.SetID("Standard Size Get Validation Node")

	//Set start node for the search workflow
	ssGetWorkflow.AddDecisionNode(validationNode,
		ssGetNode, ssErrorNode)

	//Set start node for the search workflow
	ssGetWorkflow.SetStartNode(validationNode)

	//Assign the workflow definition to the Orchestrator
	ssGetOrchestrator.Create(ssGetWorkflow)

	logger.Info(ssGetOrchestrator.String())
	logger.Info("Standard Size Get Pipeline Created")
	return *ssGetOrchestrator
}

// Init initializes the API
func (a *GetStandardSizeApi) Init() {

}

// GetHealthCheck -> Returns HealthCheckInterface
func (a *GetStandardSizeApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(StandardSizeGetHealthCheck)
}
