package common

import (
	"common/ResourceFactory"
	"common/notification"
	"common/utils"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

//converts any struct to map, takes out seqids from it and gets data
//for those ids from db
func GetDataForInsertedIds(resp interface{}) ([]Schema, error) {
	//converting resp struct to map
	respData, errs := utils.ConvertStructArrToMapArr(resp)
	if errs != nil {
		logger.Error(fmt.Sprintf("Error while converting data []struct to []map[string]interface{} :%s", errs.Error()))
		return nil, errs
	}
	var ids []int
	for _, v := range respData {
		if v["seqId"] != 0 {
			ids = append(ids, int(v["seqId"].(float64)))
		}
	}
	//if nothing was updated in mongo,return
	if len(ids) == 0 {
		return nil, nil
	}
	//getting data for inserted sellers
	data, er := GetData(ids)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while getting data for inserted sellers : %s", er.Error()))
		return nil, errs
	}
	return data, nil
}

//gets data for the seqIds passed from mongo
func GetData(ids []int) ([]Schema, error) {
	bsonMap := make(map[string]interface{})
	m := make(map[string]interface{})
	m["$in"] = ids
	bsonMap["seqId"] = m
	data, err := GetDetailsFromMongo(bsonMap, len(ids), 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting data from search by ids :%s", err.Error()))
		return nil, err
	}
	return data, nil
}

//takes map[string]interface{} and checks whether the string key exists
//in mongo or not,returns true if it exists, false otherwise along with the found object.
func CheckIfKeyExists(bsonMap map[string]interface{}) (bool, Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_CHK_KEY_EXISTS)
	defer func() {
		logger.EndProfile(profiler, MONGO_CHK_KEY_EXISTS)
	}()
	mgoSession := ResourceFactory.GetMongoSession(SELLERS)
	mgoObj := mgoSession.SetCollection(SELLERS_COLLECTION)
	defer mgoSession.Close()
	orgData := Schema{}
	if val, ok := bsonMap["seqId"]; ok {
		val = val.(int)
	}
	err := mgoObj.Find(bsonMap).One(&orgData)
	if orgData.SeqId != 0 {
		return true, orgData, err
	}
	return false, orgData, nil
}

//takes Schema as input and inserts it if the seqId existed
//in db, returns err otherwise
func Insert(mgoObj *mgo.Collection, data Schema, upsertValue bool) (Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_INSERT)
	defer func() {
		logger.EndProfile(profiler, MONGO_INSERT)
	}()
	var orgInfo Schema
	upsertVal := upsertValue
	deleteVal := false
	returnNew := true
	updatedVal := bson.M{"$set": data}
	findCriteria := bson.M{"seqId": data.SeqId}
	change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
	_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &orgInfo)
	return orgInfo, err
}

func PushDataToMemcache(data []Schema, command string, comment string) {
	for _, v := range data {
		aliceData := new(AliceStruct)
		aliceData.SeqId = strconv.Itoa(v.SeqId)
		aliceData.Status = v.Status
		aliceData.SellerId = v.SellerId
		m := map[string]interface{}{strconv.Itoa(v.SeqId): aliceData}
		alicedata, _ := utils.JSONMarshal(m, true)
		stringTime := strconv.FormatInt(time.Now().UnixNano(), 10)
		version := stringTime[0 : len(stringTime)-4]

		sqldriver, sqlerr := ResourceFactory.GetMySqlDriver(SELLER_UPDATE)
		if sqlerr != nil {
			logger.Error(fmt.Sprintf("(s Update)#PushToMemcache(): %s", sqlerr.Error()))
			notification.SendNotification("Error while getting mysql driver", fmt.Sprintf("#PushToMemcache(): %s", sqlerr.Error()), nil, "error")
		}
		sql := `INSERT INTO alice_message
			 (timestamp, data, command, caller, comment, type)
			 VALUES (?,?,?,?,?,?)`

		_, serr := sqldriver.Execute(sql,
			version, alicedata, command, "org service", comment, "suppliers",
		)
		if serr != nil {
			logger.Error(fmt.Sprintf("(s Update)#PushToMemcache()2: %s", serr.DeveloperMessage))
			notification.SendNotification("Error while executing alice_message", fmt.Sprintf("#PushToMemcache()2: %s", serr.DeveloperMessage), nil, "error")
		}
	}
}
