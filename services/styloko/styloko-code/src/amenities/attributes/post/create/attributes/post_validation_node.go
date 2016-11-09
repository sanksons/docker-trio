package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"common/utils"
	"encoding/json"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
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
	return "Validate Attribute Node For Insertion"
}

func (d Validation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_POST)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_POST)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(ATTRIBUTE_INSERT, "Attribute insertion validation started")
	logger.Info("Validate Attribute Set started")

	data, err := utils.GetPostData(io)
	if err != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Post Data",
			DeveloperMessage: "Error while getting post data"}
	}

	formData := new(Attribute)
	err = json.Unmarshal(data, &formData)
	if err != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Json",
			DeveloperMessage: "Error while unmarshalling json"}
	}

	errMap, ok := d.ValidateData(formData)
	if !ok {
		io.IOData.Set(FAILURE_DATA, errMap)
		return false, nil
	}
	if formData.IsGlobal == 0 && formData.Set != "" {
		setData, setErr, ok := d.CheckIfSetExists(formData.Set)
		if !ok {
			io.IOData.Set(FAILURE_DATA, setErr)
			return false, nil
		}
		io.IOData.Set(POST_SET_DATA, *setData)
	}
	io.IOData.Set(INSERT_DATA, formData)
	return true, nil
}

func (d Validation) ValidateData(wAttr *Attribute) (
	florest_constants.AppErrors, bool) {
	var errorMap []string
	errors := florest_constants.AppErrors{}
	if wAttr.Name == nil {
		errorMap = append(errorMap, "Name cannot be null")
	}

	if wAttr.Label == nil {
		errorMap = append(errorMap, "Label cannot be null")
	}

	if wAttr.Mandatory == nil {
		errorMap = append(errorMap, "Mandatory field cannot be null")
	}

	ok := d.CheckIfExists(*wAttr.Name, wAttr.Set, wAttr.IsGlobal)
	if ok {
		errorMap = append(errorMap, "Attribute Already exists")
	}

	if wAttr.IsGlobal == 0 && wAttr.Set == "" {
		errorMap = append(errorMap, "Please enter specific set or global flag")
	}

	if wAttr.IsGlobal == 1 && wAttr.Set != "" {
		errorMap = append(errorMap, "Attribute can be created in specific set or global flag")
	}

	ok = d.stringInSlice(wAttr.ProductType, []string{"config", "simple", "source"})
	if !ok {
		errorMap = append(errorMap, "Please enter valid product type")
	}

	ok = d.stringInSlice(wAttr.AliceExport, []string{"no", "meta", "attribute"})
	if !ok {
		errorMap = append(errorMap, "Please enter valid value for alice export")
	}

	if wAttr.PetMode != "" {
		ok = d.stringInSlice(wAttr.PetMode, []string{"edit", "display", "invisible", "edit_on_create"})
		if !ok {
			errorMap = append(errorMap, "Please enter alice export value")
		}
	}

	ok = d.stringInSlice(wAttr.PetType, []string{"textfield", "textarea", "numberfield", "datefield", "datetime", "checkbox", "dropdown", "multiselect", "combo", "multicombo", "pricecost"})
	if !ok {
		errorMap = append(errorMap, "Please enter valid value for pet type")
	} else if wAttr.PetType == "textfield" && wAttr.MaxLength == nil {
		errorMap = append(errorMap, "Please enter maxlength")
	} else if wAttr.PetType == "numberfield" && (wAttr.MaxLength == nil || wAttr.DecimalPlaces == nil) {
		errorMap = append(errorMap, "Please enter maxlength and decimal places for numberfield")
	}

	if wAttr.Validation != nil {
		ok = d.stringInSlice(*wAttr.Validation, []string{"decimal", "integer", "percent", "email", "url", "letters", "lettersnumbers"})
		if !ok {
			errorMap = append(errorMap, "Please enter valid value for validation")
		}
	}

	ok = d.stringInSlice(wAttr.AttributeType, []string{"system", "option", "multi_option", "value", "custom"})
	if !ok {
		errorMap = append(errorMap, "Please enter valid value for attribute type")
	}

	if wAttr.UniqueValue != nil {
		ok = d.stringInSlice(*wAttr.UniqueValue, []string{"global", "config"})
		if !ok {
			errorMap = append(errorMap, "Please enter a valid value for Unique Value")
		}
	}

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

func (d Validation) CheckIfExists(name string, setName string, isGlobal int) bool {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_CHECK_IF_EXISTS)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, POST_CHECK_IF_EXISTS)
		mongoDriver.Close()
	}()
	var attribute Attribute
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	if isGlobal == 1 {
		err := attributeObj.Find(bson.M{"isGlobal": 1, "name": name}).One(&attribute)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while getting attributes from mongo :%s", err.Error()))
		}
	} else {
		err := attributeObj.Find(bson.M{"isGlobal": 0, "name": name, "set.name": setName}).One(&attribute)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while getting attributes from mongo :%s", err.Error()))
		}
	}

	if attribute.Name != nil {
		return true
	}
	return false
}

func (d Validation) stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (d Validation) CheckIfSetExists(name string) (
	*Set, florest_constants.AppErrors, bool) {
	errors := florest_constants.AppErrors{}
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_CHECK_SET_EXISTS)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, POST_CHECK_SET_EXISTS)
		mongoDriver.Close()
	}()
	var attributeSet Set
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(bson.M{"name": name}).One(&attributeSet)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while finding attribute set in mongo :%s", err.Error()))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Error while getting from mongo for name : " + name,
			DeveloperMessage: err.Error(),
		})
		return nil, errors, false
	}
	if attributeSet.Name == nil {
		logger.Error("Set does not exist with provided name : " + name)
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Set does not exist with provided name : " + name,
			DeveloperMessage: "Attribute-Set does not exist",
		})
		return nil, errors, false
	}
	return &attributeSet, errors, true
}
