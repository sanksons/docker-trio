package sellercenter

import (
	proUtil "amenities/products/common"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
)

type FetchDataNodeSC struct {
	id string
}

func (cs *FetchDataNodeSC) SetID(id string) {
	cs.id = id
}

func (cs FetchDataNodeSC) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs FetchDataNodeSC) Name() string {
	return "FetchDataNodeSC"
}

func (cs FetchDataNodeSC) GetSSR(io workflow.WorkFlowData) (SellerSearchRequest, error) {
	r, err := io.IOData.Get(proUtil.REQUEST_STRUCT)
	if err != nil {
		return SellerSearchRequest{}, fmt.Errorf("(cs FetchDataNodeSC) GetSSR: %s", err.Error())
	}
	req, ok := r.(SellerSearchRequest)
	if !ok {
		return SellerSearchRequest{}, fmt.Errorf("(cs FetchDataNodeSC) Execute(): Invalid IO Data")
	}
	return req, nil
}

func (cs FetchDataNodeSC) Execute(io workflow.WorkFlowData) (
	workflow.WorkFlowData, error) {

	logger.Debug("Enter FetchData SC")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_SC_FETCH_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_SC_FETCH_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Fetch Data SC")

	req, err := cs.GetSSR(io)
	if err != nil {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: err.Error(),
		}
	}

	//check if its a reset counter request
	if req.ResetCounter {
		io.ExecContext.SetDebugMsg("key", "Its a Reset Request")
		//reset counters
		proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).ResetSSRCounter()
		return io, nil
	}

	products, er := cs.GetProducts(req)
	if er != nil {
		logger.Error(er)
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs FetchDataNodeSC) Execute(): Unable to Fetch Data" + er.Error(),
			DeveloperMessage: er.Error(),
		}
	}
	io.ExecContext.SetDebugMsg("count of products", strconv.Itoa(len(products)))
	io.IOData.Set(proUtil.IODATA, products)
	logger.Debug("Exit FetchData SC")
	return io, nil
}

//
// Get products based on the supplied condition.
//
func (cs FetchDataNodeSC) GetProducts(req SellerSearchRequest) (
	[]proUtil.Product, error) {

	products, err := proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).GetProductsForSeller(
		req.SellerIds,
		req.Limit,
		req.Offset,
		req.LastSCId,
	)
	if err != nil {
		return products, fmt.Errorf("(cs FetchDataNodeSC)#GetProducts: %s", err.Error())
	}
	return products, nil
}
