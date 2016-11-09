package bootstrap

import (
	proUtil "amenities/products/common"
	"amenities/products/jbus"
	factory "common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

type BootstrapNode struct {
	id string
}

func (cs *BootstrapNode) SetID(id string) {
	cs.id = id
}

func (cs BootstrapNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs BootstrapNode) Name() string {
	return "BootstrapNode"
}

func (cs BootstrapNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter bootstrap node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.BOOTSTRAP_NODE)
	defer logger.EndProfile(profiler, proUtil.BOOTSTRAP_NODE)

	transactionId := jbus.GenerateTransactionId()
	header, _ := utils.GetRequestHeader(io, KILL_JOB_HEADER)
	if strings.ToLower(header) == "true" {
		cs.setSolrFlag(KILL_JOB_FLAG, "true")
		cs.deleteSolrFlag(PREVIOUS_JOB_FLAG)
		io.IOData.Set(constants.RESULT, "Kill signal sent to running job.")
		return io, nil
	}
	previousJob := cs.getSolrFlag(PREVIOUS_JOB_FLAG)
	if !previousJob {
		cs.setSolrFlag(PREVIOUS_JOB_FLAG, "true")
		go func() {
			defer proUtil.RecoverHandler("BootstrapJob")
			cs.StartBootStrapping(transactionId)
		}()
		io.IOData.Set(constants.RESULT, "Bootstrapping products")
	} else {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "A bootstrap job is already running.", DeveloperMessage: "Cannot run multiple bootstrap jobs"}
	}
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Bootstrap Node")
	logger.Info("Exit bootstrap node")
	return io, nil
}

func (cs BootstrapNode) StartBootStrapping(transactionId string) {
	var limit = poolQueueSize
	var offset int
	firstRun := true
	processedCount.resetCounter()
	for true {
		killJob := cs.getSolrFlag(KILL_JOB_FLAG)
		if killJob {
			cs.deleteSolrFlag(KILL_JOB_FLAG)
			break
		}
		count := processedCount.getCounter()
		if limit == count || firstRun {
			processedCount.resetCounter()
			firstRun = false
			pros, err := cs.FetchProducts(limit, offset)
			offset += limit
			if err == mgo.ErrNotFound {
				break
			}
			if err != nil {
				logger.Error(err)
				break
			}
			for _, pro := range pros {
				workerPool.StartJob(Job{
					ProductId:        pro.Id,
					TransactionId:    transactionId,
					PublishInvisible: false,
				},
				)
			}
			if len(pros) < limit {
				break
			}
		} else {
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}
	cs.deleteSolrFlag(PREVIOUS_JOB_FLAG)
	logger.Info("Added all products to worker pool.")
}

func (cs BootstrapNode) FetchProducts(limit int, offset int) ([]proUtil.ProductSmall, error) {
	logger.Warning(fmt.Sprintf("fetching products %d, %d", offset, limit))
	pros := []proUtil.ProductSmall{}
	mgoSess := factory.GetMongoSession("BOOTSTRAP")
	defer mgoSess.Close()
	err := mgoSess.SetCollection(
		proUtil.PRODUCT_COLLECTION,
	).Find(
		proUtil.M{
			"status":      proUtil.STATUS_ACTIVE,
			"petApproved": 1,
		},
	).Select(
		proUtil.M{
			"seqId": 1, "sku": 1, "_id": 0,
		},
	).Sort("-createdAt").Skip(offset).Limit(limit).All(&pros)
	if err != nil && err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf(
			"(cs BootstrapNode)#FetchProducts[%d, %d] : %s",
			offset,
			limit,
			err.Error(),
		)
	}
	return pros, nil
}

// Redis key management functions below

// getSolrFlag returns bool for provided key from redis hMap
func (cs BootstrapNode) getSolrFlag(key string) bool {
	redisAdapter, err := factory.GetRedisDriver(proUtil.BOOTSTRAP_DRIVER)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(cs BootstrapNode)#getSolrFlag: Cannot acquire redis driver: %s",
			err.Error(),
		))
	}
	adapterResponse := redisAdapter.HGetAllMap(proUtil.BOOTSTRAP_HASH_MAP)
	flags, err := adapterResponse.Result()
	if _, ok := flags[key]; ok {
		return true
	}
	return false
}

// setSolrFlag sets redis flag for solr
func (cs BootstrapNode) setSolrFlag(key, value string) {
	redisAdapter, err := factory.GetRedisDriver(proUtil.BOOTSTRAP_DRIVER)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(cs BootstrapNode)#setSolrFlag(): %s",
			err.Error(),
		))
	}
	redisAdapter.HSetNX(proUtil.BOOTSTRAP_HASH_MAP, key, value)
}

// deleteSolrFlag deletes the key from redis hMap
func (cs BootstrapNode) deleteSolrFlag(key string) {
	redisAdapter, err := factory.GetRedisDriver(proUtil.BOOTSTRAP_DRIVER)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(cs BootstrapNode)#deleteSolrFlag: %s",
			err.Error(),
		))
	}
	redisAdapter.HDel(proUtil.BOOTSTRAP_HASH_MAP, key)
}

const (
	KILL_JOB_FLAG     = "killJob"
	PREVIOUS_JOB_FLAG = "previousJob"
	KILL_JOB_HEADER   = "Kill-Job"
)
