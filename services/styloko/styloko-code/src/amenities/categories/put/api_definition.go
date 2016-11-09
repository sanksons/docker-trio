package put

import (
	"amenities/categories/common"
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

// CategoryUpdatePool -> Pool object for job dispatcher
var categoryUpdatePool pool.Safe

// CategoryAPI -> Basic orchestrator struct
type CategoryAPI struct {
}

// GetVersion -> Version Manager boilerplate
func (a *CategoryAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.CATEGORY_API,
		Version:  "V1",
		Action:   "PUT",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Orchestrator Boilerplate
func (a *CategoryAPI) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Categories Update Creation begin")

	categoriesUpdateOrchestrator := new(orchestrator.Orchestrator)
	categoriesUpdateWorkflow := new(orchestrator.WorkFlowDefinition)
	categoriesUpdateWorkflow.Create()

	//Creation of the nodes in the workflow definition

	// Update Node is where mongo update happens
	categoriesUpdateNode := new(CategoriesUpdate)
	categoriesUpdateNode.SetID("Categories Update Node")
	eerr := categoriesUpdateWorkflow.AddExecutionNode(categoriesUpdateNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Category check node finds ID from mongo and gets the POST data if
	// ID is found in Mongo.
	// Struct validation happens here.
	categoriesCheckNode := new(CategoryCheckNode)
	categoriesCheckNode.SetID("Categories Check Node")
	eerr = categoriesUpdateWorkflow.AddExecutionNode(categoriesCheckNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Empty node is required for false edge of decision node.
	// As no process needs to run in that case.
	emptyNode := new(EmptyNode)
	emptyNode.SetID("Empty Node")
	eerr = categoriesUpdateWorkflow.AddExecutionNode(emptyNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Decision nodes below
	// Path based decision node. Returns false if length is not 1.
	validationNode := new(ValidationNode)
	validationNode.SetID("Path decision Node")

	categoriesUpdateWorkflow.AddDecisionNode(validationNode, categoriesCheckNode, emptyNode)
	categoriesUpdateWorkflow.AddConnection(categoriesCheckNode, categoriesUpdateNode)
	//Set start node for the search workflow
	categoriesUpdateWorkflow.SetStartNode(validationNode)

	//Assign the workflow definition to the Orchestrator
	categoriesUpdateOrchestrator.Create(categoriesUpdateWorkflow)

	logger.Info(categoriesUpdateOrchestrator.String())
	logger.Info("Categories Update Pipeline Updated")
	return *categoriesUpdateOrchestrator
}

// GetHealthCheck -> healthcheck Boilerplate
func (a *CategoryAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CategoryUpdateHealthCheck)
}

// Init -> Basic Init function.
func (a *CategoryAPI) Init() {
	logger.Info("Starting worker pool for Category PUT API")
	categoryUpdatePool = pool.NewWorkerSafe(constants.CATEGORY_UPDATE, common.CATEGORY_UPDATE_POOL_SIZE, common.CATEGORY_UPDATE_QUEUE_SIZE, common.CATEGORY_UPDATE_RETRY_COUNT, common.CATEGORY_UPDATE_WAIT_TIME)
	categoryUpdatePool.StartWorkers(common.CategoryUpdateWorker)
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
