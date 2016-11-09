package post

import (
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"common/mongodb"
	"fmt"

	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"gopkg.in/mgo.v2/bson"
)

// CreateStandardSizeApi -> Basic orchestrator struct
type CreateStandardSizeApi struct {
}

// GetVersion -> Version Manager boilerplate
func (a *CreateStandardSizeApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.STANDARDSIZE_API,
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

// GetOrchestrator -> Orchestrator Boilerplate
func (a *CreateStandardSizeApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Standard Size Create Creation begin")

	standardSizeCreateOrchestrator := new(orchestrator.Orchestrator)
	standardSizeCreateWorkflow := new(orchestrator.WorkFlowDefinition)
	standardSizeCreateWorkflow.Create()

	//Creation of the nodes in the workflow definition
	standardSizeCreateNode := new(StandardSizeCreate)
	standardSizeCreateNode.SetID("StandardSize Create node")
	eerr := standardSizeCreateWorkflow.AddExecutionNode(standardSizeCreateNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	standardSizeErrorNode := new(StandardSizeErrorResponse)
	standardSizeErrorNode.SetID("Standard Size Create Error Response Node")
	eerr = standardSizeCreateWorkflow.AddExecutionNode(standardSizeErrorNode)
	if eerr != nil {
		logger.Error(fmt.Sprintln(eerr))
	}

	validationNode := new(StandardSizeCreateValidation)
	validationNode.SetID("Standard Size Create Validation Node")

	//Set start node for the search workflow
	standardSizeCreateWorkflow.AddDecisionNode(validationNode,
		standardSizeCreateNode, standardSizeErrorNode)

	//Set start node for the search workflow
	standardSizeCreateWorkflow.SetStartNode(validationNode)

	//Assign the workflow definition to the Orchestrator
	standardSizeCreateOrchestrator.Create(standardSizeCreateWorkflow)

	logger.Info(standardSizeCreateOrchestrator.String())
	logger.Info("Standard Size Create Pipeline Created")
	return *standardSizeCreateOrchestrator
}

// GetHealthCheck -> healthcheck Boilerplate
func (a *CreateStandardSizeApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(StandardSizeCreateHealthCheck)
}

// Init -> Basic Init function.
func (a *CreateStandardSizeApi) Init() {
	mgoSession := mongoFactory.GetMongoSession(common.STANDARDSIZE_CREATE)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection("counters")
	var counter mongodb.CounterInfo
	//initialize counter for standard size collection
	err := mgoObj.Find(bson.M{"_id": common.STANDARDSIZE_COLLECTION}).One(&counter)
	if err != nil {
		mgoSession.Insert("counters", mongodb.CounterInfo{Id: common.STANDARDSIZE_COLLECTION, SeqId: 0})
	}
	//initialize counter for standard size error collection
	err = mgoObj.Find(bson.M{"_id": common.STANDARDSIZEERROR_COLLECTION}).One(&counter)
	if err != nil {
		mgoSession.Insert("counters", mongodb.CounterInfo{Id: common.STANDARDSIZEERROR_COLLECTION, SeqId: 0})
	}
}
