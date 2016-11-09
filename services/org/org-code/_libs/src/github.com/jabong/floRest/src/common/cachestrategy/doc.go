package cachestrategy

import ()

/** Two types of cache strategy
1. DB First Cache Strategy - Sets the data to DB first (use DB Adapter)and then sets it to cache using a background thread.
 							 Always gets the data from cache. So, if data is not available in cache,
							 It will be fetched from DB

2. Cache First Cache Strategy - Sets the data to Cache first and then sets it to DB (using db adapter) using a background thread.
								Always gets the data from cache, so, if data is not available in cache,
								it will return nil

DB adapter must be registered with DBAdapterManager and should be specified in the cache strategy config to use
it.

Check the usage in "cache_strategy_user.go" file

Sample Config for DB First Strategy :
	"CacheStrategyConfig":{
   		"Strategy": "dbFirstStrategy",
   		"DBAdapterType": "sample",
   		"Cache": {
   			"Platform": "centralCache",
   			"Host": "http://localhost:8080/cache/api/v1/buckets",
   			"KeyPrefix": "default"
   		},
   		"ThreadPool": {
   			"NThreads": 5,
   			"TaskQueueSize": 10
   		}
   	}

Sample Config for Cache First Strategy :
	"CacheStrategyConfig":{
   		"Strategy": "cacheFirstStrategy",
   		"DBAdapterType": "sample",
   		"Cache": {
   			"Platform": "centralCache",
   			"Host": "http://localhost:8080/cache/api/v1/buckets",
   			"KeyPrefix": "default"
   		},
   		"ThreadPool": {
   			"NThreads": 5,
   			"TaskQueueSize": 10
   		}
   	}
*/
