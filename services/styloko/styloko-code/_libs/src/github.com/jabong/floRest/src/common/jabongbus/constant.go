package jabongbus

import ()

// subscriber constants
const (
	EXCEPTION              = "-exception"
	SUBSCRIBER             = "-sub"
	SUBSCRIBER_COUNT       = "-subcount"
	SUBSTATUS              = "-substatus"
	DEADQ_SUBSCRIBER       = "-deadq_sub"
	DEADQ_SUBSCRIBER_COUNT = "-deadq_subcount"
	DEADQ_SUBSTATUS        = "-deadq_substatus"
	MSG_NOT_ACKED          = "last message is not acked"
	SUB_STATUS_DURATION    = 300
	MAX_ACTIVE_REDIS_CON   = 2
	DEFAULT_RETRY_COUNT    = 0
)
