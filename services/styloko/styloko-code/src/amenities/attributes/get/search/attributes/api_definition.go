package attributes

import (
	"common/constants"
	"fmt"

	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type GetAttributesApi struct {
}

type M map[string]interface{}

func (a *GetAttributesApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.ATTRIBUTEAPI,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *GetAttributesApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Attributes Search begin")

	attributesWorkflow := new(orchestrator.WorkFlowDefinition)
	attributesWorkflow.Create()

	//Creation of the nodes in the workflow definition
	getOneNode := new(GetAttribute)
	getOneNode.SetID("Attributes Get One")
	eerr := attributesWorkflow.AddExecutionNode(getOneNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	getAllNode := new(GetAllAttributes)
	getAllNode.SetID("Attribute Get All")
	err := attributesWorkflow.AddExecutionNode(getAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	attrSearchNode := new(SearchAttribute)
	attrSearchNode.SetID("Attribute Search")
	err = attributesWorkflow.AddExecutionNode(attrSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	pathParamNode := new(PathParamDecision)
	pathParamNode.SetID("Attribute Path Param")
	err = attributesWorkflow.AddDecisionNode(pathParamNode, getOneNode, getAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	queryParamNode := new(QueryParamDecision)
	queryParamNode.SetID("Attribute Query Param")
	err = attributesWorkflow.AddDecisionNode(queryParamNode, attrSearchNode, pathParamNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}
	//Set start node for the search workflow
	attributesWorkflow.SetStartNode(queryParamNode)

	attributesOrchestrator := new(orchestrator.Orchestrator)
	attributesOrchestrator.Create(attributesWorkflow)
	logger.Info(attributesOrchestrator.String())
	return *attributesOrchestrator
}

func (a *GetAttributesApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(AttributesHealthCheck)
}

func (a *GetAttributesApi) Init() {
}
