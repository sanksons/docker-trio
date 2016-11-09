package cache

import (
	"github.com/jabong/floRest/src/common/logger"
	"testing"
)

//startTest does initialisation that is required before each test
func startTest() {
	//TO DO:- Have seperate config for the test
	logger.Initialise("../../../config/logger/logger_dev.json")
}

//endTest cleans up all resources that are allocated or initialised
//during startTest
func endTest() {
	//memcacheImpl = nil
}

func Test_Get(t *testing.T) {
	startTest()

	c := getTestMemCacheConfig()
	m := new(MemcacheImpl)
	m.Init(c.Host, c.KeyPrefix, c.DumpFilePath)
	//m.keyPrefix = c.KeyPrefix
	ci, err := Get(c)
	if err != nil {
		t.Errorf("Test cache_manager Get failed %v", err)
	}
	tm, ok := ci.(*MemcacheImpl)
	if !ok {
		t.Errorf("Test cache_manager Get failed. Not returned memcacheimpl type", tm)
	}
	if m.keyPrefix != tm.keyPrefix /*|| reflect.DeepEqual(m.client, tm.client) == false*/ {
		t.Errorf("Test cache_manager Get failed. Returned different cache imp %v", tm)
	}

	c.Platform = "123232"
	ci, err = Get(c)
	if err == nil {
		t.Error("Test cache_manager Get failed. Unknown cache dao returned")
	}

	endTest()
}

func Test_newMemcache(t *testing.T) {
	startTest()

	c := getTestMemCacheConfig()
	m := new(MemcacheImpl)
	m.Init(c.Host, c.KeyPrefix, c.DumpFilePath)

	tm, err := newMemcache(getTestMemCacheConfig())
	if err != nil || m.keyPrefix != tm.keyPrefix /*|| reflect.DeepEqual(m, tm) == false*/ {
		t.Error("Test newMemcache failed")
	}
	endTest()
}

//TODO move this to a test config
func getTestMemCacheConfig() Config {
	c := Config{}
	c.Platform = "memcache"
	c.KeyPrefix = "core"
	c.Host = "localhost:11211"
	return c
}
