package exactquery

import (
	proUtil "amenities/products/common"
	search "amenities/products/get/search"
	"common/utils"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//
// Checks if the query we have is of Exact type.
//
type IsExactQuery struct {
	id string
}

func (cs *IsExactQuery) SetID(id string) {
	cs.id = id
}

func (cs IsExactQuery) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs IsExactQuery) Name() string {
	return "IsExactQuery"
}

func (cs IsExactQuery) GetDecision(io workflow.WorkFlowData) (bool, error) {
	logger.Debug("Enter IsExactQuery Decision node")

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(
		proUtil.DEBUG_KEY_NODE,
		"IsExactQuery Decision Node",
	)

	isExact := cs.IsExactQuery(io)
	if isExact {
		io.IOData.Set(search.QUERY_TYPE, search.QUERY_TYPE_EXACT)
	} else {
		io.IOData.Set(search.QUERY_TYPE, search.QUERY_TYPE_FILTER)
	}
	logger.Debug("Exit IsExactQuery Decision node")
	return isExact, nil
}

//Get Query type from IO data
func (cs IsExactQuery) GetQueryType(io workflow.WorkFlowData) string {
	q, err := io.IOData.Get(search.QUERY_TYPE)
	if err != nil {
		logger.Error(err)
		return search.QUERY_TYPE_FILTER
	}
	data, ok := q.(string)
	if !ok {
		return search.QUERY_TYPE_FILTER
	}
	return data
}

//Check if the query id of exact type
func (cs IsExactQuery) IsExactQuery(io workflow.WorkFlowData) bool {
	_, err := utils.GetPathParams(io)
	if err == nil {
		//Its an exact query type
		return true
	}
	request, _ := utils.GetRequestFromIO(io)
	//check if we have a limit parameter supplied
	_, err = utils.GetQueryParam(request, proUtil.PARAM_LIMIT)
	if err == nil {
		//Its a filter query type
		return false
	}
	//check if we have sku in param
	_, err = utils.GetQueryParam(request, proUtil.PARAM_SKU)
	if err == nil {
		//Its exact query type
		return true
	}
	//check if we have ID in param
	_, err = utils.GetQueryParam(request, proUtil.PARAM_ID)
	if err == nil {
		//Its exact query type
		return true
	}
	return false
}
