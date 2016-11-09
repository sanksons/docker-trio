package cache

import (
	"encoding/json"
	"fmt"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	h "net/http"
	"time"
)

// CentralCacheImpl denotes Central Cache. To know more about central cache please refer
// https://wiki.jira.rocket-internet.de/display/INDFAS/Cache+Service+Cluster
type CentralCacheImpl struct {
	baseUrl      string
	keyPrefix    string
	dumpFilePath string
	expirySec    int32
	timeOut      time.Duration
}

var default_headers map[string]string

//centralCacheErrorHeaderStatus is a list of all error header status returned by Central Cache
var centralCacheErrorHeaderStatus = map[constants.HttpCode]bool{
	constants.HttpStatusBadRequestCode:          true,
	constants.HttpStatusNotFound:                true,
	constants.HttpStatusInternalServerErrorCode: true,
}

//Init initialiases c to connect to server with all keys stored under the bucket keyPrefix
func (c *CentralCacheImpl) Init(conf Config) {
	c.baseUrl = conf.Host
	c.keyPrefix = conf.KeyPrefix
	c.dumpFilePath = conf.DumpFilePath
	c.expirySec = conf.ExpirySec
	c.timeOut = time.Duration(time.Millisecond * conf.TimeOut)
	default_headers = make(map[string]string)
	default_headers["Accept-Encoding"] = "gzip, deflate, sdch"
}

//Get gets key from central cache
func (c *CentralCacheImpl) Get(key string, serialize bool, compress bool) (*Item, error) {
	url := c.getServerUrl(key, true)
	res, err := http.HttpGet(url, default_headers, c.timeOut)
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
	cr := new(CentralCacheGetResponse)
	err = json.Unmarshal(res.Body, cr)
	if err != nil {
		return nil, ErrInvalidCacheResponse
	}
	item := new(Item)
	item.Key = cr.Data.Key
	item.Value = cr.Data.Value
	return item, nil
}

//GetBatch gets the list of values stored in central cache indexed by keys. If a key does not exist in central
//cache then it has nil set for that key value in the returned list
func (c *CentralCacheImpl) GetBatch(keys []string, serialize bool, compress bool) (map[string]*Item, error) {
	keycount := len(keys)
	if keycount == 0 {
	    return nil, ErrInvalidData
	} 
	resMap := make(map[string]*Item, keycount)
	keyString := ""
	for i, v := range keys {
		if v != "" {
		    resMap[v] = nil
			if i == (keycount - 1) {
				keyString = fmt.Sprintf("%skey=%s", keyString, v)
			} else {
				keyString = fmt.Sprintf("%skey=%s&", keyString, v)
			}
		}
	}
	if keyString == "" {
		return nil, ErrInvalidData
	}
	url := c.getServerUrl(fmt.Sprintf("bulk?%s", keyString), false)
	res, err := http.HttpGet(url, default_headers, c.timeOut)
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
	cr := new(CentralCacheGetBatchResponse)
	err = json.Unmarshal(res.Body, cr)
	if err != nil {
		return nil, ErrInvalidCacheResponse
	}
	for _, v := range cr.Data {
		item := new(Item)
		item.Key = v.Key
		item.Value = v.Value
		resMap[item.Key] = item
	}
	return resMap, nil
}

//setWithTTL stores an item in central cache for either expirySec or provided TTL value
func (c *CentralCacheImpl) setWithTTL(item Item, serialize bool, compress bool, ttl int32) error {
	req := new(CentralCachePutRequest)
	req.Key = item.Key
	val, ok := item.Value.([]byte)
	if !ok {
		v, ok := item.Value.(string)
		if !ok {
			return ErrNotSupportedFormat
		}
		req.Value = v
	} else {
		req.Value = string(val)
	}
	req.Ttl = ttl
	jsonEn, _ := json.Marshal(req)
	url := c.getServerUrl(item.Key, false)
	res, err := http.HttpPut(url, nil, string(jsonEn), c.timeOut)
	if err != nil {
		return err
	}
	if res.HttpStatus != h.StatusCreated {
		errResp, err := getCCErrorResponse(res.Body)
		if err != nil {
			return err
		}
		logger.Error(errResp.Message)
		return getCacheErrorType(CentralCache, errResp.Code)
	}
	return nil
}

//Set stores an item in central cache for expirySec
func (c *CentralCacheImpl) Set(item Item, serialize bool, compress bool) error {
	return c.setWithTTL(item, serialize, compress, c.expirySec)
}

// SetWithTimeout stores an item in central cache for provided TTL value
func (c *CentralCacheImpl) SetWithTimeout(item Item, serialize bool, compress bool, ttl int32) error {
	return c.setWithTTL(item, serialize, compress, ttl)
}

// Delete deletes an item from central cache
func (c *CentralCacheImpl) Delete(key string) error {
	url := c.getServerUrl(key, false)
	res, err := http.HttpDelete(url, nil, "", c.timeOut)
	if err != nil {
		return err
	}
	if res.HttpStatus != h.StatusNoContent {
		errResp, err := getCCErrorResponse(res.Body)
		if err != nil {
			return err
		}
		logger.Error(errResp.Message)
		return getCacheErrorType(CentralCache, errResp.Code)
	}
	return nil
}

// DeleteBatch deletes an array of keys from central cache
// Blitz will return error when any one key isn't found.
func (c *CentralCacheImpl) DeleteBatch(keys []string) error {
	data := make(map[string]interface{}, 1)
	var keyArr []map[string]string
	for _, x := range keys {
		key := map[string]string{"key": x}
		keyArr = append(keyArr, key)
	}
	data["data"] = keyArr
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	url := c.getServerUrl("bulkDelete", false)
	res, err := http.HttpDelete(url, nil, string(jsonData), c.timeOut)
	if err != nil {
		return err
	}
	if res.HttpStatus != h.StatusOK {
		errResp, err := getCCErrorResponse(res.Body)
		if err != nil {
			return err
		}
		logger.Error(errResp.Message)
		return getCacheErrorType(CentralCache, errResp.Code)
	}
	return nil
}

//getServerUrl returns the central cache server url from the supplied key
func (c *CentralCacheImpl) getServerUrl(key string, isSingleGetReq bool) string {
	if key == "bulkDelete" {
		return fmt.Sprintf("%s/%s/%s", c.baseUrl, c.keyPrefix, "bulk")
	}
	if isSingleGetReq {
		return fmt.Sprintf("%s/%s/entity/%s", c.baseUrl, c.keyPrefix, key)
	}
	return fmt.Sprintf("%s/%s/entities/%s", c.baseUrl, c.keyPrefix, key)
}

func (c *CentralCacheImpl) Dump(key string, value []byte) error {
	return nil
}

//newCentralCache returns a new instance of CentralCacheImpl from the supplied conf
func newCentralCache(conf Config) (*CentralCacheImpl, error) {
	c := new(CentralCacheImpl)
	c.Init(conf)
	return c, nil
}
