package cachestrategy

import (
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/threadpool"
)

type Config struct {
	Strategy      string
	DBAdapterType string
	Cache         cache.Config
	ThreadPool    threadpool.Config
}
