package delete

import (
	proUtil "amenities/products/common"
	"common/utils"
	"errors"
	"strconv"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// DeleteCacheNode is a basic node struct
type DeleteCacheNode struct {
	id string
}

// SetID sets the ID for the node.
func (cs *DeleteCacheNode) SetID(id string) {
	cs.id = id
}

// GetID returns ID for the Node
func (cs DeleteCacheNode) GetID() (id string, err error) {
	return cs.id, nil
}

// Name returns name of the Node
func (cs DeleteCacheNode) Name() string {
	return "DeleteCacheNode"
}

// Execute executes the workflow
func (cs DeleteCacheNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter Delete Cache node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.DELETE_CACHE_NODE)
	defer logger.EndProfile(profiler, proUtil.DELETE_CACHE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Delete Cache Node")

	// Get delete struct for SKU and ID
	delStruct, ty, err := cs.PrepareDeleteQuery(io)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "Invalid Query Params",
			DeveloperMessage: err.Error(),
		}
	}

	switch ty {
	case "sku":
		// Making it async to reduce response time
		go func() {
			defer proUtil.RecoverHandler("DeleteCache#sku")
			cacheMngr.DeleteBySku(delStruct.Skus, true)
		}()
		break
	case "id":
		// Making it async to reduce response time
		go func() {
			defer proUtil.RecoverHandler("DeleteCache#Id")
			cacheMngr.DeleteById(delStruct.Ids, true)
		}()
		break
	}

	//Set data for next node
	logger.Info("Exit Delete Cache node")
	io.IOData.Set(constants.RESULT, "success")
	return io, nil
}

// PrepareDeleteQuery -> prepared delete query
// Accepted query strings ?ids=[1,2,3] or ?sku=[AWP12321INDFAS,AWP31241INDFAS]
func (cs DeleteCacheNode) PrepareDeleteQuery(io workflow.WorkFlowData) (Query, string, error) {
	delStruct := Query{}
	skuFlag := false
	idFlag := false

	httpReq, _ := utils.GetRequestFromIO(io)
	tmp, _ := utils.GetQueryParam(httpReq, "ids")
	ids, _ := tmp.([]string)
	if len(ids) > 0 {
		idFlag = true
		for x := 0; x < len(ids); x++ {
			id, err := strconv.Atoi(ids[0])
			if err != nil {
				return delStruct, "", errors.New("Invalid data in Query. Integer conversion failed.")
			}
			delStruct.Ids = append(delStruct.Ids, id)
		}
	}

	tmp, _ = utils.GetQueryParam(httpReq, "skus")
	skus, _ := tmp.([]string)
	if len(skus) > 0 {
		skuFlag = true
		delStruct.Skus = skus
	}
	if idFlag && skuFlag {
		return delStruct, "", errors.New("Cannot use both IDs and SKUs together")
	}

	if !idFlag && !skuFlag {
		return delStruct, "", errors.New("No IDs or SKUs provided.")
	}

	if idFlag {
		return delStruct, "id", nil
	}
	return delStruct, "sku", nil
}
