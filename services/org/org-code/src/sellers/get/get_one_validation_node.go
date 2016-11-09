package get

import (
	"common/appconstant"
	"common/utils"
	"fmt"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type GetOneDecision struct {
	id string
}

func (d *GetOneDecision) SetID(id string) {
	d.id = id
}

func (d GetOneDecision) GetID() (id string, err error) {
	return d.id, nil
}

func (d GetOneDecision) Name() string {
	return "GetOneDecision node for GET"
}

func (d GetOneDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_GET, "Single seller get decision started")
	//checking if path params exist for not
	//if path params exists returns true else false
	pathParams, err := utils.GetPathParams(io)
	if err != nil {
		logger.Info("pathparams :", pathParams)
		logger.Error(fmt.Sprintf("Error while getting path parameters: %s", err.Error()))
		return false, florest_constants.AppError{Code: appconstant.InvalidPath, Message: "Invalid Path", DeveloperMessage: "HTTP Path not formed correctly"}
	}
	if pathParams != nil {
		logger.Info("id:", pathParams[0])
		io.IOData.Set(GET_ONE, pathParams[0])
		return true, nil
	}
	return false, nil
}
