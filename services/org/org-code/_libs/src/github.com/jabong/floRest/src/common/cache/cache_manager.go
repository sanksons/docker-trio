package cache

import (
	"errors"

	"github.com/jabong/floRest/src/common/utils/logger"
)

func Get(conf Config) (CacheInterface, error) {
	switch conf.Platform {
	case CentralCache:
		centralCacheImpl, err := newCentralCache(conf)
		if err != nil {
			logger.Info("Error in Getting Centralcache Impl " + err.Error())
			return nil, err
		}
		logger.Info("Centralcache Dao Initialised")
		return centralCacheImpl, nil
	case CentralConfigCache:
		centralConfigCacheImpl, err := newCentralConfigCache(conf)
		if err != nil {
			logger.Info("Error in Getting CentralConfigcache Impl " + err.Error())
			return nil, err
		}
		logger.Info("CentralConfigcache Dao Initialised")
		return centralConfigCacheImpl, nil
	case CentralCacheTest:
		centralCacheTestImpl, err := newCentralCacheTest(conf)
		if err != nil {
			logger.Info("Error in Getting Test Centralcache Impl " + err.Error())
			return nil, err
		}
		logger.Info("CentralcacheTest Dao Initialised")
		return centralCacheTestImpl, nil
	}
	return nil, errors.New("Unknown Cache Dao Type" + conf.Platform)
}
