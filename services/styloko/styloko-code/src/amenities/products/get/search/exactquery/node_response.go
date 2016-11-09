package exactquery

import (
	proUtil "amenities/products/common"
	search "amenities/products/get/search"
	"errors"
	"fmt"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/monitor"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ResponseNodeFE struct {
	id string
}

func (cs *ResponseNodeFE) SetID(id string) {
	cs.id = id
}

func (cs ResponseNodeFE) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ResponseNodeFE) Name() string {
	return "ResponseNodeFE"
}

func (cs ResponseNodeFE) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Debug("Enter Response Node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_PRO_RESPONSE_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_PRO_RESPONSE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "FE Response Node")

	products, err := LoadDataExactQuery{}.GetProductData(io)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs ResponseNodeFE) Execute(): Cannot get products data",
			DeveloperMessage: err.Error(),
		}
	}
	//Get Query
	query, err := PrepareExactQuery{}.GetQuery(io)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs ResponseNodeFE) Execute(): Cannot get products data",
			DeveloperMessage: err.Error(),
		}
	}
	var response interface{}
	var perr error
	var errorCode constants.AppErrorCode
	cacheControl := true
	if query.IsSingle {
		response, perr = cs.PreparePDPResponse(
			products,
			query.Expanse,
			query.Visibility,
		)
		errorCode = constants.InvalidRequestUri
	} else {
		response, cacheControl, perr = cs.PrepareMultiResponse(
			products,
			query.Expanse,
			query.Visibility,
		)
		errorCode = constants.ResourceErrorCode
	}
	if perr != nil {
		logger.Error(perr)
		return io, &constants.AppError{
			Code:             errorCode,
			Message:          "Not found",
			DeveloperMessage: perr.Error(),
		}
	}
	if cacheControl {
		io.IOData.Set(constants.RESPONSE_HEADERS_CONFIG, config.GlobalAppConfig.ResponseHeaders)
	}
	io.IOData.Set(constants.RESULT, response)
	logger.Debug("Exit Response Node")
	return io, nil
}

func (cs ResponseNodeFE) PrepareMultiResponse(
	productData []search.ProductData,
	expanse string,
	visibilityType string,
) (
	interface{}, bool, error,
) {
	type Response struct {
		Visibility bool        `json:"visible"`
		Identifier string      `json:"identifier"`
		Data       interface{} `json:"data"`
	}
	response := []Response{}
	var cacheCount int
	for _, data := range productData {
		resp := Response{}
		resp.Visibility = data.Visibility
		resp.Identifier = data.Identifier
		if data.Cache {
			cacheCount += 1
			logger.Warning("[cache]Serving From Cache: " + data.Identifier)
			resp.Data = data.Data
		} else {
			logger.Warning("[mongo]Serving From Mongo: " + data.Identifier)
			product, ok := data.Data.(proUtil.Product)
			if !ok {
				logger.Error("(cs ResponseNodeFE)#PrepareMultiResponse: unable to assert product")
				return nil, false, errors.New("Product Assertion failed")
			}
			presp, cacheTTL, setNoCache := cs.PrepareResponse(product, expanse, visibilityType)
			resp.Data = presp
			//set cache
			data2Cache := proUtil.ProductCache{}
			data2Cache.Data = presp
			data2Cache.Visibility = data.Visibility

			go func(ttl int, setNoCache bool) {
				defer proUtil.RecoverHandler("SetCache")
				if setNoCache {
					monitor.GetInstance().Count(
						"counter_ProductGetMultiCacheFailed", 1, []string{"styloko"}, 1)
					return
				}
				logger.Warning("[set]Set in cache: " + product.SKU)
				search.Pc.SetBySku(product.SKU, expanse, visibilityType, data2Cache, ttl)
				search.Pc.Set(product.SeqId, expanse, visibilityType, data2Cache, ttl)
			}(cacheTTL, setNoCache)

		}
		response = append(response, resp)
	}
	logger.Warning(fmt.Sprintf("[count]Found [%d] skus in cache out of %d", cacheCount, len(productData)))
	if len(response) < 1 {
		return response, false, nil
	}
	return response, true, nil
}

func (cs ResponseNodeFE) PreparePDPResponse(
	productData []search.ProductData,
	expanse string,
	visibilityType string,
) (
	interface{}, error,
) {
	if len(productData) == 0 {
		return nil, errors.New("Product Does not Exist")
	}
	data := productData[0]
	if data.Cache {
		//in cache
		if !data.Visibility {
			return nil, errors.New("Product Not Visible")
		}
		return data.Data, nil
	}
	//not in cache
	product, ok := data.Data.(proUtil.Product)
	if !ok {
		logger.Error("(cs ResponseNodeFE)#PreparePDPResponse: unable to assert product")
		return nil, errors.New("Product Assertion failed")
	}
	resp, cacheTTL, setNoCache := cs.PrepareResponse(product, expanse, visibilityType)

	//set cache
	data2Cache := proUtil.ProductCache{}
	data2Cache.Data = resp
	data2Cache.Visibility = data.Visibility

	go func(ttl int, setNoCache bool) {
		if setNoCache {
			monitor.GetInstance().Count(
				"counter_ProductGetPDPCacheFailed", 1, []string{"styloko"}, 1)
			return
		}
		defer proUtil.RecoverHandler("SetCache")
		search.Pc.SetBySku(product.SKU, expanse, visibilityType, data2Cache, ttl)
		search.Pc.Set(product.SeqId, expanse, visibilityType, data2Cache, ttl)
	}(cacheTTL, setNoCache)

	if !data.Visibility {
		return nil, errors.New("Product Not Visible")
	}
	return resp, nil
}

func (cs ResponseNodeFE) PrepareResponse(p proUtil.Product,
	expanse string, visibilityType string) (interface{}, int, bool) {
	response := proUtil.ProductResponse{}
	response.Product = &p
	response.VisibilityType = visibilityType

	return response.GetResponse(expanse), response.CacheTTL, response.SetNoCache
}
