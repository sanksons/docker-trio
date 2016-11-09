package cachestrategy

import (
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/cachestrategy"
)

type SampleDBAdapter struct {
}

func (adapter SampleDBAdapter) Init(conf cachestrategy.Config) {
	fmt.Println("Sample DB Adapter initialized")
}

func (adapter SampleDBAdapter) ExecuteRead(key string) (item *cache.Item, err error) {
	fmt.Println("Sample DB Adapter - In Execute Read method - Key : " + key)
	item = new(cache.Item)
	item.Key = "key"
	item.Value = "value"
	return item, nil
}

func (adapter SampleDBAdapter) ExecuteReadBulk(keys []string) (items map[string]*cache.Item, err error) {
	fmt.Println("Sample DB Adapter - In Execute Read Bulk method - Keys : %v", keys)
	item := new(cache.Item)
	item.Key = "key"
	item.Value = "value"
	items["key"] = item
	return items, nil
}

func (adapter SampleDBAdapter) ExecuteWrite(item cache.Item) error {
	fmt.Println("Sample DB Adapter - In Execute write method")
	return nil
}

func (adapter SampleDBAdapter) ExecuteDelete(key string) error {
	fmt.Println("Sample DB Adapter - In Execute delete method - Key : " + key)
	return nil
}
