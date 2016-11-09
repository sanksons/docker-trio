package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"

	"time"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type UpdateAttribute struct {
	id string
}

func (u *UpdateAttribute) SetID(id string) {
	u.id = id
}

func (u UpdateAttribute) GetID() (id string, err error) {
	return u.id, nil
}

func (u UpdateAttribute) Name() string {
	return "Attribute Update api"
}

func (u UpdateAttribute) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_UPDATE)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_UPDATE)
	}()
	io.ExecContext.SetDebugMsg(ATTRIBUTE_UPDATE, "Attribute Update Execute")

	data, _ := io.IOData.Get(PATH_PARAMETERS)
	seqId, ok := data.(int)

	d, _ := io.IOData.Get(UPDATE_DATA)
	formData, ok := d.(*Attribute)
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

func (u UpdateAttribute) Update(seqId int, formData *Attribute) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_UPDATE)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_UPDATE)
		mongoDriver.Close()
	}()
	formData.UpdatedAt = time.Now()
	updateQuery := M{"$set": formData}
	query := M{"seqId": seqId}
	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTES_COLLECTION,
		updateQuery, query, false)
	if err != nil {
		return res, err
	}

	return res, nil
}
