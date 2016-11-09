package put

import (
	"amenities/brands/common"
	"common/utils"
	"encoding/json"
	"errors"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	s "strings"
)

// Decision basic struct for decision node
type Decision struct {
	id string
}

// SetID sets ID for node
func (d *Decision) SetID(id string) {
	d.id = id
}

// GetID returns ID for node
func (d Decision) GetID() (id string, err error) {
	return d.id, nil
}

// Name returns name for node
func (d Decision) Name() string {
	return "UPDATE brand by id"
}

// GetDecision return bool on condition
func (d Decision) GetDecision(io workflow.WorkFlowData) (bool, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.VALIDATE_PUT)
	defer func() {
		logger.EndProfile(profiler, common.VALIDATE_PUT)
	}()

	//Enable logs for Debugging
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(common.BRAND_OPERATION, "Brand update validation started")

	logger.Info("Validate brand started")
	var errors []map[string]interface{}
	flag := false
	data, _ := utils.GetPostData(io)
	brandUpdate := new(common.BrandData)
	err := json.Unmarshal(data, &brandUpdate)
	if err != nil {
		errorMap := make(map[string]interface{})
		errorMap["JSON"] = "Unmarshal failure. Bad data."
		errors = append(errors, errorMap)
		io.IOData.Set(common.FAILURE_DATA, errors)
		return false, nil
	}

	//calling validate to check seqId
	for _, v := range brandUpdate.Branddata {
		errMap := d.Validate(v)
		if len(errMap) != 0 {
			flag = true
		}
		errors = append(errors, errMap)
	}

	//checking if flag was not unset during validation
	if flag == false {
		for k, v := range brandUpdate.Branddata {
			status, err := validateStatus(v)
			if err != nil {
				return false, err
			}
			brandUpdate.Branddata[k].Status = status
		}
		io.IOData.Set(common.BRAND_UPDATE_DATA, brandUpdate.Branddata)
		logger.Info("Brand data extracted")
		return true, nil
	}
	//if flag was set, setting errors in FAILURE_DATA
	io.IOData.Set(common.FAILURE_DATA, errors)
	return false, err
}

// Validate returns validation failure errors in map[string]interface{}
func (d Decision) Validate(brandInfo common.Brand) map[string]interface{} {
	errorMap := make(map[string]interface{})
	if brandInfo.SeqId == 0 {
		errorMap["seqId"] = "Sequence Id is Mandatory."
	}
	if brandInfo.Name != "" {
		found, brandData, _ := common.CheckIfKeyExists(bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: brandInfo.Name, Options: "i"}}})
		if found && len(brandData) > 1 && brandInfo.SeqId != brandData[0].SeqId {
			errorMap["name"] = "Brand Name already exists"
		}
	}

	if brandInfo.UrlKey != "" {
		found, brandData, _ := common.CheckIfKeyExists(bson.M{"urlKey": bson.M{"$regex": bson.RegEx{Pattern: brandInfo.UrlKey, Options: "i"}}})
		if found && len(brandData) > 1 && brandInfo.SeqId != brandData[0].SeqId {
			errorMap["urlKey"] = "Url Key already exists"
		}
	}
	return errorMap
}

func validateStatus(brandInfo common.Brand) (string, error) {
	status := s.ToLower(brandInfo.Status)
	if status == "active" || status == "deleted" || status == "inactive" {
		return status, nil
	}
	return brandInfo.Status, errors.New("Invalid status provided.")
}
