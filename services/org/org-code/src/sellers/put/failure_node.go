package put

import (
	"common/appconstant"
	"common/notification"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type Failure struct {
	id string
}

func (f *Failure) SetID(id string) {
	f.id = id
}

func (f Failure) GetID() (id string, err error) {
	return f.id, nil
}

func (f Failure) Name() string {
	return "UPDATE seller by id"
}

func (f Failure) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PUT_FAILURE)
	defer func() {
		logger.EndProfile(profiler, PUT_FAILURE)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+f.Name(), rc)
	io.ExecContext.SetDebugMsg(PUT_FAILURE, "Failure Node execution started")
	res, _ := io.IOData.Get(FAILURE_DATA)
	io.IOData.Set(florest_constants.RESULT, res)
	notification.SendNotification("Validation failure during Seller Update", fmt.Sprintf("Mandatory fields missing : %v", res), nil, "error")
	return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Validation Failure"}
}
