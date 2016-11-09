package utils

import (
	"errors"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"path"
	"strings"
)

// GetPathParams -> Returns parameter array for request path.
// Builds on the logic that first three words are API endpoint, rest are params.
func GetPathParams(io orchestrator.WorkFlowData) (params []string, err error) {
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		return nil, errors.New(MalformedRequest)
	}
	urlString := appHTTPReq.PathParameter()
	remQuery := strings.Split(urlString, "?")
	finURL := path.Clean(remQuery[0])
	allSubDirs := strings.Split(finURL, "/")
	if allSubDirs[0] == "" {
		allSubDirs = allSubDirs[1:]
	}
	if len(allSubDirs) <= 3 {
		if len(allSubDirs) == 2 {
			return nil, errors.New(PathNoParams)
		}
		return nil, nil
	}
	return allSubDirs[3:], nil
}

// GetQueryParams -> Returns query param data for qs.
// Will return false if data doesn't exist or req failure.
func GetQueryParams(io orchestrator.WorkFlowData, query string) (string, bool) {
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		logger.Error("Bad request. Fail at GetQueryParams.")
		return "", false
	}
	rawReq := appHTTPReq.OriginalRequest
	q := rawReq.FormValue(query)
	if q == "" {
		logger.Info("No valid query params found.")
		return "", false
	}
	return q, true
}

// GetRawQueryMap -> Returns a map[string][]string of all query values.
func GetRawQueryMap(io orchestrator.WorkFlowData) (map[string][]string, error) {
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		logger.Error("Bad request. Fail at GetRawQueryMap.")
		return nil, errors.New(MalformedRequest)
	}
	rawReq := appHTTPReq.OriginalRequest

	// Forced typeCast to map[string][]string
	qParams := map[string][]string(rawReq.Form)
	return qParams, nil
}

// GetSearchQueries -> Returns search params from the q value in queryparams.
// queryDict is stored inside constants. Must be passed manually from the workflow.
// sample query str => orgType.eq~SELLER___status.eq~ACTIVE___address.city.eq~gurgaon
func GetSearchQueries(io orchestrator.WorkFlowData) (searchMap []map[string]interface{}, flag bool, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			searchMap, flag, err = nil, false, errors.New(SearchQueryFail)
		}
	}()
	httpReq, _ := io.IOData.Get(florest_constants.REQUEST)
	appHTTPReq, ok := httpReq.(*utilhttp.Request)
	if !ok || appHTTPReq == nil {
		logger.Error("Bad request. Fail at GetSearchQueries.")
		return nil, false, errors.New(MalformedRequest)
	}
	rawReq := appHTTPReq.OriginalRequest
	q := rawReq.FormValue("q")
	if q == "" {
		logger.Info("Cannot find q in query params.")
		return nil, false, nil
	}
	searchMap = buildSearchMap(q)
	return searchMap, true, nil
}

// buildSearchMap => Generates the map for search queries if found.
func buildSearchMap(query string) []map[string]interface{} {
	baseParams := strings.Split(query, "___")
	var searchArray []map[string]interface{}
	for _, value := range baseParams {
		nest := strings.Split(value, ".")
		data := recursiveMapper(nest)
		searchArray = append(searchArray, data)
	}
	return searchArray
}

// recursiveMapper => creates nested mapstring interfaces for search query
func recursiveMapper(original []string) map[string]interface{} {
	if len(original) > 1 {
		tmp := make(map[string]interface{})
		tmp[original[0]] = recursiveMapper(original[1:])
		return tmp
	}
	keyVal := strings.Split(original[0], "~")
	tmp := make(map[string]interface{})
	tmp["operator"] = keyVal[0]
	tmp["value"] = keyVal[1]
	return tmp
}
