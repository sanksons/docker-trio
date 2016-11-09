package sqldb

import ()

type SqlDbError struct {
	ErrCode          string
	DeveloperMessage string
}

const (
	ERR_NO_DRIVER       = "Driver not found"
	ERR_INITIALIZATION  = "Initialization failed"
	ERR_QUERY_FAILURE   = "Failure in Query() method"
	ERR_EXECUTE_FAILURE = "Failure in Execute() method"
	ERR_PING_FAILURE    = "Failure in Ping() method"
	ERR_GETTXN_FAILURE  = "Failure in GetTxnObj() method"
	ERR_CLOSE_FAILURE   = "Failure in Close() method"
)

// getErrObj returns error object with given details
func getErrObj(errCode string, developerMessage string) (ret *SqlDbError) {
	ret = new(SqlDbError)
	ret.ErrCode = errCode
	ret.DeveloperMessage = developerMessage
	return ret
}
