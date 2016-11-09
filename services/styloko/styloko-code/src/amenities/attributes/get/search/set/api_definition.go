package set

import (
	"common/constants"
	"fmt"

	"common/appconfig"

	"github.com/jabong/floRest/src/common/cache"
	conf "github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

var cacheObj cache.CacheInterface

type GetAttributesSetsApi struct {
}

type M map[string]interface{}

func (a *GetAttributesSetsApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.ATTRIBUTESETAPI,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *GetAttributesSetsApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Attributes Set Creation begin")

	attributesSetWorkflow := new(orchestrator.WorkFlowDefinition)
	attributesSetWorkflow.Create()

	//Creation of the nodes in the workflow definition
	getOneNode := new(GetAttributeSet)
	getOneNode.SetID("Attributes Set Get One")
	eerr := attributesSetWorkflow.AddExecutionNode(getOneNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	getAllNode := new(GetAllAttributeSet)
	getAllNode.SetID("Attributes Set Get All")
	err := attributesSetWorkflow.AddExecutionNode(getAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	setSearchNode := new(SearchAttributeSet)
	setSearchNode.SetID("Attributes Set Search")
	err = attributesSetWorkflow.AddExecutionNode(setSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	pathParamNode := new(PathParamDecision)
	pathParamNode.SetID("Attribute Set Path Param")
	err = attributesSetWorkflow.AddDecisionNode(pathParamNode, getOneNode, getAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	queryParamNode := new(QueryParamDecision)
	queryParamNode.SetID("Attribute Set Query Param")
	err = attributesSetWorkflow.AddDecisionNode(queryParamNode, setSearchNode, pathParamNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}
	//Set start node for the search workflow
	attributesSetWorkflow.SetStartNode(queryParamNode)

	attributesSetOrchestrator := new(orchestrator.Orchestrator)
	attributesSetOrchestrator.Create(attributesSetWorkflow)
	logger.Info(attributesSetOrchestrator.String())
	return *attributesSetOrchestrator
}

func (a *GetAttributesSetsApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(AttributesSetHealthCheck)
}

func (a *GetAttributesSetsApi) Init() {
	c := conf.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(c.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
