package post

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
	logger.Debug("Enter Response node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.POST_RESPONSE_NODE)
	defer logger.EndProfile(profiler, proUtil.POST_RESPONSE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Response Node")

	d, _ := io.IOData.Get(proUtil.IODATA)
	iodata, ok := d.([]*ProIOData)
	if !ok {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: "(cs ResponseNode)#Execute(): Invalid IO Data",
		}
	}
	finalResponse := make([]map[string]interface{}, 0)
	for _, v := range iodata {
		response := make(map[string]interface{})
		response["status"] = v.Status
		if v.Error != nil {
			response["error"] = v.Error
		} else {
			var simpleIndex int
			videohashMap := make(map[string]bool, 0)
			for index, s := range v.Product.Simples {
				if v.ReqData.SellerSKU == s.SellerSKU {
					simpleIndex = index
				}
			}

			proResponse := make(map[string]interface{})
			proResponse["configId"] = v.Product.SeqId
			proResponse["id"] = v.Product.Simples[simpleIndex].Id
			proResponse["sku"] = v.Product.Simples[simpleIndex].SKU
			proResponse["sellerSku"] = v.Product.Simples[simpleIndex].SellerSKU
			proResponse["updatedAt"] = proUtil.ToMySqlTime(&v.Product.Simples[simpleIndex].UpdatedAt)
			//prepare video hash
			if len(v.Product.Videos) > 0 {
				for _, video := range v.Product.Videos {
					videohashMap[video.Hash] = true
				}
			}
			proResponse["videoProcessInfo"] = videohashMap
			response["product"] = proResponse
			response["isDuplicate"] = v.IsDuplicate
		}
		finalResponse = append(finalResponse, response)
	}
	logger.Debug("Exit Response node")
	io.IOData.Set(constants.RESULT, finalResponse)
	return io, nil
}
