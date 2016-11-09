package get

import (
	"common/constants"
	"fmt"

	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

// MigrationAPI -> Base struct for MigrationAPI
type MigrationAPI struct {
}

// GetVersion -> Returns versionmanager instance
func (a *MigrationAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.MIGRATION_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Returns the Orchestrator
func (a *MigrationAPI) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Migration API workflow begin")

	migrationOrchestrator := new(orchestrator.Orchestrator)
	migrationWorkflow := new(orchestrator.WorkFlowDefinition)
	migrationWorkflow.Create()

	// migrationGetNode fetches category by ID
	migrationGetNode := new(MigrationGet)
	migrationGetNode.SetID("Migrations Get Node")
	eerr := migrationWorkflow.AddExecutionNode(migrationGetNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	statusCheckNode := new(StatusCheck)
	statusCheckNode.SetID("Statuc Check Node")
	eerr = migrationWorkflow.AddExecutionNode(statusCheckNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	// pathDecision fetches category by ID
	pathDecision := new(PathDecision)
	pathDecision.SetID("Migrations Path Decision Node")
	eerr = migrationWorkflow.AddDecisionNode(pathDecision, migrationGetNode, statusCheckNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the migration workflow
	migrationWorkflow.SetStartNode(pathDecision)

	migrationOrchestrator.Create(migrationWorkflow)
	logger.Info(migrationOrchestrator.String())
	logger.Info("Migration pipeline created")
	return *migrationOrchestrator
}

// GetHealthCheck -> Returns HealthCheckInterface
func (a *MigrationAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(MigrationHealthCheck)
}

// Init -> API initialization function.
func (a *MigrationAPI) Init() {
	//api initialization should come here
}
