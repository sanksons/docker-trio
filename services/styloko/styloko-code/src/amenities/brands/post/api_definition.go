package post

import (
	"amenities/brands/common"
	"common/constants"
	"common/pool"
	"fmt"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

//Base struct for Brand Get API
type BrandAPI struct {
}

//Pool object for job dispatcher
var brandCreatePool pool.Safe

//Function to return Version manager instance
func (a *BrandAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.BRAND_API,
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

//Function to return the Orchestrator
func (sc *BrandAPI) GetOrchestrator() orchestrator.Orchestrator {

	logger.Info("brand Create API pipeline begins")

	brandCreateWorkflow := new(orchestrator.WorkFlowDefinition)
	brandCreateWorkflow.Create()

	brandMongoNode := new(Insert)
	brandMongoNode.SetID("brand create node")

	err := brandCreateWorkflow.AddExecutionNode(brandMongoNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("brand create failure Node")
	err = brandCreateWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	decisionNode := new(Validate)
	decisionNode.SetID("brand create validate")

	//Describing connection between Nodes
	err = brandCreateWorkflow.AddDecisionNode(decisionNode, brandMongoNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	brandCreateWorkflow.SetStartNode(decisionNode)

	brandCreateOrchestrator := new(orchestrator.Orchestrator)
	brandCreateOrchestrator.Create(brandCreateWorkflow)
	logger.Info(brandCreateOrchestrator.String())

	return *brandCreateOrchestrator
}

//Function that returns HealthCheckInterface
func (sc *BrandAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(BrandCreateHealthCheck)
}

//Function to initialise API
func (sc *BrandAPI) Init() {
	logger.Info("Starting worker pool for Brand POST API")
	brandCreatePool = pool.NewWorkerSafe(constants.BRAND_CREATE, common.BRAND_CREATE_POOL_SIZE, common.BRAND_CREATE_QUEUE_SIZE, common.BRAND_CREATE_RETRY_COUNT, common.BRAND_CREATE_WAIT_TIME)
	brandCreatePool.StartWorkers(common.BrandCreateWorker)
}
