package get

import (
	"common/appconfig"
	//	"common/pool"
	sizeUtil "amenities/sizechart/common"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

var cacheObj cache.CacheInterface

type SizeChartGetApi struct {
}

func (a *SizeChartGetApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: sizeUtil.SizeChartResource,
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (a *SizeChartGetApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("GetAPI sizechart Pipeline Begins.")

	sizeChartOrchestrator := new(orchestrator.Orchestrator)
	sizeChartWorkflow := new(orchestrator.WorkFlowDefinition)
	sizeChartWorkflow.Create()

	sizeChartNode1 := new(SizeChartGet)
	sizeChartNode1.SetID("SizeChartGet")
	err := sizeChartWorkflow.AddExecutionNode(sizeChartNode1)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sizeChartWorkflow.SetStartNode(sizeChartNode1)

	//Assign the workflow definition to the Orchestrator
	sizeChartOrchestrator.Create(sizeChartWorkflow)

	logger.Info("GetAPI sizechart Pipeline Created")
	return *sizeChartOrchestrator
}

func (a *SizeChartGetApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(SizeChartHealthCheck)
}

func (a *SizeChartGetApi) Init() {
	//api initialization should come here
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
