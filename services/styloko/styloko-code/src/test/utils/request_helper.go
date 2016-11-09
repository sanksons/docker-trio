package utils

import (
	"bytes"

	"io/ioutil"
	"net/http"
)

// CreateTestRequest -> generates a test request, mostly useful for GET requests.
func CreateTestRequest(httpMethod string, urlString string) *http.Request {
	if string(urlString[0]) != "/" {
		urlString = "/" + urlString
	}
	request, _ := http.NewRequest(httpMethod, urlString, nil)
	request.RequestURI = urlString

	return request
}

// CreateTestRequestWithBody -> generates test request for POST & PUT request methods.
func CreateTestRequestWithBody(httpMethod string, urlString string, reqBody string) *http.Request {
	if string(urlString[0]) != "/" {
		urlString = "/" + urlString
	}
	request, _ := http.NewRequest(httpMethod, urlString, nil)
	request.RequestURI = urlString
	b := bytes.NewReader([]byte(reqBody))
	bCloser := ioutil.NopCloser(b)
	request.Body = bCloser
	return request
}

func SetHeadersInRequest(headers map[string]string, request *http.Request) *http.Request {
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return request
}
