package http

import (
	"net/http"
	"strconv"
)

type RequestHeader struct {
	ContentType   string
	Accept        string
	UserId        string
	SessionId     string
	AuthToken     string
	TransactionId string
	RequestId     string
	Timestamp     string
	UserAgent     string
	Referrer      string
	BucketsList   string
	Debug         bool
	ClientAppId   string
}

func GetReqHeader(req *http.Request) RequestHeader {
	return RequestHeader{
		ContentType:   req.Header.Get("Content-Type"),
		Accept:        req.Header.Get("Accept"),
		UserId:        req.Header.Get("X-Jabong-UserId"),
		SessionId:     req.Header.Get("X-Jabong-SessionId"),
		AuthToken:     req.Header.Get("X-Jabong-Token"),
		TransactionId: req.Header.Get("X-Jabong-Tid"),
		RequestId:     req.Header.Get("X-Jabong-Reqid"),
		Timestamp:     req.Header.Get("ts"),
		UserAgent:     req.Header.Get("User-Agent"),
		Referrer:      req.Header.Get("Referer"),
		BucketsList:   req.Header.Get("bucket"),
		Debug:         getBoolHeaderField(req, "X-Jabong-Debug"),
		ClientAppId:   req.Header.Get("X-Jabong-Appid"),
	}
}

func getBoolHeaderField(req *http.Request, headerKey string) bool {
	value, err := strconv.ParseBool(req.Header.Get(headerKey))
	if err != nil {
		return false
	}
	return value
}
