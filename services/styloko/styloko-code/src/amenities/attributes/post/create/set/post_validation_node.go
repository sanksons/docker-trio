package set

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
	validator "gopkg.in/go-playground/validator.v8"
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
	return "Insert Attribute Set by id"
}

func (d Validation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_PUT)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_PUT)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(SET_INSERT, "Attribute Set insertion validation started")
	logger.Info("Validate Attribute Set started")
	data, _ := utils.GetPostData(io)
	setData := new(SetRequestJson)
	err := json.Unmarshal(data, &setData)
	if err != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Json",
			DeveloperMessage: "Error while unmarshalling json"}
	}
	errr, ok := d.CheckIfSetExists(setData.Name)
	if ok {
		io.IOData.Set(FAILURE_DATA, errr)
		return false, nil
	}
	//Validate the request
	errMsg, ok := d.ValidateJson(setData)
	if ok {
		io.IOData.Set(INSERT_DATA, setData)
		return true, nil
	}
	io.IOData.Set(FAILURE_DATA, errMsg)
	return false, nil
}

func (d Validation) ValidateJson(reqBody *SetRequestJson) (
	florest_constants.AppErrors, bool) {
	errors := florest_constants.AppErrors{}
	isValidationSuccess := true
	errs := validate.Struct(reqBody)
	if errs != nil {
		isValidationSuccess = false
		validationErrors := errs.(validator.ValidationErrors)
		msgs := d.PrepareErrorMessages(validationErrors)
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             florest_constants.IncorrectDataErrorCode,
			Message:          "Validation Failed For Input Json",
			DeveloperMessage: strings.Join(msgs, "**"),
		})
	}
	if !isValidationSuccess {
		return errors, false
	}
	return errors, true
}

func (d Validation) PrepareErrorMessages(errs validator.ValidationErrors) []string {
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

func (d Validation) CheckIfSetExists(name string) (
	florest_constants.AppErrors, bool) {
	errors := florest_constants.AppErrors{}
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_CHECK_IF_SET_EXISTS)
	mongoDriver = factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, POST_CHECK_IF_SET_EXISTS)
		mongoDriver.Close()
	}()
	var attributeSet AttributeSet
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(bson.M{"name": name}).One(&attributeSet)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while finding attribute set in mongo :%s", err.Error()))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Error while getting from mongo for name : " + name,
			DeveloperMessage: err.Error(),
		})
		return errors, true
	}
	if attributeSet.Name != "" {
		logger.Error("Attibute Set already exists having name : " + name)
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Attribute set already exist having name : " + name,
			DeveloperMessage: "Attribute set exist",
		})
		return errors, true
	}
	return errors, false
}
