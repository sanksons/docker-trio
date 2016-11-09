package jabongbus

import (
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/utils/http"
	"time"
)

// structure for publisher
type jPublish struct {
	url string
}

// PubRequest request structure for client to publish message
type PubRequest struct {
	Headers struct {
		UserId        string `json:"userId"`
		SessionId     string `json:"sessionId"`
		TransactionId string `json:"tId"`
		RequestId     string `json:"reqId"`
	}
	Body Message
}

// PubResponse response structure for client after publishing message
type PubResponse struct {
	BusResponse   []byte
	DebugResponse struct {
		Url string
	}
}

// GetPublisher get publisher implementation
func GetPublisher(conf *PublisherConfig) Publisher {
	obj := new(jPublish)
	obj.url = conf.Url // set url
	return obj
}

// PublishMessage publish the message sent by client
func (pub *jPublish) PublishMessage(preq *PubRequest, timeOut time.Duration, debug bool) (ret *PubResponse, err error) {
	// get body string
	ary, jerr := messageToString(&preq.Body)
	if jerr != nil {
		return nil, jerr
	}
	// set headers
	headers := make(map[string]string)
	headers[constants.JABONG_USER_ID] = preq.Headers.UserId
	headers[constants.JABONG_SESSION_ID] = preq.Headers.SessionId
	headers[constants.JABONG_TRANSACTION_ID] = preq.Headers.TransactionId
	headers[constants.JABONG_REQUEST_ID] = preq.Headers.RequestId
	// get return object
	ret = new(PubResponse)
	if debug { // set debug values
		ret.DebugResponse.Url = pub.url
	}
	// call http post method
	httpResponse, err := http.HttpPost(pub.url, headers, string(ary), timeOut)
	// error in bus response
	if err != nil {
		return ret, ErrHttpPost
	}
	ret.BusResponse = httpResponse.Body

	// return
	return ret, nil
}
