package exactquery

import (
	proUtil "amenities/products/common"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type PublishNode struct {
	id string
}

func (cs *PublishNode) SetID(id string) {
	cs.id = id
}

func (cs PublishNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs PublishNode) Name() string {
	return "PublishNode"
}

func (cs PublishNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	//check if we need to publish to bus

	if !cs.doPublish(io) {
		return io, nil
	}

	logger.Debug("Enter publish node")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_publish_node")
	defer logger.EndProfile(profiler, "product_publish_node")

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "publish node")
	products, err := LoadDataExactQuery{}.GetProductData(io)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs PublishNode)#Execute(): Cannot load data",
			DeveloperMessage: err.Error(),
		}
	}
	for _, v := range products {
		p, ok := v.Data.(proUtil.Product)
		if !ok {
			return io, &constants.AppError{
				Code:             constants.ResourceErrorCode,
				Message:          "(cs VisibilityCheck) Execute(): Assertion failed",
				DeveloperMessage: err.Error(),
			}
		}
		p.Publish("", true)
		p.PushToMemcache("Custom Push")
	}
	logger.Debug("Exit visibility check node")
	return io, nil
}

func (cs PublishNode) doPublish(io workflow.WorkFlowData) bool {
	rp, _ := io.IOData.Get(constants.REQUEST)
	appHttpReq, pOk := rp.(*utilhttp.Request)
	if !pOk || appHttpReq == nil {
		return false
	}
	publish := appHttpReq.OriginalRequest.Header.Get(proUtil.HEADER_PUBLISH)
	if publish == "true" {
		return true
	}
	return false
}
