package cachestrategy

import (
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/cachestrategy"
	"github.com/jabong/floRest/src/common/config"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type CacheStrategyUser struct {
	id string
}

func (n *CacheStrategyUser) SetID(id string) {
	n.id = id
}

func (n CacheStrategyUser) GetID() (id string, err error) {
	return n.id, nil
}

func (a CacheStrategyUser) Name() string {
	return "CacheStrategyUser"
}

func (a CacheStrategyUser) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	adapter := new(SampleDBAdapter)
	cachestrategy.DBAdapterMgr.RegisterDBAdapter("sample", adapter)
	cacheStrategy, err := cachestrategy.Get(config.GlobalAppConfig.CacheStrategyConfig)
	if err != nil {
		logger.Error("Error in getting the cache strategy --- " + err.Error())
		return io, nil
	}

	item := cache.Item{"somekey", "somevalue"}
	err1 := cacheStrategy.Set(item, false, false)
	if err1 != nil {
		logger.Error("Error in setting the item using the cache strategy --- " + err1.Error())
		return io, nil
	}
	item1, err2 := cacheStrategy.Get("somekey", false, false)

	if err2 != nil {
		logger.Error("Error in getting the item using the cache strategy --- " + err2.Error())
		return io, nil
	}

	fmt.Println("KEY : " + item1.Key)
	fmt.Println("VALUE : " + a.convertValueToString(item1.Value))
	return io, nil
}

func (a CacheStrategyUser) convertValueToString(value interface{}) string {
	val, ok := value.([]byte)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	return string(val)
}
