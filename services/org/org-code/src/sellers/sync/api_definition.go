package sync

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"sellers/common"
)

type Sync struct {
}

func (s *Sync) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.SYNC,
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (s *Sync) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Sync API pipeline begins")

	getSyncWorkflow := new(orchestrator.WorkFlowDefinition)
	getSyncWorkflow.Create()

	getSyncNode := new(StartSync)
	getSyncNode.SetID("Seller Sync Node")
	err := getSyncWorkflow.AddExecutionNode(getSyncNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	getSyncWorkflow.SetStartNode(getSyncNode)

	getSyncOrchestrator := new(orchestrator.Orchestrator)
	getSyncOrchestrator.Create(getSyncWorkflow)
	logger.Info(getSyncOrchestrator.String())
	return *getSyncOrchestrator
}

func (s *Sync) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(getSyncHealthCheck)
}

func (s *Sync) Init() {
	//api initialization should come here
}
