package mongodb

import (
	"errors"
)

var ErrNotFound error = errors.New("Data Not Found")

type MongodbError struct {
	ErrCode          string
	DeveloperMessage string
}

func (e *MongodbError) Error() string {
	if e == nil {
		return ""
	}
	return e.DeveloperMessage
}

const (
	ERR_INITIALIZATION          = "Initialization failed"
	ERR_FINDONE_FAILURE         = "Failure in FindOne() method"
	ERR_FINDALL_FAILURE         = "Failure in FindAll() method"
	ERR_INSERT_FAILURE          = "Failure in Insert() method"
	ERR_UPDATE_FAILURE          = "Failure in Update() method"
	ERR_REMOVE_FAILURE          = "Failure in Remove() method"
	ERR_FIND_AND_MODIFY_FAILURE = "Failure in FindAndModify() method"
)

// getErrObj returns error object with given details
func getErrObj(errCode string, developerMessage string) (ret *MongodbError) {
	ret = new(MongodbError)
	ret.ErrCode = errCode
	ret.DeveloperMessage = developerMessage
	return ret
}
