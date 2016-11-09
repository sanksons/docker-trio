package examples

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
		Resource: "HELLO",
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (a *HelloApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Error("Hello World Pipeline Creation begin")

	helloWorldOrchestrator := new(orchestrator.Orchestrator)
	helloWorldWorkflow := new(orchestrator.WorkFlowDefinition)
	helloWorldWorkflow.Create()

	//Creation of the nodes in the workflow definition
	helloWorldNode := new(HelloWorld)
	helloWorldNode.SetID("hello world node 1")
	eerr := helloWorldWorkflow.AddExecutionNode(helloWorldNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	//Set start node for the search workflow
	helloWorldWorkflow.SetStartNode(helloWorldNode)

	//Assign the workflow definition to the Orchestrator
	helloWorldOrchestrator.Create(helloWorldWorkflow)

	logger.Info(helloWorldOrchestrator.String())
	logger.Info("Hello World Pipeline Created")
	return *helloWorldOrchestrator
}

func (a *HelloApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(HelloWorldHealthCheck)
}

func (a *HelloApi) Init() {
	//api initialization should come here
}
