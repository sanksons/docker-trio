package post

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type UpdateErpApi struct {
}

func (u *UpdateErpApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: ERP,
		Version:  "V1",
		Action:   "POST",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (u *UpdateErpApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Update Erp API pipeline begins")

	erpUdateOrchestrator := new(orchestrator.Orchestrator)
	erpUpdateWorkflow := new(orchestrator.WorkFlowDefinition)
	erpUpdateWorkflow.Create()

	//Creation of the nodes in the workflow definition
	erpUdateNode := new(ErpUpdate)
	erpUdateNode.SetID("ERP Update Node")
	eerr := erpUpdateWorkflow.AddExecutionNode(erpUdateNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	erpUpdateWorkflow.SetStartNode(erpUdateNode)

	//Assign the workflow definition to the Orchestrator
	erpUdateOrchestrator.Create(erpUpdateWorkflow)

	logger.Info(erpUdateOrchestrator.String())
	return *erpUdateOrchestrator
}

func (a *UpdateErpApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(ErpUdateHealthCheck)
}

func (a *UpdateErpApi) Init() {
	//api initialization should come here
}
