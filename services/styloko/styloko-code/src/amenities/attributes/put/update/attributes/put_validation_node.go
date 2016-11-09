package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"common/utils"
	"encoding/json"
	_ "errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
)

type Validation struct {
	id string
}

func (d *Validation) SetID(id string) {
	d.id = id
}

func (d Validation) GetID() (id string, err error) {
	return d.id, nil
}

func (d Validation) Name() string {
	return "Validate Attribute by id"
}

func (d Validation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_PUT)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_PUT)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(ATTRIBUTE_UPDATE, "Attribute update validation started")
	logger.Info("Validate Attribute started")

	pathParams, _ := utils.GetPathParams(io)
	seqId, patheErr, ok := d.GetPathParmas(pathParams)
	if !ok {
		io.IOData.Set(FAILURE_DATA, patheErr)
		return false, nil
	}

	io.IOData.Set(PATH_PARAMETERS, seqId)

	//validation skipped altogther
	rset, err, ok := d.CheckIfExists(seqId)
	if !ok {
		io.IOData.Set(FAILURE_DATA, err)
		return false, nil
	}

	data, _ := utils.GetPostData(io)
	formData := new(Attribute)
	dataEr := json.Unmarshal(data, &formData)
	if dataEr != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Json",
			DeveloperMessage: "Error while unmarshalling json"}
	}

	errMap, ok := d.Validate(formData, rset)
	if ok {
		io.IOData.Set(UPDATE_DATA, formData)
		return true, nil
	}

	io.IOData.Set(FAILURE_DATA, errMap)
	return false, nil
}

func (d Validation) Validate(uAttr *Attribute, ca *Attribute) (
	florest_constants.AppErrors, bool) {
	var errorMap []string
	errors := florest_constants.AppErrors{}
	// if uAttr.Name != "" && uAttr.Name != ca.Name {
	// 	errorMap = append(errorMap, "Name cannot be updated")
	// }

	// if uAttr.IsGlobal != ca.IsGlobal {
	// 	errorMap = append(errorMap, "Global Flag cannot be updated")
	// }

	// if ca.AttributeType == "option" || ca.AttributeType == "value" {
	// 	if uAttr.PetMode != "" {
	// 		errorMap = append(errorMap, "Pet Mode cannot be updated")
	// 	}
	// }
	// if ca.AttributeType == "system" {
	// 	if uAttr.PetMode != "" {
	// 		errorMap = append(errorMap, "Pet mode cannot be updated for system attributes")
	// 	}
	// 	if uAttr.DefaultValue != nil {
	// 		errorMap = append(errorMap, "Default value cannot be updated for system attributes")
	// 	}
	// 	if uAttr.Validation != nil {
	// 		errorMap = append(errorMap, "Validation cannot be updated for system attributes")
	// 	}
	// 	if uAttr.UniqueValue != nil {
	// 		errorMap = append(errorMap, "Unique value cannot be updated for system attributes")
	// 	}
	// 	if uAttr.Mandatory != nil {
	// 		errorMap = append(errorMap, "Mandatory cannot be updated for system attributes")
	// 	}
	// 	if uAttr.MandatoryImport != nil {
	// 		errorMap = append(errorMap, "Mandatory import cannot be updated for system attributes")
	// 	}
	// 	if uAttr.AliceExport != "" {
	// 		errorMap = append(errorMap, "Alice export cannot be updated for system attributes")
	// 	}
	// 	if uAttr.SolrSearchable != nil || uAttr.SolrSuggestions != nil || uAttr.SolrFilter != nil {
	// 		errorMap = append(errorMap, "Solr fields cannot be updated for system attributes")
	// 	}

	// }
	// if ca.UniqueValue != nil && *ca.UniqueValue == "config" {
	// 	if uAttr.AliceExport != "" {
	// 		errorMap = append(errorMap, "Alice export cannot be updated with unique value as config")
	// 	}
	// 	if uAttr.MandatoryImport != nil {
	// 		errorMap = append(errorMap, "Mandatory import cannot be updated with unique value as config")
	// 	}
	// 	if uAttr.Mandatory != nil {
	// 		errorMap = append(errorMap, "Mandatory cannot be updated with unique value as config")
	// 	}
	// }

	// //explicilty putting as "" so it is ignored in omitempty tag and never updated
	// if uAttr.ProductType != "" {
	// 	uAttr.ProductType = ""
	// }

	// if uAttr.AliceExport != "" {
	// 	ok := d.stringInSlice(uAttr.AliceExport, []string{"no", "meta", "attribute"})
	// 	if !ok {
	// 		errorMap = append(errorMap, "Please enter valid value for alice export")
	// 	}
	// }

	// if uAttr.PetMode != "" {
	// 	ok := d.stringInSlice(uAttr.PetMode, []string{"edit", "display", "invisible", "edit_on_create"})
	// 	if !ok {
	// 		errorMap = append(errorMap, "Please enter alice export value")
	// 	}
	// }

	// if uAttr.PetType != nil {
	// 	uAttr.PetType = nil
	// }

	// if uAttr.Validation != nil {
	// 	uAttr.Validation = nil
	// }

	// if uAttr.AttributeType != "" {
	// 	uAttr.AttributeType = ""
	// }

	// if uAttr.UniqueValue != nil {
	// 	ok := d.stringInSlice(*uAttr.UniqueValue, []string{"global", "config"})
	// 	if !ok {
	// 		errorMap = append(errorMap, "Please enter a valid value for Unique Value")
	// 	}
	// }

	if len(errorMap) > 0 {
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Validation Failed For Input Json",
			DeveloperMessage: strings.Join(errorMap, "**"),
		})
		return errors, false
	}
	return errors, true
}

func (d Validation) CheckIfExists(seqId int) (*Attribute,
	florest_constants.AppErrors, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PUT_CHECK_IF_EXISTS)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, PUT_CHECK_IF_EXISTS)
		mongoDriver.Close()
	}()
	errors := florest_constants.AppErrors{}
	var attribute Attribute
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	err := attributeObj.Find(bson.M{"seqId": seqId}).One(&attribute)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attributes from mongo :%s", err.Error()))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data not found",
			DeveloperMessage: "Unable to fetch from mongo for id " + strconv.Itoa(seqId),
		})
		return nil, errors, false
	}
	if attribute.Name == "" {
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Data exist check failed",
			DeveloperMessage: "Attibute does not exist having id " + strconv.Itoa(seqId),
		})
		return nil, errors, false
	}
	return &attribute, errors, true
}

func (d Validation) GetPathParmas(pathParams []string) (
	int, florest_constants.AppErrors, bool) {
	var seqId int
	msg := "Path Parmeters are not valid"
	errors := florest_constants.AppErrors{}
	if len(pathParams) == 1 {
		seqId, err := strconv.Atoi(pathParams[0])
		if err != nil {
			msg = "Attribute Id cannot be converted into int " + pathParams[0]
		}
		return seqId, errors, true
	}
	errors.Errors = append(errors.Errors, florest_constants.AppError{
		Code:             appconstant.BadRequestCode,
		Message:          "Path Parmeters Validation Failed",
		DeveloperMessage: msg,
	})
	return seqId, errors, false
}

func (d Validation) stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
