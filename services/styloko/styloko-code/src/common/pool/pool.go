package pool

import (
	"fmt"

	"github.com/jabong/floRest/src/common/utils/logger"
)

/*
*
* Pool package must be initialized in the Init function
* The StartWorkers method required an anonymous function as
* an argument which in turn takes an interface{} as its argument.
*
* Implement a function with signature :
* func work(args interface{})
*
 */

// Pool -> Basic Pool struct
type Pool struct {
	name     string
	poolSize int
	channel  chan interface{}
}

// NewWorker -> Creates a new worker and returns a Pool object
func NewWorker(name string, poolSize int, queueMax int) Pool {
	pool := Pool{
		name:     name,
		poolSize: poolSize,
		channel:  make(chan interface{}, queueMax),
	}
	logger.Info(fmt.Sprintf("Thread pool started for %s, pool size %d, queue size %d", name, poolSize, queueMax))
	return pool
}

// StartWorkers -> Starts workers based on pool size
func (p Pool) StartWorkers(job func(interface{})) {
	for w := 0; w <= p.poolSize; w++ {
		go p.worker(job)
	}
	logger.Info(fmt.Sprintf("%d workers started in waiting state.", p.poolSize))
}

// StartJob -> Use this function to start a new job, send data in interface
func (p Pool) StartJob(jobName interface{}) {
	p.channel <- jobName
}

// worker -> runs the given job with the interface in channel as argument
func (p Pool) worker(job func(interface{})) {
	for j := range p.channel {
		logger.Info(fmt.Sprintf("Job starting for %s", p.name))
		func(name string) {
			defer recoverHandler(name)
			job(j)
		}(p.name)
		logger.Info(fmt.Sprintf("Job finished for %s", p.name))
	}
}

// Close -> closes the pool channel. Will finish pending jobs before closing.
func (p Pool) Close() {
	close(p.channel)
}
