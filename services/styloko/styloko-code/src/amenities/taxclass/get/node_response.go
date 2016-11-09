package get

import (
	proUtil "amenities/products/common"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ResponseNode struct {
	id string
}

func (cs *ResponseNode) SetID(id string) {
	cs.id = id
}

func (cs ResponseNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ResponseNode) Name() string {
	return "ResponseNode"
}

func (cs ResponseNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter Response node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_TAXCLASS_RESPONSE_NODE)
	defer logger.EndProfile(profiler, GET_TAXCLASS_RESPONSE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(DEBUG_KEY_NODE, "ResponseNode")

	d, _ := io.IOData.Get(IODATA)
	tClass, ok := d.([]proUtil.TaxClass)
	if !ok {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: "(cs ResponseNode)#Execute(): Invalid IO Data"}
	}
	response := []TaxClassResponse{}
	for _, v := range tClass {
		taxclass := TaxClassResponse{}
		taxclass.Id = v.Id
		taxclass.IsDefault = false
		if v.IsDefault == 1 {
			taxclass.IsDefault = true
		}
		taxclass.Name = v.Name
		taxclass.TaxPercent = v.TaxPercent
		response = append(response, taxclass)
	}
	logger.Info("Exit Response node")
	io.IOData.Set(constants.RESULT, response)
	return io, nil
}
