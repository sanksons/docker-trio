package set

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	mongodb "common/mongodb"
	"common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type SearchAttributeSet struct {
	id string
}

func (s *SearchAttributeSet) SetID(id string) {
	s.id = id
}

func (s SearchAttributeSet) GetID() (id string, err error) {
	return s.id, nil
}

func (s SearchAttributeSet) Name() string {
	return "Search Attribute Set"
}

func (s SearchAttributeSet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_SEARCH)
	defer func() {
		logger.EndProfile(profiler, GET_SEARCH)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(GET_SEARCH, "Attribute search execution started")
	data, _ := io.IOData.Get(GET_SEARCH)
	isGlobal, err := strconv.Atoi(data.(string))
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.BadRequestCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Global flag should be int"}
	}
	if isGlobal == 1 {
		v, ok := s.AllAttributes(isGlobal)
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "No Data Found",
			DeveloperMessage: "Attributes not found for the passed query"}
	}
	if isGlobal == 0 {
		v, ok := s.SetWithAttibutes()
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "No Data Found",
			DeveloperMessage: "Attribute Sets with attributes not found for the passed query"}
	}
	return io, &florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          "Invalid Parameter",
		DeveloperMessage: "Param value in query is Invalid"}
}

//get all attributes acc to global flag
func (s SearchAttributeSet) AllAttributes(isGlobal int) ([]Attribute, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ALL_ATTRIBUTES)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ALL_ATTRIBUTES)
		mongoDriver.Close()
	}()
	var attributes []Attribute
	var query mongodb.Query
	query.Sort = []string{"seqId"}
	criteria := M{"isGlobal": isGlobal}
	query.Criteria = criteria
	err := mongoDriver.FindAll(constants.ATTRIBUTES_COLLECTION, query, &attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Error In getting all atrributes from mongo : %s", err.Error()))
		return nil, false
	}
	if len(attributes) == 0 {
		return nil, false
	}
	attributes, err = s.ChangeCatalogTyAttributeForSC(attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to get ty attribute. Error: %s", err.Error()))
		return nil, false
	}
	return attributes, true
}

//get all attributes with set data
func (s SearchAttributeSet) SetWithAttibutes() ([]AttributeSets, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_SET_WITH_ATTRIBUTES)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_SET_WITH_ATTRIBUTES)
		mongoDriver.Close()
	}()
	var attributeSets []AttributeSet
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(nil).All(&attributeSets)
	if err != nil {
		logger.Error(fmt.Sprintf("Error In getting all atrribute sets from mongo : %s", err.Error()))
		return nil, false
	}
	if len(attributeSets) == 0 {
		return nil, false
	}
	attributes := s.PrepareSet(attributeSets)
	if attributes == nil {
		return nil, false
	}
	return attributes, true
}

//prepare data set for sellercenter
func (s SearchAttributeSet) PrepareSet(attrSets []AttributeSet) []AttributeSets {
	var attributeSetsArr []AttributeSets
	for _, attrSet := range attrSets {
		attributes := s.GetAttributesBySetId(attrSet.Name)
		if attributes == nil {
			return nil
		}
		var temp AttributeSets
		temp.SeqId = attrSet.SeqId
		temp.GlobalIdentifier = attrSet.GlobalIdentifier
		temp.Label = attrSet.Label
		temp.Name = attrSet.Name
		temp.Attributes = attributes
		attributeSetsArr = append(attributeSetsArr, temp)
	}
	return attributeSetsArr
}

//get attribute by set and attribute id
func (s SearchAttributeSet) GetAttributesBySetId(setName string) []Attribute {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ATTRIBUTES_BY_SET_ID)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ATTRIBUTES_BY_SET_ID)
		mongoDriver.Close()
	}()
	var attributes []Attribute
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	err := attributeObj.Find(bson.M{"set.name": setName}).All(&attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Error In getting all atrributes by set id from mongo : %s", err.Error()))
		return nil
	}
	//
	// This is a PATCH for SC as SC does not behaves correctly for
	// old variation and size attributes. This should be fixed at SC
	// but at the time its not fixed I am handling it here.
	//
	//Remove unwanted attributes.
	var newAttributes []Attribute
	var variationSets = []string{"bags", "beauty", "fragrances", "home", "toys"}
	var sizeSets = []string{"sports_equipment"}

	for _, a := range attributes {
		if ((a.Name == "variation") && (utils.InArrayString(variationSets, setName))) ||
			((a.Name == "size") && (utils.InArrayString(sizeSets, setName))) {
			continue
		}
		newAttributes = append(newAttributes, a)
	}
	return newAttributes
}

// Returns the attributes list with modified ty attribute (system->option)
func (s SearchAttributeSet) ChangeCatalogTyAttributeForSC(attributes []Attribute) ([]Attribute, error) {
	var newAttributes []Attribute
	falseVar := false
	zeroVar := 0
	for _, attribute := range attributes {
		if attribute.Name == "ty" {
			mysqlDriver, err := factory.GetMySqlDriver("tyAttribute")
			if err != nil {
				return nil, err
			}
			rows, qErr := mysqlDriver.Query("SELECT id_catalog_ty, name from catalog_ty")
			if qErr != nil {
				return nil, fmt.Errorf(qErr.DeveloperMessage)
			}
			defer rows.Close()
			var option Option
			var tyOptions []Option
			for rows.Next() {
				var name string
				var id int
				sErr := rows.Scan(&id, &name)
				if sErr != nil {
					return nil, sErr
				}
				option.Name = &name
				option.IsDefault = &falseVar
				option.Position = &zeroVar
				option.SeqId = id
				tyOptions = append(tyOptions, option)
			}
			attribute.AttributeType = "option"
			attribute.Options = tyOptions
		}
		newAttributes = append(newAttributes, attribute)
	}
	return newAttributes, nil
}
