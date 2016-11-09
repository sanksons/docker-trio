package examples

import (
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

type HelloWorldMultiError struct {
	id string
}

func (n *HelloWorldMultiError) SetID(id string) {
	n.id = id
}

func (n HelloWorldMultiError) GetID() (id string, err error) {
	return n.id, nil
}

func (a HelloWorldMultiError) Name() string {
	return "HelloWord"
}

func (a HelloWorldMultiError) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	//Business Logic

	err := new(constants.AppErrors)
	err.Errors = append(err.Errors, constants.AppError{
		Code:             constants.ParamsInSufficientErrorCode,
		Message:          "In Sufficient Params Error Code",
		DeveloperMessage: "In Sufficient Params Error Code",
	})
	err.Errors = append(err.Errors, constants.AppError{
		Code:             constants.ResourceErrorCode,
		Message:          "Resource Error Code",
		DeveloperMessage: "Resource Error Code",
	})
	return io, err
}
