package http

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/constants"
)

//HttpGet makes an http get request with default header parameters
func HttpGet(url string, headers map[string]string,
	timeOut time.Duration) (ret *APIResponse, err error) {
	var resp *http.Response // response
	defer recoverFromPanic("HttpGet()", resp, &err)
	req, err := getReqWithoutBody("GET", url)
	if err != nil {
		return nil, err
	}
	return httpExecuter(req, resp, headers, timeOut)
}

//HttpPost makes an http post request with given parameters
func HttpPost(url string, headers map[string]string, body string,
	timeOut time.Duration) (ret *APIResponse, err error) {
	var resp *http.Response // response
	defer recoverFromPanic("HttpPost()", resp, &err)
	req, err := getReqWithBody("POST", url, body)
	if err != nil {
		return nil, err
	}
	return httpExecuter(req, resp, headers, timeOut)
}

//HttpPut makes an http put request with given parameters
func HttpPut(url string, headers map[string]string, body string,
	timeOut time.Duration) (ret *APIResponse, err error) {
	var resp *http.Response // response
	defer recoverFromPanic("HttpPut()", resp, &err)
	req, err := getReqWithBody("PUT", url, body)
	if err != nil {
		return nil, err
	}
	return httpExecuter(req, resp, headers, timeOut)
}

//HttpDelete makes an http delete request with given parameters
func HttpDelete(url string, headers map[string]string, body string,
	timeOut time.Duration) (ret *APIResponse, err error) {
	var resp *http.Response // response
	defer recoverFromPanic("HttpDelete()", resp, &err)
	req, err := getReqWithBody("DELETE", url, body)
	if err != nil {
		return nil, err
	}
	return httpExecuter(req, resp, headers, timeOut)
}

//HttpPatch makes an http patch request with given parameters
func HttpPatch(url string, headers map[string]string, body string,
	timeOut time.Duration) (ret *APIResponse, err error) {
	var resp *http.Response // response
	defer recoverFromPanic("HttpPatch()", resp, &err)
	req, err := getReqWithBody("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	return httpExecuter(req, resp, headers, timeOut)
}

// recoverFromPanic closes any open http response, recovers with panic details
func recoverFromPanic(name string, resp *http.Response, err *error) {
	if resp != nil {
		resp.Body.Close()
	}
	if r := recover(); r != nil {
		*err = errors.New(name + ":" + fmt.Sprintf("%s", r))
	}
}

// httpExecuter executes the http call with given headers and timeout
func httpExecuter(req *http.Request, resp *http.Response, headers map[string]string, timeOut time.Duration) (*APIResponse, error) {
	ret := new(APIResponse)
	var err error

	// add client headers
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	var client *http.Client
	// set client
	if isPoolSet() {
		if err = incNumCon(); err != nil {
			return ret, err
		}
		defer func() {
			if derr := decNumCon(); derr != nil {
				err = derr // set decrement error
			}
		}()
		client = &http.Client{Transport: poolObj.transport, Timeout: timeOut}
	} else {
		client = &http.Client{Timeout: timeOut}
	}
	resp, err = client.Do(req)
	if err != nil {
		return ret, err
	}
	// read http status
	ret.HttpStatus = constants.HttpCode(resp.StatusCode)

	// read headers
	ret.Headers = make(map[string]string)
	for h, v := range resp.Header {
		ret.Headers[h] = v[0]
	}
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	// read body
	body, berr := ioutil.ReadAll(reader)
	if berr != nil {
		return nil, berr
	}
	ret.Body = body
	// return
	return ret, err
}

// getReqWithBody returns http request for given url and body
func getReqWithBody(name string, url string, body string) (req *http.Request, err error) {
	req, err = http.NewRequest(name, url, strings.NewReader(body))
	if err == nil {
		req.Close = true
		// must for body
		req.Header.Set("Content-Type", "application/json")
	}
	// return
	return req, err
}

// getReqWithoutBody returns http request for given url
func getReqWithoutBody(name string, url string) (req *http.Request, err error) {
	req, err = http.NewRequest(name, url, nil)
	if err == nil {
		req.Close = true
	}
	// return
	return req, err
}
