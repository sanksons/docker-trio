package pool

import (
	"common/notification"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

/*
*
* Pool package must be initialized in the Init function
* The StartWorkers method required an anonymous function as
* an argument which in turn takes an interface{} as its argument.
* Return for the job must be a bool confirming fail or success.
*
* Implement a function with signature :
* func work(args interface{})(bool)
*
 */

// Safe -> struct for retry based pool
// retryCount: is the number of time job will be retried if Failed
// sleepTime: is the time for which job will sleep before retrying
type Safe struct {
	name       string
	poolSize   int
	channel    chan interface{}
	retryCount int
	sleepTime  int64
	job        func(interface{}) error
	queueMax   int
}

// NewWorkerSafe -> Creates a new worker and returns a Pool object. This implements retry for failed jobs
func NewWorkerSafe(name string, poolSize int, queueMax int, retryCount int, sleepTime int64) Safe {
	pool := Safe{
		name:       name,
		poolSize:   poolSize,
		channel:    make(chan interface{}, queueMax),
		retryCount: retryCount,
		sleepTime:  sleepTime,
		queueMax:   queueMax,
	}
	logger.Info(fmt.Sprintf("Safe thread pool started for %s, pool size %d, queue size %d", name, poolSize, queueMax))
	logger.Info(fmt.Sprintf("Retries: %d, Sleep time: %d", retryCount, sleepTime))
	return pool
}

// StartWorkers -> Starts workers based on pool size
func (p Safe) StartWorkers(job func(interface{}) error) {
	p.job = job
	for w := 0; w <= p.poolSize; w++ {
		go p.worker(job)
	}
	logger.Info(fmt.Sprintf("%d workers started in waiting state.", p.poolSize))
}

// StartJob -> Use this function to start a new job, send data in interface
func (p Safe) StartJob(jobData interface{}) {
	p.channel <- jobData
}

// worker -> runs the given job with the interface in channel as argument
func (p Safe) worker(job func(interface{}) error) {
	defer recoverHandler(p.name)
	for j := range p.channel {
		logger.Info(fmt.Sprintf("Job starting for %s", p.name))
		retryCount := p.retryCount
		err := func(name string) error {
			defer recoverHandler(name)
			er := job(j)
			return er
		}(p.name)
		for i := 0; i < retryCount; i++ {
			if err != nil {
				logger.Error(fmt.Sprintf("Failure occured for job in API: %s", p.name))
				logger.Error(err.Error())
				logger.Info(fmt.Sprintf("Sleeping for %d milliseconds", p.sleepTime))
				time.Sleep(time.Duration(p.sleepTime) * time.Millisecond)
				err = func(name string) error {
					defer recoverHandler(name)
					er := job(j)
					return er
				}(p.name)
			} else {
				break
			}
		}
		if err != nil {
			logger.Error("Failure occured. All retries failed.")
			title := fmt.Sprintf("Worker failure in %s. All retries (%d) have failed.", p.name, p.retryCount)
			text := fmt.Sprintf("Failure reson for worker %s: %s", p.name, err.Error())
			tags := []string{"pool-error", "worker-failure", p.name}
			notification.SendNotification(title, text, tags, "error")
			continue
		}
		logger.Info("Job successful.")
	}
}

// Close -> closes the pool channel. Will finish pending jobs before closing.
func (p Safe) Close() {
	close(p.channel)
	logger.Info("Channel has been closed")
}

// restartWorkers -> re generates the channel and restarts workers
func (p Safe) restartWorkers() {
	p.channel = make(chan interface{}, p.queueMax)
	p.StartWorkers(p.job)
}
