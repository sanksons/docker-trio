package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// CentralCacheTestImpl denotes Central Cache for Test cases.
type CentralCacheTestImpl struct {
	baseUrl      string
	keyPrefix    string
	dumpFilePath string
	expirySec    int32
	file         []byte
}

//Init initialiases c to connect to server with all keys stored under the bucket keyPrefix
func (c *CentralCacheTestImpl) Init(conf Config) {
	c.baseUrl = conf.Host
	c.keyPrefix = conf.KeyPrefix
	c.dumpFilePath = conf.DumpFilePath
	c.expirySec = conf.ExpirySec
	c.file = nil
}

func (c *CentralCacheTestImpl) InitializeFromFile() {
	pwd, osErr := os.Getwd()
	if osErr != nil {
		panic(fmt.Sprintf("Error in getting file path-  %s", osErr))
	}
	filePath := pwd + "/" + c.dumpFilePath
	fmt.Println(fmt.Sprintf("Central Cache File:  %+v", filePath))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading Central Test Cache file %s \n %s", filePath, err))
	}
	c.file = file
}

//Get gets key from central cache
func (c *CentralCacheTestImpl) Get(key string, serialize bool, compress bool) (*Item, error) {
	keyArray := make([]string, 1)
	keyArray[0] = key
	conf, err := c.GetBatch(keyArray, false, false)
	if err != nil {
		return nil, err
	}
	item := new(Item)
	item.Key = key
	itemValue, ok := conf[key].Value.(string)
	if !ok {
		return nil, errors.New("Error: Cannot convert to type string")
	}
	item.Value = itemValue
	return item, nil
}

//GetBatch gets the list of values stored in central cache indexed by keys. If a key does not exist in central
//cache then it has nil set for that key value in the returned list
func (c *CentralCacheTestImpl) GetBatch(keys []string, serialize bool, compress bool) (map[string]*Item, error) {
	resMap := make(map[string]*Item, len(keys))
	for _, v := range keys {
		resMap[v] = nil
	}

	body := c.file
	if body == nil {
		return nil, errors.New("Cache file response error")
	}

	cr := new(CentralCacheGetBatchResponse)
	err := json.Unmarshal(body, cr)
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

//Set stores an item in central cache for expirySec
func (c *CentralCacheTestImpl) Set(item Item, serialize bool, compress bool) error {
	return nil
}

// SetWithTimeout stores an item in central cache for provided TTL
func (c *CentralCacheTestImpl) SetWithTimeout(item Item, serialize bool, compress bool, ttl int32) error {
	return nil
}

// Delete deletes an item from central cache
func (c *CentralCacheTestImpl) Delete(key string) error {
	return nil
}

// DeleteBatch deletes an array of keys from central cache
func (c *CentralCacheTestImpl) DeleteBatch(key []string) error {
	return nil
}

// Dump dumps data to central cache
func (c *CentralCacheTestImpl) Dump(key string, value []byte) error {
	return nil
}

//newCentralCache returns a new instance of CentralCacheTestImpl from the supplied conf
func newCentralCacheTest(conf Config) (*CentralCacheTestImpl, error) {
	c := new(CentralCacheTestImpl)
	c.Init(conf)
	c.InitializeFromFile()
	return c, nil
}
