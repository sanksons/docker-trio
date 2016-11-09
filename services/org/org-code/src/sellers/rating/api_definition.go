package rating

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

var cacheObj cache.CacheInterface
var prodUpdatePool pool.Safe

type UploadRatingApi struct {
}

func (u *UploadRatingApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: RATING,
		Version:  "V1",
		Action:   "POST",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (u *UploadRatingApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Upload Rating API pipeline begins")

	uploadRatingOrchestrator := new(orchestrator.Orchestrator)
	uploadRatingWorkflow := new(orchestrator.WorkFlowDefinition)
	uploadRatingWorkflow.Create()

	//Creation of the nodes in the workflow definition
	uploadRatingNode := new(UploadRating)
	uploadRatingNode.SetID("Upload Rating Node")
	eerr := uploadRatingWorkflow.AddExecutionNode(uploadRatingNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	uploadRatingWorkflow.SetStartNode(uploadRatingNode)

	//Assign the workflow definition to the Orchestrator
	uploadRatingOrchestrator.Create(uploadRatingWorkflow)

	logger.Info(uploadRatingOrchestrator.String())
	return *uploadRatingOrchestrator
}

func (a *UploadRatingApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(UploadRatingHeathCheck)
}

func (a *UploadRatingApi) Init() {
	//api initialization should come here
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting cache object :%s", err.Error()))
	}
	logger.Info("Starting worker pool for PRODUCT INVALIDATION")
	prodUpdatePool = pool.NewWorkerSafe(PRODUCT_INAVLIDATION, PRODUCT_INAVLIDATION_POOL_SIZE, PRODUCT_INAVLIDATION_QUEUE_SIZE, PRODUCT_INAVLIDATION_RETRY_COUNT, PRODUCT_INAVLIDATION_WAIT_TIME)
	prodUpdatePool.StartWorkers(common.InvalidateProductsForUpdatedSellers)
}
