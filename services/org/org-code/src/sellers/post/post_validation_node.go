package post

import (
	"common/appconstant"
	"common/notification"
	"common/utils"
	"encoding/json"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sellers/common"
)

type Validate struct {
	id string
}

func (v *Validate) SetID(id string) {
	v.id = id
}

func (v Validate) GetID() (id string, err error) {
	return v.id, nil
}

func (v Validate) Name() string {
	return "Validate node for POST"
}

func (v Validate) GetDecision(io workflow.WorkFlowData) (bool, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, VALIDATE_SELLER_POST)
	defer func() {
		logger.EndProfile(profiler, VALIDATE_SELLER_POST)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+v.Name(), rc)
	io.ExecContext.SetDebugMsg(VALIDATE_SELLER_POST, "Seller update decision started")

	var errors []map[string]interface{}
	flag := false
	data, errs := utils.GetPostData(io)
	if errs != nil {
		notification.SendNotification("Error while getting Post Data for Seller Create", errs.Error(), nil, "error")
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Post Data for Seller Create", DeveloperMessage: errs.Error()}
	}
	orgCreate := new(common.Org)
	err := json.Unmarshal(data, &orgCreate)

	if err != nil {
		notification.SendNotification("Seller Create Json Incorrect", "Data type mismatch", nil, "error")
		return false, &florest_constants.AppError{Code: appconstant.DataTypeMismatch, Message: "Json Incorrect", DeveloperMessage: err.Error()}
	}

	for i := 0; i < len(orgCreate.Orgdata); i++ {
		errMap := v.ValidateData(&orgCreate.Orgdata[i])
		if len(errMap) != 0 {
			flag = true
		}
		errors = append(errors, errMap)
	}

	if flag == false {
		io.IOData.Set(ORG_DATA, orgCreate.Orgdata)
		return true, nil
	}

	io.IOData.Set(FAILURE_DATA, errors)
	logger.Info("Seller data extracted")
	return false, nil
}

//valides data for mandatory and non-mandatory fields in post request data
func (v Validate) ValidateData(orgInfo *common.Schema) map[string]interface{} {
	errorMap := make(map[string]interface{})
	if orgInfo.SeqId > 0 {
		errorMap["seqId"] = "Forceful insertion of Sequence Id is not allowed"
	}
	orgInfo.Sync = false

	if orgInfo.SellerId == "" {
		errorMap["sellerId"] = "Seller Id is Missing"
	} else {
		bsonMap := make(map[string]interface{})
		bsonMap["slrId"] = orgInfo.SellerId
		ok, _, _ := common.CheckIfKeyExists(bsonMap)
		if ok {
			errorMap["sellerId"] = "Seller Id already exists.Should be Unique."
		}
	}

	if orgInfo.OrderEmail == "" {
		errorMap["orderEmail"] = "Email Id is Missing"
	} else {
		bsonMap := make(map[string]interface{})
		bsonMap["ordrEml"] = orgInfo.OrderEmail
		ok, _, _ := common.CheckIfKeyExists(bsonMap)
		if ok {
			errorMap["orderEmail"] = "OrderEmail Id already exists.Should be Unique."
		}
	}

	if orgInfo.SellerName == "" {
		errorMap["sellerName"] = "Seller Name is Missing"
	} else {
		bsonMap := make(map[string]interface{})
		bsonMap["slrName"] = orgInfo.SellerName
		ok, _, _ := common.CheckIfKeyExists(bsonMap)
		if ok {
			errorMap["sellerName"] = "Seller Name already exists.Should be Unique."
		}
	}

	if orgInfo.OrgName == "" {
		errorMap["orgName"] = "Org Name is Missing"
	}

	if orgInfo.ContactName == "" {
		errorMap["contactName"] = "Contact Name is Missing"
	}

	if orgInfo.Phone == "" {
		errorMap["phone"] = "Phone Number is Missing"
	}

	if orgInfo.City == "" {
		errorMap["city"] = "City is Missing "
	}

	if orgInfo.Address1 == "" {
		errorMap["address1"] = "Address is Missing "
	}

	if orgInfo.Postcode == 0 {
		errorMap["postcode"] = "Postcode is Missing"
	}

	if orgInfo.CountryCode == "" {
		errorMap["country"] = "Country Code is Missing"
	}
	return errorMap
}
