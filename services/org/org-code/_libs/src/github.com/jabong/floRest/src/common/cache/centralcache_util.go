package cache

import (
	"encoding/json"
	"errors"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//getCCErrorResponse returns a CentralCacheError from a json encoded response body
func getCCErrorResponse(res []byte) (*CentralCacheError, error) {
	cerror := new(CentralCacheErrorResponse)
	if err := json.Unmarshal(res, cerror); err != nil {
		logger.Error("Json unmarshal or error response failed with error - " + err.Error())
		return nil, err
	}
	if len(cerror.Errors) <= 0 {
		errMsg := "Error Response returned zero errors"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	cerr := cerror.Errors[0]
	return &cerr, nil
}
