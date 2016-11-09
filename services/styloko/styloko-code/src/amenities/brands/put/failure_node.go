package put

import (
	"amenities/brands/common"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// Failure basic struct
type Failure struct {
	id string
}

// SetID sets ID for node
func (f *Failure) SetID(id string) {
	f.id = id
}

// GetID returns node ID
func (f Failure) GetID() (id string, err error) {
	return f.id, nil
}

// Name returns name of the node
func (f Failure) Name() string {
	return "UPDATE Brand by id"
}

// Execute runs the node workflow
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
	listError := new(florest_constants.AppErrors)
	errs := res.([]map[string]interface{})
	for x := range errs {
		for key, value := range errs[x] {
			s, _ := value.(string)
			err := common.GenError(key, s)
			listError.Errors = append(listError.Errors, err)
		}
	}
	return io, listError
}
