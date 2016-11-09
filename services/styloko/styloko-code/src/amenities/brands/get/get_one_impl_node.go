package get

import (
	"amenities/brands/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"encoding/json"
	//"encoding/json"
	"fmt"
	"strconv"

	"github.com/jabong/floRest/src/common/cache"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//Struct for Get Brand by ID(node based)
type GetBrand struct {
	id string
}

//Function to SetID for current node from orchestrator
func (b *GetBrand) SetID(id string) {
	b.id = id
}

//Function that returns current node ID to orchestrator
func (b GetBrand) GetID() (id string, err error) {
	return b.id, nil
}

//Function that returns node name to orchestrator
func (b GetBrand) Name() string {
	return "GET brand by id"
}

func (b GetBrand) One(id int) (common.Brand, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_GET)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_GET)
		mgoSession.Close()
	}()
	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	var brand common.Brand
	er := mgoObj.Find(bson.M{"seqId": id}).One(&brand)
	if er != nil {
		return brand, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          er.Error(),
			DeveloperMessage: "Brand Id doesnot exist in mongo ",
		}
	}
	return brand, er
}

func (b GetBrand) OneByName(name string) (common.Brand, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_GET)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_GET)
		mgoSession.Close()
	}()
	//Establishing Mongo Session
	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	var brand common.Brand
	er := mgoObj.Find(bson.M{"name": name}).One(&brand)
	return brand, er
}

//Function to start node execution
func (b GetBrand) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_GET)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_GET)
	}()

	//Enable logs for Debugging
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+b.Name(), rc)
	io.ExecContext.SetDebugMsg(common.BRAND_GET, "Single brand get execution started")

	//Retrieving BrandID
	data, _ := io.IOData.Get(common.BRAND_GET)
	id, err := strconv.Atoi(data.(string))
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting id to int :%s", err.Error()))
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Invalid Brand Id datatype",
			DeveloperMessage: "Invalid Brand Id",
		}

	}

	//Check cache for the searched ID
	v, err := b.GetFromCache(data.(string))
	if err == nil {
		io.IOData.Set(florest_constants.RESULT, v)
		return io, nil
	}

	//Gets data from mongo for the id passed and returns error if ID was not found
	brand, err := b.One(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting brand : %s", err.Error()))
		return io, err
	}

	io.IOData.Set(florest_constants.RESULT, brand)
	//Setting searched ID to cache
	go b.SetCache(data.(string), brand)
	return io, nil
}

//Function to set brand information to cache
func (b GetBrand) SetCache(key string, data interface{}) {
	e, _ := json.Marshal(data)
	i := cache.Item{
		Key:   fmt.Sprintf("%s-%s", common.BRANDS, key),
		Value: string(e),
	}
	err := cacheObj.Set(i, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in setting cache: %s", err.Error()))
	}
}

//Function to get brand information from cache
func (b GetBrand) GetFromCache(key string) (interface{}, error) {
	item, err := cacheObj.Get(fmt.Sprintf("%s-%s", common.BRANDS, key), false, false)
	if err != nil {
		logger.Warning(fmt.Sprintf("%s %s", key, err.Error()))
		return nil, err
	}

	var v interface{}
	json.Unmarshal([]byte(item.Value.(string)), &v)
	return v, nil
}
