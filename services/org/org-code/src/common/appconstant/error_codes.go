package appconstant

import (
	florest_Constant "github.com/jabong/floRest/src/common/constants"
)

const (
	InconsistantDataStateErrorCode       florest_Constant.AppErrorCode = 17200
	InvalidIdCode                        florest_Constant.AppErrorCode = 17201
	FunctionalityNotImplementedErrorCode florest_Constant.AppErrorCode = 17202
	IncorrectParameters                  florest_Constant.AppErrorCode = 17203
	BadRequestCode                       florest_Constant.AppErrorCode = 17204
	ResourceNotFoundCode                 florest_Constant.AppErrorCode = 17205
	InvalidPath                          florest_Constant.AppErrorCode = 17206
	DataTypeMismatch                     florest_Constant.AppErrorCode = 17207
	MigrationError                       florest_Constant.AppErrorCode = 17208
)

const (
	HttpStatusNotImplementedErrorCode florest_Constant.HttpCode = 501
	HttpBadRequestCode                florest_Constant.HttpCode = 400
	HTTPResourceNotFound              florest_Constant.HttpCode = 404
	UnableToUnmarshall                florest_Constant.HttpCode = 500
)

var AppErrorCodeToHttpCodeMap = map[florest_Constant.AppErrorCode]florest_Constant.HttpCode{
	InconsistantDataStateErrorCode:       florest_Constant.HttpStatusInternalServerErrorCode,
	FunctionalityNotImplementedErrorCode: HttpStatusNotImplementedErrorCode,
	BadRequestCode:                       HttpBadRequestCode,
	ResourceNotFoundCode:                 HTTPResourceNotFound,
	InvalidPath:                          HttpBadRequestCode,
	DataTypeMismatch:                     UnableToUnmarshall,
	MigrationError:                       UnableToUnmarshall,
}
