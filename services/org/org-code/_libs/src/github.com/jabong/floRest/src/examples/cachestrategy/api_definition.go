package cachestrategy

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type HelloApi struct {
}

func (a *HelloApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: "CACHESTRATEGY",
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *HelloApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Cache Strategy Pipeline Creation begin")

	cacheStrategyOrchestrator := new(orchestrator.Orchestrator)
	cacheStrategyWorkflow := new(orchestrator.WorkFlowDefinition)
	cacheStrategyWorkflow.Create()

	//Creation of the nodes in the workflow definition
	cacheStrategyNode := new(CacheStrategyUser)
	cacheStrategyNode.SetID("cache strategy user node 1")
	eerr := cacheStrategyWorkflow.AddExecutionNode(cacheStrategyNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	cacheStrategyWorkflow.SetStartNode(cacheStrategyNode)

	//Assign the workflow definition to the Orchestrator
	cacheStrategyOrchestrator.Create(cacheStrategyWorkflow)

	logger.Info(cacheStrategyOrchestrator.String())
	logger.Info("Cache Strategy Pipeline Created")
	return *cacheStrategyOrchestrator
}

func (a *HelloApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return nil
}

func (a *HelloApi) Init() {
	//api initialization should come here
}
