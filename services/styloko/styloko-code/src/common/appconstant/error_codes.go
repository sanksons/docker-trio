package appconstant

import (
	florest_Constant "github.com/jabong/floRest/src/common/constants"
)

const (
	InconsistantDataStateErrorCode       florest_Constant.AppErrorCode = 1407
	FunctionalityNotImplementedErrorCode florest_Constant.AppErrorCode = 1408
	DataNotFoundErrorCode                florest_Constant.AppErrorCode = 17400
	InvalidDataErrorCode                 florest_Constant.AppErrorCode = 17401
	FailedToCreateErrorCode              florest_Constant.AppErrorCode = 17402
	BadRequestCode                       florest_Constant.AppErrorCode = 17403
	ResourceNotFoundCode                 florest_Constant.AppErrorCode = 17404
	ServiceFailureCode                   florest_Constant.AppErrorCode = 17405
)

// These are HTTP Status codes. Please do not put values like 1700, etc here.
// Visit : https://en.wikipedia.org/wiki/List_of_HTTP_status_codes , for a list
// of valid HTTP status codes if you wish to add any new code.
const (
	HttpServiceFailureErrorCode       florest_Constant.HttpCode = 500
	HttpStatusNotImplementedErrorCode florest_Constant.HttpCode = 501
	HttpDataNotFoundErrorCode         florest_Constant.HttpCode = 404
	HttpInvalidDataErrorCode          florest_Constant.HttpCode = 400
	HttpCreateFailErrorCode           florest_Constant.HttpCode = 400
	HttpBadRequestCode                florest_Constant.HttpCode = 400
)

// AppErrorCodeToHttpCodeMap is a map which makes a map of HTTP codes and AppError codes.
// If you wish to add any new AppErrorCode then please add a new Key: value pair to this map.
var AppErrorCodeToHttpCodeMap = map[florest_Constant.AppErrorCode]florest_Constant.HttpCode{
	InconsistantDataStateErrorCode:       florest_Constant.HttpStatusInternalServerErrorCode,
	FunctionalityNotImplementedErrorCode: HttpStatusNotImplementedErrorCode,
	DataNotFoundErrorCode:                HttpDataNotFoundErrorCode,
	InvalidDataErrorCode:                 HttpInvalidDataErrorCode,
	FailedToCreateErrorCode:              HttpCreateFailErrorCode,
	BadRequestCode:                       HttpBadRequestCode,
	ServiceFailureCode:                   HttpServiceFailureErrorCode,
	ResourceNotFoundCode:                 HttpDataNotFoundErrorCode,
}
