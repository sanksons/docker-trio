package nodes

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type UpdateMongo struct {
	id string
}

func (um *UpdateMongo) SetID(id string) {
	um.id = id
}

func (um UpdateMongo) GetID() (string, error) {
	return um.id, nil
}

func (um UpdateMongo) Name() string {
	return "UpdateMongo"
}

func (um UpdateMongo) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter Update Mongo")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.PUT_UPDATE_NODE)
	defer logger.EndProfile(profiler, proUtil.PUT_UPDATE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Update Node")

	d, _ := io.IOData.Get(proUtil.IODATA)
	ioDataArr, ok := d.([]*put.ProIOData)
	if !ok {
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "(um UpdateMongo)#Execute(): []*ProIOData assertion failed",
			DeveloperMessage: "(um UpdateMongo)#Execute(): []*ProIOData assertion failed",
		}
	}
	um.Update(ioDataArr)
	io.IOData.Set(proUtil.IODATA, ioDataArr)
	logger.Info("Exit Update Mongo")
	return io, nil
}

//
// Generalized update method for all update cases.
//
func (um UpdateMongo) Update(proIOData []*put.ProIOData) {
	for i, v := range proIOData {
		if v.Status == proUtil.STATUS_FAILURE {
			continue
		}
		profiler := logger.NewProfiler()
		logger.StartProfile(profiler, "Update")
		defer func() {
			logger.EndProfile(profiler, "Update")
		}()
		//acquire lock on the resource.
		if v.ReqData.Lock() {
			p, er := v.ReqData.Update()
			if er != nil {
				logger.Error("(um UpdateMongo)#Update: " + er.Error())
				proIOData[i].SetFailure(constants.AppError{
					Code:             constants.IncorrectDataErrorCode,
					Message:          "Mongo Updation Failed",
					DeveloperMessage: er.Error(),
				})
			} else {
				proIOData[i].SetSuccess(&p)
			}
			//release the resource
			v.ReqData.UnLock()
		}
	}
}
