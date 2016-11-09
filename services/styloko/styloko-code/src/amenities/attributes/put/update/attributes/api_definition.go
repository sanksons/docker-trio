package attributes

import (
	"common/appconfig"
	"common/constants"
	"common/pool"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	conf "github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	validator "gopkg.in/go-playground/validator.v8"
)

type AttributesUpdateApi struct {
}

// attributeUpdatePool -> Pool object for job dispatcher
var attributeUpdatePool pool.Safe

var validate *validator.Validate

type M map[string]interface{}

var cacheObj cache.CacheInterface

func (a *AttributesUpdateApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: constants.ATTRIBUTEAPI,
		Version:  "V1",
		Action:   "PUT",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *AttributesUpdateApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Attributes Updation begin")

	attributeWorkflow := new(orchestrator.WorkFlowDefinition)
	attributeWorkflow.Create()

	updateOptionNode := new(UpdateOption)
	updateOptionNode.SetID("Attribute option update node")
	err := attributeWorkflow.AddExecutionNode(updateOptionNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	updateNode := new(UpdateAttribute)
	updateNode.SetID("Attribute update node")
	err1 := attributeWorkflow.AddExecutionNode(updateNode)
	if err1 != nil {
		logger.Error(fmt.Sprintln(err1))
	}

	failureNode := new(Failure)
	failureNode.SetID("Attribute put failure node")
	err2 := attributeWorkflow.AddExecutionNode(failureNode)
	if err2 != nil {
		logger.Error(fmt.Sprintln(err2))
	}

	failureNode1 := new(Failure)
	failureNode1.SetID("Attribute put failure node1")
	err3 := attributeWorkflow.AddExecutionNode(failureNode1)
	if err3 != nil {
		logger.Error(fmt.Sprintln(err3))
	}

	emptyNode1 := new(EmptyNode)
	emptyNode1.SetID("Empty Node1")
	eerr9 := attributeWorkflow.AddExecutionNode(emptyNode1)
	if eerr9 != nil {
		logger.Error(fmt.Sprintln(eerr9))
	}

	emptyNode2 := new(EmptyNode)
	emptyNode2.SetID("Empty Node2")
	eerr4 := attributeWorkflow.AddExecutionNode(emptyNode2)
	if eerr4 != nil {
		logger.Error(fmt.Sprintln(eerr4))
	}

	pathDecisionNode := new(PathParamDecision)
	pathDecisionNode.SetID("Attribute path param decision")

	validateOptionNode := new(ValidateOption)
	validateOptionNode.SetID("Attribute validate Option decision node")

	validationNode := new(Validation)
	validationNode.SetID("Attribute put decision node")

	err6 := attributeWorkflow.AddDecisionNode(validateOptionNode,
		updateOptionNode, failureNode)
	if err6 != nil {
		logger.Error(fmt.Sprintln(err6))
	}

	err8 := attributeWorkflow.AddDecisionNode(validationNode,
		updateNode, failureNode1)
	if err8 != nil {
		logger.Error(fmt.Sprintln(err8))
	}

	eerr10 := attributeWorkflow.AddConnection(emptyNode1, validateOptionNode)
	if eerr10 != nil {
		logger.Error(fmt.Sprintln(eerr10))
	}

	eerr7 := attributeWorkflow.AddConnection(emptyNode2, validationNode)
	if eerr7 != nil {
		logger.Error(fmt.Sprintln(eerr7))
	}

	err5 := attributeWorkflow.AddDecisionNode(pathDecisionNode,
		emptyNode1, emptyNode2)
	if err5 != nil {
		logger.Error(fmt.Sprintln(err5))
	}

	//Set start node for the search workflow
	attributeWorkflow.SetStartNode(pathDecisionNode)
	//Assign the workflow definition to the Orchestrator
	attributesOrchestrator := new(orchestrator.Orchestrator)
	attributesOrchestrator.Create(attributeWorkflow)

	logger.Info(attributesOrchestrator.String())
	logger.Info("Attributes Pipeline Created")
	return *attributesOrchestrator
}

func (a *AttributesUpdateApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(AttributesHealthCheck)
}

func (a *AttributesUpdateApi) Init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
	attributeUpdatePool = pool.NewWorkerSafe(ATTRIBUTE_UPDATE, ATTRIBUTE_UPDATE_POOL_SIZE, ATTRIBUTE_UPDATE_QUEUE_SIZE, ATTRIBUTE_UPDATE_RETRY_COUNT, ATTRIBUTE_UPDATE_WAIT_TIME)
	attributeUpdatePool.StartWorkers(optionUpdateWorker)

	// initialize cache config for attribute
	c := conf.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(c.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
