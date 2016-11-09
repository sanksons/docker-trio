package common

import (
	"net/http/httptest"
	"test/servicetest"
	"test/utils"
)

// Get -> Wrapper for GET request on CreateTestRequest
func Get(urlString string) *httptest.ResponseRecorder {
	return servicetest.GetResponse(utils.CreateTestRequest("GET", urlString))
}

// Post -> Wrapper for POST request on CreateTestRequestWithBody
func Post(urlString string, body string) *httptest.ResponseRecorder {
	return servicetest.GetResponse(utils.CreateTestRequestWithBody("POST", urlString, body))
}

// Put -> Wrapper for PUT request on CreateTestRequestWithBody
func Put(urlString string, body string) *httptest.ResponseRecorder {
	return servicetest.GetResponse(utils.CreateTestRequestWithBody("PUT", urlString, body))
}
