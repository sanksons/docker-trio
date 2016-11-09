package post

import (
	proUtils "amenities/products/common"
	sizeUtil "amenities/sizechart/common"
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
)

var (
	cacheMngr           proUtils.CacheManager
	conf                *appconfig.AppConfig
	dbAdapterName       string
	sizechartCreatePool pool.Safe
	sizechartProdPool   pool.Safe
)

type SizeChartCreateApi struct {
}

func (a *SizeChartCreateApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: sizeUtil.SizeChartResource,
		Version:  "V1",
		Action:   "POST",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (a *SizeChartCreateApi) GetOrchestrator() orchestrator.Orchestrator {
	//           SCTypeDecisionNode
	//           /			    \
	//       SkuSCValidate     StdSCValidate
	//               \		 /
	//	 			  SCCreate
	//					|
	//				MapSizeChartProduct
	logger.Info("Upload Size Chart Pipeline Creation begin")

	sizeChartOrchestrator := new(orchestrator.Orchestrator)
	sizeChartWorkflow := new(orchestrator.WorkFlowDefinition)
	sizeChartWorkflow.Create()

	sizeChartNode1 := new(SizeChartTypeDecision)
	sizeChartNode1.SetID("SizeChartType")

	sizeChartNode2 := new(SkuSizeChartValidate)
	sizeChartNode2.SetID("SkuSizeChartValidate")
	err := sizeChartWorkflow.AddExecutionNode(sizeChartNode2)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sizeChartNode4 := new(SizeChartValidate)
	sizeChartNode4.SetID("BrandBrickSizeChartValidate")
	err = sizeChartWorkflow.AddExecutionNode(sizeChartNode4)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sizeChartNode5 := new(SizeChartCreate)
	sizeChartNode5.SetID("SizeChartCreate")
	err = sizeChartWorkflow.AddExecutionNode(sizeChartNode5)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sizeChartNode6 := new(MapSizeChartProduct)
	sizeChartNode6.SetID("MapSizeChartProduct")
	err = sizeChartWorkflow.AddExecutionNode(sizeChartNode6)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	err = sizeChartWorkflow.AddDecisionNode(sizeChartNode1, sizeChartNode2, sizeChartNode4)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	err = sizeChartWorkflow.AddConnection(sizeChartNode2, sizeChartNode5)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	err = sizeChartWorkflow.AddConnection(sizeChartNode4, sizeChartNode5)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	err = sizeChartWorkflow.AddConnection(sizeChartNode5, sizeChartNode6)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sizeChartWorkflow.SetStartNode(sizeChartNode1)

	//Assign the workflow definition to the Orchestrator
	sizeChartOrchestrator.Create(sizeChartWorkflow)

	logger.Info("Upload Size Chart Pipeline Created")
	return *sizeChartOrchestrator
}

func (a *SizeChartCreateApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(SizeChartHealthCheck)
}

func (a *SizeChartCreateApi) Init() {
	//api initialization should come here

	logger.Info("Starting worker pool for Sizechart POST API")
	sizechartCreatePool = pool.NewWorkerSafe(sizeUtil.SizeChartCreate, sizeUtil.PoolSize,
		sizeUtil.QueueSize, sizeUtil.RetryCount,
		sizeUtil.WaitTime)
	sizechartProdPool = pool.NewWorkerSafe(sizeUtil.SizeChartProdCreate, sizeUtil.PoolSize,
		sizeUtil.QueueSize, sizeUtil.RetryCount,
		sizeUtil.WaitTime)
	sizechartCreatePool.StartWorkers(sizeUtil.CreateSizeChartWorker)
	sizechartProdPool.StartWorkers(sizeUtil.CreateSizeChartToProductWorker)

	// cache initialisation
	var err error
	conf = config.ApplicationConfig.(*appconfig.AppConfig)
	dbAdapterName = conf.DbAdapter
	//Initialize cache object
	cacheMngr.CacheObj, err = cache.Get(conf.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
