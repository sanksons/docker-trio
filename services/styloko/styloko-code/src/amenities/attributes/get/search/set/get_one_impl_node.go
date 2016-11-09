package set

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	mongodb "common/mongodb"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"common/utils"

	"github.com/jabong/floRest/src/common/cache"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type GetAttributeSet struct {
	id string
}

func (s *GetAttributeSet) SetID(id string) {
	s.id = id
}

func (s GetAttributeSet) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetAttributeSet) Name() string {
	return "Get Attribute Set by id"
}

func (s GetAttributeSet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ONE)
	defer func() {
		logger.EndProfile(profiler, GET_ONE)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(GET_ONE, "Single attribute set get execution started")

	data, _ := io.IOData.Get(PATH_PARAMETERS)
	params, ok := data.(Parameters)
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.BadRequestCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Parameters are not valid"}
	}
	setId, err := strconv.Atoi(*params.SetId)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.BadRequestCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Set is not valid"}
	}

	switch params.Count {
	case 1:
		v, ok := s.SetById(setId)
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data not found",
			DeveloperMessage: "Attribute Set does not exist"}
	case 2:
		endPt := *params.EndPoint
		if endPt != ATTRIBUTES {
			return io, &florest_constants.AppError{
				Code:             appconstant.BadRequestCode,
				Message:          "Invalid path params",
				DeveloperMessage: "End point should be attributes"}
		}
		v, ok := s.AttributesBySetId(setId)
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data not found",
			DeveloperMessage: "Attributes for this set do not exist"}
	case 3:
		endPt := *params.EndPoint
		if endPt != ATTRIBUTES {
			return io, &florest_constants.AppError{
				Code:             appconstant.BadRequestCode,
				Message:          "Invalid path params",
				DeveloperMessage: "End point should be attributes"}
		}
		attrId, err := strconv.Atoi(*params.AttrIdName)
		var (
			v  *Attribute
			ok bool
		)
		if err != nil {
			v, ok = s.AttributeByName(setId, *params.AttrIdName)
		} else {
			v, ok = s.AttributeById(setId, attrId)
		}
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data not found",
			DeveloperMessage: "Attribute does not exist"}
	}

	return io, &florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          "Invalid Parameter",
		DeveloperMessage: "Param value in query is Invalid"}
}

func (s GetAttributeSet) SetById(setId int) (*AttributeSet, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_SET_BY_ID)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_SET_BY_ID)
		mongoDriver.Close()
	}()
	var attributeSet AttributeSet
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(bson.M{"seqId": setId}).One(&attributeSet)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attribute set from mongo :%s", err.Error()))
		return nil, false
	}
	if attributeSet.SeqId == 0 {
		return nil, false
	}
	return &attributeSet, true
}

func (s GetAttributeSet) AttributesBySetId(setId int) (*AttributeSets, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ATTRIBUTES_BY_SET_ID)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ATTRIBUTES_BY_SET_ID)
		mongoDriver.Close()
	}()
	var attributeSets AttributeSets
	attributeSetObj := mongoDriver.SetCollection(constants.ATTRIBUTESETS_COLLECTION)
	err := attributeSetObj.Find(bson.M{"seqId": setId}).One(&attributeSets)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attributes by setid from mongo :%s", err.Error()))
		return nil, false
	}
	if attributeSets.SeqId == 0 {
		return nil, false
	}
	attributeSets.Attributes = s.GetAttributes(attributeSets.Name)

	return &attributeSets, true
}

func (s GetAttributeSet) GetAttributes(setName string) []Attribute {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ATTRIBUTES)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ATTRIBUTES)
		mongoDriver.Close()
	}()
	var attributes []Attribute
	var query mongodb.Query
	query.Sort = []string{"seqId"}
	criteria := M{"set.name": setName}
	query.Criteria = criteria
	err := mongoDriver.FindAll(constants.ATTRIBUTES_COLLECTION, query, &attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Error In getting attributes from mongo :%s", err.Error()))
		return nil
	}
	return attributes
}

func (s GetAttributeSet) AttributeById(setId int, attrId int) (*Attribute, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ATTRIBUTE_BY_ID)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ATTRIBUTE_BY_ID)
		mongoDriver.Close()
	}()

	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	var (
		attribute Attribute
		err       error
	)
	if setId == 0 {
		err = attributeObj.Find(bson.M{"seqId": attrId, "isGlobal": 1}).One(&attribute)
	} else {
		err = attributeObj.Find(bson.M{"seqId": attrId, "set.seqId": setId}).One(&attribute)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attributes by id from mongo :%s", err.Error()))
		return nil, false
	}
	if attribute.SeqId == 0 {
		return nil, false
	}
	return &attribute, true
}

func (s GetAttributeSet) AttributeByName(setId int, attrName string) (*Attribute, bool) {
	key := fmt.Sprintf(constants.ATTR_CACHE_KEY_FORMAT_CRITERIA,
		attrName, ATTRIBUTE_PRODUCT_TYPE, setId)
	val, errCache := s.GetFromCache(key)
	if errCache == nil {
		return val, true
	}
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ATTRIBUTE_BY_NAME)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ATTRIBUTE_BY_NAME)
		mongoDriver.Close()
	}()

	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	var (
		attribute Attribute
		err       error
	)
	if setId == 0 {
		err = attributeObj.Find(bson.M{"name": attrName, "isGlobal": 1}).One(&attribute)
	} else {
		err = attributeObj.Find(bson.M{"name": attrName, "set.seqId": setId}).One(&attribute)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attributes by id from mongo :%s", err.Error()))
		return nil, false
	}
	if attribute.SeqId == 0 {
		return nil, false
	}
	go s.SetInCache(attribute)
	return &attribute, true
}

func (s GetAttributeSet) GetFromCache(key string) (*Attribute, error) {
	item, err := cacheObj.Get(key, false, false)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	d, ok := item.Value.(string)
	if !ok {
		return nil, errors.New("Assertion error")
	}
	attr := Attribute{}
	err = json.Unmarshal([]byte(d), &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

func (s GetAttributeSet) SetInCache(attr Attribute) {
	defer utils.RecoverHandler("SetInCache")

	e, err := json.Marshal(attr)
	if err != nil {
		logger.Error(fmt.Sprintf("(#SetInCache)Error in JSON marshalling:%v", err.Error()))
		return
	}
	val := string(e)
	i1 := cache.Item{
		Key: strings.ToLower(fmt.Sprintf(
			constants.ATTR_CACHE_KEY_FORMAT_CRITERIA,
			attr.Name, attr.ProductType, attr.AttrSet.SeqId,
		)),
		Value: val,
	}
	err = cacheObj.SetWithTimeout(i1, false, false, int32(constants.ATTR_CACHE_EXPIRY))
	if err != nil {
		logger.Error(fmt.Sprintf("(#SetInCache)Error in setting Cache:%v", err.Error()))
	}
}
