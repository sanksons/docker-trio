package jabongbus

import ()

// ProcessMessage interface to process message, should be implemented by client
type ProcessMessage interface {
	Process(*Message, error)
}

// PeristentSubscriber persistent subscriber interface
type PeristentSubscriber interface {
	init(*Subscriberconfig) error
	SetProcessMsg(ProcessMessage)
	Get(timeout int)
	Ack() error
	NoAck() error
	StopSub()
}

// NonPersistentSubscriber non persistent subscriber interface
type NonPersistentSubscriber interface {
	init(*Subscriberconfig)
	SetProcessMsg(ProcessMessage)
	Get(timeout int)
	StopSub()
}

// ExceptionSubscriber exception subscriber interface
type DeadQSubscriber interface {
	init(*Subscriberconfig) error
	SetProcessMsg(ProcessMessage)
	Get(timeout int)
	Ack() error
	StopSub()
}
