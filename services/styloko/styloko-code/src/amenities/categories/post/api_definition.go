package post

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

// CategoryAPI -> Basic orchestrator struct
type CategoryAPI struct {
}

var cacheObj cache.CacheInterface

// categoryCreatePool -> Pool object for job dispatcher
var categoryCreatePool pool.Safe

// GetVersion -> Version Manager boilerplate
func (a *CategoryAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.CATEGORY_API,
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Orchestrator Boilerplate
func (a *CategoryAPI) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Categories Create Creation begin")

	categoriesCreateOrchestrator := new(orchestrator.Orchestrator)
	categoriesCreateWorkflow := new(orchestrator.WorkFlowDefinition)
	categoriesCreateWorkflow.Create()

	//Creation of the nodes in the workflow definition
	categoriesCreateNode := new(CategoriesCreate)
	categoriesCreateNode.SetID("Categories Create Node")
	eerr := categoriesCreateWorkflow.AddExecutionNode(categoriesCreateNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Empty node is required for false edge of decision node.
	// As no process needs to run in that case.
	emptyNode := new(EmptyNode)
	emptyNode.SetID("Empty Node")
	eerr = categoriesCreateWorkflow.AddExecutionNode(emptyNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Data Validation node.
	validationNode := new(CategoriesCreateValidation)
	validationNode.SetID("Categories Validation Node")

	//Set start node for the search workflow
	categoriesCreateWorkflow.AddDecisionNode(validationNode, categoriesCreateNode, emptyNode)

	categoriesCreateWorkflow.SetStartNode(validationNode)

	//Assign the workflow definition to the Orchestrator
	categoriesCreateOrchestrator.Create(categoriesCreateWorkflow)

	logger.Info(categoriesCreateOrchestrator.String())
	logger.Info("Categories Create Pipeline Created")
	return *categoriesCreateOrchestrator
}

// GetHealthCheck -> healthcheck Boilerplate
func (a *CategoryAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CategoryCreateHealthCheck)
}

// Init -> Basic Init function.
func (a *CategoryAPI) Init() {
	logger.Info("Starting worker pool for Category POST API")
	categoryCreatePool = pool.NewWorkerSafe(constants.CATEGORY_CREATE, common.CATEGORY_CREATE_POOL_SIZE, common.CATEGORY_CREATE_QUEUE_SIZE, common.CATEGORY_CREATE_RETRY_COUNT, common.CATEGORY_CREATE_WAIT_TIME)
	categoryCreatePool.StartWorkers(common.CategoryCreateWorker)
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
