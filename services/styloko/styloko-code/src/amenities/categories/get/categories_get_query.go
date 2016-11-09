package get

import (
	"amenities/categories/common"
	"common/appconstant"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// CategoriesGetQuery -> struct for node based data
type CategoriesGetQuery struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesGetQuery) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesGetQuery) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesGetQuery) Name() string {
	return "CategoriesGetQuery"
}

// Execute -> Starts node execution.
// If no query params are found, then all active categories are returned.
// TODO: Caching strategy.
func (cs CategoriesGetQuery) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_QUERY)
	defer func() {
		logger.EndProfile(profiler, common.GET_QUERY)
	}()
	q, _ := io.IOData.Get(common.CATEGORY_QUERY_PARAMS)
	if query, ok := q.(string); ok {
		cf := CategoriesGetAll{}
		ctgryStruct, ok, err := cf.findByStatus(query)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "No Data found.", DeveloperMessage: err.Error()}
		}
		if !ok {
			return io, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode, Message: "No Data found.", DeveloperMessage: "Cannot find categories for provided status."}
		}
		io.IOData.Set(florest_constants.RESULT, ctgryStruct)
		return io, nil
	}
	return io, &florest_constants.AppError{Code: appconstant.FailedToCreateErrorCode, Message: "Failure in getting data for query.", DeveloperMessage: "Type assertion failure. Cannot assert to string."}
}
