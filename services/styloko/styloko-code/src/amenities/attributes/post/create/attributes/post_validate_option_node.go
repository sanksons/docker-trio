package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

type ValidateOption struct {
	id string
}

func (d *ValidateOption) SetID(id string) {
	d.id = id
}

func (d ValidateOption) GetID() (id string, err error) {
	return d.id, nil
}

func (d ValidateOption) Name() string {
	return "Validate Attribute-option Insertion"
}

func (d ValidateOption) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_POST)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_POST)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(ATTRIBUTE_OPTION_INSERT, "Attribute-option insertion validation started")
	logger.Info("Validate Attribute-option started")

	pathParams, _ := io.IOData.Get(PATH_PARAMETERS)
	params := pathParams.(Parameters)

	attrData, err, ok := d.CheckIfExists(params)
	if !ok {
		io.IOData.Set(FAILURE_DATA, err)
		return false, nil
	}

	data, _ := utils.GetPostData(io)
	formData := []Option{}
	errData := json.Unmarshal(data, &formData)
	if errData != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Json",
			DeveloperMessage: "Error while unmarshalling json"}
	}
	//Validate the request
	errMsg, ok := d.ValidateJson(formData)
	if !ok {
		io.IOData.Set(FAILURE_DATA, errMsg)
		return false, nil
	}
	errMap, ok := d.ValidateData(params.AttrId, formData, *attrData)
	if ok {
		io.IOData.Set(INSERT_DATA, formData)
		return true, nil
	}
	io.IOData.Set(FAILURE_DATA, errMap)
	return false, nil
}

func (d ValidateOption) ValidateJson(reqBody []Option) (
	florest_constants.AppErrors, bool) {
	errors := florest_constants.AppErrors{}
	isValidationSuccess := true
	for _, x := range reqBody {
		errs := validate.Struct(x)
		if errs != nil {
			isValidationSuccess = false
			validationErrors := errs.(validator.ValidationErrors)
			msgs := d.PrepareErrorMessages(validationErrors)
			errors.Errors = append(errors.Errors, florest_constants.AppError{
				Code:             appconstant.InvalidDataErrorCode,
				Message:          "Validation Failed For Input Json",
				DeveloperMessage: strings.Join(msgs, "**"),
			})
		}
	}
	if !isValidationSuccess {
		return errors, false
	}
	return errors, true
}

func (d ValidateOption) PrepareErrorMessages(errs validator.ValidationErrors) []string {
	var msgs []string
	for _, err := range errs {
		var msg string
		switch err.Tag {
		case "required":
			msg = err.Field + ": Is Required."
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

func (d ValidateOption) ValidateData(attrId int, data []Option, attrInfo CheckOptions) (
	florest_constants.AppErrors, bool) {
	var errorMap []string
	errors := florest_constants.AppErrors{}
	if !(attrInfo.AttributeType == "multi_option" || attrInfo.AttributeType == "option") {
		errorMap = append(errorMap, "Options cannot be added in this attribute")
	}
	for _, x := range data {
		if x.Value == "" {
			errorMap = append(errorMap, "Value can not be updated")
		}
	}

	//check option value exist or not
	for _, x := range data {
		err, ok := d.CheckValueExist(attrId, x)
		if !ok {
			errorMap = append(errorMap, err.Error())
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
func (d ValidateOption) CheckValueExist(attrId int, data Option) (error, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_CHECK_VALUE_EXIST)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, POST_CHECK_VALUE_EXIST)
		mongoDriver.Close()
	}()
	var options CheckOptions
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	err := attributeObj.
		Find(bson.M{"seqId": attrId,
			"options": bson.M{
				"$elemMatch": bson.M{
					"value": data.Value,
				},
			},
		}).Select(bson.M{
		"_id": 0,
		"options": bson.M{
			"$elemMatch": bson.M{
				"value": data.Value,
			},
		},
	}).One(&options)

	if err != nil && err.Error() != "not found" {
		logger.Error(fmt.Sprintf("Error while getting attribute from mongo :%s", err.Error()))
		return err, false
	}

	if len(options.Options) > 0 {
		return errors.New("Option already exists."), false
	}

	return nil, true
}

func (d ValidateOption) CheckIfExists(params Parameters) (*CheckOptions,
	florest_constants.AppErrors, bool) {
	var res CheckOptions
	errors := florest_constants.AppErrors{}
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_CHECK_IF_EXISTS)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, POST_CHECK_IF_EXISTS)
		mongoDriver.Close()
	}()
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	err := attributeObj.
		Find(bson.M{"seqId": params.AttrId}).
		Select(bson.M{
			"_id":           0,
			"attributeType": 1,
			"options":       1,
		}).One(&res)

	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attribute from mongo :%s", err.Error()))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Error while getting id : " + strconv.Itoa(params.AttrId),
			DeveloperMessage: err.Error(),
		})
		return nil, errors, false
	}

	if res.AttributeType == "" {
		logger.Error("Attribute does not exist having id : " + strconv.Itoa(params.AttrId))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Attribute does not exist having id : " + strconv.Itoa(params.AttrId),
			DeveloperMessage: "Attribute does not exist",
		})
		return nil, errors, false
	}
	return &res, errors, true
}
