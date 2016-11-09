package common

import (
	"common/appconstant"

	florest_constants "github.com/jabong/floRest/src/common/constants"
)

func GenError(err string, msg string) florest_constants.AppError {
	tmp := florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          err,
		DeveloperMessage: msg,
	}
	return tmp
}
