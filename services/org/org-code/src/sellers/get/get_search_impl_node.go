package get

import (
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"os"
	"sellers/common"
	"strconv"
	"strings"
)

type SearchSeller struct {
	id string
}

func (s *SearchSeller) SetID(id string) {
	s.id = id
}

func (s SearchSeller) GetID() (id string, err error) {
	return s.id, nil
}

func (s SearchSeller) Name() string {
	return "Search seller"
}

func (s SearchSeller) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_SEARCH)
	defer func() {
		logger.EndProfile(profiler, SELLER_SEARCH)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_SELLER_GET_SEARCH)
	rc, err := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting context: %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "", DeveloperMessage: err.Error()}
	}
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_SEARCH, "Seller search execution started")
	data, err := io.IOData.Get(GET_SEARCH)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting search request data: %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting search request data", DeveloperMessage: err.Error()}
	}
	// parse query
	dataMap := data.(map[string]interface{})
	if len(dataMap) == 0 {
		logger.Error("Error while getting search params.")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String"}
	}
	//getting limit and offset
	limit, offset, err := s.GetLimitOffset(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing limit and offset %s:", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid limit and offset", DeveloperMessage: err.Error()}
	}
	//getting headers
	header, err := utils.GetRequestHeader(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting headers %s:", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Headers", DeveloperMessage: err.Error()}
	}
	//geting data for search parameters
	resp, err := s.GetSearchData(dataMap, limit, offset, header)

	if err != nil {
		logger.Error(fmt.Sprintf("Error in Getting Search Data : %s", err.Error()))
		// return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error in Getting Search Data", DeveloperMessage: err.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, resp)
	//getting count for all sellers according to the search criteria
	count := 0
	var e error
	val, ok := utils.GetQueryParams(io, "count")
	if ok && val == "true" {
		count, e = s.GetSellerDataCount(dataMap)
		if e != nil {
			logger.Error(fmt.Sprintf("Error while getting seller count from mongo : %v", e))
			return io, &florest_constants.AppError{Code: appconstant.ResourceNotFoundCode, Message: "Could not get Count from Mongo", DeveloperMessage: e.Error()}
		}
	} else {
		count = len(resp)
	}
	//setting count in META
	info := common.SetCountInMeta(io, count)
	io.IOData.Set(florest_constants.RESPONSE_META_DATA, info)
	return io, nil
}

// get seller details
func (s SearchSeller) GetIdDetails(ids []int, limit int, offset int) ([]common.Schema, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_SEARCH_IDS)
	defer func() {
		logger.EndProfile(profiler, MONGO_SEARCH_IDS)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	var org []common.Schema
	query := map[string]interface{}{"seqId": bson.M{"$in": ids}}
	value := os.Getenv("SELLER_STATUS")
	if value != "" && value == "active" {
		query["status"] = "active"
	}
	err := mgoObj.Find(query).Sort("crtdAt").Limit(limit).Skip(offset).All(&org)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting id details from mongo: %s", err.Error()))
		return nil, err
	}
	if len(org) == 0 {
		logger.Error("No seller found against passed seq id's")
	}
	return org, nil
}

// get batch ids from the cache
func (s SearchSeller) GetBatch(ids []int) ([]int, []common.Schema) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET_BATCH_CACHE)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET_BATCH_CACHE)
	}()
	var keys []string
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("%s-%d", common.SELLERS, id))
	}
	items, err := cacheObj.GetBatch(keys, false, false)
	if err != nil {
		logger.Error(err.Error())
		return ids, nil
	}
	var missed []int
	var found []common.Schema
	for _, val := range items {
		if val.Value == "" {
			i, _ := strconv.Atoi(strings.Split(val.Key, "-")[1])
			missed = append(missed, i)
		} else {
			var v common.Schema
			json.Unmarshal([]byte(val.Value.(string)), &v)
			found = append(found, v)
		}
	}
	return missed, found
}

// set seller information to cache
func (s SearchSeller) SetMulti(org []common.Schema) {
	gs := GetSeller{}
	for _, o := range org {
		gs.SetCache(strconv.Itoa(o.SeqId), o)
	}
}

//function searches in the cache first, if missed then brings data from mongo
func (s SearchSeller) SearchByIds(ids []int, limit int, offset int) ([]common.Schema, error) {
	// fetch from cache
	missed, found := s.GetBatch(ids)
	if len(missed) == 0 {
		return found, nil
	}
	// get missed id's from cache
	resp, err := s.GetIdDetails(missed, limit, offset)
	//if all ids were found in cache
	if err != nil && len(found) == 0 {
		logger.Error(fmt.Sprintf("Error while getting Id Details : %s", err.Error()))
		return nil, err
	}
	// merge found & cached items
	if len(found) != 0 && len(resp) != 0 {
		resp = append(resp, found...)
	}
	//if no data found of ids, return
	if len(resp) == 0 {
		logger.Error(fmt.Sprintf("No records found against search criteria %d", ids))
		return nil, err
	}
	//setting in cache
	s.SetMulti(resp)
	return resp, nil
}

