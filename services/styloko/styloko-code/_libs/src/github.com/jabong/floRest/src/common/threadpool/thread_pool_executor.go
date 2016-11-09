package threadpool

import (
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
)

type ThreadPoolType chan chan Task

type ThreadPoolExecutor struct {
	nThreads      int
	taskQueueSize int
	taskQueue     chan Task
	threadPool    ThreadPoolType
}

// Creates a new thread pool executor and initializes it.
func NewThreadPoolExecutor(conf Config) (threadPoolExecutor *ThreadPoolExecutor, errObj error) {
	defer func() {
		// do recovery if error
		if r := recover(); r != nil {
			errObj = errors.New(fmt.Sprintf("Failed to create threadpool executor, Error:%s", r))
			return
		}
	}()

	threadPoolExecutor = new(ThreadPoolExecutor)
	threadPoolExecutor.nThreads = conf.NThreads
	threadPoolExecutor.taskQueueSize = conf.TaskQueueSize

	// First, initialize the channel we are going to put the thread's thread channels into.
	threadPoolExecutor.threadPool = make(ThreadPoolType, conf.NThreads)

	threadPoolExecutor.taskQueue = make(chan Task, conf.TaskQueueSize)

	// Now, create all the threads.
	for i := 0; i < conf.NThreads; i++ {
		logger.Debug("Starting thread "+strconv.Itoa(i+1), false)
		thread := NewThread(i+1, threadPoolExecutor.threadPool)
		thread.Start()
	}

	go func() {
		defer func() {
			// do recovery if error
			if r := recover(); r != nil {
				logger.Error(fmt.Sprintf("Error in starting the thread pool, Error:%s", r))
				return
			}
		}()
		for {
			select {
			case task := <-threadPoolExecutor.taskQueue:
				logger.Debug("Received task requeust", false)
				threadChan := <-threadPoolExecutor.threadPool
				logger.Debug("Dispatching task to the thread", false)
				threadChan <- task
			}
		}
	}()

	return threadPoolExecutor, nil
}

// Dispatches the task to an available thread to execute it
func (threadPoolExecutor *ThreadPoolExecutor) ExecuteTask(t Task) {
	if threadPoolExecutor.taskQueue == nil {
		panic(fmt.Sprintf("ThreadPoolExecutor is not initalized. Use threadPool.NewThreadPoolExecutor method to create & initialize threadPoolExecutor"))
	}
	threadPoolExecutor.taskQueue <- t
}
