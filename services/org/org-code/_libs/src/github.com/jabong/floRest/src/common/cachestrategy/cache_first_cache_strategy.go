package cachestrategy

import (
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/threadpool"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type CacheFirstCacheStrategy struct {
	AbstractCacheStrategy
}

func (cacheStrategy *CacheFirstCacheStrategy) Get(key string, serialize bool, compress bool) (item *cache.Item, err error) {
	item, err = cacheStrategy.cacheImpl.Get(key, serialize, compress)
	if err != nil {
		logger.Error("Error fetching data from cache, Key : " + key + ", Error : " + err.Error())
		return nil, err
	} else if item == nil {
		logger.Error("Data missing in cache, Key : " + key)
		return nil, nil
	}
	return item, nil
}

func (cacheStrategy *CacheFirstCacheStrategy) Set(item cache.Item, serialize bool, compress bool) error {
	err := cacheStrategy.cacheImpl.Set(item, serialize, compress)
	if err != nil {
		logger.Error("Error setting data to cache, Error : " + err.Error() + ", Key : " + item.Key + ", Value : " + cacheStrategy.convertValueToString(item.Value))
		return err
	} else {
		task := threadpool.Task{cacheStrategy, "WriteToDBAsynchronously", []interface{}{item}}
		cacheStrategy.threadPoolExecutor.ExecuteTask(task)
	}
	return nil
}

func (cacheStrategy *CacheFirstCacheStrategy) Delete(key string) error {
	err := cacheStrategy.cacheImpl.Delete(key)
	if err != nil {
		logger.Error("Error deleting data from cache, Error : " + err.Error() + ", Key : " + key)
		return err
	} else {
		task := threadpool.Task{cacheStrategy, "DeleteFromDBAsynchronously", []interface{}{key}}
		cacheStrategy.threadPoolExecutor.ExecuteTask(task)
	}
	return nil
}

func (cacheStrategy *CacheFirstCacheStrategy) GetBatch(keys []string, serialize bool, compress bool) (items map[string]*cache.Item, err error) {
	items, err = cacheStrategy.cacheImpl.GetBatch(keys, serialize, compress)
	if err != nil {
		logger.Error("Error fetching bulk data from cache, Keys : %v, Error : "+err.Error(), keys)
		return nil, err
	} else if items == nil {
		logger.Error("All keys are missing in the cache, Keys : %v", keys)
		return nil, nil
	}
	return items, nil
}

func (cacheStrategy *CacheFirstCacheStrategy) WriteToDBAsynchronously(item cache.Item) {
	err := cacheStrategy.dbAdapterImpl.ExecuteWrite(item)
	if err != nil {
		logger.Error("Error writing data to DB, Error : " + err.Error() + ", Key : " + item.Key + ", Value : " + cacheStrategy.convertValueToString(item.Value))
	}
}

func (cacheStrategy *CacheFirstCacheStrategy) DeleteFromDBAsynchronously(key string) {
	err := cacheStrategy.dbAdapterImpl.ExecuteDelete(key)
	if err != nil {
		logger.Error("Error deleting data from DB, Error : " + err.Error() + ", Key : " + key)
	}
}
