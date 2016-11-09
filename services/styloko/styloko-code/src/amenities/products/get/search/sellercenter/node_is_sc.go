package sellercenter

import (
	proUtil "amenities/products/common"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type IsSC struct {
	id string
}

func (cs *IsSC) SetID(id string) {
	cs.id = id
}

func (cs IsSC) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs IsSC) Name() string {
	return "IsSC"
}

func (cs IsSC) GetDecision(io workflow.WorkFlowData) (bool, error) {
	logger.Debug("Enter Decision node")

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Decision Node")
	rp, _ := io.IOData.Get(florestConsts.REQUEST)
	appHttpReq, pOk := rp.(*utilhttp.Request)
	if !pOk || appHttpReq == nil {
		return false, &florestConsts.AppError{
			Code:    florestConsts.IncorrectDataErrorCode,
			Message: "invalid request"}
	}
	platform := appHttpReq.OriginalRequest.Header.Get("Request-Platform")
	logger.Debug("Exit Decision node")
	if platform == proUtil.SELLER_CENTER {
		return true, nil
	}
	return false, nil
}
