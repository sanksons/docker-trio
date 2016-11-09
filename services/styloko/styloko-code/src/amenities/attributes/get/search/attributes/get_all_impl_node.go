package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	mongodb "common/mongodb"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type GetAllAttributes struct {
	id string
}

func (s *GetAllAttributes) SetID(id string) {
	s.id = id
}

func (s GetAllAttributes) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetAllAttributes) Name() string {
	return "GET all Attributes"
}

func (s GetAllAttributes) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ALL)
	defer func() {
		logger.EndProfile(profiler, GET_ALL)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(GET_ALL, "Attribute get all execution started")
	v, ok := s.All()
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data set is empty",
			DeveloperMessage: "No data for given parameter"}
	}
	io.IOData.Set(florest_constants.RESULT, v)
	return io, nil
}

//fucntion to get all attributes
func (s GetAllAttributes) All() ([]Attribute, bool) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ALL)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ALL)
		mongoDriver.Close()
	}()
	var query mongodb.Query
	var attributes []Attribute
	query.Sort = []string{"seqId"}
	err := mongoDriver.FindAll(constants.ATTRIBUTES_COLLECTION, query, &attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting all attributes from mongo :%s", err.Error()))
		return nil, false
	}
	return attributes, true
}
