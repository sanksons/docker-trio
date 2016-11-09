package post

import (
	"common/pool"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"sellers/common"
)

var sellerInsertPool pool.Safe

type CreateSellerApi struct {
}

func (sc *CreateSellerApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.SELLERS,
		Version:  "V1",
		Action:   "POST",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (sc *CreateSellerApi) GetOrchestrator() orchestrator.Orchestrator {

	logger.Info("Seller Create API pipeline begins")

	sellerCreateWorkflow := new(orchestrator.WorkFlowDefinition)
	sellerCreateWorkflow.Create()

	sellerMongoNode := new(Insert)
	sellerMongoNode.SetID("Seller Create Inset Node")
	err := sellerCreateWorkflow.AddExecutionNode(sellerMongoNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Seller Create Failure Node")
	err = sellerCreateWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	decisionNode := new(Validate)
	decisionNode.SetID("Seller Create Validation Node")
	err = sellerCreateWorkflow.AddDecisionNode(decisionNode, sellerMongoNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	sellerCreateWorkflow.SetStartNode(decisionNode)

	sellerCreateOrchestrator := new(orchestrator.Orchestrator)
	sellerCreateOrchestrator.Create(sellerCreateWorkflow)
	logger.Info(sellerCreateOrchestrator.String())
	return *sellerCreateOrchestrator
}

func (sc *CreateSellerApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(SellerCreateHealthCheck)
}

func (sc *CreateSellerApi) Init() {
	//api initialization should come here
	logger.Info("Starting worker pool for SELLER POST API")
	sellerInsertPool = pool.NewWorkerSafe(SELLER_INSERT, SELLER_INSERT_POOL_SIZE, SELLER_INSERT_QUEUE_SIZE, SELLER_INSERT_RETRY_COUNT, SELLER_INSERT_WAIT_TIME)
	sellerInsertPool.StartWorkers(common.SyncWithMysql)
}
