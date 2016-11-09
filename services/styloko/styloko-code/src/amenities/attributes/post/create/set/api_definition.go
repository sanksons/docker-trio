package set

import (
	"common/constants"
	mongodb "common/mongodb"
	"fmt"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	validator "gopkg.in/go-playground/validator.v8"
)

type AttributeSetCreateApi struct {
}

var mongoDriver *mongodb.MongoDriver

var validate *validator.Validate

func (a *AttributeSetCreateApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.ATTRIBUTESETAPI,
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *AttributeSetCreateApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Attributes Set Insertion begin")

	attributesSetWorkflow := new(orchestrator.WorkFlowDefinition)
	attributesSetWorkflow.Create()

	insertNode := new(InsertSet)
	insertNode.SetID("Attribute Set insert node")
	err := attributesSetWorkflow.AddExecutionNode(insertNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	failureNode := new(Failure)
	failureNode.SetID("Attribute Set post failure node")
	err = attributesSetWorkflow.AddExecutionNode(failureNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//AddDecisionNode(decisionNode,yesNode,noNode)
	validationNode := new(Validation)
	validationNode.SetID("Attribute Set post decision node")
	err = attributesSetWorkflow.AddDecisionNode(validationNode, insertNode, failureNode)
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

func (a *AttributeSetCreateApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(AttributesSetHealthCheck)
}

func (a *AttributeSetCreateApi) Init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}
