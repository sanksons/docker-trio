package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func getAppRequest(r *http.Request) (*utilhttp.Request, error) {
	req, rerr := utilhttp.GetRequest(r)
	if rerr != nil {
		return nil, rerr
	}

	return &req, nil
}

func getBucketsMap(bucketsList string) map[string]string {
	buckets := strings.Split(bucketsList, constants.FIELD_SEPARATOR)
	bucketMap := make(map[string]string, len(buckets))
	for _, v := range buckets {
		bKV := strings.Split(v, constants.KEY_VALUE_SEPARATOR)
		if len(bKV) < 2 { //invalid bucket
			continue
		}
		bucketMap[bKV[0]] = bKV[1]
	}
	return bucketMap
}

//Get the Service WorkFlow Data
func GetData(r *http.Request) (*orchestrator.WorkFlowData, error) {
	serviceInputOutput := new(orchestrator.WorkFlowIOInMemoryImpl)
	appReq, rerr := getAppRequest(r)
	if rerr != nil {
		return nil, rerr
	}

	serviceInputOutput.Set(constants.URI, appReq.URI)
	serviceInputOutput.Set(constants.HTTPVERB, appReq.HTTPVerb)
	serviceInputOutput.Set(constants.REQUEST, appReq)

	logger.Info(fmt.Sprintf("Service Input Output %v", serviceInputOutput))

	serviceEcContext := new(orchestrator.WorkFlowECInMemoryImpl)
	serviceEcContext.Set(constants.USER_ID, appReq.Headers.UserId)
	serviceEcContext.Set(constants.SESSION_ID, appReq.Headers.SessionId)
	serviceEcContext.Set(constants.AUTHTOKEN, appReq.Headers.AuthToken)
	serviceEcContext.Set(constants.USER_AGENT, appReq.Headers.UserAgent)
	serviceEcContext.Set(constants.HTTP_REFERRER, appReq.Headers.Referrer)
	serviceEcContext.Set(constants.REQUEST_ID, appReq.Headers.RequestId)
	serviceEcContext.SetBuckets(getBucketsMap(appReq.Headers.BucketsList))
	serviceEcContext.SetDebugFlag(appReq.Headers.Debug)

	serviceEcContext.Set(constants.REQUEST_CONTEXT,
		utilhttp.RequestContext{
			AppName:       config.GlobalAppConfig.AppName,
			UserId:        appReq.Headers.UserId,
			SessionId:     appReq.Headers.SessionId,
			RequestId:     appReq.Headers.RequestId,
			TransactionId: appReq.Headers.TransactionId,
			URI:           appReq.URI,
			ClientAppId:   appReq.Headers.ClientAppId,
		})

	logger.Info(fmt.Sprintf("Service Execution Context %v", serviceInputOutput))

	serviceWorkFlowData := new(orchestrator.WorkFlowData)
	serviceWorkFlowData.Create(serviceInputOutput, serviceEcContext)

	return serviceWorkFlowData, nil
}
