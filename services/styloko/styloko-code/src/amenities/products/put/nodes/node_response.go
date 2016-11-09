package nodes

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ResponseNode struct {
	id string
}

func (rn *ResponseNode) SetID(id string) {
	rn.id = id
}

func (rn ResponseNode) GetID() (string, error) {
	return rn.id, nil
}

func (rn ResponseNode) Name() string {
	return "ResponseNode"
}

func (rn ResponseNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Entered response node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.PUT_RESPONSE_NODE)
	defer logger.EndProfile(profiler, proUtil.PUT_RESPONSE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "In Response node")

	d, _ := io.IOData.Get(proUtil.IODATA)
	ioDataArr, ok := d.([]*put.ProIOData)
	if !ok {
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "(rn ResponseNode)#Execute(): []*ProIOData assertion failed",
			DeveloperMessage: "(rn ResponseNode)#Execute(): []*ProIOData assertion failed",
		}
	}
	var response interface{}

	updateType, _ := io.IOData.Get(proUtil.HEADER_UPDATE_TYPE)
	switch updateType {
	case proUtil.UPDATE_TYPE_IMAGEDEL:
		ioData := ioDataArr[0]
		response = ioData.ReqData.Response(nil)

	case proUtil.UPDATE_TYPE_IMAGEADD:
		response = rn.ImageAdd(ioDataArr)

	case proUtil.UPDATE_TYPE_VIDEO_STATUS:
		response = rn.VideoStatus(ioDataArr)

	default:
		response = rn.Response(ioDataArr)
	}
	io.IOData.Set(constants.RESULT, response)
	logger.Info("Exit response node")
	return io, nil
}

//
// Prepare appropriate response for video status update case.
//
func (rn ResponseNode) VideoStatus(ioDataArr []*put.ProIOData) []map[string]interface{} {
	finalResponse := make([]map[string]interface{}, 0)
	for _, v := range ioDataArr {
		response := make(map[string]interface{})
		response["status"] = v.Status
		response["videoId"] = v.ReqData.Response(v.Product)
		finalResponse = append(finalResponse, response)
	}
	return finalResponse
}

//
// Prepare appropriate response for image Addition case.
//
func (rn ResponseNode) ImageAdd(ioDataArr []*put.ProIOData) []interface{} {
	finalResponse := make([]interface{}, 0)
	for _, v := range ioDataArr {
		p := &proUtil.Product{}
		if v.Product != nil {
			p = v.Product
		}
		response := v.ReqData.Response(p)
		resp, ok := response.(map[string]interface{})
		if ok && (v.Error != nil) {
			resp["error"] = v.Error
			finalResponse = append(finalResponse, resp)
			continue
		}
		finalResponse = append(finalResponse, response)
	}
	return finalResponse
}

//
// This is a general function for response creation, so incase the
// update type does not belong to a specialized case this function will
// be used.
//
func (rn ResponseNode) Response(ioDataArr []*put.ProIOData) []map[string]interface{} {
	finalResponse := make([]map[string]interface{}, 0)
	for _, v := range ioDataArr {
		response := make(map[string]interface{})
		response["status"] = v.Status
		if v.Error != nil {
			response["error"] = v.Error
		} else {
			response["product"] = v.ReqData.Response(v.Product)
		}
		finalResponse = append(finalResponse, response)
	}
	return finalResponse
}
