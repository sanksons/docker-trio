package exactquery

import (
	proUtil "amenities/products/common"
	"amenities/products/common/prodsizechart"
	search "amenities/products/get/search"
	"errors"
	_ "reflect"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type LoadDataExactQuery struct {
	id string
}

func (cs *LoadDataExactQuery) SetID(id string) {
	cs.id = id
}

func (cs LoadDataExactQuery) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs LoadDataExactQuery) Name() string {
	return "LoadDataExactQuery"
}

func (cs LoadDataExactQuery) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Info("Enter load data node")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_LOAD_DATA_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_LOAD_DATA_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "LoadDataExactQuery")
	query, err := PrepareExactQuery{}.GetQuery(io)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs LoadDataExactQuery) Execute(): Cannot load query data",
			DeveloperMessage: err.Error(),
		}

	}
	cachedata, _ := CacheGet{}.GetCacheData(io)

	productsData, err := cs.ProcessExactQuery(query, cachedata)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs LoadDataExactQuery) Execute(): Cannot prepare products data",
			DeveloperMessage: err.Error(),
		}
	}
	io.IOData.Set(search.PRODUCTDATA, productsData)
	logger.Info("Exit load data node")
	return io, nil
}

func (cs LoadDataExactQuery) GetProductData(io workflow.WorkFlowData) ([]search.ProductData, error) {
	q, err := io.IOData.Get(search.PRODUCTDATA)
	if err != nil {
		return nil, err
	}
	data, ok := q.([]search.ProductData)
	if !ok {
		return nil, errors.New("(cs LoadDataExactQuery)#GetProductData:Assertion failed")
	}
	return data, nil
}

func (cs LoadDataExactQuery) ProcessExactQuery(
	query search.ExactQuery,
	cachedata proUtil.ProductCacheCollection,
) ([]search.ProductData, error) {

	var response []search.ProductData
	var err error
	if query.Id != nil {
		response, err = cs.prepareByIds(query.Id, cachedata)
	} else {
		response, err = cs.prepareBySkus(query.Sku, cachedata)
	}
	return response, err
}

func (cs LoadDataExactQuery) prepareByIds(ids []int, cachedata proUtil.ProductCacheCollection) ([]search.ProductData, error) {
	var response []search.ProductData
	var tobeloaded []int
	for _, id := range ids {
		if _, ok := cachedata[strconv.Itoa(id)]; !ok {
			tobeloaded = append(tobeloaded, id)
		}
	}
	procoll := &proUtil.ProductCollection{}
	if len(tobeloaded) > 0 {
		//load from DB
		err := procoll.LoadByIds(tobeloaded)
		if err != nil {
			return response, err
		}
	}
	//parse Ids to load data from cache or db
	for _, id := range ids {
		pData := search.ProductData{}
		//check if exists in cache
		if val, ok := cachedata[strconv.Itoa(id)]; ok {
			//load from cache
			pData.Data = val.Data
			pData.Cache = true
			pData.Visibility = val.Visibility
			pData.Identifier = strconv.Itoa(id)
			response = append(response, pData)
		} else {
			//load from db
			for _, p := range procoll.Products {
				if p.SeqId == id {
					// Manipulate sizechart and return the result
					prodsizechart.CalculateSizeChart(p)
					pData.Data = p
					pData.Cache = false
					pData.Identifier = strconv.Itoa(id)
					response = append(response, pData)
					break
				}
			}
		}
	}
	return response, nil
}

func (cs LoadDataExactQuery) prepareBySkus(skus []string, cachedata proUtil.ProductCacheCollection) ([]search.ProductData, error) {
	var response []search.ProductData
	var tobeloaded []string
	for _, sku := range skus {
		if _, ok := cachedata[strings.ToLower(sku)]; !ok {
			tobeloaded = append(tobeloaded, sku)
		}
	}
	procoll := &proUtil.ProductCollection{}
	if len(tobeloaded) > 0 {
		//load from DB
		err := procoll.LoadBySkus(tobeloaded)
		if err != nil {
			return response, err
		}
	}
	for _, sku := range skus {
		pData := search.ProductData{}
		//check if exists in cache
		if val, ok := cachedata[strings.ToLower(sku)]; ok {
			//load from cache
			pData.Data = val.Data
			pData.Cache = true
			pData.Visibility = val.Visibility
			pData.Identifier = sku
			response = append(response, pData)
		} else {
			//load from db
			for _, p := range procoll.Products {
				if p.SKU == sku {
					// Manipulate sizechart and return the result
					prodsizechart.CalculateSizeChart(p)
					pData.Data = p
					pData.Cache = false
					pData.Identifier = sku
					response = append(response, pData)
					break
				}
			}
		}
	}
	return response, nil
}
