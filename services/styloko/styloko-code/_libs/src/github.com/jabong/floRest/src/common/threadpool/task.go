package threadpool

import ()

type Task struct {
	Instance   interface{}
	MethodName string
	Args       []interface{}
}
