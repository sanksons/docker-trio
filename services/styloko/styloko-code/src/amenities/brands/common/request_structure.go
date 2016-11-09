package common

import (
	utilHttp "github.com/jabong/floRest/src/common/utils/http"
)

type RequestParams struct {
	RequestContext utilHttp.RequestContext
	QueryParams    QueryParams
}

type QueryParams struct {
	id int
}
