package get

import (
	"amenities/brands/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"fmt"
	"strconv"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//Struct for GetAll Brands(node based)
type GetAllBrands struct {
	id string
}

//Function to SetID for current node from orchestrator
func (b *GetAllBrands) SetID(id string) {
	b.id = id
}

//Function that returns current node ID to orchestrator
func (b GetAllBrands) GetID() (id string, err error) {
	return b.id, nil
}

//Function that returns node name to orchestrator
func (b GetAllBrands) Name() string {
	return "GET all Brands"
}

//Function to start node execution
func (b GetAllBrands) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRANDS_GET_ALL)
	defer func() {
		logger.EndProfile(profiler, common.BRANDS_GET_ALL)
	}()

	//Enable logs for Debugging
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+b.Name(), rc)
	io.ExecContext.SetDebugMsg(common.BRANDS_GET_ALL, "Brand get all execution started")

	//Reading path params
	pathParams, _ := utils.GetPathParams(io)
	if len(pathParams) == 0 {
		logger.Info("No path params, getAll started. Looking for query params")
		limit, offset, err := b.GetLimitOffset(io)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while parsing limit and offset %s:", err.Error()))
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid limit and offset", DeveloperMessage: err.Error()}
		}
		//Reading Query params on the basis of status
		query, ok := utils.GetQueryParams(io, common.STATUS)
		if !ok {

			//Finding All brands.
			brandStruct, ok := b.BrandFindByStatus(common.ALL, limit, offset)
			if !ok {
				return io, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode, Message: "No data found.", DeveloperMessage: "Cannot find any active categories."}
			}

			io.IOData.Set(florest_constants.RESULT, brandStruct)
		} else {
			//Finding brands on the basis of status
			brandStruct, ok := b.BrandFindByStatus(query, limit, offset)
			if !ok {
				return io, &florest_constants.AppError{
					Code:             appconstant.InvalidDataErrorCode,
					Message:          "Invalid Brand Status",
					DeveloperMessage: "Status param value in query is Invalid"}
			}
			io.IOData.Set(florest_constants.RESULT, brandStruct)
		}
		//getting count for all brands
		count, _ := b.BrandDataCount()
		//setting count in META
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
		io.IOData.Set(florest_constants.RESPONSE_META_DATA, info)
		return io, nil
	}

	return io, &florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          "Invalid path provided",
		DeveloperMessage: "Invalid path provided"}

}

//Function to find Brands based on the status field
func (b GetAllBrands) BrandFindByStatus(status string, limit int, offset int) ([]common.Brand, bool) {

	//Enable profiling
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRANDS_GET_ALL)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRANDS_GET_ALL)
		mgoSession.Close()
	}()

	var brandStruct []common.Brand
	if status == common.ACTIVE || status == common.INACTIVE || status == common.DELETED {
		logger.Info("Get " + status + "brands finished")
		mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
		if limit == 0 {
			mgoObj.Find(bson.M{"status": status}).Sort("-seqId").All(&brandStruct)
		} else {
			mgoObj.Find(bson.M{"status": status}).Sort("-seqId").Limit(limit).Skip(offset).All(&brandStruct)
		}
		return brandStruct, true
	}
	if status == common.ALL {
		logger.Info("Get " + status + "brand finished")
		mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
		if limit == 0 {
			mgoObj.Find(nil).Sort("-seqId").All(&brandStruct)
		} else {
			mgoObj.Find(nil).Sort("-seqId").Limit(limit).Skip(offset).All(&brandStruct)
		}
		return brandStruct, true
	}
	return brandStruct, false
}

func (b GetAllBrands) BrandDataCount() (int, error) {
	//Enable profiling
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRANDS_GET_ALL)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRANDS_GET_ALL)
		mgoSession.Close()
	}()

	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	count, err := mgoObj.Find(nil).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b GetAllBrands) GetLimitOffset(io workflow.WorkFlowData) (int, int, error) {
	//getting limit from query params
	//else setting default at 1000
	lim, ok := utils.GetQueryParams(io, "limit")
	if !ok {
		logger.Info("No limit passed. Setting default limit.")
		lim = "0"
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
