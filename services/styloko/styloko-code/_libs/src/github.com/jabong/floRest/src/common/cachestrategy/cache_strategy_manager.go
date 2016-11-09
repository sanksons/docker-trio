package cachestrategy

import (
	"errors"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// Creates, initializes and returns the cache strategy instance based on the given config
func Get(conf Config) (CacheStrategyInterface, error) {
	var cacheStrategy CacheStrategyInterface
	if conf.Strategy == CacheFirstStrategy {
		cacheStrategy = new(CacheFirstCacheStrategy)
	} else if conf.Strategy == DBFirstStrategy {
		cacheStrategy = new(DBFirstCacheStrategy)
	} else {
		errMsg := "Cache strategy - " + conf.Strategy + " is not supported"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	cacheStrategy.Init(conf)
	return cacheStrategy, nil
}
