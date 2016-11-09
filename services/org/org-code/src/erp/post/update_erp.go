package post

import (
	"common/appconstant"
	"common/notification"
	"common/utils"
	"encoding/json"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sellers/common"
	"sellers/get"
)

type ErpUpdate struct {
	id string
}

func (e *ErpUpdate) SetID(id string) {
	e.id = id
}

func (e ErpUpdate) GetID() (id string, err error) {
	return e.id, nil
}

func (e ErpUpdate) Name() string {
	return "UPDATE details in ERP"
}

func (e ErpUpdate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, UPDATE_ERP)
	defer func() {
		logger.EndProfile(profiler, UPDATE_ERP)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+e.Name(), rc)
	io.ExecContext.SetDebugMsg(UPDATE_ERP, "Erp_Update Node execution started")

	data, err := utils.GetPostData(io)
	if err != nil {
		notification.SendNotification("Error while getting Post Data for Erp Update", err.Error(), nil, "error")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Post Data for Erp Update", DeveloperMessage: err.Error()}
	}

	erpData, err := e.GetDetailsForData(data)
	if err != nil {
		notification.SendNotification("Error while getting data for passed Ids in Erp Update", err.Error(), nil, "error")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting data for passed Ids", DeveloperMessage: err.Error()}
	}

	err = common.SendDataToErp(erpData)
	if err != nil {
		notification.SendNotification("Error while sending data to ERP", err.Error(), nil, "error")
		logger.Error(fmt.Sprintf("Error while sending data to ERP : %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while sendin data to ERP", DeveloperMessage: err.Error()}
	}
	return io, nil
}

//gets data from the ids passed from mongo and populates it in a format to be consumed by the orchestrator
func (e ErpUpdate) GetDetailsForData(data interface{}) ([]common.ErpData, error) {
	feResponse := new(FeResponse)
	err := json.Unmarshal(data.([]byte), &feResponse)
	if err != nil {
		return nil, err
	}
	gs := get.SearchSeller{}
	//sending default limit as lenth of array of seqIds and offset as 0
	res, err := gs.GetIdDetails(feResponse.SeqIds, len(feResponse.SeqIds), 0)
	if err != nil {
		return nil, err
	}
	erpData := make([]common.ErpData, 1)
	erpData[0] = common.ErpData{
		Method: "ERP.insertSeller",
		Params: common.ErpSellerData{
			SellerData: res,
		},
	}
	return erpData, nil
}
