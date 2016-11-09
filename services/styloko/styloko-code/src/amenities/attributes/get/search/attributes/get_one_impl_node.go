package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type GetAttribute struct {
	id string
}

func (s *GetAttribute) SetID(id string) {
	s.id = id
}

func (s GetAttribute) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetAttribute) Name() string {
	return "Get Attribute Set by id"
}

func (s GetAttribute) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ONE)
	defer func() {
		logger.EndProfile(profiler, GET_ONE)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(GET_ONE, "Attribute get execution started")
	data, _ := io.IOData.Get(PATH_PARAMETERS)
	pathParams, ok := data.([]string)
	if ok {
		attrId, err := strconv.Atoi(pathParams[0])
		if err != nil {
			return io, &florest_constants.AppError{
				Code:             appconstant.InvalidDataErrorCode,
				Message:          "Error in reading param in int",
				DeveloperMessage: "Param value in query is Invalid"}
		}
		v, ok := s.ById(attrId)
		if ok {
			io.IOData.Set(florest_constants.RESULT, v)
			return io, nil
		}
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Error while getting attribute",
			DeveloperMessage: "Attribute does not exist"}
	}
	return io, &florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          "Invalid Attribute Id",
		DeveloperMessage: "Param value in query is Invalid"}
}

//get attribute by attribute id
func (s GetAttribute) ById(attrId int) (*Attribute, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_BY_ID)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, GET_BY_ID)
		mongoDriver.Close()
	}()
	var attribute Attribute
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	err := attributeObj.Find(bson.M{"seqId": attrId}).One(&attribute)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting attribute from mongo :%s", err.Error()))
		return nil, false
	}
	if attribute.SeqId == 0 {
		return nil, false
	}
	return &attribute, true
}
