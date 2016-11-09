package post

import (
	"amenities/brands/common"
	"common/utils"
	"encoding/json"
	"errors"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	s "strings"
	"unicode"
)

//Struct for Failure (node based)
type Validate struct {
	id string
}

//Function to SetID for current node from orchestrator
func (v *Validate) SetID(id string) {
	v.id = id
}

//Function that returns current node ID to orchestrator
func (v Validate) GetID() (id string, err error) {
	return v.id, nil
}

//Function that returns node name to orchestrator
func (v Validate) Name() string {
	return "Validate node for POST"
}

//Function to start node execution
func (v Validate) GetDecision(io workflow.WorkFlowData) (bool, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.VALIDATE_BRAND_POST)
	defer func() {
		logger.EndProfile(profiler, common.VALIDATE_BRAND_POST)
	}()

	var errors []map[string]interface{}
	flag := false
	data, _ := utils.GetPostData(io)

	brandCreate := new(common.BrandData)
	err := json.Unmarshal(data, &brandCreate)

	if err != nil {
		return false, err
	}

	for _, v := range brandCreate.Branddata {
		errMap := validate(v)
		if len(errMap) != 0 {
			flag = true
		}
		errors = append(errors, errMap)
	}
	if flag == false {
		for k, v := range brandCreate.Branddata {
			status, err := validateStatus(v)
			if err != nil {
				return false, err
			}
			brandCreate.Branddata[k].Status = status
			brandCreate.Branddata[k].UrlKey = cleanUp(brandCreate.Branddata[k].UrlKey)
		}
		io.IOData.Set(common.BRAND_DATA, brandCreate.Branddata)
		return true, nil
	}
	io.IOData.Set(common.FAILURE_DATA, errors)
	logger.Info("Brand data extracted")

	return false, err
}

//Function to ensure presence of mandatory fields
func validate(brandInfo common.Brand) map[string]interface{} {
	errorMap := make(map[string]interface{})

	if brandInfo.Status == "" {
		errorMap["status"] = "Status is Missing"
	}

	if brandInfo.Name == "" {
		errorMap["name"] = "Brand Name is Missing"
	} else {
		found, brandData, _ := common.CheckIfKeyExists(bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: brandInfo.Name, Options: "i"}}})
		if found && len(brandData) >= 1 {
			errorMap["name"] = "Brand Name already exists"
		}
	}

	if brandInfo.UrlKey == "" {
		errorMap["urlKey"] = "UrlKey is Missing"
	} else {
		found, brandData, _ := common.CheckIfKeyExists(bson.M{"urlKey": bson.M{"$regex": bson.RegEx{Pattern: cleanUp(brandInfo.UrlKey), Options: "i"}}})
		if found && len(brandData) >= 1 {
			errorMap["urlKey"] = "Url Key already exists"
		}
	}

	if brandInfo.BrandClass == "" {
		errorMap["brandClass"] = "Brand Class is Missing"
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

func cleanUp(url string) string {
	var newUrl string
	for _, v := range url {
		if string(v) == " " {
			newUrl += "-"
			continue
		}
		if unicode.IsLetter(v) || unicode.IsNumber(v) {
			newUrl += string(v)
		}
	}
	newUrl = s.Replace(newUrl, "--", "-", -1)
	return s.Trim(newUrl, "-")
}
