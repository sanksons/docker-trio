package put

import (
	"amenities/categories/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// CategoryCheckNode -> struct for node based data
type CategoryCheckNode struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoryCheckNode) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoryCheckNode) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoryCheckNode) Name() string {
	return "CategoryCheckNode"
}

// Execute -> Decides which node to run next. Here its a validation node.
func (cs CategoryCheckNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.PUT_CHECK_NODE)
	defer func() {
		logger.EndProfile(profiler, common.PUT_CHECK_NODE)
	}()
	logger.Info("Fetching category ID from constants")
	ctgryID, _ := io.IOData.Get(common.CATEGORY_ID)
	id, ok := ctgryID.(int)
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.InconsistantDataStateErrorCode, Message: "Invalid Data.", DeveloperMessage: "Type assertion failure for ID."}
	}
	if found := cs.checkID(id); !found {
		return io, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode, Message: "Category doesn't exist.", DeveloperMessage: "Couldn't fetch ID."}
	}
	data, _ := utils.GetPostData(io)
	ctgryUpdate := new(common.CategoryUpdate)
	err := json.Unmarshal(data, &ctgryUpdate)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.InconsistantDataStateErrorCode, Message: "Invalid Data."}
	}
	appErrors, flag := common.ValidateV2(*ctgryUpdate)
	if !flag {
		return io, &appErrors
	}
	logger.Debug("Category Update struct is: ")
	logger.Debug(ctgryUpdate)
	io.IOData.Set(common.CATEGORY_VALID_DATA, ctgryUpdate)
	return io, nil
}

func (cs CategoryCheckNode) checkID(id int) bool {
	logger.Info("Get category by ID started")
	var ctgry common.CategoryGet
	mgoSession := mongo.GetMongoSession(constants.CATEGORY_UPDATE)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(constants.CATEGORY_COLLECTION)
	mgoObj.Find(bson.M{"seqId": id}).One(&ctgry)
	if ctgry.CategoryId != 0 {
		logger.Info("Category ID found")
		return true
	}
	logger.Info("Category ID not found")
	return false
}
