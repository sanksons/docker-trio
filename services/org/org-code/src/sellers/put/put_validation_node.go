package put

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
)

type Decision struct {
	id string
}

func (d *Decision) SetID(id string) {
	d.id = id
}

func (d Decision) GetID() (id string, err error) {
	return d.id, nil
}

func (d Decision) Name() string {
	return "UPDATE seller by id"
}

func (d Decision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_SELLER_PUT)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_SELLER_PUT)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+d.Name(), rc)
	io.ExecContext.SetDebugMsg(VALIDATE_SELLER_PUT, "Seller update decision started")
	logger.Info("Validate seller started")

	data, err := utils.GetPostData(io)
	if err != nil {
		notification.SendNotification("Error while getting Post Data for Seller Update", err.Error(), nil, "error")
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Post Data for Seller Update", DeveloperMessage: err.Error()}
	}
	orgUpdate := new(common.Org)
	logger.Info(data)
	err = json.Unmarshal(data, &orgUpdate)
	if err != nil {
		notification.SendNotification("Seller Update Json Incorrect", "Data type mismatch", nil, "error")
		return false, &florest_constants.AppError{Code: appconstant.DataTypeMismatch, Message: "Json Incorrect", DeveloperMessage: err.Error()}
	}

	var errors []map[string]interface{}
	flag := false
	//calling validate to check seqId, sellerId and orderEmail
	for i := 0; i < len(orgUpdate.Orgdata); i++ {
		errMap := d.ValidateData(&orgUpdate.Orgdata[i])
		if len(errMap) != 0 {
			flag = true
		}
		errors = append(errors, errMap)
	}
	//checking if flag was not unset during validation
	if flag == false {
		io.IOData.Set(ORG_UPDATE_DATA, orgUpdate.Orgdata)
		logger.Info("Seller data extracted")
		return true, nil
	}
	//if flag was set, setting errors in FAILURE_DATA
	for k, v := range errors {
		logger.Error(fmt.Sprintf("%d -> %s", k, v))
	}
	io.IOData.Set(FAILURE_DATA, errors)
	return false, nil
}

//function takes common.Schema as input and returns a map[string]interface{}
//of errors in the input.It checks that seqId should be provided in the input and if
//sellerId is sent,that it be unique.
func (d Decision) ValidateData(orgInfo *common.Schema) map[string]interface{} {

	errorMap := make(map[string]interface{})
	if orgInfo.SeqId == 0 {
		errorMap["seqId"] = "Sequence Id is Mandatory."
	}

	if orgInfo.Status != "" {
		if orgInfo.Status != "active" && orgInfo.Status != "inactive" && orgInfo.Status != "deleted" {
			errorMap["status"] = "Only active,inactive and deleted values are allowed for status"
		}
	}

	if orgInfo.SellerId != "" {
		m := map[string]interface{}{"slrId": orgInfo.SellerId}
		ok, data, err := common.CheckIfKeyExists(m)
		if err != nil {
			errorMap["sellerId"] = "Error while checking for existing Seller Id"
		}
		if ok && data.SellerId != "" {
			//no updations to be made in this case
			orgInfo.SellerId = ""
		}
	}

	return errorMap
}
