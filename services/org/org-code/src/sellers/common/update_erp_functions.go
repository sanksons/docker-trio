package common

import (
	mongo "common/ResourceFactory"
	"common/appconfig"
	"encoding/json"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//sends prepared data to ERP Orechestrator and calls transform request which
//transforms the response to get correctly accepted ids by ERP
func SendDataToErp(erpData []ErpData) error {
	jsonData, err := json.Marshal(erpData)
	if err != nil {
		logger.Error("Error while marshalling erpData", err)
		return err
	}
	logger.Info(fmt.Sprintf("Json Data sent to Orchestrator:%s", string(jsonData)))
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	headers := make(map[string]string)
	headers["X-Custom-Header"] = "ERP Data"
	logger.Info(fmt.Sprintf("Config Erp Url:%s", config.Erp.Url))
	resp, err := http.HttpPost(config.Erp.Url, headers, string(jsonData), 100*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending POST Request,probably timeout from ERP Server %s", err.Error()))
		return err
	}
	TransformRequest(resp)
	return nil
}

//This fuunction takes http response got from ERP Orchestrator and
//gets correctly and wrongly inserted ids and sets sync flag as true for correctly inserted ids
func TransformRequest(resp *http.APIResponse) {
	if resp.HttpStatus != 200 {
		logger.Error(fmt.Sprintf("Error In ERP Response : %s", string(resp.Body)))
		return
	}
	logger.Info(string(resp.Body))
	res1 := []ErpResponse{}

	err := json.Unmarshal(resp.Body, &res1)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling json to ERP Response : %s", err.Error()))
	}
	res := res1[0]
	insSlrIds, notInsSlrIds := GetInsertedAndNotInsertedSellerIds(res.Data)
	logger.Info(fmt.Sprintf("Inserted sellerIds are : %s", insSlrIds))
	logger.Info(fmt.Sprintf("Not Inserted sellerIds are : %s", notInsSlrIds))

	err = SetSyncFlagForInsertedSlrIds(insSlrIds)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while setting sync flag for inserted sellerIds in Erp : %s", err.Error()))
	}
	logger.Info("Transforming Request Complete")
}

//this function takes in an []map[string]interface{} and gives
//array of correctly inserted ids and not correctly inserted ids
func GetInsertedAndNotInsertedSellerIds(data []map[string]interface{}) ([]string, []string) {
	var arrNotInsSlrIds []string
	var arrInsSlrIds []string
	for _, v := range data {
		if _, ok := v["slrId"]; !ok {
			logger.Error("Error while recieveing response from ERP : SellerId is missing")
			return nil, nil
		}
		slrId := v["slrId"].(string)
		if v["result"] == false {
			logger.Error(fmt.Sprintf("Error recived for sellerId %s from ERP ADAPTER %v :", slrId, v["errorMsg"]))
			arrNotInsSlrIds = append(arrNotInsSlrIds, slrId)
		}
		if v["result"] == true {
			arrInsSlrIds = append(arrInsSlrIds, slrId)
		}
	}
	return arrInsSlrIds, arrNotInsSlrIds
}

//This function takes the correctly inserted ids as parameter,
// loops over it and sets calls SetSyncInMongo()
func SetSyncFlagForInsertedSlrIds(sellerIds []string) error {
	for _, v := range sellerIds {
		err := SetSyncTrueInMongo(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while setting sync flag in mongo : %v", err))
			return err
		}
	}
	return nil
}

//This function sets sync flag as true for the passed seller id in mongo
func SetSyncTrueInMongo(sellerId string) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_SYNC_INSERT)
	defer func() {
		logger.EndProfile(profiler, MONGO_SYNC_INSERT)
	}()
	mgoSession := mongo.GetMongoSession(SELLERS)
	mgoObj := mgoSession.SetCollection(SELLERS_COLLECTION)
	defer mgoSession.Close()
	var orgInfo Schema
	upsertVal := false
	deleteVal := false
	returnNew := true
	updatedVal := bson.M{"$set": bson.M{"sync": true}}
	findCriteria := bson.M{"slrId": sellerId}
	change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
	_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &orgInfo)
	logger.Info(fmt.Sprintf("SellerId %s updated to %v", sellerId, orgInfo))
	return err
}
