package get

import (
	"amenities/services/attributes"
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"strconv"
	"strings"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// StandardSizeGetValidation -> struct for node based data
type StandardSizeGetValidation struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (ss *StandardSizeGetValidation) SetID(id string) {
	ss.id = id
}

// GetID -> returns current node ID to orchestrator
func (ss StandardSizeGetValidation) GetID() (id string, err error) {
	return ss.id, nil
}

// Name -> Returns node name to orchestrator
func (ss StandardSizeGetValidation) Name() string {
	return "StandardSizeGetValidation"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (ss StandardSizeGetValidation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.GET_VALIDATE)
	defer func() {
		logger.EndProfile(profiler, common.GET_VALIDATE)
	}()

	param, err := utils.GetRawQueryMap(io)
	queryArr, ok := param["q"]
	if !ok {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "No query passed"})
		return false, nil
	}
	limitArr, ok := param["limit"]
	if !ok {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "No limit passed"})
		return false, nil
	}
	offsetArr, ok := param["offset"]
	if !ok {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "No offset passed"})
		return false, nil
	}
	if len(queryArr) == 0 || len(limitArr) == 0 || len(offsetArr) == 0 {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "No limit/offset/query passed"})
		return false, nil
	}
	qArr := strings.Split(queryArr[0], "~")
	if len(qArr) < 2 {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "No query value passed"})
		return false, nil
	}
	query := strings.Split(qArr[0], ".")
	if len(query) < 2 || query[0] != "attributeSet" || query[1] != "eq" {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "Invalid query passed"})
		return false, nil
	}

	mgoSession := mongoFactory.GetMongoSession(common.STANDARDSIZE_SEARCH)
	defer mgoSession.Close()
	attrSet := attributes.GetByName(qArr[1], mgoSession)
	if attrSet.SeqId == 0 {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: "Invalid attribute set passed"})
		return false, nil
	}

	limit, err := strconv.Atoi(limitArr[0])
	if err != nil {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: err.Error()})
		return false, nil
	}
	offset, err := strconv.Atoi(offsetArr[0])
	if err != nil {
		io.IOData.Set(common.STANDARDSIZE_GET_ERROR, florest_constants.AppError{Code: appconstant.BadRequestCode,
			Message: common.DATA_VALIDATION_FAIL, DeveloperMessage: err.Error()})
		return false, nil
	}
	validData := make(map[string]int, 0)
	validData["limit"] = limit
	validData["offset"] = offset
	validData["attrSetId"] = attrSet.SeqId
	io.IOData.Set(common.STANDARDSIZE_GET_VALID_DATA, validData)
	return true, nil
}
