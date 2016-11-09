package post

import (
	"amenities/brands/common"
	"common/appconstant"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//Struct for Failure (node based)
type Failure struct {
	id string
}

//Function to SetID for current node from orchestrator
func (f *Failure) SetID(id string) {
	f.id = id
}

//Function that returns current node ID to orchestrator
func (f Failure) GetID() (id string, err error) {
	return f.id, nil
}

//Function that returns node name to orchestrator
func (f Failure) Name() string {
	return "CREATE Brand by id"
}

//Function to start node execution
func (f Failure) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.FAILURE)
	defer func() {
		logger.EndProfile(profiler, common.FAILURE)
	}()

	//Enable logs for Debugging
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+f.Name(), rc)
	io.ExecContext.SetDebugMsg(common.FAILURE, "Failure Node execution started")

	res, _ := io.IOData.Get(common.FAILURE_DATA)
	io.IOData.Set(florest_constants.RESULT, res)

	return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Validation Failure", DeveloperMessage: "Mandatory Fields Missing."}
}
