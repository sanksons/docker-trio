package post

import (
	"amenities/categories/common"
	"common/appconstant"
	"fmt"
	"reflect"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// CategoriesErrorResponse -> struct for node based data
type CategoriesErrorResponse struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesErrorResponse) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesErrorResponse) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesErrorResponse) Name() string {
	return "CategoriesErrorResponse"
}

// Execute -> Starts node execution.
func (cs CategoriesErrorResponse) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	res, _ := io.IOData.Get(common.CATEGORY_ERROR)
	errs, ok := res.(florest_constants.AppErrors)
	fmt.Println(reflect.TypeOf(errs))
	if !ok {
		str, _ := res.(string)
		return io, &florest_constants.AppError{Code: appconstant.FunctionalityNotImplementedErrorCode, Message: str}
	}
	fmt.Printf("%v\n", errs.Errors)
	return io, &errs
}
