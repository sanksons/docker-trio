package exactquery

import (
	proUtil "amenities/products/common"
	search "amenities/products/get/search"
	"common/utils"
	"errors"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//
// Prepare data for exact query
//
type PrepareExactQuery struct {
	id string
}

func (cs *PrepareExactQuery) SetID(id string) {
	cs.id = id
}

func (cs PrepareExactQuery) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs PrepareExactQuery) Name() string {
	return "PrepareExactQuery"
}

func (cs PrepareExactQuery) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Debug("Enter Prepare exact query node")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_PREPARE_QUERY_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_PREPARE_QUERY_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Prepare exact query")

	query, err := cs.PrepareQuery(io)
	if err != nil {
		logger.Error(err)
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "Invalid request",
			DeveloperMessage: err.Error(),
		}
	}
	cs.setCustomDatadogMetrics(io, query)

	io.IOData.Set(search.QUERY, query)
	io.ExecContext.SetDebugMsg("preparedquery", query.ToString())
	logger.Debug("Exit Prepare exact query node")
	return io, nil
}

func (cs PrepareExactQuery) PrepareQuery(io workflow.WorkFlowData) (search.ExactQuery, error) {
	query := search.ExactQuery{}
	param, err := utils.GetPathParams(io)
	if err == nil {
		id, _ := strconv.Atoi(param[0])
		if id > 0 {
			query.Id = []int{id}
		} else {
			query.Sku = []string{param[0]}
		}
		query.Visibility = cs.GetVisibility(io)
		query.Expanse = cs.GetExpanse(io)
		query.IsSingle = true
		return query, nil
	}
	request, _ := utils.GetRequestFromIO(io)
	sku, err := utils.GetQueryParam(request, proUtil.PARAM_SKU)
	if err == nil {
		skus, ok := sku.([]string)
		if !ok {
			return query, errors.New(
				"(cs PrepareExactQuery)#PrepareQuery(): Assertion failed")
		}
		if len(skus) <= 0 {
			return query, errors.New("Please supply atleast one sku")
		}
		query.Sku = skus
		query.Visibility = cs.GetVisibility(io)
		query.Expanse = cs.GetExpanse(io)
		return query, nil
	}
	id, err := utils.GetQueryParam(request, proUtil.PARAM_ID)
	if err == nil {
		idsStr, ok := id.([]string)
		if !ok {
			return query, errors.New("(cs PrepareExactQuery)#PrepareQuery(): Assertion failed")
		}
		ids := []int{}
		for _, v := range idsStr {
			id, _ := strconv.Atoi(v)
			if id > 0 {
				ids = append(ids, id)
			}
		}
		if len(ids) <= 0 {
			return query, errors.New("Please supply atleast one productid")
		}
		query.Id = ids
		query.Visibility = cs.GetVisibility(io)
		query.Expanse = cs.GetExpanse(io)
		return query, nil
	}
	//if do not fall in above cases definitely an error
	return query, errors.New("(cs PrepareExactQuery)#PrepareQuery: Cannot parse query")
}

func (cs PrepareExactQuery) GetVisibility(io workflow.WorkFlowData) string {
	request, _ := utils.GetRequestFromIO(io)
	visibility := request.OriginalRequest.Header.Get(proUtil.HEADER_VISIBILITY_TYPE)
	allVisibility := proUtil.GetAllVisibility()
	for _, v := range allVisibility {
		if strings.ToLower(v) == strings.ToLower(visibility) {
			return v
		}
	}
	return search.DEFAULT_VISIBILITY
}

func (cs PrepareExactQuery) GetExpanse(io workflow.WorkFlowData) string {
	request, _ := utils.GetRequestFromIO(io)
	expanse := request.OriginalRequest.Header.Get(proUtil.HEADER_EXPANSE)
	allExpanse := proUtil.GetAllExpanse()
	for _, v := range allExpanse {
		if strings.ToLower(v) == strings.ToLower(expanse) {
			return v
		}
	}
	return search.DEFAULT_EXPANSE
}

func (cs PrepareExactQuery) GetQuery(io workflow.WorkFlowData) (search.ExactQuery, error) {
	q, err := io.IOData.Get(search.QUERY)
	eq := search.ExactQuery{}
	if err != nil {
		return eq, errors.New("(cs CacheGet)#GetExactQuery():" + err.Error())
	}
	query, ok := q.(search.ExactQuery)
	if !ok {
		return eq, errors.New("(cs CacheGet)#GetExactQuery() Assertion:" + err.Error())
	}
	return query, nil
}

func (cs PrepareExactQuery) setCustomDatadogMetrics(io workflow.WorkFlowData, query search.ExactQuery) {
	if query.IsSingle {
		io.ExecContext.Set(constants.MONITOR_CUSTOM_METRIC_PREFIX, cs.setHeadersForMetrics("SINGLE", query.Expanse, query.Visibility))
		return
	}
	io.ExecContext.Set(constants.MONITOR_CUSTOM_METRIC_PREFIX, cs.setHeadersForMetrics("MULTI", query.Expanse, query.Visibility))
}

func (cs PrepareExactQuery) setHeadersForMetrics(ty, expanse, visibility string) string {
	metricName := "_CUSTOM_PRODUCTS_GET_" + ty + "_" + expanse + "_" + visibility + "_"
	return metricName
}
