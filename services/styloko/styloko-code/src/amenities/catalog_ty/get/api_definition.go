package get

import (
	"common/constants"

	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/versionmanager"
)

// CatalogTyApi base struct
type CatalogTyApi struct {
}

// GetVersion returns version number
func (a *CatalogTyApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.CATALOG_TY_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator returns orchestrator
func (a *CatalogTyApi) GetOrchestrator() orchestrator.Orchestrator {
	responseNode := new(CatalogTyGet)
	responseNode.SetID("CatalogTyNode")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	workflow.AddExecutionNode(responseNode)

	workflow.SetStartNode(responseNode)

	orchestrator := new(orchestrator.Orchestrator)
	orchestrator.Create(workflow)
	return *orchestrator
}

// Init initializes the API
func (a *CatalogTyApi) Init() {

}

// GetHealthCheck -> Returns HealthCheckInterface
func (a *CatalogTyApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(CatalogTyHealthCheck)
}
