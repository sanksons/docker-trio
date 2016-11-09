package set

import (
	"common/constants"
	"fmt"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	validator "gopkg.in/go-playground/validator.v8"
)

type AttributeSetUpdateApi struct {
}

var validate *validator.Validate

func (a *AttributeSetUpdateApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.ATTRIBUTESETAPI,
		Version:  "V1",
		Action:   "PUT",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *AttributeSetUpdateApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Attributes Set Updation begin")

	attributesSetWorkflow := new(orchestrator.WorkFlowDefinition)
	attributesSetWorkflow.Create()

	updateNode := new(UpdateSet)
	updateNode.SetID("Attribute Set update node")
	err := attributesSetWorkflow.AddExecutionNode(updateNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Attribute Set put failure node")
	err = attributesSetWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//AddDecisionNode(decisionNode,yesNode,noNode)
	validationNode := new(Validation)
	validationNode.SetID("Attribute Set put decision node")
	err = attributesSetWorkflow.AddDecisionNode(validationNode, updateNode, failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	attributesSetWorkflow.SetStartNode(validationNode)

	//Assign the workflow definition to the Orchestrator
	attributesSetOrchestrator := new(orchestrator.Orchestrator)
	attributesSetOrchestrator.Create(attributesSetWorkflow)

	logger.Info(attributesSetOrchestrator.String())
	logger.Info("Attributes Set Pipeline Created")
	return *attributesSetOrchestrator
}

func (a *AttributeSetUpdateApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(AttributesSetHealthCheck)
}

func (a *AttributeSetUpdateApi) Init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}
