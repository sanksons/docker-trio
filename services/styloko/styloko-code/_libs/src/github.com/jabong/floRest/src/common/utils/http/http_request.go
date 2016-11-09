package http

import (
	"errors"
	"net/http"
	"strings"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	PUT    HTTPMethod = "PUT"
	POST   HTTPMethod = "POST"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

type Request struct {
	HTTPVerb        HTTPMethod
	URI             string
	OriginalRequest *http.Request
	Headers         RequestHeader
}

func getHTTPMethod(method string) (HTTPMethod, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return GET, nil
	case "PUT":
		return PUT, nil
	case "POST":
		return POST, nil
	case "DELETE":
		return DELETE, nil
	case "PATCH":
		return PATCH, nil
	}
	return "", errors.New("Incorrect HTTP Method")
}

func GetRequest(r *http.Request) (Request, error) {
	httpVerb, verr := getHTTPMethod(r.Method)
	if verr != nil {
		return Request{}, verr
	}

	return Request{
		HTTPVerb:        httpVerb,
		URI:             r.URL.String(),
		OriginalRequest: r,
		Headers:         GetReqHeader(r)}, nil
}

func (req *Request) BodyParameter() (string, error) {
	return getBodyParam(req.OriginalRequest)
}

func (req *Request) PathParameter() string {
	return req.OriginalRequest.URL.String()
}
