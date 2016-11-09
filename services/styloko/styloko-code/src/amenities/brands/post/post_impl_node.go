package post

import (
	"amenities/brands/common"
	"amenities/brands/put"
	mongo "common/ResourceFactory"
	"common/constants"
	"fmt"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//Struct for Failure (node based)
type Insert struct {
	id string
}

//Function to SetID for current node from orchestrator
func (b *Insert) SetID(id string) {
	b.id = id
}

//Function that returns current node ID to orchestrator
func (b Insert) GetID() (id string, err error) {
	return b.id, nil
}

//Function that returns node name to orchestrator
func (b Insert) Name() string {
	return "Insert Brand by id"
}

//Function to start node execution
func (b Insert) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_INSERT)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_INSERT)
	}()

	//Enable logs for Debugging
	io.ExecContext.SetDebugMsg(common.BRAND_INSERT, "Brand insert execution started")
	logger.Info("Insert Brand Started")
	p, _ := io.IOData.Get(common.BRAND_DATA)

	//parsing data into []common.Brand from interface{}
	data, pOk := p.([]common.Brand)

	if !pOk || data == nil {
		logger.Error("BrandCreate invalid type of columns")
		return io, &florest_constants.AppError{Code: florest_constants.ParamsInValidErrorCode, Message: "Invalid Parameters"}
	}

	//Calling insert function to insert data into mongo
	resp := b.CheckAndInsert(data)

	//Setting value recieved from update() into RESULT
	io.IOData.Set(florest_constants.RESULT, resp)
	logger.Info("Insert brand finished")

	return io, nil
}

//This function return a newly generated seqId by incrementing one
//from the seqId found in counters collection and checks if it already exists
//sends 0 and err if id exists, id and nil otherwise
func (b Insert) GetNewSeqId() (int, error) {
	//Enable profiling
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_INSERT)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_INSERT)
		mgoSession.Close()
	}()

	id := mgoSession.GetNextSequence(constants.BRAND_COLLECTION)
	ok, err := b.CheckIfSeqIdExists(id)
	if ok {
		return 0, err
	}
	return id, nil
}

//Checks if sequence Id already exists, if not data is inserted in mongo
// and and []Response is returned with Brand Name and newly created sequence Id
func (b Insert) CheckAndInsert(data []common.Brand) []Response {
	var arrInsert []Response
	d := put.Update{}
	for _, v := range data {
		var temp Response
		id, err := b.GetNewSeqId()
		if err == nil {
			v.SeqId = id
			brandInfo, err := d.Insert(v, true)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while inserting into Mongo %s", err.Error()))
				temp.Name = v.Name
				temp.SeqId = 0
				temp.Error = err.Error()
				arrInsert = append(arrInsert, temp)
				continue
			}
			temp.Name = brandInfo.Name
			temp.SeqId = brandInfo.SeqId
			arrInsert = append(arrInsert, temp)

			// SQL data preperation for worker
			sqlData := b.prepData(brandInfo.SeqId, brandInfo)

			// Worker start for MySQL data push in background.
			brandCreatePool.StartJob(sqlData)
			continue
		}
		temp.Name = v.Name
		temp.SeqId = id
		arrInsert = append(arrInsert, temp)
	}
	return arrInsert
}

//This function takes as input a seqId and returns true
//or false if the id exists or not respectively
func (b Insert) CheckIfSeqIdExists(seqId int) (bool, error) {
	//Enable profiling
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_INSERT)
	mgoSession := mongo.GetMongoSession(common.BRAND_OPERATION)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_INSERT)
		mgoSession.Close()
	}()

	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	brandData := common.BrandUpdate{}
	err := mgoObj.Find(bson.M{"seqId": seqId}).One(&brandData)

	if err == nil {
		return true, err
	}

	return false, nil
}

//prepares data in a map[string]interface{} to send to the parallel worker
func (b Insert) prepData(id int, brandStruct common.Brand) map[string]interface{} {
	sqlData := make(map[string]interface{})
	sqlData["brandInfo"] = brandStruct

	return sqlData
}
