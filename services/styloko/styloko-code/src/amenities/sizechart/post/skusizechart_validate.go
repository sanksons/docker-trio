package post

import (
	sizeUtils "amenities/sizechart/common"
	"common/appconstant"
	utils "common/utils"
	"encoding/json"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type SkuSizeChartValidate struct {
	id string
}

func (n *SkuSizeChartValidate) SetID(id string) {
	n.id = id
}

func (n SkuSizeChartValidate) GetID() (id string, err error) {
	return n.id, nil
}

func (a SkuSizeChartValidate) Name() string {
	return "SizeChartValidationNode"
}

func (a SkuSizeChartValidate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	reqBody, err := utils.GetPostData(io)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Unable to get Request Body data",
			DeveloperMessage: err.Error(),
		}
	}
	// Log the request body
	logger.Warning(string(reqBody))

	var sizechart sizeUtils.SizeChart
	err = json.Unmarshal(reqBody, &sizechart)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          sizeUtils.InvalidJson,
			DeveloperMessage: err.Error(),
		}
	}
	// validates the sku sizechart and return the failed sku and brandwise sizchart data
	brandWiseScData, failedSkus, successSkus := validateSkuSizeChart(sizechart)

	if len(brandWiseScData) == 0 {
		return io, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          "Validation for all skus failed",
			DeveloperMessage: sizeUtils.FailedValidation}
	}
	io.IOData.Set(sizeUtils.SizeChartBrandWise, brandWiseScData)
	io.IOData.Set(sizeUtils.SizeChartInput, sizechart)
	io.IOData.Set(sizeUtils.FailedSkus, failedSkus)
	io.IOData.Set(sizeUtils.SuccessSkus, successSkus)

	return io, nil
}
