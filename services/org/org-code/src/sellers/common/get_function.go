package common

import (
	mongo "common/ResourceFactory"
	"errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
)

//gets details from mongo based on the bson map passed
func GetDetailsFromMongo(bsonMap map[string]interface{}, limit int, offset int) ([]Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_SEARCH)
	defer func() {
		logger.EndProfile(profiler, MONGO_SEARCH)
	}()
	mgoSession := mongo.GetMongoSession(SELLERS)
	mgoObj := mgoSession.SetCollection(SELLERS_COLLECTION)
	defer mgoSession.Close()
	var org []Schema
	err := mgoObj.Find(bsonMap).Sort("-updtdAt").Limit(limit).Skip(offset).All(&org)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting seller details from mongo: %s", err.Error()))
		return nil, err
	}
	if len(org) == 0 {
		logger.Error("No seller found against passed search criteria")
		return nil, errors.New("No seller found against passed search criteria")
	}
	return org, nil
}

//sets count in response meta data for florest io data
func SetCountInMeta(io workflow.WorkFlowData, count int) *http.ResponseMetaData {
	var info *http.ResponseMetaData

	infoVal, merr := io.IOData.Get(florest_constants.RESPONSE_META_DATA)
	if merr != nil {
		info = http.NewResponseMetaData()
	} else if infoObj, ok := infoVal.(*http.ResponseMetaData); ok {
		info = infoObj
	} else {
		info = http.NewResponseMetaData()
	}
	info.ApiMetaData["count"] = count
	return info
}

//Generic function to parse queryMap into bsonMap that can be used to search in mongo
func GetBsonMapFromSearchMap(queryMap []map[string]interface{}) (map[string]interface{}, error) {
	bsonMap := make(map[string]interface{})
	for _, v := range queryMap {
		for x, y := range v {
			xNew := GetSellerMapping(x)
			val := y.(map[string]interface{})
			if x == "" {
				return nil, errors.New("Key field missing")
			}
			if val["value"] == nil {
				return nil, errors.New("Value field missing")
			}
			switch val["operator"] {
			case "eq":
				if x == "sync" {
					if val["value"] == "true" {
						bsonMap[xNew] = true
						break
					}
					bsonMap[xNew] = false
					break
				}
				bsonMap[xNew] = val["value"]
				break
			case "in":
				bsonMap[xNew] = bson.M{"$in": val["value"]}
				break
			default:
				return nil, errors.New("Incorrect value for operator")
			}
		}
	}
	return bsonMap, nil
}

// parse search query
func ParseQuery(id string) []int {
	id = strings.Replace(id, "\"", "", -1)
	id = strings.Replace(id, "[", "", -1)
	id = strings.Replace(id, "]", "", -1)
	idSlice := strings.Split(id, ",")
	var arrId []int
	for _, v := range idSlice {
		id, _ := strconv.Atoi(v)
		arrId = append(arrId, id)
	}
	return arrId
}

//mapping for search string to mongo names
func GetSellerMapping(key string) string {
	sellerMapping := make(map[string]string)
	sellerMapping["id"] = "seqId"
	sellerMapping["sellerId"] = "slrId"
	sellerMapping["name"] = "slrName"
	sellerMapping["status"] = "status"
	sellerMapping["sync"] = "sync"
	sellerMapping["city"] = "city"
	sellerMapping["categories"] = "categories"
	return sellerMapping[key]
}

//gets data from mongo for the sequence Id passed
//added profile name to distinguish for two different metrics on datadog
//as it was unable to distinguish in calles for commission and get by id
func GetById(id int, profilerName string) (Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, profilerName)
	defer func() {
		logger.EndProfile(profiler, profilerName)
	}()
	mgoSession := mongo.GetMongoSession(SELLERS)
	mgoObj := mgoSession.SetCollection(SELLERS_COLLECTION)
	defer mgoSession.Close()
	var org Schema
	err := mgoObj.Find(bson.M{"seqId": id}).One(&org)
	return org, err
}
