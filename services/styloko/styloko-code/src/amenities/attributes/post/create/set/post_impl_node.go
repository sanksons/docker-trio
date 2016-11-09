package set

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	_ "fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	_ "strconv"
	"time"
)

type InsertSet struct {
	id string
}

func (u *InsertSet) SetID(id string) {
	u.id = id
}

func (u InsertSet) GetID() (id string, err error) {
	return u.id, nil
}

func (u InsertSet) Name() string {
	return "AttributesSet"
}

func (u InsertSet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SET_INSERT)
	defer func() {
		logger.EndProfile(profiler, SET_INSERT)
	}()
	io.ExecContext.SetDebugMsg(SET_INSERT, "Attribute Set Update Execute")
	d, _ := io.IOData.Get(INSERT_DATA)
	formData, ok := d.(*SetRequestJson)
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Form Data is not valid"}
	}

	res, uerr := u.Insert(formData)
	if uerr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Parameters",
			DeveloperMessage: uerr.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, res)
	return io, nil
}

func (u InsertSet) Insert(formData *SetRequestJson) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SET_INSERT)
	mongoDriver = factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, SET_INSERT)
		mongoDriver.Close()
	}()
	data := u.TransformData(formData)
	updateQuery := bson.M{"$set": data}
	query := bson.M{"seqId": data.SeqId}
	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTESETS_COLLECTION, updateQuery, query, true)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (u InsertSet) TransformData(formData *SetRequestJson) AttributeSet {
	var attr AttributeSet
	attr.SeqId = mongoDriver.GetNextSequence(constants.ATTRIBUTESETS_COLLECTION)
	attr.Name = formData.Name
	attr.Label = formData.Label
	attr.Identifier = formData.Identifier
	attr.UpdatedAt = time.Now()
	attr.CreatedAt = time.Now()
	return attr
}
