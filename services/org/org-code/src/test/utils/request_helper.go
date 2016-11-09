package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func CreateTestRequest(httpMethod string, urlString string) *http.Request {
	request, _ := http.NewRequest(httpMethod, urlString, nil)
	request.RequestURI = urlString
	return request
}

func CreateTestRequestWithBody(httpMethod string, urlString string, reqBody string) *http.Request {
	request, _ := http.NewRequest(httpMethod, urlString, nil)
	request.RequestURI = urlString
	b := bytes.NewReader([]byte(reqBody))
	bCloser := ioutil.NopCloser(b)
	request.Body = bCloser
	return request
}
