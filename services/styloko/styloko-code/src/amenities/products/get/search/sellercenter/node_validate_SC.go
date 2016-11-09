package sellercenter

import (
	proUtil "amenities/products/common"
	"common/utils"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
)

type ValidateNodeSC struct {
	id string
}

func (cs *ValidateNodeSC) SetID(id string) {
	cs.id = id
}

func (cs ValidateNodeSC) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ValidateNodeSC) Name() string {
	return "ValidateNodeSC"
}

func (cs ValidateNodeSC) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Debug("Enter Validate node SC")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_SC_VALIDATE_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_SC_VALIDATE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Validate SC")

	request, err := utils.GetRequestFromIO(io)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:    constants.IncorrectDataErrorCode,
			Message: err.Error()}
	}
	ssr, err := cs.PrepareSellerSearchRequest(request)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:    constants.IncorrectDataErrorCode,
			Message: err.Error()}
	}
	io.IOData.Set(proUtil.REQUEST_STRUCT, ssr)
	io.ExecContext.SetDebugMsg("data", ssr.ToString())
	logger.Debug("Exit Validate node SC")
	return io, nil
}

//prepare SellerSearchRequest
func (cs ValidateNodeSC) PrepareSellerSearchRequest(httpReq *utilhttp.Request) (
	SellerSearchRequest, error) {
	var request SellerSearchRequest
	_, err := utils.GetQueryParam(httpReq, "reset")
	if err == nil {
		request.ResetCounter = true
	} else {
		//set default value
		request.ResetCounter = false
	}

	lastScId, err := utils.GetQueryParam(httpReq, "lastProductSetId")
	if err == nil {
		request.LastSCId, _ = strconv.Atoi(lastScId.(string))
	} else {
		//set default value
		request.LastSCId = 0
	}

	limit, err := utils.GetQueryParam(httpReq, proUtil.PARAM_LIMIT)
	if err == nil {
		request.Limit, _ = strconv.Atoi(limit.(string))
	} else {
		//set default value
		request.Limit = 10
	}
	offset, err := utils.GetQueryParam(httpReq, proUtil.PARAM_OFFSET)
	if err == nil {
		request.Offset, _ = strconv.Atoi(offset.(string))
	} else {
		//set default value
		request.Offset = 0
	}
	var sellerInts []int
	sellerS, err := utils.GetQueryParam(httpReq, proUtil.PARAM_SELLERS)
	if err == nil {
		sellerSlice, ok := sellerS.([]string)
		if ok {
			for _, v := range sellerSlice {
				i, _ := strconv.Atoi(v)
				sellerInts = append(sellerInts, i)
			}
		} else {
			i, _ := strconv.Atoi(sellerS.(string))
			sellerInts = append(sellerInts, i)
		}
	}
	request.SellerIds = sellerInts
	return request, nil
}
