package cache

import ()

type CentralCacheResponseStatus struct {
	Success bool
}

type CentralCacheError struct {
	Code             int
	Message          string
	DeveloperMessage string
}

type CentralCacheResponse struct {
	Status CentralCacheResponseStatus
}

//CentralCacheErrorResponse denotes the structure of a response in case there is an error
type CentralCacheErrorResponse struct {
	CentralCacheResponse
	Errors []CentralCacheError
}

type CentralCachedData struct {
	Key   string
	Value string
}

//CentralCacheGetResponse denotes the structure of a get response
type CentralCacheGetResponse struct {
	CentralCacheResponse
	Data CentralCachedData
}

//CentralCacheGetBatchResponse denotes the structure of a getBatch response
type CentralCacheGetBatchResponse struct {
	CentralCacheResponse
	Data []CentralCachedData
}

//CentralCachePutRequest denotes the structure for making a PUT request to add an item
//in cache
type CentralCachePutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`

	//Ttl denotes the time for which an item should persist in cache
	Ttl int32 `json:"ttl"`
}
