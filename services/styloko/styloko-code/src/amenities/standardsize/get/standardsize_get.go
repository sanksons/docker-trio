package get

import (
	"amenities/services/attributes"
	"amenities/services/brands"
	"amenities/services/categories"
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"common/appconstant"
	"common/mongodb"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// StandardSizeGet -> struct for node based data
type StandardSizeGet struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *StandardSizeGet) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs StandardSizeGet) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs StandardSizeGet) Name() string {
	return "StandardSizeGet"
}

// Execute -> Starts node execution
func (cs StandardSizeGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_SEARCH)
	defer func() {
		logger.EndProfile(profiler, common.GET_SEARCH)
	}()

	res, _ := io.IOData.Get(common.STANDARDSIZE_GET_VALID_DATA)
	resData := res.(map[string]int)
	limit := resData["limit"]
	offset := resData["offset"]
	attrSetId := resData["attrSetId"]
	results, err := cs.getAllStandardSize(limit, offset, attrSetId)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.DataNotFoundErrorCode,
			Message: "Failure in getting data", DeveloperMessage: err.Error()}
	}
	finalRes := common.FinalResult{}
	finalRes.Count = len(results)
	finalRes.SizesMapping = results
	io.IOData.Set(florest_constants.RESULT, finalRes)
	return io, nil
}

func (cs StandardSizeGet) getAllStandardSize(limit, offset,
	attrSetId int) ([]common.StandardSizeResult, error) {

	mgoSession := mongoFactory.GetMongoSession(common.STANDARDSIZE_SEARCH)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(common.STANDARDSIZE_COLLECTION)

	var stores []common.StandardSizeStore
	query := bson.M{"attrbtStId": attrSetId}
	err := mgoObj.Find(query).Sort("seqId").
		Skip(offset).Limit(limit).All(&stores)
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot get standard size:%v", err.Error()))
		return nil, err
	}
	var result []common.StandardSizeResult
	for _, val := range stores {
		res := cs.transformStore(val, mgoSession)
		result = append(result, res)
	}
	return result, nil
}

func (cs StandardSizeGet) transformStore(ip common.StandardSizeStore,
	mgoSession *mongodb.MongoDriver) common.StandardSizeResult {
	attrSet := attributes.GetAttributeSetById(ip.AttributeSetId, mgoSession)
	brnd, _ := brands.ById(ip.BrandId)
	lfCtg := categories.ById(ip.LeafCategoryId)

	var op common.StandardSizeResult
	op.SeqId = ip.SeqId
	op.Brand = brnd.Name
	op.AttributeSet = attrSet.Name
	op.LeafCategory = lfCtg.Name
	op.Size = ip.Size
	return op
}
