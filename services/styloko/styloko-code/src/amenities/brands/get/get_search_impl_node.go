package get

import (
	"amenities/brands/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"errors"
	"fmt"
	"strconv"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//Struct for Search Brand(node based)
type SearchBrand struct {
	id string
}

//Function to SetID for current node from orchestrator
func (b *SearchBrand) SetID(id string) {
	b.id = id
}

//Function that returns current node ID to orchestrator
func (b SearchBrand) GetID() (id string, err error) {
	return b.id, nil
}

//Function that returns node name to orchestrator
func (b SearchBrand) Name() string {
	return "Search Brand"
}

//Function to start node execution
func (b SearchBrand) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	rc, err := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting context: %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "", DeveloperMessage: err.Error()}
	}
	logger.Info("entered "+b.Name(), rc)
	io.ExecContext.SetDebugMsg(common.BRAND_SEARCH, "Brand search execution started")
	data, err := io.IOData.Get(common.BRAND_SEARCH)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting search request data: %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting search request data", DeveloperMessage: err.Error()}
	}
	// parse query
	dataMap := data.(map[string]interface{})
	if len(dataMap) == 0 {
		logger.Error("Error while getting search params.")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String"}
	}
	//getting limit and offset
	limit, offset, err := b.GetLimitOffset(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing limit and offset %s:", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid limit and offset", DeveloperMessage: err.Error()}
	}
	//geting data for search parameters
	resp, err := b.GetDetailsFromMongo(dataMap, limit, offset)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in Getting Search Data : %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error in Getting Search Data", DeveloperMessage: err.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, resp)

	//getting count for all brands according to the search criteria
	count, e := b.GetBrandDataCount(dataMap)
	if e != nil {
		logger.Error(fmt.Sprintf("Error while getting seller count from mongo : %v", e))
		return io, &florest_constants.AppError{Code: appconstant.ResourceNotFoundCode, Message: "Could not get Count from Mongo", DeveloperMessage: e.Error()}
	}
	//setting count in META
	info := b.SetCountInMeta(io, count)
	io.IOData.Set(florest_constants.RESPONSE_META_DATA, info)
	return io, nil
}

//gets details from mongo based on the bson map passed
func (b SearchBrand) GetDetailsFromMongo(bsonMap map[string]interface{}, limit int, offset int) ([]common.Brand, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.MONGO_SEARCH)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.MONGO_SEARCH)
		mgoSession.Close()
	}()
	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	var brandInfo []common.Brand
	err := mgoObj.Find(bsonMap).Sort("-updtdAt").Limit(limit).Skip(offset).All(&brandInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting brand details from mongo: %s", err.Error()))
		return nil, err
	}
	if len(brandInfo) == 0 {
		logger.Error("No brand found against passed search criteria")
		return nil, errors.New("No brand found against passed search criteria")
	}
	return brandInfo, nil
}

//gets limit and offset if passed, if not sets limit=1000 and offset=0
func (b SearchBrand) GetLimitOffset(io workflow.WorkFlowData) (int, int, error) {
	//getting limit from query params
	//else setting default at 1000
	lim, ok := utils.GetQueryParams(io, "limit")
	if !ok {
		logger.Info("No limit passed. Setting default limit.")
		lim = "1000"
	}
	//converting limit from string to int
	limit, err := strconv.Atoi(lim)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing limit %v:", lim))
		return 0, 0, err
	}
	//getting offset from query params
	//else setting default at 0
	skip, ok := utils.GetQueryParams(io, "offset")
	if !ok {
		logger.Info("No offset passed.Setting default offset")
		skip = "0"
	}
	//converting offset to int
	offset, err := strconv.Atoi(skip)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing offset %v:", skip))
		return 0, 0, err
	}
	return limit, offset, nil
}

//Function gets brand data count for the datamap passed
func (b SearchBrand) GetBrandDataCount(dataMap map[string]interface{}) (int, error) {
	var err error
	var count int
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.MONGO_GET_BRAND_SEARCH_COUNT)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.MONGO_GET_BRAND_SEARCH_COUNT)
		mgoSession.Close()
	}()
	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	if _, ok := dataMap["seqId"]; ok {
		return 0, nil
	}
	for k, v := range dataMap {
		count, err = mgoObj.Find(bson.M{k: v}).Count()
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

//sets count in response meta data for florest io data
func (b SearchBrand) SetCountInMeta(io workflow.WorkFlowData, count int) *http.ResponseMetaData {
	var info *http.ResponseMetaData

	infoVal, merr := io.IOData.Get(florest_constants.RESPONSE_META_DATA)
	if merr != nil {
		info = http.NewResponseMetaData()
	} else if infoObj, ok := infoVal.(*http.ResponseMetaData); ok {
		info = infoObj
	} else {
		info = http.NewResponseMetaData()
	}
	info.ApiMetaData["count"] = count
	return info
}
