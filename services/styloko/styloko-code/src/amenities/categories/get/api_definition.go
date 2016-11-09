package get

import (
	"common/appconfig"
	"common/constants"
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

// CategoryAPI -> Base struct for CategoryAPI
type CategoryAPI struct {
}

// GetVersion -> Returns versionmanager instance
func (a *CategoryAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.CATEGORY_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Returns the Orchestrator
func (a *CategoryAPI) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Categories Search Creation begin")

	categoriesSearchOrchestrator := new(orchestrator.Orchestrator)
	categoriesSearchWorkflow := new(orchestrator.WorkFlowDefinition)
	categoriesSearchWorkflow.Create()

	// categoriesGetID fetches category by ID
	categoriesGetID := new(CategoriesGetID)
	categoriesGetID.SetID("Categories Get By ID Node")
	eerr := categoriesSearchWorkflow.AddExecutionNode(categoriesGetID)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// categoriesGetAll fetches all categories
	categoriesGetAll := new(CategoriesGetAll)
	categoriesGetAll.SetID("Categories Get All Node")
	eerr = categoriesSearchWorkflow.AddExecutionNode(categoriesGetAll)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// categoriesGetQuery fetches all categories by status type
	categoriesGetQuery := new(CategoriesGetQuery)
	categoriesGetQuery.SetID("Categories Get Query Node")
	eerr = categoriesSearchWorkflow.AddExecutionNode(categoriesGetQuery)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Decision nodes below
	// Path based decision node. Returns false if length is not 1.
	pathDecisionNode := new(PathDecision)
	pathDecisionNode.SetID("Path decision Node")

	// Checks for query params. Returns false if no query params found.
	queryDecisionNode := new(QueryDecision)
	queryDecisionNode.SetID("Query Decision Node")

	// emptyNode is an empty node that does nothing. Added due to decision
	// node limitation.
	emptyNode := new(EmptyNode)
	emptyNode.SetID("Empty Node")
	eerr = categoriesSearchWorkflow.AddExecutionNode(emptyNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Path based decision node.
	eerr = categoriesSearchWorkflow.AddDecisionNode(pathDecisionNode, categoriesGetID, emptyNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	/*
		Empty node is required here becuase floRest does a cycle check and passing a decision node as an
		argument to the function, causes a cycle detection and code crashes. Please refrain from adding
		decsion nodes as an edge to a decision node itself.
	*/

	// Query decision node.
	eerr = categoriesSearchWorkflow.AddDecisionNode(queryDecisionNode, categoriesGetQuery, categoriesGetAll)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// Empty Node to act as a proxy.
	eerr = categoriesSearchWorkflow.AddConnection(emptyNode, queryDecisionNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	categoriesSearchWorkflow.SetStartNode(pathDecisionNode)

	//Assign the workflow definition to the Orchestrator
	categoriesSearchOrchestrator.Create(categoriesSearchWorkflow)

	logger.Info(categoriesSearchOrchestrator.String())
	logger.Info("Categories Search Pipeline Created")
	return *categoriesSearchOrchestrator
}

// GetHealthCheck -> Returns HealthCheckInterface
func (a *CategoryAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CategorySearchHealthCheck)
}

// Init -> API initialization function.
func (a *CategoryAPI) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
