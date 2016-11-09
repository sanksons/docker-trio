package http

import (
	"github.com/jabong/floRest/src/common/constants"
	"github.com/twinj/uuid"
)

//GetHTTPHeaders returns a map of required headers
//GetHTTPHeaders reads the header values from input request context. A new transaction id is
//created for each call to this method
func GetHTTPHeaders(rc *RequestContext) map[string]string {
	if rc == nil {
		return nil
	}

	headerMap := make(map[string]string, 4)
	headerMap[constants.JABONG_USER_ID] = rc.UserId
	headerMap[constants.JABONG_SESSION_ID] = rc.SessionId
	headerMap[constants.JABONG_REQUEST_ID] = rc.RequestId
	headerMap[constants.JABONG_TRANSACTION_ID] = GetTransactionId()
	return headerMap
}

//GetTransactionId returns a new v4 UUID
func GetTransactionId() string {
	return uuid.NewV4().String()
}
