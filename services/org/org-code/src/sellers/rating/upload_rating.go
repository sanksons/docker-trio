package rating

import (
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/notification"
	"common/utils"
	"encoding/json"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sellers/common"
)

type UploadRating struct {
	id string
}

func (r *UploadRating) SetID(id string) {
	r.id = id
}

func (r UploadRating) GetID() (id string, err error) {
	return r.id, nil
}

func (r UploadRating) Name() string {
	return "UPLOAD RATING for seller"
}

func (r UploadRating) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, UPLOAD_RATING)
	defer func() {
		logger.EndProfile(profiler, UPLOAD_RATING)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+r.Name(), rc)
	io.ExecContext.SetDebugMsg(UPLOAD_RATING, "Upload Rating Node execution started")

	data, err := utils.GetPostData(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while ugetting post data : %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Post Data for Uploading Rating", DeveloperMessage: err.Error()}
	}

	req := new(Request)
	err = json.Unmarshal(data, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling json to Upload Request : %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while unmarshalling json to Upload Request", DeveloperMessage: err.Error()}
	}

	ids, errMap := r.UpdateRating(req.Data)
	if len(errMap) != 0 {
		io.IOData.Set(florest_constants.RESULT, errMap)
	}
	if len(ids) != 0 {
		sellerData, err := common.GetData(ids)
		if err != nil {
			notification.SendNotification("Error while getting data for updated ids after Seller Rating Upload", err.Error(), nil, "error")
			logger.Error(fmt.Sprintf("Error while getting data for updated ids:%s", err.Error()))
		}
		prodUpdatePool.StartJob(sellerData)
	}
	return io, nil
}

//This function updates seller rating in mongo by ranging over the array of seller data passed
func (r UploadRating) UpdateRating(reqData []common.Schema) ([]int, map[string]interface{}) {
	errorMap := make(map[string]interface{})
	var ids []int
	for _, v := range reqData {
		id, err := r.UpdateRatingInMongo(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while Updating Rating in mongo : %v", err))
			errorMap[v.SellerId] = err.Error()
		}
		if id != 0 {
			ids = append(ids, id)
		}
	}
	return ids, errorMap
}

//This function inserts rating in mongo for the seller data passed by updating it
func (r UploadRating) UpdateRatingInMongo(data common.Schema) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, UPLOAD_RATING_MONGO)
	defer func() {
		logger.EndProfile(profiler, UPLOAD_RATING_MONGO)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	var orgInfo common.Schema
	upsertVal := false
	deleteVal := false
	returnNew := true
	updatedVal := bson.M{"$set": bson.M{"rating": data.Rating}}
	findCriteria := bson.M{"slrId": data.SellerId}
	change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
	_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &orgInfo)
	if err == nil {
		go r.Delete(orgInfo.SeqId)
		logger.Info(fmt.Sprintf("SellerId %s updated to %v", data.SellerId, orgInfo))
		return orgInfo.SeqId, nil
	}
	return 0, err
}

//Deleting key from cache
func (r UploadRating) Delete(key int) {
	cacheObj.Delete(fmt.Sprintf("%s-%d", common.SELLERS, key))
}
