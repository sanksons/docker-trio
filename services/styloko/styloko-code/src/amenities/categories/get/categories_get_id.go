package get

import (
	"amenities/categories/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// CategoriesGetID -> struct for node based data
type CategoriesGetID struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesGetID) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesGetID) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesGetID) Name() string {
	return "CategoriesGetID"
}

// Execute -> Starts node execution.
// If no query params are found, then all active categories are returned.
// TODO: Caching strategy.
func (cs CategoriesGetID) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_ONE)
	defer func() {
		logger.EndProfile(profiler, common.GET_ONE)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_CATEGORY_GET_ONE)
	path, _ := io.IOData.Get(common.CATEGORY_PATH_PARAMS)
	if pathParams, ok := path.([]string); ok {
		ctgry, ok, err := cs.FindById(pathParams[0])
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "No Data found.", DeveloperMessage: err.Error()}
		}
		if !ok {
			return io, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode, Message: "Invalid Category Id", DeveloperMessage: "Id param value in query is Invalid"}
		}
		io.IOData.Set(florest_constants.RESULT, ctgry)
		return io, nil
	}
	return io, &florest_constants.AppError{Code: appconstant.FailedToCreateErrorCode, Message: "Failure in getting category by ID", DeveloperMessage: "Type assertion failure. Cannot assert to []string."}
}

// FindById -> returns mongo document with given ID string.
func (cs CategoriesGetID) FindById(idStr string) (common.CategoryGetVerbose, bool, error) {
	mgoSession := mongo.GetMongoSession(constants.CATEGORY_SEARCH)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(constants.CATEGORY_COLLECTION)

	var ctgry common.CategoryGetVerbose
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Debug(err)
		return ctgry, false, nil
	}
	err = mgoObj.Find(bson.M{"seqId": id}).One(&ctgry)
	if err != nil {
		return ctgry, false, err
	}
	if ctgry.CategoryId != 0 {
		logger.Info("ID found in Mongo")
		logger.Debug(ctgry)
		return ctgry, true, nil
	}
	logger.Info("ID not found in Mongo")
	return ctgry, false, nil
}
