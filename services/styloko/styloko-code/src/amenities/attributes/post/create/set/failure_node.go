package set

import (
	_ "common/appconstant"
	_ "fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	_ "reflect"
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
	return "Failure attribute set by id"
}

func (f Failure) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, FAILURE)
	defer func() {
		logger.EndProfile(profiler, FAILURE)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+f.Name(), rc)
	io.ExecContext.SetDebugMsg(FAILURE, "Failure Node execution started")

	res, _ := io.IOData.Get(FAILURE_DATA)
	err := res.(florest_constants.AppErrors)
	return io, &err
}
