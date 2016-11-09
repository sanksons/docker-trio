package cachestrategy

import (
	"github.com/jabong/floRest/src/common/cache"
)

type DBAdapterInterface interface {
	// initialize the DB Adapter by creating connection/session to DB
	Init(conf Config)

	// Gets an item from DB indexed with key
	ExecuteRead(key string) (item *cache.Item, err error)

	// Gets a list of all items from DB indexed with keys
	ExecuteReadBulk(keys []string) (items map[string]*cache.Item, err error)

	// Sets the given item in DB
	ExecuteWrite(item cache.Item) error

	// Deletes the item from DB indexed with key
	ExecuteDelete(key string) error
}
