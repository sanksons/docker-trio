package appconstant

import (
	florest_Constant "github.com/jabong/floRest/src/common/constants"
)

const (
	InconsistantDataStateErrorCode       florest_Constant.AppErrorCode = 1407
	FunctionalityNotImplementedErrorCode florest_Constant.AppErrorCode = 1408
)

const (
	HttpStatusNotImplementedErrorCode florest_Constant.HttpCode = 501
)

var AppErrorCodeToHttpCodeMap = map[florest_Constant.AppErrorCode]florest_Constant.HttpCode{
	InconsistantDataStateErrorCode:       florest_Constant.HttpStatusInternalServerErrorCode,
	FunctionalityNotImplementedErrorCode: HttpStatusNotImplementedErrorCode,
}
