package get

import (
	"amenities/categories/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

// CategoriesGetAll -> struct for node based data
type CategoriesGetAll struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesGetAll) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesGetAll) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesGetAll) Name() string {
	return "CategoriesGetAll"
}

// Execute -> Starts node execution.
// If no query params are found, then all active categories are returned.
// TODO: Caching strategy.
func (cs CategoriesGetAll) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_ALL)
	defer func() {
		logger.EndProfile(profiler, common.GET_ALL)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_CATEGORY_GET_ALL)
	ctgryStruct, _, err := cs.findByStatus("all")
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "No Data found.", DeveloperMessage: err.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, ctgryStruct)
	return io, nil
}

// findByStatus -> Returns all categories by status code.
func (cs CategoriesGetAll) findByStatus(status string) ([]common.CategoryGet, bool, error) {
	var ctgryStruct []common.CategoryGet
	mgoSession := mongo.GetMongoSession(constants.CATEGORY_SEARCH)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(constants.CATEGORY_COLLECTION)
	if status == common.ACTIVE || status == common.DELETED || status == common.INACTIVE {
		err := mgoObj.Find(bson.M{"status": status}).All(&ctgryStruct)
		logger.Info("Get " + status + "categories finished")
		logger.Debug(err)
		return ctgryStruct, true, err
	}
	if status == common.ALL {
		err := mgoObj.Find(nil).All(&ctgryStruct)
		logger.Info("Get " + status + "categories finished")
		logger.Debug(err)
		return ctgryStruct, true, err
	}
	return ctgryStruct, false, nil
}

// FindByIds returns categories by IDS
func (cs CategoriesGetAll) FindByIds(categoryIds []int) ([]common.CategoryGetVerbose, error) {
	if len(categoryIds) == 0 {
		return nil, nil
	}
	mongoSession := mongo.GetMongoSession(constants.CATEGORY_SEARCH)
	defer mongoSession.Close()
	mgoObj := mongoSession.SetCollection(constants.CATEGORY_COLLECTION)
	var cat []common.CategoryGetVerbose
	err := mgoObj.Find(bson.M{"seqId": bson.M{"$in": categoryIds}}).Sort("seqId").All(&cat)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting id details from mongo: %s", err.Error()))
		return nil, err
	}
	if len(cat) == 0 {
		return nil, errors.New("No categories found")
	}
	return cat, nil
}
