package exactquery

import (
	proUtil "amenities/products/common"
	search "amenities/products/get/search"
	"common/utils"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/monitor"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
	"strings"
)

//
// Get data from cache based on the query
//
type CacheGet struct {
	id string
}

func (cs *CacheGet) SetID(id string) {
	cs.id = id
}

func (cs CacheGet) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs CacheGet) Name() string {
	return "CacheGetNode"
}

func (cs CacheGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Debug("Enter Cache Get Node")

	//Check if cache to be used.
	if cs.IsNoCache(io) {
		return io, nil
	}

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_CACHE_GET_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_CACHE_GET_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Cache Get Node")

	peq := PrepareExactQuery{}
	query, er := peq.GetQuery(io)
	if er != nil {
		logger.Error(er)
		return io, &constants.AppError{
			Code:    constants.IncorrectDataErrorCode,
			Message: er.Error(),
		}
	}
	cacheData := proUtil.ProductCacheCollection{}
	var totalPros int
	if query.Id != nil {
		totalPros = len(query.Id)
		data, err := search.Pc.GetAll(query.Id, query.Expanse, query.Visibility)
		if err == nil {
			for k, v := range data {
				if v != nil {
					cacheData[strconv.Itoa(k)] = *v
				}
			}
		}
	} else if query.Sku != nil {
		totalPros = len(query.Sku)
		data, err := search.Pc.GetAllBySku(query.Sku, query.Expanse, query.Visibility)
		if err == nil {
			for k, v := range data {
				if v != nil {
					cacheData[k] = *v
				}
			}
		}
	} else {
		//wrong parameters
		return io, &constants.AppError{
			Code:    constants.IncorrectDataErrorCode,
			Message: "(cs CacheGet) Execute(): Invalid Data"}
	}
	//Metrics for Cache get count.
	go func(total int, cached int) {
		defer utils.RecoverHandler("Publish metric: counter_ProductGet")
		monitor.GetInstance().Count(
			"counter_ProductGet", int64(total), []string{"styloko"}, 1,
		)
		monitor.GetInstance().Count(
			"counter_ProductGetCache", int64(cached), []string{"styloko"}, 1,
		)
	}(totalPros, len(cacheData))
	//Metrics for Cache get count [ends].

	io.IOData.Set(search.CACHE_DATA, cacheData)
	io.ExecContext.SetDebugMsg("cache", cacheData.ToStringKeys())
	logger.Debug("Exit Cache Get Node")
	return io, nil
}

func (cs CacheGet) GetCacheData(io workflow.WorkFlowData) (proUtil.ProductCacheCollection, error) {
	q, err := io.IOData.Get(search.CACHE_DATA)
	if err != nil {
		return proUtil.ProductCacheCollection{},
			fmt.Errorf("(cs CacheGet)#GetCacheData: %s", err.Error())
	}
	data, ok := q.(proUtil.ProductCacheCollection)
	if !ok {
		return proUtil.ProductCacheCollection{}, errors.New("(cs CacheGet)#GetCacheData:Assertion failed")
	}
	return data, nil
}

func (cs CacheGet) IsNoCache(io workflow.WorkFlowData) bool {
	rp, _ := io.IOData.Get(constants.REQUEST)
	appHttpReq, pOk := rp.(*utilhttp.Request)
	if !pOk || appHttpReq == nil {
		return false
	}
	nocache := appHttpReq.OriginalRequest.Header.Get(proUtil.HEADER_NOCACHE)
	publish := appHttpReq.OriginalRequest.Header.Get(proUtil.HEADER_PUBLISH)
	if strings.ToLower(nocache) == "true" || strings.ToLower(publish) == "true" {
		return true
	}
	return false
}
