package put

import (
	"common/appconfig"
	"common/pool"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"sellers/common"
)

var (
	cacheObj         cache.CacheInterface
	sellerUpdatePool pool.Safe
	updateErpPool    pool.Safe
	prodUpdatePool   pool.Safe
)

type UpdateSellerApi struct {
}

func (su *UpdateSellerApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.SELLERS,
		Version:  "V1",
		Action:   "PUT",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (su *UpdateSellerApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Seller Update API pipeline begins")

	sellerUpdateWorkflow := new(orchestrator.WorkFlowDefinition)
	sellerUpdateWorkflow.Create()

	//Creation of the nodes in the workflow definition

	updateNode := new(Update)
	updateNode.SetID("Seller Update Node")
	err := sellerUpdateWorkflow.AddExecutionNode(updateNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Update Update Failure Node")
	err = sellerUpdateWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//AddDecisionNode(decisionNode,yesNode,noNode)
	decisionNode := new(Decision)
	decisionNode.SetID("Seller Update Validation Node")
	err = sellerUpdateWorkflow.AddDecisionNode(decisionNode, updateNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	sellerUpdateWorkflow.SetStartNode(decisionNode)

	sellerUpdateOrchestrator := new(orchestrator.Orchestrator)
	sellerUpdateOrchestrator.Create(sellerUpdateWorkflow)
	logger.Info(sellerUpdateOrchestrator.String())
	return *sellerUpdateOrchestrator
}

func (su *UpdateSellerApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(sellerUpdateHealthCheck)
}

func (su *UpdateSellerApi) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting cache object: %s", err.Error()))
	}
	logger.Info("Starting worker pool for SELLER PUT API")
	sellerUpdatePool = pool.NewWorkerSafe(SELLER_UPDATE, SELLER_UPDATE_POOL_SIZE, SELLER_UPDATE_QUEUE_SIZE, SELLER_UPDATE_RETRY_COUNT, SELLER_UPDATE_WAIT_TIME)
	sellerUpdatePool.StartWorkers(common.SyncWithMysql)
	logger.Info("Starting worker pool for UPDATE ERP")
	updateErpPool = pool.NewWorkerSafe(UPDATE_ERP, UPDATE_ERP_POOL_SIZE, UPDATE_ERP_QUEUE_SIZE, UPDATE_ERP_RETRY_COUNT, UPDATE_ERP_WAIT_TIME)
	updateErpPool.StartWorkers(UpdateOnErp)
	logger.Info("Starting worker pool for PRODUCT INVALIDATION")
	prodUpdatePool = pool.NewWorkerSafe(PRODUCT_INAVLIDATION, PRODUCT_INAVLIDATION_POOL_SIZE, PRODUCT_INAVLIDATION_QUEUE_SIZE, PRODUCT_INAVLIDATION_RETRY_COUNT, PRODUCT_INAVLIDATION_WAIT_TIME)
	prodUpdatePool.StartWorkers(common.InvalidateProductsForUpdatedSellers)
}
