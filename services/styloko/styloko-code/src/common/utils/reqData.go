package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
)

var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	}}

// GetPostData returns []byte for easier unmarshal.
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

// GetRequestHeader returns headers for requested string
func GetRequestHeader(io orchestrator.WorkFlowData, name string) (string, error) {
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		return "", errors.New(MalformedRequest)
	}
	return appHTTPReq.OriginalRequest.Header.Get(name), nil
}

func CreateRequestWithBody(httpMethod,
	urlString, reqBody string) (*http.Request, error) {

	request, err := http.NewRequest(httpMethod, urlString, nil)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader([]byte(reqBody))
	bCloser := ioutil.NopCloser(b)
	request.Body = bCloser
	return request, nil
}
