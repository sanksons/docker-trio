package set

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	mongodb "common/mongodb"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type GetAllAttributeSet struct {
	id string
}

func (s *GetAllAttributeSet) SetID(id string) {
	s.id = id
}

func (s GetAllAttributeSet) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetAllAttributeSet) Name() string {
	return "GET all Attribute Sets"
}

func (s GetAllAttributeSet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ALL)
	defer func() {
		logger.EndProfile(profiler, GET_ALL)
	}()

	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(GET_ALL, "AttributeSet get all execution started")

	v, ok := s.AllSets()
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.ResourceNotFoundCode,
			Message:          "Data set is empty",
			DeveloperMessage: "No data for given parameter"}
	}
	io.IOData.Set(florest_constants.RESULT, v)
	return io, nil
}

func (s GetAllAttributeSet) AllSets() ([]AttributeSet, bool) {
	var query mongodb.Query
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_ALL_SETS)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, GET_ALL_SETS)
		mongoDriver.Close()
	}()
	var attributeSet []AttributeSet
	query.Sort = []string{"seqId"}
	err := mongoDriver.FindAll(constants.ATTRIBUTESETS_COLLECTION, query, &attributeSet)
	if err != nil {
		logger.Error("Error in fetching data from database")
		return nil, false
	}
	if len(attributeSet) == 0 {
		return nil, false
	}
	return attributeSet, true
}
