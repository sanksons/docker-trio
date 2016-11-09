package cache

import (
	"encoding/json"
)

//CentralConfigCacheGetResponse denotes the structure of a get config response
type CentralConfigCacheGetResponse struct {
	CentralConfigCacheResponse
	Data CentralConfigCachedData
}

type CentralConfigCacheResponse struct {
	Status CentralConfigCacheResponseStatus
}

type CentralConfigCacheResponseStatus struct {
	Success bool
}

type CentralConfigCachedData struct {
	Key   string
	Value json.RawMessage
}
