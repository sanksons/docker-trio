package api

import (
	proUtils "amenities/products/common"
	put "amenities/products/put"
	nodes "amenities/products/put/nodes"
	"common/appconfig"
	"common/constants"
	"common/pool/tasker"
	"strconv"
	"time"

	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/versionmanager"
	validator "gopkg.in/go-playground/validator.v8"
)

//
// API definitionn for Product Update
//
type ProductAPI struct {
}

func (a *ProductAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: proUtils.PRODUCT_API,
		Version:  "V1",
		Action:   "PUT",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *ProductAPI) GetOrchestrator() orchestrator.Orchestrator {
	validateUpdate := new(nodes.ValidateUpdate)
	validateUpdate.SetID("RequestValidator")

	updateMongo := new(nodes.UpdateMongo)
	updateMongo.SetID("UpdateMongo")

	invalidateCache := new(nodes.InvalidateCache)
	invalidateCache.SetID("InvalidateCache")

	responseNode := new(nodes.ResponseNode)
	responseNode.SetID("ResponseNode")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()

	//
	// WorkFlow:
	//
	// Validate --> Update --> InvalidateCache --> Response
	//
	//set start node
	workflow.AddExecutionNode(validateUpdate)
	workflow.AddExecutionNode(updateMongo)
	workflow.AddExecutionNode(invalidateCache)
	workflow.AddExecutionNode(responseNode)

	err := workflow.AddConnection(validateUpdate, updateMongo)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(updateMongo, invalidateCache)
	if err != nil {
		logger.Error(err)
	}
	err = workflow.AddConnection(invalidateCache, responseNode)
	if err != nil {
		logger.Error(err)
	}
	workflow.SetStartNode(validateUpdate)

	Orchestrator := new(orchestrator.Orchestrator)
	Orchestrator.Create(workflow)
	return *Orchestrator
}

func (a *ProductAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(ProductUpdateHealthCheck)
}

func (a *ProductAPI) Init() {
	//Store config into global object
	var err error
	put.Conf = config.ApplicationConfig.(*appconfig.AppConfig)
	put.DbAdapterName = put.Conf.DbAdapter
	//Initialize cache object
	put.CacheMngr.CacheObj, err = cache.Get(put.Conf.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
	//Validator config
	config := &validator.Config{TagName: "validate"}
	put.Validate = validator.New(config)
	proUtils.ReadAttributeFile()

	//Initiate Tasker
	//@todo: enable tasker.
	a.InitTasker()
}

//
// This method is called from make Test.
//
func (a *ProductAPI) InitTest() {
	//Store config into global object
	var err error
	put.Conf = config.ApplicationConfig.(*appconfig.AppConfig)
	put.DbAdapterName = put.Conf.DbAdapter
	//Initialize cache object
	put.CacheMngr.CacheObj, err = cache.Get(put.Conf.Cache)
	if err != nil {
		logger.Error(err.Error())
	}
	//Validator config
	config := &validator.Config{TagName: "validate"}
	put.Validate = validator.New(config)
	proUtils.ReadTestAttributeFile()
}

func (a *ProductAPI) InitTasker() {
	//Get config for Tasker
	var sleepTime time.Duration
	var fetchLimit int
	maxRoutinesStr := put.Conf.ProductSyncing.MaxRoutines
	maxRoutines, _ := strconv.Atoi(maxRoutinesStr)
	if maxRoutines > 0 {
		fetchLimit = maxRoutines
	} else {
		fetchLimit = 10
	}
	sleepTimeStr := put.Conf.ProductSyncing.SleepTime
	st, _ := strconv.Atoi(sleepTimeStr)
	if st > 0 {
		sleepTime = time.Duration(st)
	} else {
		sleepTime = 10
	}
	//initiateTasker
	tasker.Tasker{
		UUID:         "Product",
		ResourceType: constants.PRODUCT_RESOURCE_NAME,
		SleepTime:    sleepTime,
		FetchLimit:   fetchLimit,
	}.Initiate()
}
