package put

import (
	"amenities/brands/common"
	"amenities/services/products"
	mongo "common/ResourceFactory"
	"common/constants"
	"fmt"
	"strconv"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Struct for Update Node
type Update struct {
	id string
}

//Function to SetID for current node from orchestrator
func (b *Update) SetID(id string) {
	b.id = id
}

//Function that returns current node ID to orchestrator
func (b Update) GetID() (id string, err error) {
	return b.id, nil
}

//Function that returns node name to orchestrator
func (b Update) Name() string {
	return "Update brand by id"
}

//Function to start node execution
func (b Update) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_UPDATE)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_UPDATE)
	}()

	//Enable logs for Debugging
	io.ExecContext.SetDebugMsg(common.BRAND_UPDATE, "Brand update execution started")
	logger.Info("Update Brand Started")

	p, _ := io.IOData.Get(common.BRAND_UPDATE_DATA)

	//parsing data into []common.Brand from interface{}
	data, pOk := p.([]common.Brand)
	if !pOk || data == nil {
		logger.Error("BrandUpdate invalid type of columns")
		return io, &florest_constants.AppError{Code: florest_constants.ParamsInValidErrorCode, Message: "Json Incorrect"}
	}

	//calling update function to update data into mongo
	resp := b.Update(data)

	//setting value recieved from update() into RESULT
	io.IOData.Set(florest_constants.RESULT, resp)
	logger.Info("Update brand finished")
	go b.deleteCache(data)
	b.purgeProductsCache(data)
	return io, nil
}

func (b Update) purgeProductsCache(brands []common.Brand) {
	var brandIds []int
	for _, brand := range brands {
		brandIds = append(brandIds, brand.SeqId)
	}
	products.PurgeCacheByBrands(brandIds)
}

// deleteCache -> deletes cache for all the brands updated
func (b Update) deleteCache(brands []common.Brand) {
	var ids []string
	for x := range brands {
		tmp := strconv.Itoa(brands[x].SeqId)
		tmp = common.BRANDS + "-" + tmp
		ids = append(ids, tmp)
	}
	err := cacheObj.DeleteBatch(ids)
	if err != nil {
		logger.Error("Error in deleting brand from cache")
	}
}

//Function takes []common.BrandUpdate and returns []Response with
// seqId as the id of the updated brand and result as true or false
//if the brand update was success or failure respectively.
func (b Update) Update(data []common.Brand) []Response {
	var arrUpdate []Response
	for _, v := range data {
		var temp Response

		brandInfo, err := b.Insert(v, false)
		if err != nil {
			temp.SeqId = v.SeqId
			temp.Result = false
			arrUpdate = append(arrUpdate, temp)
			logger.Error(fmt.Sprintf("Error while updating into Mongo.%s", err.Error()))
		} else {
			temp.SeqId = brandInfo.SeqId
			temp.Result = true
			arrUpdate = append(arrUpdate, temp)
		}
		// SQL data preperation for worker
		sqlData := b.prepData(brandInfo.SeqId, brandInfo)
		// Worker start for MySQL data push in background.
		brandUpdatePool.StartJob(sqlData)
		// Invalidate cache
		b.Delete(brandInfo.SeqId)
	}

	return arrUpdate
}

//Function that takes common.Brand as input and inserts it if the seqId existed
//in db, returns err otherwise
func (b Update) Insert(data common.Brand, upsertValue bool) (common.Brand, error) {
	//Enable profiling
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.VALIDATE_PUT)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.VALIDATE_PUT)
		mgoSession.Close()
	}()

	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	var brandInfo common.Brand
	upsertVal := upsertValue
	deleteVal := false
	returnNew := true
	updatedVal := bson.M{"$set": data}
	findCriteria := bson.M{"seqId": data.SeqId}
	change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
	_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &brandInfo)
	return brandInfo, err
}

//prepares data in a map[string]interface{} to send to the parallel worker
func (b Update) prepData(id int, brandStruct common.Brand) map[string]interface{} {
	sqlData := make(map[string]interface{})
	sqlData["brandInfo"] = brandStruct

	return sqlData
}

//Deleting key from cache
func (b Update) Delete(key int) {
	cacheObj.Delete(fmt.Sprintf("%s-%d", common.BRANDS, key))
}
