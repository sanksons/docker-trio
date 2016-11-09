package get

import (
	"common/appconfig"

	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

const (
	category_tree_api = "CATEGORYTREE"
)

var cacheObj cache.CacheInterface

// CategoryTreeApi base struct
type CategoryTreeApi struct {
}

// GetVersion returns version number
func (a *CategoryTreeApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: category_tree_api,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator returns orchestrator
func (a *CategoryTreeApi) GetOrchestrator() orchestrator.Orchestrator {
	responseNode := new(CategoryTreeGet)
	responseNode.SetID("CategoryTreeNode")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	workflow.AddExecutionNode(responseNode)

	workflow.SetStartNode(responseNode)

	orchestrator := new(orchestrator.Orchestrator)
	orchestrator.Create(workflow)
	return *orchestrator
}

// Init initializes the API
func (a *CategoryTreeApi) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}

// GetHealthCheck -> Returns HealthCheckInterface
func (a *CategoryTreeApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CategoryTreeHealthCheck)
}
