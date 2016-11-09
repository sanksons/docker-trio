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
	return "UPDATE Attribute Set by id"
}

func (d Validation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_PUT)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_PUT)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(SET_UPDATE, "Attribute Set update validation started")
	logger.Info("Validate Attribute Set started")

	pathParams, _ := utils.GetPathParams(io)
	setId, pthErr, ok := d.GetPathParmas(pathParams)
	if !ok {
		io.IOData.Set(FAILURE_DATA, pthErr)
		return false, nil
	}

	rset, err, ok := d.CheckIfSetExists(setId)
	if !ok {
		io.IOData.Set(FAILURE_DATA, err)
		return false, nil
	}
	io.IOData.Set(PATH_PARAMETERS, setId)

	data, _ := utils.GetPostData(io)
	setData := new(AttributeSet)
	jsonErr := json.Unmarshal(data, &setData)
	if jsonErr != nil {
		return false, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Invalid Json",
			DeveloperMessage: "Error while unmarshalling json"}
	}
	//Validate the request
	errMsg, ok := d.ValidateJson(setData)
	if !ok {
		io.IOData.Set(FAILURE_DATA, errMsg)
		return false, nil
	}

	//validate values
	derr, ok := d.ValidateData(setData, rset)
	if ok {
		io.IOData.Set(UPDATE_DATA, setData)
		return true, nil
	}
	io.IOData.Set(FAILURE_DATA, derr)
	return false, nil
}

func (d Validation) ValidateJson(reqBody *AttributeSet) (
	florest_constants.AppErrors, bool) {
	errors := florest_constants.AppErrors{}
	isValidationSuccess := true
	errs := validate.Struct(reqBody)
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
	if !isValidationSuccess {
		return errors, false
	}
	return errors, true
}

func (d Validation) ValidateData(wAttr *AttributeSet, rAttr *AttributeSet) (
	florest_constants.AppErrors, bool) {
	var errorMap []string
	errors := florest_constants.AppErrors{}
	if wAttr.Name != "" && rAttr.Name != wAttr.Name {
		errorMap = append(errorMap, "Name can not be updated")
	}
	if wAttr.Label == "" {
		errorMap = append(errorMap, "Label cannot be updated")
	}
	if wAttr.Identifier == "" {
		errorMap = append(errorMap, "Identifier cannot be updated")
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

func (d Validation) CheckIfSetExists(setId int) (*AttributeSet,
	florest_constants.AppErrors, bool) {
	profiler := logger.NewProfiler()
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	logger.StartProfile(profiler, PUT_CHECK_IF_SET_EXISTS)
	defer func() {
		logger.EndProfile(profiler, PUT_CHECK_IF_SET_EXISTS)
		mongoDriver.Close()
	}()
	errors := florest_constants.AppErrors{}
	var attributeSet AttributeSet
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(bson.M{"seqId": setId}).One(&attributeSet)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attribute set details from mongo :%s", err.Error()))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Attibute set does not exist having id : " + strconv.Itoa(setId),
			DeveloperMessage: "Attribute set does not exist in database",
		})
		return nil, errors, false
	}
	if attributeSet.Name == "" {
		logger.Error("Attibute set does not exist having id : " + strconv.Itoa(setId))
		errors.Errors = append(errors.Errors, florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Attibute set does not exist having id : " + strconv.Itoa(setId),
			DeveloperMessage: "Attribute set does not exist in database",
		})
		return nil, errors, false
	}
	return &attributeSet, errors, true
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

func (d Validation) GetPathParmas(pathParams []string) (int,
	florest_constants.AppErrors, bool) {
	var seqId int
	msg := "Path Parmeters are not valid"
	errors := florest_constants.AppErrors{}
	if len(pathParams) == 1 {
		seqId, err := strconv.Atoi(pathParams[0])
		if err != nil {
			msg = "Set Id cannot convert into int " + pathParams[0]
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
