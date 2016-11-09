package cache

//Different type of Cache Implementations
const (
	Memcache           string = "memcache"
	FileCache          string = "file"
	CentralCache       string = "centralCache"
	InmemoryCache      string = "inmemoryCache"
	CentralConfigCache string = "centralconfigcache"
	CentralCacheTest   string = "centralCacheTest"
)

//Rest API Verbs used in cache implementations
const (
	httpGet    string = "GET"
	httpPut    string = "PUT"
	httpDelete string = "DELETE"
)
