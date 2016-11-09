package cachestrategy

import (
	"github.com/jabong/floRest/src/common/cache"
)

type CacheStrategyInterface interface {
	// Init initialises the cache and db adapter
	Init(conf Config)

	// Get gets an item from a cache store indexed with key. serialize and compress indicates if the cache implementation
	// has to undergo some serialization or compression before returning the item
	Get(key string, serialize bool, compress bool) (item *cache.Item, err error)

	// Set sets an item in both cache store and DB. serialise and compress indicates if the cache implementation
	// has to undergo some serialization or compression before setting the item in cache
	Set(item cache.Item, serialize bool, compress bool) error

	// Delete deletes a Key from both cache and DB
	Delete(key string) error

	// GetBatch gets a list of all items indexed with keys. serialize and compress indicates if the
	// cache implementation has to undergo some serialization or compression before returning the items
	GetBatch(keys []string, serialize bool, compress bool) (items map[string]*cache.Item, err error)
}
