package stock

import (
	productUpdaters "amenities/products/put/updaters"
	"common/redis"
	"common/utils"
	"encoding/json"
	_ "fmt"
	constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"io/ioutil"
	"strings"
)

type StockNode struct {
	id string
}

func (cs *StockNode) SetID(id string) {
	cs.id = id
}

func (cs StockNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs StockNode) Name() string {
	return "StockNode"
}

func (cs StockNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	httpReq, err := utils.GetRequestFromIO(io)
	if err != nil {
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "cannot read request data",
			DeveloperMessage: err.Error(),
		}
	}
	var stockReq StockRequest
	httpReqOrig := httpReq.OriginalRequest
	data, _ := ioutil.ReadAll(httpReqOrig.Body)
	logger.Warning(string(data))

	err = json.Unmarshal(data, &stockReq)
	if err != nil {
		logger.Error(err.Error())
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "Data UnMarshalling Failed",
			DeveloperMessage: err.Error(),
		}
	}
	client, err := redis.GetDriver()
	if err != nil {
		logger.Error(err.Error())
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "Cannot connect to Redis",
			DeveloperMessage: err.Error(),
		}
	}
	newStock := stockReq.GetStock()
	oldStock, rerr := client.GetSetInt(REDIS_STOCK_KEY+stockReq.Simple, newStock)
	if rerr != nil {
		logger.Error(rerr.Error())
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "Cannot fetch old stock from Redis",
			DeveloperMessage: rerr.Error(),
		}
	}
	//check if we need to bypass cache invalidate.
	var bypass bool
	if (newStock > 0) && (oldStock > 0) {
		bypass = true
	}

	cacheInv := productUpdaters.CacheInvalidate{
		Id:   stockReq.Simple,
		Type: productUpdaters.CACHE_INV_SIMPLEID,
	}
	errSlice := cacheInv.Validate()
	if errSlice != nil {
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "validation failed",
			DeveloperMessage: strings.Join(errSlice, ";"),
		}
	}
	//Invalidate cache only on stock in and out.
	if !bypass {
		err = cacheInv.InvalidateCache()
		if err != nil {
			logger.Error(err)
			return io, &constants.AppError{Code: constants.ResourceErrorCode,
				Message:          "Cache Invalidation Failed",
				DeveloperMessage: err.Error(),
			}
		}
	}
	//Publish message to bus.
	err = cacheInv.Publish()
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "Product Publish Failed",
			DeveloperMessage: err.Error(),
		}
	}
	io.IOData.Set(constants.RESULT, "success")
	return io, nil
}
