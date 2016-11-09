package get

import (
	categoryService "amenities/services/categories"
	prodUtils "amenities/services/products"
	sizeUtils "amenities/sizechart/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	utils "common/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jabong/floRest/src/common/cache"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type SizeChartGet struct {
	id string
}

func (n *SizeChartGet) SetID(id string) {
	n.id = id
}

func (n SizeChartGet) GetID() (id string, err error) {
	return n.id, nil
}

func (a SizeChartGet) Name() string {
	return "SizeChartGetNode"
}

type ProductResponse struct {
	SizeChart interface{}
	Sku       string
	Id        int
}

func (a SizeChartGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	reqParamsArray, err := utils.GetPathParams(io)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Unable to get Request Parameters",
			DeveloperMessage: err.Error(),
		}
	}
	reqParam := reqParamsArray[0]
	nocache, hErr := utils.GetRequestHeader(io, sizeUtils.CacheControlHeader)
	if hErr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Unable to get Request Header",
			DeveloperMessage: hErr.Error(),
		}
	}
	if strings.ToLower(nocache) != "true" {
		sizechart, cErr := a.GetfromCache(reqParam)
		if cErr == nil {
			io.IOData.Set(florest_constants.RESULT, sizechart)
			return io, nil
		}
	}

	productResponse, gErr := a.GetSizechartByProduct(reqParam)
	if gErr != nil {
		logger.Error(fmt.Sprintf("Error while getting sizechart : %s", gErr.Error()))
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Unable to get sizechart.",
			DeveloperMessage: gErr.Error(),
		}
	}
	// Process the sizechart to check active/inactive/deleted simples
	mErr := a.ManipulateSizeChart(productResponse.Sku, productResponse.SizeChart)
	if mErr != nil {
		logger.Error(fmt.Sprintf("%s", mErr.Error()))
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Unable to recalculate sizechart for inactive/deleted simples.",
			DeveloperMessage: mErr.Error(),
		}
	}

	io.IOData.Set(florest_constants.RESULT, productResponse.SizeChart)
	// Set data in cache
	go a.SetCache(productResponse.Sku, productResponse.SizeChart)
	go a.SetCache(utils.ToString(productResponse.Id), productResponse.SizeChart)

	return io, nil
}

func (a SizeChartGet) ManipulateSizeChart(sku string, sizechart interface{}) error {
	if sizechart == nil {
		return nil
	}
	inactiveDeletedSimpleSizes, err := prodUtils.GetInactiveDeletedSimpleSizes(sku, "mongo")
	if err != nil {
		return fmt.Errorf("#ManipulateSizeChart: could not get inactive/deleted simple sizes %s", err.Error())
	}
	if len(inactiveDeletedSimpleSizes) == 0 {
		return nil
	}
	var sortedData [][]interface{}

	var count int
	var result map[string]interface{}
	result = make(map[string]interface{})

	sizeChartMap := sizechart.(map[string]interface{})
	sizeData := sizeChartMap["sizes"]
	sizeDataMap := sizeData.(map[string]interface{})
	sortedData = make([][]interface{}, len(sizeDataMap))

	// convert map to sorted datastructure
	for mapIndex, mapData := range sizeDataMap {
		intIndex, err := utils.GetInt(mapIndex)
		if err != nil {
			return fmt.Errorf("ManipulateSizeChart: err in get int %s", err.Error())
		}
		sortedData[intIndex] = mapData.([]interface{})
	}
	for _, rowData := range sortedData {
		// check if this size is to be retained/remove from sizechart
		if !utils.InArrayString(inactiveDeletedSimpleSizes, rowData[0].(string)) {
			result[utils.ToString(count)] = rowData
			count++
		}
	}
	sizeChartMap["sizes"] = result
	return nil
}

func (a SizeChartGet) GetSizechartByProduct(reqParam interface{}) (ProductResponse, error) {
	var sku string
	mongoSession := mongo.GetMongoSession(sizeUtils.SizeChartGetAPI)
	defer func() {
		mongoSession.Close()
	}()
	mgoObj := mongoSession.SetCollection(sizeUtils.ProductsCollection)

	type Product struct {
		Sku             string      `bson:"sku"`
		SeqId           int         `bson:"seqId"`
		Data            interface{} `bson:"sizeChart"`
		PrimaryCategory int         `bson:"primaryCategory"`
	}
	var prod Product
	configId, err := utils.GetInt(reqParam)
	if err != nil {
		sku = utils.ToString(reqParam)
		err := mgoObj.Find(bson.M{"sku": sku}).
			Select(bson.M{"sku": 1, "seqId": 1, "sizeChart": 1, "primaryCategory": 1, "_id": 0}).
			One(&prod)
		if err != nil {
			return ProductResponse{nil, "", 0}, err
		}
	} else {
		// Check product by Id
		err := mgoObj.Find(bson.M{"seqId": configId}).
			Select(bson.M{"sku": 1, "seqId": 1, "sizeChart": 1, "primaryCategory": 1, "_id": 0}).
			One(&prod)
		if err != nil {
			return ProductResponse{nil, "", 0}, err
		}
	}
	// Prepare sizechrt data from mongoresult
	type ProductSizeChart struct {
		Id   int         `json:"id"`
		Data interface{} `json:"data"`
	}
	var sizechart ProductSizeChart
	mapBytes, mErr := json.Marshal(prod.Data)
	if mErr != nil {
		return ProductResponse{nil, "", 0}, mErr
	}
	uErr := json.Unmarshal(mapBytes, &sizechart)
	if uErr != nil {
		return ProductResponse{nil, "", 0}, uErr
	}

	// Check if sizechart is active over categry to which product belongs
	catData := categoryService.ById(prod.PrimaryCategory)
	if catData.SizeChartAcive == 1 {
		if sizechart.Data == nil {
			return ProductResponse{sizechart.Data, prod.Sku, prod.SeqId}, fmt.Errorf("Sizechart doesnot exist for product")
		}
		return ProductResponse{sizechart.Data, prod.Sku, prod.SeqId}, nil
	}
	// sizechart is not to be served, but sku exists.
	return ProductResponse{nil, prod.Sku, prod.SeqId}, fmt.Errorf("Sizechart is not active for product")
}

func (a SizeChartGet) GetfromCache(key string) (interface{}, error) {
	item, err := cacheObj.Get(fmt.Sprintf("%s-%s", sizeUtils.SizeChartCacheKey, key), false, false)
	if err != nil {
		logger.Warning(fmt.Sprintf("%s %s", key, err.Error()))
		return nil, err
	}

	var v interface{}
	json.Unmarshal([]byte(item.Value.(string)), &v)
	return v, nil
}

func (a SizeChartGet) SetCache(key string, data interface{}) {
	e, _ := json.Marshal(data)
	i := cache.Item{
		Key:   fmt.Sprintf("%s-%s", sizeUtils.SizeChartCacheKey, key),
		Value: string(e),
	}
	err := cacheObj.Set(i, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in setting cache: %s", err.Error()))
	}
}
