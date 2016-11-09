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

type UpdateSet struct {
	id string
}

func (u *UpdateSet) SetID(id string) {
	u.id = id
}

func (u UpdateSet) GetID() (id string, err error) {
	return u.id, nil
}

func (u UpdateSet) Name() string {
	return "AttributesSet"
}

func (u UpdateSet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SET_UPDATE)
	defer func() {
		logger.EndProfile(profiler, SET_UPDATE)
	}()
	io.ExecContext.SetDebugMsg(SET_UPDATE, "Attribute Set Update Execute")

	data, _ := io.IOData.Get(PATH_PARAMETERS)
	seqId, ok := data.(int)

	d, _ := io.IOData.Get(UPDATE_DATA)
	formData, ok := d.(*AttributeSet)
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Form Data is not valid"}
	}

	res, uerr := u.Update(seqId, formData)
	if uerr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Parameters",
			DeveloperMessage: uerr.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, res)
	return io, nil
}

func (u UpdateSet) Update(seqId int, formData *AttributeSet) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SET_UPDATE)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTESETAPI)
	defer func() {
		logger.EndProfile(profiler, SET_UPDATE)
		mongoDriver.Close()
	}()
	updatedAt := time.Now()
	updateQuery := bson.M{"$set": bson.M{
		"label":      formData.Label,
		"identifier": formData.Identifier,
		"updatedAt":  updatedAt},
	}
	query := bson.M{"seqId": seqId}
	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTESETS_COLLECTION, updateQuery, query, false)
	if err != nil {
		return res, err
	}
	return res, nil
}
