package exactquery

import (
	proUtil "amenities/products/common"
	search "amenities/products/get/search"
	_ "errors"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"reflect"
	_ "strconv"
)

type VisibilityCheck struct {
	id string
}

func (cs *VisibilityCheck) SetID(id string) {
	cs.id = id
}

func (cs VisibilityCheck) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs VisibilityCheck) Name() string {
	return "VisibilityCheck"
}

func (cs VisibilityCheck) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Info("Enter visibility check node")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_VISIBILITY_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_VISIBILITY_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "visibility node")
	products, err := LoadDataExactQuery{}.GetProductData(io)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs VisibilityCheck) Execute(): Cannot load data",
			DeveloperMessage: err.Error(),
		}
	}
	query, err := PrepareExactQuery{}.GetQuery(io)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.ResourceErrorCode,
			Message:          "(cs VisibilityCheck) Execute(): Cannot load query",
			DeveloperMessage: err.Error(),
		}
	}
	for i, v := range products {
		if v.Cache == true {
			continue
		}
		p, ok := v.Data.(proUtil.Product)
		if !ok {
			return io, &constants.AppError{
				Code:             constants.ResourceErrorCode,
				Message:          "(cs VisibilityCheck) Execute(): Assertion failed",
				DeveloperMessage: reflect.TypeOf(v.Data).Name(),
			}
		}
		if cs.IsVisible(p, query.Visibility) {
			products[i].Visibility = true
			continue
		}
		products[i].Visibility = false
	}
	io.IOData.Set(search.PRODUCTDATA, products)
	logger.Info("Exit visibility check node")
	return io, nil
}

func (cs VisibilityCheck) IsVisible(p proUtil.Product, visibility string) bool {
	vc := proUtil.VisibilityChecker{}
	vc.Product = p
	vc.VisibilityType = visibility
	return vc.IsVisible()
}
