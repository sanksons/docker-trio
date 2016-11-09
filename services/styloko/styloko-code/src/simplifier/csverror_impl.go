package simplifier

import (
	"common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type csvError struct {
	id string
}

func (n *csvError) SetID(id string) {
	n.id = id
}

func (n csvError) GetID() (id string, err error) {
	return n.id, nil
}

func (a csvError) Name() string {
	return "csvError"
}

func (a csvError) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, CSV_ERRORS)
	defer func() {
		logger.EndProfile(profiler, CSV_ERRORS)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, CSV_ERRORS)
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+a.Name(), rc)
	io.ExecContext.SetDebugMsg(CSV_ERRORS, "Csv Error execution started")

	val, ok := utils.GetQueryParams(io, "jobName")
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid JobName", DeveloperMessage: "No query param passed for jobName"}
	}
	stringErr, err := a.GetErrorsByJobname(val)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting error by jobname:%s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting error by jobname", DeveloperMessage: err.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, stringErr)
	return io, nil
}

func (a csvError) GetErrorsByJobname(jobName string) (string, error) {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON_ERRORS)
	defer mgoSession.Close()
	var errDocs []ErrorStruct
	err := mgoObj.Find(bson.M{"jobname": jobName}).All(&errDocs)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting documents by JobName from mongo:%s", err.Error()))
		return "", err
	}
	if len(errDocs) == 0 {
		return "", errors.New("No document was found for passed JobName")
	}
	var stringErr string
	for _, v := range errDocs {
		stringErr += fmt.Sprintf("%s\n", v.ErrMsg)
	}
	return stringErr, nil
}
