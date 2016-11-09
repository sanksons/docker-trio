package get

import (
	"amenities/brands/common"
	"common/appconfig"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

var cacheObj cache.CacheInterface

//Base struct for Brand Get API
type BrandAPI struct {
}

//Function to return Version manager instance
func (a *BrandAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.BRAND_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

//Function to return the Orchestrator
func (a *BrandAPI) GetOrchestrator() orchestrator.Orchestrator {

	//Pictorial description of Nodes

	//       startDecisionNode
	// 			   /\
	// brandGetNode  searchOrGetAllDecisionNode
	// 					/\
	// 	  brandGetAllNode brandSearchNode

	logger.Info("Brand Get API pipeline begins")

	brandGetWorkflow := new(orchestrator.WorkFlowDefinition)
	brandGetWorkflow.Create()

	//Creation of the nodes in the workflow definition
	brandGetNode := new(GetBrand)
	brandGetNode.SetID("Brand getSingle api")
	err := brandGetWorkflow.AddExecutionNode(brandGetNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	brandSearchNode := new(SearchBrand)
	brandSearchNode.SetID("Brand search api")
	err = brandGetWorkflow.AddExecutionNode(brandSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	brandGetAllNode := new(GetAllBrands)
	brandGetAllNode.SetID("Brand getAll api")
	err = brandGetWorkflow.AddExecutionNode(brandGetAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	searchOrGetAllDecisionNode := new(GetAllDecision)
	searchOrGetAllDecisionNode.SetID("Brand get All or Search decision node")

	//Describing connection between Nodes
	err = brandGetWorkflow.AddDecisionNode(searchOrGetAllDecisionNode, brandGetAllNode, brandSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Path based decision, Starting Node
	startDecisionNode := new(GetOneDecision)
	startDecisionNode.SetID("Brand get api")

	//Describing connection between Nodes
	err = brandGetWorkflow.AddDecisionNode(startDecisionNode, brandGetNode, searchOrGetAllDecisionNode)

	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	brandGetWorkflow.SetStartNode(startDecisionNode)
	brandGetOrchestrator := new(orchestrator.Orchestrator)

	//Assign the workflow definition to the Orchestrator
	brandGetOrchestrator.Create(brandGetWorkflow)
	logger.Info(brandGetOrchestrator.String())
	logger.Info("Brand Search Pipeline Created")
	return *brandGetOrchestrator
}

//Function that returns HealthCheckInterface
func (bg *BrandAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(BrandsGetHealthCheck)
}

//Function to initialise API
func (bg *BrandAPI) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
