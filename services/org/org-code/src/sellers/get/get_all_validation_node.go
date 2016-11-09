package get

import (
	"common/appconstant"
	"common/utils"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sellers/common"
)

type GetAllDecision struct {
	id string
}

func (d *GetAllDecision) SetID(id string) {
	d.id = id
}

func (d GetAllDecision) GetID() (id string, err error) {
	return d.id, nil
}

func (d GetAllDecision) Name() string {
	return "GetAllDecision node for GET"
}

func (d GetAllDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_GET, "Seller get all decision started")
	//checks query params,if they exist return true else false
	queryMap, found, err := utils.GetSearchQueries(io)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String", DeveloperMessage: err.Error()}
	}
	//gets raw query map. Used to make endpont exlusive
	qmap, err := utils.GetRawQueryMap(io)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String", DeveloperMessage: err.Error()}
	}
	//if no query params foumd or any raw query map used call made to get all sellers
	if found == false && len(qmap) == 0 {
		return true, nil
	}
	//gets bson map from search map passed
	bsonMap, err := common.GetBsonMapFromSearchMap(queryMap)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String", DeveloperMessage: err.Error()}
	}
	//if no query return true
	if len(bsonMap) == 0 {
		return true, nil
	}
	//else set bsonMap to get in next node
	io.IOData.Set(GET_SEARCH, bsonMap)
	return false, nil
}
