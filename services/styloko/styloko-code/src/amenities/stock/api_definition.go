package stock

import (
	_ "common/constants"
	_ "fmt"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	_ "github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type StockAPI struct {
}

func (a *StockAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: "STOCK",
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Returns the Orchestrator
func (a *StockAPI) GetOrchestrator() orchestrator.Orchestrator {
	stockNode := new(StockNode)
	stockNode.SetID("stockNode")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	workflow.AddExecutionNode(stockNode)

	workflow.SetStartNode(stockNode)

	orchestrator := new(orchestrator.Orchestrator)
	orchestrator.Create(workflow)
	return *orchestrator

}

func (a *StockAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(StockHealthCheck)
}

func (a *StockAPI) Init() {
}