//gets limit and offset if passed, if not sets limit=1000 and offset=0
func (s SearchSeller) GetLimitOffset(io workflow.WorkFlowData) (int, int, error) {
	//getting limit from query params
	//else setting default at 1000
	lim, ok := utils.GetQueryParams(io, "limit")
	if !ok {
		logger.Info("No limit passed. Setting default limit.")
		lim = "1000"
	}
	//converting limit from string to int
	limit, err := strconv.Atoi(lim)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing limit %v:", lim))
		return 0, 0, err
	}
	//getting offset from query params
	//else setting default at 0
	skip, ok := utils.GetQueryParams(io, "offset")
	if !ok {
		logger.Info("No offset passed.Setting default offset")
		skip = "0"
	}
	//converting offset to int
	offset, err := strconv.Atoi(skip)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing offset %v:", skip))
		return 0, 0, err
	}
	return limit, offset, nil
}

//generic function to get seller data for any datamap passed
func (s SearchSeller) GetSearchData(dataMap map[string]interface{}, limit int, offset int, header map[string]interface{}) ([]common.Schema, error) {
	resp := make([]common.Schema, 0)
	var err error
	var stringIds string
	//search by ids(different because caching is involved)
	//for the rest gets data from mongo directly
	dataVal, ok := dataMap["seqId"]
	if !ok {
		//search by all other params
		resp, err = common.GetDetailsFromMongo(dataMap, limit, offset)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in GetDetailsFromMongo :%s", err.Error()))
			return nil, err
		}
	} else {
		//search by ids
		x := dataVal.(bson.M)
		intIds := common.ParseQuery(x["$in"].(string))
		resp, err = s.SearchByIds(intIds, limit, offset)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in SearchByIds:%s", err.Error()))
			return nil, err
		}
	}
	//if header is set,get commission value
	if _, ok := header["Commissions"]; ok {
		headerArr := header["Commissions"].([]interface{})
		_, ok = headerArr[0].(string)
		if ok {
			//get update commission value from slave db of SC
			resp, err = s.GetSellerWithCommissions(resp, stringIds)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in GetUpdateCommission :%s", err.Error()))
				return resp, err
			}
		}
	}
	return resp, nil
}

//takes input schema without update commiission values, gets values from mysql
//appends and returns schema with update commission values
func (s SearchSeller) GetSellerWithCommissions(slrData []common.Schema, ids string) ([]common.Schema, error) {
	comMap, err := s.GetCommissions(ids)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting commission: %s", err.Error()))
		return slrData, err
	}
	//ranging over existing data to append commission values
	for key, _ := range slrData {
		if val, ok := comMap[slrData[key].SeqId]; ok {
			slrData[key].UpdateCommission = val
		}
	}
	return slrData, nil
}

//takes an id and gets commission from the seller centre slave database for the passed ids
func (s SearchSeller) GetCommissions(ids string) (map[int][]common.Commission, error) {
	comData := []common.GetCommission{}
	ids = strings.Replace(ids, "[", "", -1)
	ids = strings.Replace(ids, "]", "", -1)
	sql := `SELECT
                          seller.src_id AS sellerId,
                          catalog_category.src_id AS categoryId,
                          seller_commission.percentage
                        FROM seller_commission
                        INNER JOIN catalog_category
                          ON catalog_category.id_catalog_category = seller_commission.fk_catalog_category
                        INNER JOIN seller
                          ON seller.id_seller=seller_commission.fk_seller
                        WHERE seller.src_id IN  (` + ids + `)`
	driver, errs := mongo.GetMySqlDriverSC(UPDATE_COMMISSIONS)
	if errs != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", errs.Error()))
		return nil, errs
	}
	rows, err := driver.Query(sql)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing query: %v", err))
		return nil, errors.New("Error while executing query to get commission value")
	}
	for rows.Next() {
		com := common.GetCommission{}
		err := rows.Scan(&com.SeqId,
			&com.CategoryId,
			&com.Percentage)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while scanning commiission values: %v", err))
			return nil, err
		}
		comData = append(comData, com)
	}
	comMap, er := s.GetCommissionMap(comData)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while getting map of commission data: %s", er.Error()))
		return nil, er
	}
	return comMap, nil
}

//populates the recieved commission into struct to be visible in api response
func (s SearchSeller) GetCommissionMap(comData []common.GetCommission) (map[int][]common.Commission, error) {
	comMap := make(map[int][]common.Commission)
	for _, v := range comData {
		comArr := common.Commission{}
		comArr.CategoryId = v.CategoryId
		comArr.Percentage = v.Percentage
		comMap[v.SeqId] = append(comMap[v.SeqId], comArr)
	}
	return comMap, nil
}

//Function gets seller data count for the datamap passed
func (s SearchSeller) GetSellerDataCount(dataMap map[string]interface{}) (int, error) {
	var err error
	var count int
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, MONGO_GET_SELLER_SEARCH_COUNT)
	defer func() {
		logger.EndProfile(profiler, MONGO_GET_SELLER_SEARCH_COUNT)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	if _, ok := dataMap["seqId"]; ok {
		return 0, nil
	}
	for k, v := range dataMap {
		count, err = mgoObj.Find(bson.M{k: v}).Count()
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
