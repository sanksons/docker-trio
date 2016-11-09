package put

import (
	"amenities/categories/common"
	"amenities/services/products"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// CategoriesUpdate -> struct for node based data
type CategoriesUpdate struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesUpdate) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesUpdate) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesUpdate) Name() string {
	return "CategoriesUpdate"
}

// Execute -> Starts node execution.
func (cs CategoriesUpdate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.PUT)
	defer func() {
		logger.EndProfile(profiler, common.PUT)
	}()
	data, _ := io.IOData.Get(common.CATEGORY_VALID_DATA)
	ctgryUpdate, ok := data.(*common.CategoryUpdate)
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.InconsistantDataStateErrorCode, Message: "Invalid Data.", DeveloperMessage: "Type assertion failure for category struct."}
	}
	// ID needs to be refetched from path params because only one constant can be get.
	// Error checks have been disabled since other nodes have already done this.
	pathParams, _ := utils.GetPathParams(io)
	id, _ := strconv.Atoi(pathParams[0])

	// Segment insertion into struct
	segIds := common.InsertSegments(ctgryUpdate.SegIds)
	ctgryUpdate.CategorySeg = segIds
	res, ok, err := cs.update(id, *ctgryUpdate)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "Mongo Error.", DeveloperMessage: err.Error()}
	}
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.FailedToCreateErrorCode, Message: "Update failed.", DeveloperMessage: "Update to MongoDB failed."}
	}
	go cs.Delete(id)
	logger.Debug("Result")
	logger.Debug(res)
	io.IOData.Set(florest_constants.RESULT, res)
	return io, nil
}

// // update -> Update the given struct in Mongo.
func (cs CategoriesUpdate) update(id int, ctgryStruct common.CategoryUpdate) (interface{}, bool, error) {
	logger.Info("Mongo update started")
	mgoSession := mongo.GetMongoSession(constants.CATEGORY_UPDATE)
	defer mgoSession.Close()
	updateQuery := bson.M{"$set": ctgryStruct}
	query := bson.M{"seqId": id}
	res, err := mgoSession.FindAndModify(constants.CATEGORY_COLLECTION, updateQuery, query, false)
	if err != nil {
		logger.Info("Mongo update failure")
		logger.Debug(err)
		return nil, false, err
	}

	// Invalidate product cache for category
	categoryIDs := make([]int, 1)
	categoryIDs[0] = id
	products.PurgeCacheByCategories(categoryIDs)

	// SQL data preperation for worker
	sqlData := cs.prepData(id, ctgryStruct)
	// Worker start for MySQL data push in background.
	categoryUpdatePool.StartJob(sqlData)

	logger.Info("Mongo update success")
	return res, true, nil
}

// prepData -> prepares data in a map[string]interface{} to send to the parallel worker
func (cs CategoriesUpdate) prepData(id int, ctgryStruct common.CategoryUpdate) map[string]interface{} {
	sqlData := make(map[string]interface{})
	sqlData["id"] = id
	sqlData["categoryUpdate"] = ctgryStruct
	sqlData["segIds"] = ctgryStruct.SegIds
	return sqlData
}

// Delete removes key from cache
func (cs CategoriesUpdate) Delete(key int) {
	err1 := cacheObj.Delete("CATEGORY_TREE")
	if err1 != nil {
		logger.Error(fmt.Sprintf("Error while deleting category tree from cache: %v", err1.Error()))
	}
}
