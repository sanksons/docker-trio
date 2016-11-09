package commissions

import (
	"common/appconstant"
	"common/utils"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sellers/common"
)

type Validate struct {
	id string
}

func (v *Validate) SetID(id string) {
	v.id = id
}

func (v Validate) GetID() (id string, err error) {
	return v.id, nil
}

func (v Validate) Name() string {
	return "Validate node for GET COMMISSION"
}

func (v Validate) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_GET_COMMISSION)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_GET_COMMISSION)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+v.Name(), rc)
	io.ExecContext.SetDebugMsg(VALIDATE_GET_COMMISSION, "Seller update decision started")

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
	//if no query params foumd or any raw query map used call made to failure node
	if found == false && len(qmap) == 0 {
		return false, nil
	}
	//gets bson map from search map passed
	bsonMap, err := common.GetBsonMapFromSearchMap(queryMap)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String", DeveloperMessage: err.Error()}
	}
	//if no query return true
	if len(bsonMap) == 0 {
		return false, nil
	}
	//else set bsonMap to get in next node
	io.IOData.Set(GET_SEARCH, bsonMap)
	return true, nil
}
