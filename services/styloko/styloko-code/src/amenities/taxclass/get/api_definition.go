package get

import (
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type TaxClassApi struct {
}

func (a *TaxClassApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: TAXCLASS_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *TaxClassApi) GetOrchestrator() orchestrator.Orchestrator {
	//Initiate required nodes
	validateNode := new(ValidateNode)
	validateNode.SetID("ValidateRequest")

	daNode := new(DataAccessorNode)
	daNode.SetID("dataAccessor")

	responseNode := new(ResponseNode)
	responseNode.SetID("Response")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	//       validate
	//           |
	//      Fetch Data
	//           |
	//      Prepare Resp
	//
	//Set up workflow

	//add execution nodes
	workflow.AddExecutionNode(validateNode)
	workflow.AddExecutionNode(responseNode)
	workflow.AddExecutionNode(daNode)

	//add connections
	err := workflow.AddConnection(validateNode, daNode)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(daNode, responseNode)
	if err != nil {
		logger.Error(err)
	}
	//set start node
	workflow.SetStartNode(validateNode)

	orchestrator := new(orchestrator.Orchestrator)
	orchestrator.Create(workflow)
	return *orchestrator
}

func (a *TaxClassApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(TaxClassHealthCheck)
}

func (a *TaxClassApi) Init() {

}
