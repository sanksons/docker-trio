package nodes

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type InvalidateCache struct {
	id string
}

func (vu *InvalidateCache) SetID(id string) {
	vu.id = id
}

func (vu InvalidateCache) GetID() (string, error) {
	return vu.id, nil
}

func (vu InvalidateCache) Name() string {
	return "InvalidateCache"
}

func (ic InvalidateCache) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Info("Entered invalidate cache")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.PUT_INVALIDATE_CACHE_NODE)
	defer logger.EndProfile(profiler, proUtil.PUT_INVALIDATE_CACHE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Invalidate Cache")

	d, _ := io.IOData.Get(proUtil.IODATA)
	ioDataArr, ok := d.([]*put.ProIOData)
	if !ok {
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "(vu InvalidateCache)#Execute(): []*ProIOData assertion failed",
			DeveloperMessage: "(vu InvalidateCache)#Execute(): []*ProIOData assertion failed",
		}
	}
	for _, v := range ioDataArr {
		if v.Status == proUtil.STATUS_FAILURE {
			continue
		}
		err := v.ReqData.InvalidateCache()
		if err != nil {
			logger.Error(err)
		}
		//Publish Info to Bus
		err = v.ReqData.Publish()
		if err != nil {
			logger.Error(err)
		}
	}
	io.IOData.Set(proUtil.IODATA, ioDataArr)
	logger.Info("Exit invalidate cache")
	return io, nil
}
