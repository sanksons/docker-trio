package cachestrategy

import (
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/threadpool"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type DBFirstCacheStrategy struct {
	AbstractCacheStrategy
}

func (cacheStrategy *DBFirstCacheStrategy) Get(key string, serialize bool, compress bool) (item *cache.Item, err error) {
	item, err = cacheStrategy.cacheImpl.Get(key, serialize, compress)
	if item == nil || err != nil {
		logger.Debug("Could not get the item from the cache. Hence reading from DB and setting it in cache")
		item, err = cacheStrategy.dbAdapterImpl.ExecuteRead(key)
		if err != nil {
			logger.Error("Error fetching data from DB, Key : " + key + ", Error : " + err.Error())
			return nil, err
		} else if item == nil {
			logger.Error("Data missing in DB as well, Key : " + key)
			return nil, nil
		} else {
			task := threadpool.Task{cacheStrategy, "WriteToCacheAsynchronously", []interface{}{*item, serialize, compress}}
			cacheStrategy.threadPoolExecutor.ExecuteTask(task)
		}
	}
	return item, nil
}

func (cacheStrategy *DBFirstCacheStrategy) Set(item cache.Item, serialize bool, compress bool) error {
	err := cacheStrategy.dbAdapterImpl.ExecuteWrite(item)
	if err != nil {
		logger.Error("Error writing data to DB, Error : " + err.Error() + ", Key : " + item.Key + ", Value : " + cacheStrategy.convertValueToString(item.Value))
		return err
	} else {
		task := threadpool.Task{cacheStrategy, "WriteToCacheAsynchronously", []interface{}{item, serialize, compress}}
		cacheStrategy.threadPoolExecutor.ExecuteTask(task)
	}
	return nil
}

func (cacheStrategy *DBFirstCacheStrategy) Delete(key string) error {
	err := cacheStrategy.dbAdapterImpl.ExecuteDelete(key)
	if err != nil {
		logger.Error("Error deleting data from DB, Error : " + err.Error() + ", Key : " + key)
		return err
	} else {
		task := threadpool.Task{cacheStrategy, "DeleteFromCacheAsynchronously", []interface{}{key}}
		cacheStrategy.threadPoolExecutor.ExecuteTask(task)
	}
	return nil
}

func (cacheStrategy *DBFirstCacheStrategy) GetBatch(keys []string, serialize bool, compress bool) (items map[string]*cache.Item, err error) {
	items, err = cacheStrategy.cacheImpl.GetBatch(keys, serialize, compress)
	if items == nil || err != nil {
		logger.Debug("Could not get the items from the cache. Hence reading from DB and setting it in cache")
		items, err = cacheStrategy.dbAdapterImpl.ExecuteReadBulk(keys)
		if err != nil {
			logger.Error("Error fetching bulk data from DB, Keys : %v , Error : "+err.Error(), keys)
			return nil, err
		} else if items == nil {
			logger.Error("Data missing in DB as well, Keys : %v", keys)
			return nil, nil
		} else {
			task := threadpool.Task{cacheStrategy, "WriteBulkToCacheAsynchronously", []interface{}{keys, items, serialize, compress}}
			cacheStrategy.threadPoolExecutor.ExecuteTask(task)
		}
	}
	return items, nil
}

func (cacheStrategy *DBFirstCacheStrategy) WriteBulkToCacheAsynchronously(keys []string, items map[string]*cache.Item, serialize bool, compress bool) {
	for _, key := range keys {
		err := cacheStrategy.cacheImpl.Set(*items[key], serialize, compress)
		if err != nil {
			logger.Error("Error setting data to cache, Error : " + err.Error() + ", Key : " + items[key].Key + ", Value : " + cacheStrategy.convertValueToString(items[key].Value))
		}
	}
}

func (cacheStrategy *DBFirstCacheStrategy) WriteToCacheAsynchronously(item cache.Item, serialize bool, compress bool) {
	err := cacheStrategy.cacheImpl.Set(item, serialize, compress)
	if err != nil {
		logger.Error("Error setting data to cache, Error : " + err.Error() + ", Key : " + item.Key + ", Value : " + cacheStrategy.convertValueToString(item.Value))
	}
}

func (cacheStrategy *DBFirstCacheStrategy) DeleteFromCacheAsynchronously(key string) {
	err := cacheStrategy.cacheImpl.Delete(key)
	if err != nil {
		logger.Error("Error deleting data from cache, Error : " + err.Error() + ", Key : " + key)
	}
}
