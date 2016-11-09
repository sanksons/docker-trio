package http

import (
	"github.com/jabong/floRest/src/common/constants"
)

type Debug struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResponseMetaData struct {
	UrlParams   map[string]interface{} `json:"urlParams"`
	ApiMetaData map[string]interface{} `json:apiMetaData`
}

func NewResponseMetaData() *ResponseMetaData {
	r := new(ResponseMetaData)
	r.UrlParams = make(map[string]interface{})
	r.ApiMetaData = make(map[string]interface{})
	return r
}

type Response struct {
	Status    constants.AppHttpStatus `json:"status"`
	Data      interface{}             `json:"data"`
	DebugData []Debug                 `json:"debugData,omitempty"`
	MetaData  *ResponseMetaData       `json:"_metaData,omitempty"`
}

type APIResponse struct {
	HttpStatus constants.HttpCode
	Headers    map[string]string
	Body       []byte
}

func NewAPIResponse() APIResponse {
	a := APIResponse{}
	a.Headers = make(map[string]string)
	return a
}
