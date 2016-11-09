package get

import (
	"amenities/standardsize/common"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// StandardSizeErrorResponse -> struct for node based data
type StandardSizeErrorResponse struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (ss *StandardSizeErrorResponse) SetID(id string) {
	ss.id = id
}

// GetID -> returns current node ID to orchestrator
func (ss StandardSizeErrorResponse) GetID() (id string, err error) {
	return ss.id, nil
}

// Name -> Returns node name to orchestrator
func (ss StandardSizeErrorResponse) Name() string {
	return "StandardSizeErrorResponse"
}

// Execute -> Starts node execution.
func (ss StandardSizeErrorResponse) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData,
	error) {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_ERROR)
	defer func() {
		logger.EndProfile(profiler, common.GET_ERROR)
	}()
	res, _ := io.IOData.Get(common.STANDARDSIZE_GET_ERROR)
	err, _ := res.(florest_constants.AppError)
	return io, &err
}
