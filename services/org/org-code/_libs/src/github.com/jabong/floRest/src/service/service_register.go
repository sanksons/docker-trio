package service

import (
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
)

var apiList []ApiInterface
var resourceBucketMapping map[string]string
var apiCustomInitFunc func()
var configEnvUpdateMap map[string]string
var globalEnvUpdateMap map[string]string

func RegisterApi(apiInstance ApiInterface) {
	apiList = append(apiList, apiInstance)
}

func RegisterConfig(applicationConfig interface{}) {
	config.ApplicationConfig = applicationConfig
}

func RegisterHttpErrors(appErrorCodeMap map[constants.AppErrorCode]constants.HttpCode) {
	constants.UpdateAppHttpError(appErrorCodeMap)
}

func RegisterResourceBucketMapping(resource, bucketId string) {
	if len(resourceBucketMapping) == 0 {
		resourceBucketMapping = make(map[string]string)
	}
	resourceBucketMapping[resource] = bucketId
}

func RegisterCustomApiInitFunc(f func()) {
	apiCustomInitFunc = f
}

func RegisterConfigEnvUpdateMap(a map[string]string) {
	configEnvUpdateMap = a
}

func RegisterGlobalEnvUpdateMap(a map[string]string) {
	globalEnvUpdateMap = a
}
