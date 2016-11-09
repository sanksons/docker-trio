package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// CentralConfigCacheImpl denotes Central Cache for Config. To know more about central config cache please refer
// https://wiki.jira.rocket-internet.de/display/INDFAS/Cache+Service+Cluster
type CentralConfigCacheImpl struct {
	baseUrl      string
	keyPrefix    string
	dumpFilePath string
	expirySec    int32
	timeOut      time.Duration
}

//Init initialises c to connect to server with all keys stored under the bucket keyPrefix
func (c *CentralConfigCacheImpl) Init(conf Config) {
	c.baseUrl = conf.Host
	c.keyPrefix = conf.KeyPrefix
	c.dumpFilePath = conf.DumpFilePath
	c.expirySec = conf.ExpirySec
	c.timeOut = time.Duration(conf.TimeOut * time.Millisecond)
}

//Get gets key from central cache
func (c *CentralConfigCacheImpl) Get(key string, serialize bool, compress bool) (*Item, error) {
	url := c.getServerUrl(key, true)
	res, err := http.HttpGet(url, nil, c.timeOut)
	if err != nil {
		return nil, err
	}
	if _, ok := centralCacheErrorHeaderStatus[res.HttpStatus]; ok {
		errResp, err := getCCErrorResponse(res.Body)
		if err != nil {
			return nil, err
		}
		logger.Error(errResp.Message)
		return nil, getCacheErrorType(CentralCache, errResp.Code)
	}
	cr := new(CentralConfigCacheGetResponse)
	err = json.Unmarshal(res.Body, cr)
	if err != nil {
		logger.Error("Failed to unmarshal the config. Error - " + err.Error())
		return nil, ErrInvalidCacheResponse
	}
	item := new(Item)
	item.Key = cr.Data.Key
	item.Value = cr.Data.Value
	return item, nil
}

//GetBatch gets the list of values stored in central cache indexed by keys. If a key does not exist in central
//cache then it has nil set for that key value in the returned list
func (c *CentralConfigCacheImpl) GetBatch(keys []string, serialize bool, compress bool) (map[string]*Item, error) {
	return nil, ErrUnsupportedOperation
}

//Set stores an item in central cache for expirySec
func (c *CentralConfigCacheImpl) Set(item Item, serialize bool, compress bool) error {
	return ErrUnsupportedOperation
}

// SetWithTimeout stores an item in central cache for provided TTL
func (c *CentralConfigCacheImpl) SetWithTimeout(item Item, serialize bool, compress bool, ttl int32) error {
	return ErrUnsupportedOperation
}

// Delete deletes an item from central cache
func (c *CentralConfigCacheImpl) Delete(key string) error {
	return ErrUnsupportedOperation
}

// DeleteBatch deletes an array of keys from central cache
func (c *CentralConfigCacheImpl) DeleteBatch(keys []string) error {
	return ErrUnsupportedOperation
}

// Dump dumps data to central cache
func (c *CentralConfigCacheImpl) Dump(key string, value []byte) error {
	return ErrUnsupportedOperation
}

//getServerUrl returns the central cache server url from the supplied key
func (c *CentralConfigCacheImpl) getServerUrl(key string, isSingleGetReq bool) string {
	return fmt.Sprintf("%s/%s", c.baseUrl, key)
}

//newCentralCache returns a new instance of CentralCacheImpl from the supplied conf
func newCentralConfigCache(conf Config) (*CentralConfigCacheImpl, error) {
	c := new(CentralConfigCacheImpl)
	c.Init(conf)
	return c, nil
}
