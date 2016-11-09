package constants

import (
	"fmt"
)

type HttpCode uint16
type AppErrorCode uint16

//Http App Status
type AppHttpStatus struct {
	HttpStatusCode HttpCode   `json:"httpStatusCode"`
	Success        bool       `json:"success"`
	Errors         []AppError `json:"errors"`
}

//App Error Status
type AppError struct {
	Code             AppErrorCode `json:"code"`
	Message          string       `json:"message"`
	DeveloperMessage string       `json:"developerMessage"`
}

func (e AppError) Error() string { return e.Message }

//List of App Errors
type AppErrors struct {
	Errors []AppError
}

func (e AppErrors) Error() string {
	var s string
	for _, appError := range e.Errors {
		s = fmt.Sprintf("%s\n%s", s, appError.Message)
	}
	return s
}

const (
	ParamsInSufficientErrorCode AppErrorCode = 1401
	ParamsInValidErrorCode      AppErrorCode = 1402
	IncorrectDataErrorCode      AppErrorCode = 1403
	InvalidUrlKeyErrorCode      AppErrorCode = 1404

	ResourceErrorCode AppErrorCode = 1501
	DbErrorCode       AppErrorCode = 1502
	IndexErrorCode    AppErrorCode = 1503
	CacheErrorCode    AppErrorCode = 1504

	InvalidRequestUri AppErrorCode = 1601

	InvalidErrorCode = 2501
)

const (
	HttpStatusSuccessCode             HttpCode = 200
	HttpStatusBadRequestCode          HttpCode = 400
	HttpStatusInternalServerErrorCode HttpCode = 500
	HttpFatalErrorCode                HttpCode = 501
	HttpStatusNotFound                HttpCode = 404
)

var appErrorCodeToHttpCodeMap = map[AppErrorCode]HttpCode{

	ResourceErrorCode: HttpStatusInternalServerErrorCode,
	DbErrorCode:       HttpStatusInternalServerErrorCode,
	IndexErrorCode:    HttpStatusInternalServerErrorCode,
	CacheErrorCode:    HttpStatusInternalServerErrorCode,

	ParamsInSufficientErrorCode: HttpStatusBadRequestCode,
	ParamsInValidErrorCode:      HttpStatusBadRequestCode,
	IncorrectDataErrorCode:      HttpStatusBadRequestCode,
	InvalidUrlKeyErrorCode:      HttpStatusBadRequestCode,
	InvalidRequestUri:           HttpStatusNotFound,

	InvalidErrorCode: HttpFatalErrorCode,
}

func GetAppHttpError(appErrors AppErrors) *AppHttpStatus {
	var httpCode HttpCode = HttpStatusSuccessCode

	//Only considering the last app error to generate the http code
	if appErrors.Errors != nil && len(appErrors.Errors) > 0 {
		lastAppError := appErrors.Errors[len(appErrors.Errors)-1]
		v, found := appErrorCodeToHttpCodeMap[lastAppError.Code]
		if !found {
			httpCode = InvalidErrorCode
		}
		httpCode = v
	}

	return getAppErrStatus(httpCode, appErrors)
}

func getAppErrStatus(status HttpCode, appErrors AppErrors) *AppHttpStatus {
	var httpStatus = &AppHttpStatus{HttpStatusCode: status}
	var apiErrors []AppError = nil

	if status != HttpStatusSuccessCode {
		apiErrors = appErrors.Errors
		httpStatus.Errors = apiErrors
		httpStatus.Success = false
	} else {
		apiErrors = appErrors.Errors
		httpStatus.Errors = apiErrors
		if apiErrors == nil || len(apiErrors) == 0 {
			httpStatus.Success = true
		} else {
			httpStatus.Success = false
		}
	}

	return httpStatus
}

func UpdateAppHttpError(appErrorCodeMap map[AppErrorCode]HttpCode) {
	for k, v := range appErrorCodeMap {
		appErrorCodeToHttpCodeMap[k] = v
	}
}
