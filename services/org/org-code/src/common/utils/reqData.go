package utils

import (
	"encoding/json"
	"errors"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
)

// GetPostData -> Returns body of post request.
// Returns []byte for easier unmarshal.
func GetPostData(io orchestrator.WorkFlowData) ([]byte, error) {
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		return nil, errors.New(MalformedRequest)
	}
	body, err := appHTTPReq.BodyParameter()
	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}

// function that return the request header value
func GetRequestHeader(io orchestrator.WorkFlowData) (map[string]interface{}, error) {
	headerMap := make(map[string]interface{})
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		return nil, errors.New(MalformedRequest)
	}
	header, err := json.Marshal(appHTTPReq.OriginalRequest.Header)
	if err != nil {
		return nil, errors.New("Error while Marshalling Header")
	}
	err = json.Unmarshal(header, &headerMap)
	if err != nil {
		return nil, errors.New("Error while Unmarshalling Header")
	}
	return headerMap, nil
}
