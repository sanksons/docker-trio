package mongodb

import ()

type MongodbError struct {
	ErrCode          string
	DeveloperMessage string
}

const (
	ERR_INITIALIZATION          = "Initialization failed"
	ERR_FINDONE_FAILURE         = "Failure in FindOne() method"
	ERR_FINDALL_FAILURE         = "Failure in FindAll() method"
	ERR_INSERT_FAILURE          = "Failure in Insert() method"
	ERR_UPDATE_FAILURE          = "Failure in Update() method"
	ERR_REMOVE_FAILURE          = "Failure in Remove() method"
	ERR_FINDMODIFY_FAILURE      = "Failure in FindAndModify() method"
	ERR_SETCOLLECTION_FAILURE   = "Failure in SetCollection() method"
	ERR_GENNXTSEQUENCE_FAILURE  = "Failure in GenNextSequence() method"
	ERR_SETCOLINCOUNTER_FAILURE = "Failure in SetCollectionInCounter() method"
)

// getErrObj returns error object with given details
func getErrObj(errCode string, developerMessage string) (ret *MongodbError) {
	ret = new(MongodbError)
	ret.ErrCode = errCode
	ret.DeveloperMessage = developerMessage
	return ret
}
