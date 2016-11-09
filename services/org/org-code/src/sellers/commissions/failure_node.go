package commissions

import (
	"common/appconstant"
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
	return "Get Commission failure node"
}

func (f Failure) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, COMMISSION_FAILURE)
	defer func() {
		logger.EndProfile(profiler, COMMISSION_FAILURE)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+f.Name(), rc)
	io.ExecContext.SetDebugMsg(COMMISSION_FAILURE, "Failure Node execution started")
	return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Failure while getting commissions", DeveloperMessage: "Incomplete Url"}
}
