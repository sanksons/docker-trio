package threadpool

import (
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"reflect"
)

// NewThread creates, and returns a new thread instance. Its only argument
// is a channel that the thread can add itself to whenever it is done its
// task.
func NewThread(id int, threadPool chan chan Task) Thread {
	// Creates and return the thread.
	thread := Thread{
		ID:         id,
		ThreadChan: make(chan Task),
		ThreadPool: threadPool,
		QuitChan:   make(chan bool)}

	return thread
}

type Thread struct {
	ID         int
	ThreadChan chan Task
	ThreadPool chan chan Task
	QuitChan   chan bool
}

// This function "starts" the thread by starting a goroutine, that is
// an infinite "for-select" loop.
func (t Thread) Start() {
	go func() {
		defer func() {
			// do recovery if error
			if r := recover(); r != nil {
				logger.Error(fmt.Sprintf("Error in starting the thread, Error:%s", r))
				return
			}
		}()
		for {
			// Add its own task channel into the thread pool.
			t.ThreadPool <- t.ThreadChan

			select {
			case task := <-t.ThreadChan:
				t.execute(task)
			case <-t.QuitChan:
				// We have been asked to stop.
				logger.Debug(fmt.Sprintf("thread%d stopping\n", t.ID), false)
				return
			}
		}
	}()
}

// Stop tells the thread to stop listening for task requests.
// Note that the thread will only stop *after* it has finished its current task.
func (t Thread) Stop() {
	t.QuitChan <- true
}

//execute executes the task to be done by the thread
func (t Thread) execute(task Task) {
	t.callMethod(task.Instance, task.MethodName, task.Args)
}

// Using reflection, it calls the method on the given instance with the given arguments
func (t Thread) callMethod(instance interface{}, methodName string, args []interface{}) {
	defer func() {
		// do recovery if error
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Error in calling the method-%s, Error:%s", methodName, r))
			return
		}
	}()

	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(instance)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(instance))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	if finalMethod.IsValid() {
		finalMethod.Call(inputs)
	} else {
		logger.Error("This method - " + methodName + " is not valid. Hence, Ignoring.")
	}
}
