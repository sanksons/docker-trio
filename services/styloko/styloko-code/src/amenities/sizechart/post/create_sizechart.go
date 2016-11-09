package post

import (
	sizeUtils "amenities/sizechart/common"
	"common/appconstant"
	"common/utils"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type SizeChartCreate struct {
	id string
}

func (n *SizeChartCreate) SetID(id string) {
	n.id = id
}

func (n SizeChartCreate) GetID() (id string, err error) {
	return n.id, nil
}

func (a SizeChartCreate) Name() string {
	return "SizeChartCreationNode"
}

func (a SizeChartCreate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.SizeChartCreation)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.SizeChartCreation)
	}()
	// check if sizechart is sku / brand-brick level
	ty, _ := utils.GetRequestHeader(io, sizeUtils.SizeChartHeader)
	// get the brand wise size chart data after successful validation
	parsedSC, err := io.IOData.Get(sizeUtils.SizeChartBrandWise)

	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          sizeUtils.CreationFailed,
			DeveloperMessage: err.Error(),
		}
	}

	inputSC, err := io.IOData.Get(sizeUtils.SizeChartInput)

	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          sizeUtils.CreationFailed,
			DeveloperMessage: err.Error(),
		}
	}
	collec, errCollec := createSizeChartCollec(inputSC.(sizeUtils.SizeChart), parsedSC, ty)
	if errCollec != "" {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          errCollec,
			DeveloperMessage: sizeUtils.CreationFailed,
		}
	}
	res := storeSizeChartMongo(collec)
	if res == false {
		return io, &florest_constants.AppError{
			Code:             appconstant.FunctionalityNotImplementedErrorCode,
			Message:          sizeUtils.CreationFailed,
			DeveloperMessage: "Record could not be inserted to mongo",
		}
	}
	// Once sizechart is saved, save default values for it, like
	// sizechart_active is 1 etc.
	err = SaveDefaultSettingForSizeChart(inputSC.(sizeUtils.SizeChart).CategoryId)
	if err != nil {
		logger.Error("Unable to store default settings for sizechart", err.Error())
	}

	// Save sizechart collection in constant
	io.IOData.Set(sizeUtils.SavedSizeChCollec, collec)
	// start the job, push the data to channel to be used by worker
	sizechartCreatePool.StartJob(collec)

	return io, nil
}
