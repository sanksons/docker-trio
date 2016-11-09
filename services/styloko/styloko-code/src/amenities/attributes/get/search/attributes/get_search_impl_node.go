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
	"strconv"
)

type SearchAttribute struct {
	id string
}

func (s *SearchAttribute) SetID(id string) {
	s.id = id
}

func (s SearchAttribute) GetID() (id string, err error) {
	return s.id, nil
}

func (s SearchAttribute) Name() string {
	return "Search Attribute Set"
}

func (s SearchAttribute) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
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

	v, ok := s.All(isGlobal)
	if ok {
		io.IOData.Set(florest_constants.RESULT, v)
		return io, nil
	}
	return io, &florest_constants.AppError{
		Code:             appconstant.ResourceNotFoundCode,
		Message:          "Data not found",
		DeveloperMessage: "Attribute does not exist for the given search criteria"}
}

//function to get all attributes on basis of global flag
func (s SearchAttribute) All(isGlobal int) ([]Attribute, bool) {
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
	criteria := M{"isGlobal": isGlobal}
	query.Criteria = criteria
	err := mongoDriver.FindAll(constants.ATTRIBUTES_COLLECTION, query, &attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("Error In reading data from mongo db : %s", err.Error()))
		return nil, false
	}
	return attributes, true
}
