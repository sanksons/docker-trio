package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"fmt"
	"time"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type InsertAttribute struct {
	id string
}

func (u *InsertAttribute) SetID(id string) {
	u.id = id
}

func (u InsertAttribute) GetID() (id string, err error) {
	return u.id, nil
}

func (u InsertAttribute) Name() string {
	return "Attribute Insert api"
}

func (u InsertAttribute) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_INSERT)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_INSERT)
	}()
	io.ExecContext.SetDebugMsg(ATTRIBUTE_INSERT, "Attribute Insert Execute")

	d, _ := io.IOData.Get(INSERT_DATA)
	formData, ok := d.(*Attribute)
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Form Data is not valid"}
	}

	err := u.validateOptions(formData.Options)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Option validation failed. Please check options array.",
			DeveloperMessage: err.Error()}
	}
	formData.Options = u.injectSeqIds(formData.Options)
	setData, _ := io.IOData.Get(POST_SET_DATA)
	set, ok := setData.(Set)
	if !ok {
		logger.Error(fmt.Println("Error while asseting set struct"))
	}
	res, uerr := u.Insert(formData, &set)
	if uerr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Parameters",
			DeveloperMessage: uerr.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, res)
	return io, nil
}

func (u InsertAttribute) Insert(formData *Attribute, set *Set) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_INSERT)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_INSERT)
		mongoDriver.Close()
	}()
	wAttr, err := u.TransformData(formData, set)
	if err != nil {
		return nil, err
	}
	updateQuery := M{"$set": wAttr}
	query := M{"seqId": wAttr.SeqId}
	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTES_COLLECTION, updateQuery, query, true)
	retData := make(map[string]interface{}, 1)
	retData["seqId"] = wAttr.SeqId
	if err != nil {
		return res, err
	}
	return retData, nil
}

// validateOptions -> validates options array
func (u InsertAttribute) validateOptions(options []Option) error {
	for _, x := range options {
		if x.Value == "" {
			return fmt.Errorf("Value cannot be empty")
		}
		if x.Position < 0 {
			return fmt.Errorf("Position cannot be negative")
		}
		if x.IsDefault > 1 || x.IsDefault < 0 {
			return fmt.Errorf("IsDefault takes only 0, 1. %d is not a valid value", x.IsDefault)
		}
	}
	return nil
}

// injectSeqIds injects seqIds to options array
func (u InsertAttribute) injectSeqIds(options []Option) []Option {
	opt := []Option{}
	for x := range options {
		options[x].SeqId = x + 1
		opt = append(opt, options[x])
	}
	return opt
}

func (u InsertAttribute) TransformData(formData *Attribute, set *Set) (WAttribute, error) {
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer mongoDriver.Close()
	var wAttr WAttribute
	wAttr.SeqId = mongoDriver.GetNextSequence(constants.ATTRIBUTES_COLLECTION)
	wAttr.IsGlobal = formData.IsGlobal
	wAttr.Name = formData.Name
	if formData.IsGlobal == 0 {
		wAttr.Set = *set
	}
	wAttr.Label = formData.Label
	wAttr.Description = formData.Description
	wAttr.ProductType = formData.ProductType
	wAttr.AttributeType = formData.AttributeType
	wAttr.MaxLength = formData.MaxLength
	wAttr.DecimalPlaces = formData.DecimalPlaces
	wAttr.DefaultValue = formData.DefaultValue
	wAttr.UniqueValue = formData.UniqueValue
	wAttr.PetType = formData.PetType
	wAttr.PetMode = formData.PetMode
	wAttr.Validation = formData.Validation
	wAttr.Mandatory = formData.Mandatory
	wAttr.MandatoryImport = formData.MandatoryImport
	wAttr.AliceExport = formData.AliceExport
	wAttr.PetQc = formData.PetQc
	wAttr.ImportConfigIdentifier = formData.ImportConfigIdentifier
	wAttr.SolrSearchable = formData.SolrSearchable
	wAttr.SolrFilter = formData.SolrFilter
	wAttr.SolrSuggestions = formData.SolrSuggestions
	wAttr.Visible = formData.Visible
	wAttr.IsActive = formData.IsActive
	wAttr.FilterType = formData.FilterType
	wAttr.CreatedAt = time.Now()
	wAttr.UpdatedAt = wAttr.CreatedAt
	if formData.AttributeType == "option" || formData.AttributeType == "multi_option" {
		wAttr.Options = formData.Options
	}
	return wAttr, nil
}
