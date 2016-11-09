package put

import (
	"amenities/brands/common"
	// "common/appconfig"
	"common/appconfig"
	"common/constants"
	"common/pool"
	"fmt"

	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

var cacheObj cache.CacheInterface

//Base struct for Brand Get API
type BrandAPI struct {
}

//Pool object for job dispatcher
var brandUpdatePool pool.Safe

//Function to return Version manager instance
func (bu *BrandAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.BRAND_API,
		Version:  "V1",
		Action:   "PUT",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

//Function to return the Orchestrator
func (bu *BrandAPI) GetOrchestrator() orchestrator.Orchestrator {

	//Pictorial Description of Nodes

	//               Decision Node
	//	                  /\
	//  Brand Update Node   Update Failure Node

	logger.Info("Brand Update API pipeline begins")

	brandUpdateWorkflow := new(orchestrator.WorkFlowDefinition)
	brandUpdateWorkflow.Create()

	//Creation of the nodes in the workflow definition
	updateNode := new(Update)
	updateNode.SetID("Brand Update Node")
	err := brandUpdateWorkflow.AddExecutionNode(updateNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Update Failure Node")
	err = brandUpdateWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Describing connection between Nodes)
	decisionNode := new(Decision)
	decisionNode.SetID("Brand Update Decision Node")
	err = brandUpdateWorkflow.AddDecisionNode(decisionNode, updateNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	brandUpdateWorkflow.SetStartNode(decisionNode)
	brandUpdateOrchestrator := new(orchestrator.Orchestrator)
	brandUpdateOrchestrator.Create(brandUpdateWorkflow)
	logger.Info(brandUpdateOrchestrator.String())

	return *brandUpdateOrchestrator
}

//Function that returns HealthCheckInterface
func (bu *BrandAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(BrandsHealthCheck)
}

//Function to initialise API
func (bu *BrandAPI) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Starting worker pool for Brand PUT API")
	brandUpdatePool = pool.NewWorkerSafe(constants.BRAND_UPDATE, common.BRAND_UPDATE_POOL_SIZE, common.BRAND_UPDATE_QUEUE_SIZE, common.BRAND_UPDATE_RETRY_COUNT, common.BRAND_UPDATE_WAIT_TIME)
	brandUpdatePool.StartWorkers(common.BrandUpdateWorker)
}
