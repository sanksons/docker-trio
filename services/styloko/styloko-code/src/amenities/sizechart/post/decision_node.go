package post

import (
	sizeUtils "amenities/sizechart/common"
	"common/utils"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// SizeChartTypeDecision -> struct for node based data
type SizeChartTypeDecision struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (s *SizeChartTypeDecision) SetID(id string) {
	s.id = id
}

// GetID -> returns current node ID to orchestrator
func (s SizeChartTypeDecision) GetID() (id string, err error) {
	return s.id, nil
}

// Name -> Returns node name to orchestrator
func (s SizeChartTypeDecision) Name() string {
	return "SizeChartTypeDecision"
}

// GetDecision -> Decides which node to run next.
// True stands for sku level, false for brand-brick sizechart
func (s SizeChartTypeDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {
	var res bool
	ty, _ := utils.GetRequestHeader(io, sizeUtils.SizeChartHeader)
	if ty == sizeUtils.BrandBrick {
		res = false
	} else if ty == sizeUtils.SKU {
		res = true
	} else {
		return false, &florest_constants.AppError{
			Code:             florest_constants.IncorrectDataErrorCode,
			Message:          "Size Chart Request Header is not valid",
			DeveloperMessage: "Bad Request"}
	}

	return res, nil
}
