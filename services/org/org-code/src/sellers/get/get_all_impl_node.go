package get

import (
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"sellers/common"
)

type GetAllSeller struct {
	id string
}

func (s *GetAllSeller) SetID(id string) {
	s.id = id
}

func (s GetAllSeller) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetAllSeller) Name() string {
	return "GET all seller"
}

func (s GetAllSeller) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET_ALL)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET_ALL)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_SELLER_GET_ALL)
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_GET_ALL, "Seller get all execution started")
	//getting limit and offset
	gs := SearchSeller{}
	limit, offset, err := gs.GetLimitOffset(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing limit and offset %s:", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid limit and offset", DeveloperMessage: err.Error()}
	}
	//getting all sellers from mongo
	resp, er := s.GetAll(limit, offset)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while getting all sellers from mongo : %v", er))
		return io, &florest_constants.AppError{Code: appconstant.ResourceNotFoundCode, Message: "Invalid Ids", DeveloperMessage: "Seller Ids do not exist"}
	}
	//setting data in RESULT
	io.IOData.Set(florest_constants.RESULT, resp)
	//getting count for all sellers
	count := 0
	var e error
	val, ok := utils.GetQueryParams(io, "count")
	if ok && val == "true" {
		count, e = s.SellerDataCount()
		if e != nil {
			logger.Error(fmt.Sprintf("Error while getting seller count from mongo : %v", e))
			return io, &florest_constants.AppError{Code: appconstant.ResourceNotFoundCode, Message: "Could not get Count from Mongo", DeveloperMessage: e.Error()}
		}
	} else {
		count = len(resp)
	}
	info := common.SetCountInMeta(io, count)
	io.IOData.Set(florest_constants.RESPONSE_META_DATA, info)
	return io, nil
}

//This function gets all seller data from Mongo given the limit and offset
func (s GetAllSeller) GetAll(limit int, offset int) ([]common.Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_GET_ALL)
	defer func() {
		logger.EndProfile(profiler, MONGO_GET_ALL)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	var org []common.Schema
	err := mgoObj.Find(bson.M{"ordrEml": bson.M{"$ne": "", "$exists": true}}).
		Sort("seqId").Limit(limit).Skip(offset).All(&org)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s GetAllSeller) SellerDataCount() (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_GET_SELLER_COUNT)
	defer func() {
		logger.EndProfile(profiler, MONGO_GET_SELLER_COUNT)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	count, err := mgoObj.Find(bson.M{"ordrEml": bson.M{"$ne": "", "$exists": true}}).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
