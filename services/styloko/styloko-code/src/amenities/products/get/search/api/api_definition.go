package api

import (
	proUtil "amenities/products/common"
	searchUtil "amenities/products/get/search"
	exactQ "amenities/products/get/search/exactquery"
	sc "amenities/products/get/search/sellercenter"
	"common/appconfig"

	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type ProductsApi struct {
}

func (a *ProductsApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: proUtil.PRODUCT_API,
		Version:  "V1",
		Action:   "GET",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *ProductsApi) GetOrchestrator() orchestrator.Orchestrator {

	//Initiate All required nodes
	isSC := new(sc.IsSC)
	isSC.SetID("Is SC call")

	validateNodeSC := new(sc.ValidateNodeSC)
	validateNodeSC.SetID("ValidateSCRequest")

	fetchDataNodeSC := new(sc.FetchDataNodeSC)
	fetchDataNodeSC.SetID("FetchDataSC")

	responseNodeSC := new(sc.ResponseNodeSC)
	responseNodeSC.SetID("ResponseSC")

	emptyNode := new(EmptyNode)
	emptyNode.SetID("EmptyNode")

	emptyNode2 := new(EmptyNode)
	emptyNode2.SetID("EmptyNode2")

	isExactQuery := new(exactQ.IsExactQuery)
	isExactQuery.SetID("IsExactQuery")

	sellerSkuApi := new(exactQ.SellerSkuApi)
	sellerSkuApi.SetID("SellerSkuApi")

	prepareExactQuery := new(exactQ.PrepareExactQuery)
	prepareExactQuery.SetID("prepareExactQuery")

	cacheGet := new(exactQ.CacheGet)
	cacheGet.SetID("cacheGet")

	loadData := new(exactQ.LoadDataExactQuery)
	loadData.SetID("loadData")

	publish := new(exactQ.PublishNode)
	publish.SetID("publish")

	visibilityCheck := new(exactQ.VisibilityCheck)
	visibilityCheck.SetID("visibility")

	responseNodeFe := new(exactQ.ResponseNodeFE)
	responseNodeFe.SetID("responseNodeFe")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	//            Is SC Request
	//             SC   |   FE
	//      validateSC            Is Exact Query
	//           |                     |
	//      Fetch Data          Prepare query
	//           |              Cache Get
	//      Prepare Resp        Load data
	//							Publish
	//                          Visibility Test
	//                          Response

	err := workflow.AddDecisionNode(isSC, validateNodeSC, emptyNode)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddDecisionNode(isExactQuery, prepareExactQuery, sellerSkuApi)
	if err != nil {
		logger.Error(err)
	}

	//add execution nodes
	workflow.AddExecutionNode(responseNodeSC)
	workflow.AddExecutionNode(fetchDataNodeSC)

	workflow.AddExecutionNode(cacheGet)
	workflow.AddExecutionNode(loadData)
	workflow.AddExecutionNode(publish)
	workflow.AddExecutionNode(visibilityCheck)
	workflow.AddExecutionNode(responseNodeFe)

	//add connections
	err = workflow.AddConnection(validateNodeSC, fetchDataNodeSC)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(fetchDataNodeSC, responseNodeSC)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(emptyNode, isExactQuery)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(prepareExactQuery, cacheGet)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(cacheGet, loadData)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(loadData, publish)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(publish, visibilityCheck)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(visibilityCheck, responseNodeFe)
	if err != nil {
		logger.Error(err)
	}
	//set start node
	workflow.SetStartNode(isSC)

	orchestrator := new(orchestrator.Orchestrator)
	orchestrator.Create(workflow)
	return *orchestrator
}

func (a *ProductsApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(ProductSearchHealthCheck)
}

func (a *ProductsApi) Init() {

	// store config into global object
	searchUtil.Conf = config.ApplicationConfig.(*appconfig.AppConfig)
	searchUtil.DbAdapterName = searchUtil.Conf.DbAdapter
	searchUtil.SellerSkuLimit = searchUtil.Conf.SellerSkuLimit
	// initialize cache object
	var err error
	searchUtil.Pc.CacheObj, err = cache.Get(searchUtil.Conf.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
}
