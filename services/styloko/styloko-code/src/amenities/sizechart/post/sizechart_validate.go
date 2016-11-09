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

type SizeChartValidate struct {
	id string
}

func (n *SizeChartValidate) SetID(id string) {
	n.id = id
}

func (n SizeChartValidate) GetID() (id string, err error) {
	return n.id, nil
}

func (a SizeChartValidate) Name() string {
	return "SizeChartValidationNode"
}

func (a SizeChartValidate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
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
	var sizeChart sizeUtils.SizeChart
	err = json.Unmarshal(reqBody, &sizeChart)
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          sizeUtils.InvalidJson,
			DeveloperMessage: err.Error(),
		}
	}
	// validates the sizechart data, and returns brandwise sizechart if data is valid
	data, invalid := validateSizeChart(sizeChart)
	if invalid != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.InvalidDataErrorCode,
			Message:          invalid.(string),
			DeveloperMessage: sizeUtils.FailedValidation,
		}
	}
	// Stores the BrandWiseSizechart and raw sizechart data for further processing
	io.IOData.Set(sizeUtils.SizeChartBrandWise, data)
	io.IOData.Set(sizeUtils.SizeChartInput, sizeChart)

	return io, nil
}
