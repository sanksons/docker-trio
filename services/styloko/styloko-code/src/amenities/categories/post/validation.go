package post

import (
	"amenities/categories/common"
	"common/appconstant"
	"common/utils"
	"encoding/json"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// CategoriesCreateValidation -> struct for node based data
type CategoriesCreateValidation struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesCreateValidation) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesCreateValidation) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesCreateValidation) Name() string {
	return "CategoriesCreateValidation"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (cs CategoriesCreateValidation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	data, err := utils.GetPostData(io)
	if err != nil {
		io.IOData.Set(common.CATEGORY_ERROR, common.NO_DATA)
		return false, nil
	}
	// errMap := make(map[string]interface{})
	ctgryCreate := new(common.CategoryCreate)
	err = json.Unmarshal(data, &ctgryCreate)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.InvalidDataErrorCode, Message: "Invalid Data in JSON", DeveloperMessage: err.Error()}
	}
	// Path based validation
	pathParams, _ := utils.GetPathParams(io)
	if len(pathParams) != 0 {
		return false, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid path provided.", DeveloperMessage: common.INVALID_PATH}
	}

	appErrors, flag := common.ValidateV2(*ctgryCreate)
	if !flag {
		return flag, &appErrors
	}
	io.IOData.Set(common.CATEGORY_VALID_DATA, ctgryCreate)
	return true, nil
}
