package migration

import (
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"migration/common"
)

type MigrationApi struct {
}

func (m *MigrationApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.MIGRATION,
		Version:  "V1",
		Action:   "POST",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (m *MigrationApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Migration API Pipeline Creation begin")

	migrationOrchestrator := new(orchestrator.Orchestrator)
	migrationWorkflow := new(orchestrator.WorkFlowDefinition)
	migrationWorkflow.Create()

	//Creation of the nodes in the workflow definition
	managerNode := new(Manager)
	managerNode.SetID("Migration Manager Node")
	eerr := migrationWorkflow.AddExecutionNode(managerNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	migrationWorkflow.SetStartNode(managerNode)

	//Assign the workflow definition to the Orchestrator
	migrationOrchestrator.Create(migrationWorkflow)

	logger.Info(migrationOrchestrator.String())
	logger.Info("Migration API Pipeline Created")
	return *migrationOrchestrator
}

func (m *MigrationApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(MigrationHealthCheck)
}

func (m *MigrationApi) Init() {
	setFlagInMongo()
}
