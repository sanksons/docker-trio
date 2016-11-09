package erp

import (
	"common/ResourceFactory"
	"common/appconfig"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sellers/common"
	"strconv"
)

var cacheObj cache.CacheInterface

//This function takes as input the post data,unmarshalls it to csv struct,
//transforms it to org schema and updates in mongo
func StartErpMigration(data []byte) ([]Response, error) {
	cStruct := make([]csvSchema, 0)
	err := initializeCacheObj()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while initialzing cache object: %s", err.Error()))
		return nil, err
	}
	err = json.Unmarshal(data, &cStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling post data into csv struct :%s", err.Error()))
		return nil, err
	}
	orgStruct, err := transformCsvStructToOrgSchema(cStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while transforming csv struct to org schema :%s", err.Error()))
		return nil, err
	}
	resp := updateInMongo(orgStruct)
	if len(resp) != 0 {
		return resp, errors.New("Error while updating in mongo.")
	}
	logger.Info("Erp Migration Successful.")
	return nil, nil
}

//This function takes as input csv type struct, populates additional info with the key slrCustInfo,
//unmarshalls into org schema struct and returns it
func transformCsvStructToOrgSchema(cStruct []csvSchema) ([]common.Schema, error) {
	orgStruct := make([]common.Schema, 0)
	cStruct = populateAdditionalInfo(cStruct)
	csvJson, err := json.Marshal(cStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while marshalling csvStruct :%v", err.Error()))
		return nil, err
	}
	err = json.Unmarshal(csvJson, &orgStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling csvJson :%v", err.Error()))
		return nil, err
	}
	return orgStruct, nil
}

//This function updates the org schema struct in mongo if the sellerId matches and invalidates cache,
//if not, it creates an array of response with sellerId and error and return reponse array
func updateInMongo(orgStruct []common.Schema) []Response {
	var ids []string
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ERP_INSERT)
	defer func() {
		logger.EndProfile(profiler, ERP_INSERT)
	}()
	mgoSession := ResourceFactory.GetMongoSession(ERP_MIGRATION)
	mgoObj := mgoSession.SetCollection(SELLERS_COLLECTION)
	defer mgoSession.Close()
	var orgInfo common.Schema
	erpMigResp := make([]Response, 0)
	upsertVal := false
	deleteVal := false
	returnNew := true
	for _, v := range orgStruct {
		logger.Info(fmt.Sprintf("Vendor Id being inserted : %s", v.SellerId))
		updatedVal := bson.M{"$set": v}
		findCriteria := bson.M{"slrId": v.SellerId}
		change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
		_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &orgInfo)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while inserting erp struct in mongo :%s", err.Error()))
			resp := Response{}
			resp.SellerId = v.SellerId
			resp.Error = err.Error()
			erpMigResp = append(erpMigResp, resp)
		}
		ids = append(ids, strconv.Itoa(orgInfo.SeqId))
	}
	go delete(ids)
	logger.Info(fmt.Sprintf("Response being returned: %v", erpMigResp))
	return erpMigResp
}

//This function takes []csvSchema, loops over all the fields for additional info
// and puts them in json tag slrCustInfo, appends it to the actual struct and returns the changes struct
func populateAdditionalInfo(csvStruct []csvSchema) []csvSchema {
	for k, v := range csvStruct {
		m := make(map[string]interface{})
		m["tinNo"] = v.TinNo
		m["panNo"] = v.PanNo
		m["paymtTrmsCode"] = v.PaymentTermsCode
		m["serTaxRegNo"] = v.SerTaxRegNo
		m["stateCode"] = v.StateCode
		m["ifscCode"] = v.IfscCode
		m["cntrctExpDate"] = v.CntrctExpDate
		m["cstNo"] = v.CstNo
		m["cinNo"] = v.CinNo
		m["natOfEntity"] = v.NatureOfEntity
		m["natOfBuis"] = v.NatureOfBuisness
		m["procTime"] = v.ProcessingTime
		m["oneshipCntrCode"] = v.OneshipCentreCode
		m["oneshipAddr"] = v.OneshipAddress
		m["oneshipCity"] = v.OneshipCity
		m["oneshipZipcode"] = v.OneshipZipcode
		m["oneshipState"] = v.OneshipState
		m["retrnProcCntrCode"] = v.ReturnProcessingCentreCode
		m["retrnProcCntrAddr"] = v.ReturnProcessingCentreAddress
		m["retrnProcCntrCity"] = v.ReturnProcessingCentreCity
		m["retrnProcCntrZipcode"] = v.ReturnProcessingCentreZipcode
		m["retrnProcCntrState"] = v.ReturnProcessingCentreState
		m["delistRsn"] = v.DelistingReason
		m["pickupTime"] = v.PickupTime
		m["rtrnPolicy"] = v.ReturnPolicy
		m["pnltyCLause"] = v.PenaltyClause
		m["pnltyClause"] = v.PenaltyClause
		m["pickupPrtnrCode"] = v.PickUpPartnerCode
		m["reversepickUpCodeRPC"] = v.ReversePickUpRPC
		m["reversepickUpCodeOSS"] = v.ReversePickUpOneship
		m["dispatchLoc"] = v.DispatchLocation
		csvStruct[k].SellerCustomInfo = m
	}
	return csvStruct
}

//Deleting key from cache
func delete(keys []string) {
	err := cacheObj.DeleteBatch(keys)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while deleting keys from cache: %s", err.Error()))
	}
}

//initializes cache object
func initializeCacheObj() error {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting cache object: %s", err.Error()))
		return err
	}
	return nil
}
