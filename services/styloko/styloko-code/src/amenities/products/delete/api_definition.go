package delete

import (
	proUtil "amenities/products/common"
	"common/appconfig"
	_ "fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

var (
	cacheMngr     proUtil.CacheManager
	conf          *appconfig.AppConfig
	dbAdapterName string
)

type DeleteProductsApi struct {
}

func (a *DeleteProductsApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: proUtil.PRODUCT_API,
		Version:  "V1",
		Action:   "DELETE",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *DeleteProductsApi) GetOrchestrator() orchestrator.Orchestrator {

	//Initiate required nodes

	DeleteCacheNode := new(DeleteCacheNode)
	DeleteCacheNode.SetID("DeleteCacheNode")

	Workflow := new(orchestrator.WorkFlowDefinition)
	Workflow.Create()

	//Set up workflow
	//add execution nodes
	Workflow.AddExecutionNode(DeleteCacheNode)

	//set start node
	Workflow.SetStartNode(DeleteCacheNode)

	Orchestrator := new(orchestrator.Orchestrator)
	Orchestrator.Create(Workflow)
	return *Orchestrator
}

func (a *DeleteProductsApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(ProductDeleteHealthCheck)
}

func (a *DeleteProductsApi) Init() {
	conf = config.ApplicationConfig.(*appconfig.AppConfig)
	//set current DB Adapter
	dbAdapterName = conf.DbAdapter
	var err error
	cacheMngr.CacheObj, err = cache.Get(conf.Cache)
	if err != nil {
		logger.Error(err.Error())
	}

}
