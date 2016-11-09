package get

import (
	"common/appconfig"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	"net/http"
	"sellers/common"
	"time"
)

var cacheObj cache.CacheInterface
var client *http.Client

type GetSellerApi struct {
}

func (sg *GetSellerApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: common.SELLERS,
		Version:  "V1",
		Action:   "GET",
		BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE, //todo - should it be a constant
	}
}

func (sg *GetSellerApi) GetOrchestrator() orchestrator.Orchestrator {
	logger.Info("Seller Get API pipeline begins")

	sellerGetWorkflow := new(orchestrator.WorkFlowDefinition)
	sellerGetWorkflow.Create()

	//Creation of the nodes in the workflow definition
	sellerGetNode := new(GetSeller)
	sellerGetNode.SetID("Seller Get By Id Node")
	err := sellerGetWorkflow.AddExecutionNode(sellerGetNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sellerSearchNode := new(SearchSeller)
	sellerSearchNode.SetID("Seller Search Node")
	err = sellerGetWorkflow.AddExecutionNode(sellerSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	sellerGetAllNode := new(GetAllSeller)
	sellerGetAllNode.SetID("Seller Get All Node ")
	err = sellerGetWorkflow.AddExecutionNode(sellerGetAllNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	getAllDecisionNode := new(GetAllDecision)
	getAllDecisionNode.SetID("Seller Get All or Search Decision Node")
	err = sellerGetWorkflow.AddDecisionNode(getAllDecisionNode, sellerGetAllNode, sellerSearchNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	getOneDecisionNode := new(GetOneDecision)
	getOneDecisionNode.SetID("Seller Get One Decision Node")
	err = sellerGetWorkflow.AddDecisionNode(getOneDecisionNode, sellerGetNode, getAllDecisionNode)
	if err != nil {
		logger.Error(fmt.Sprintln(err))
	}

	//Set start node for the search workflow
	sellerGetWorkflow.SetStartNode(getOneDecisionNode)

	sellerGetOrchestrator := new(orchestrator.Orchestrator)
	sellerGetOrchestrator.Create(sellerGetWorkflow)
	logger.Info(sellerGetOrchestrator.String())
	return *sellerGetOrchestrator
}

func (sg *GetSellerApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(SellerGetHealthCheck)
}

func (sg *GetSellerApi) Init() {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while initializing cache object :%s", err.Error()))
	}
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 200,
		},
		Timeout: time.Duration(1000) * time.Millisecond,
	}
}
