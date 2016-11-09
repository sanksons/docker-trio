package cachestrategy

import (
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/threadpool"
)

type AbstractCacheStrategy struct {
	cacheImpl          cache.CacheInterface
	dbAdapterImpl      DBAdapterInterface
	threadPoolExecutor *threadpool.ThreadPoolExecutor
}

// Initializes both DB adapter impl and cache impl using the given config
func (cacheStrategy *AbstractCacheStrategy) Init(conf Config) {
	var err error
	cacheStrategy.cacheImpl, err = cache.Get(conf.Cache)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize cache \n %s \n %s", conf.Cache, err))
	}

	cacheStrategy.dbAdapterImpl, err = DBAdapterMgr.GetDBAdapter(conf.DBAdapterType)
	if err != nil {
		panic(fmt.Sprintf("Failed to get the db adapter \n %s \n %s", conf.DBAdapterType, err))
	}
	if cacheStrategy.dbAdapterImpl == nil {
		panic(fmt.Sprintf("DB Adapter must be registered with DBAdapterManager before initialising cache strategy \n %s \n %s", conf.DBAdapterType, err))
	}

	cacheStrategy.threadPoolExecutor, err = threadpool.NewThreadPoolExecutor(conf.ThreadPool)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize thread pool \n %s \n %s", conf.ThreadPool, err))
	}

	fmt.Println("Cache strategy is initialized...")
}

// Converts the interface value to string
func (cacheStrategy *AbstractCacheStrategy) convertValueToString(value interface{}) string {
	val, ok := value.([]byte)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	return string(val)
}
