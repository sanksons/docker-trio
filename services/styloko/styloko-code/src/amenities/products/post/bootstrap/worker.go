package bootstrap

import (
	proUtil "amenities/products/common"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sync/atomic"
)

type Job struct {
	ProductId        int
	TransactionId    string
	PublishInvisible bool
}

func PublishProduct(job interface{}) error {
	processedCount.incrementCounter()
	myjob, ok := job.(Job)
	if !ok {
		return fmt.Errorf("PublishProduct(): Assertion failed")
	}
	pro, err := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GetById(myjob.ProductId)
	if err != nil {
		return fmt.Errorf("PublishProduct() [proId: %d]: %s", myjob.ProductId, err.Error())
	}
	(&pro).Publish(myjob.TransactionId, myjob.PublishInvisible)
	logger.Info(fmt.Sprintf("Product %d published.", myjob.ProductId))
	return nil
}

func (p *counters) incrementCounter() {
	atomic.AddInt64(&p.counter, 1)
}

func (p *counters) resetCounter() {
	atomic.StoreInt64(&p.counter, 0)
}

func (p *counters) getCounter() int {
	return int(atomic.LoadInt64(&p.counter))
}
