package sync

import (
	"common/appconstant"
	"common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"sellers/common"
	"strconv"
)

type StartSync struct {
	id string
}

func (s *StartSync) SetID(id string) {
	s.id = id
}

func (s StartSync) GetID() (id string, err error) {
	return s.id, nil
}

func (s StartSync) Name() string {
	return "StartSync seller by id"
}

func (s StartSync) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_START_SYNC)
	defer func() {
		logger.EndProfile(profiler, SELLER_START_SYNC)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_START_SYNC, "Seller StartSync execution started")
	logger.Info("StartSync Seller Started")

	val, ok := utils.GetQueryParams(io, "ids")
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid Query Params", DeveloperMessage: "No query param passed for ids"}
	}
	ids := common.ParseQuery(val)
	bsonMap := map[string]interface{}{"seqId": bson.M{"$in": ids}}
	slrData, err := common.GetDetailsFromMongo(bsonMap, 1000, 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in GetDetailsFromMongo :%s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error in GetDetailsFromMongo", DeveloperMessage: err.Error()}
	}
	dataMap, err := utils.ConvertStructArrToMapArr(slrData)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting data []struct to []map[string]interface{} :%s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while converting data []struct to []map[string]interface{}", DeveloperMessage: err.Error()}
	}
	errMap := make(map[string]interface{}, 0)
	for _, v := range dataMap {
		err := common.UpdateInMysql(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while updating in MySql :%s", err.Error()))
			errMap[strconv.Itoa(int(v["seqId"].(float64)))] = err.Error()
			continue
		}
		errMap[strconv.Itoa(int(v["seqId"].(float64)))] = "completed"
	}
	io.IOData.Set(florest_constants.RESULT, errMap)
	return io, nil
}
