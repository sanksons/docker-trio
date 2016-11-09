package post

import (
	proUtil "amenities/products/common"
	"common/notification"
	"common/notification/datadog"
	"common/utils"
	"encoding/json"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/monitor"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"strings"
)

//Validate Create Products call.
type ValidateNode struct {
	id string
}

func (cs *ValidateNode) SetID(id string) {
	cs.id = id
}

func (cs ValidateNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ValidateNode) Name() string {
	return "ValidateNode"
}

func (cs ValidateNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Debug("Enter Validate node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.POST_VALIDATE_NODE)
	defer logger.EndProfile(profiler, proUtil.POST_VALIDATE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Validate Node")

	httpReq, err := utils.GetRequestFromIO(io)
	if err != nil {
		logger.Error(err)
	}
	httpReqOrig := httpReq.OriginalRequest
	data, _ := ioutil.ReadAll(httpReqOrig.Body)
	logger.Warning(string(data))

	reqBody := []*ProductCreateRequestData{}
	err = json.Unmarshal(data, &reqBody)
	if err != nil {
		logger.Warning(string(data))
		logger.Error(err)

		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "cannot read request data",
			DeveloperMessage: err.Error(),
		}
	}
	// Publish metrics for product create count.
	productCount := int64(len(reqBody))
	go func(cc int64) {
		defer utils.RecoverHandler("Publish metric: counter_ProductCreate")
		monitor.GetInstance().Count(
			"counter_ProductCreate", cc, []string{"styloko"}, 1,
		)
	}(productCount)
	//Validate the request
	iodata := cs.Validate(reqBody)

	//Set data for next node
	logger.Debug("Exit Validate node")
	io.IOData.Set(proUtil.IODATA, iodata)
	return io, nil
}

func (cs ValidateNode) Validate(reqBody []*ProductCreateRequestData) []*ProIOData {

	iodata := make([]*ProIOData, len(reqBody))
	for i, v := range reqBody {
		errs := validate.Struct(v)
		iodata[i] = &ProIOData{}
		if errs != nil {
			validationErrors := errs.(validator.ValidationErrors)
			msgs := cs.PrepareErrorMessages(validationErrors)
			iodata[i].setFailure(constants.AppError{
				Code:             constants.IncorrectDataErrorCode,
				Message:          "Validation Failed : " + v.SellerSKU,
				DeveloperMessage: strings.Join(msgs, ";"),
			})
			//notify
			notification.SendNotification(
				"Product Create Failed",
				fmt.Sprintf("Message:%s, DevMessage:%s",
					iodata[i].Error.Message, iodata[i].Error.DeveloperMessage,
				),
				[]string{proUtil.TAG_PRODUCT_CREATE, proUtil.TAG_PRODUCT},
				datadog.ERROR,
			)
		} else {
			iodata[i].ReqData = v
		}
	}
	return iodata
}

func (cs ValidateNode) PrepareErrorMessages(errs validator.ValidationErrors) []string {
	var msgs []string
	for _, err := range errs {
		var msg string
		switch err.Tag {
		case "required":
			msg = err.Field + ": Is Required."
		case "ltfield":
			msg = err.Field + ": should be less than " + err.Param
		default:
			msg = "Validation failed"
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
