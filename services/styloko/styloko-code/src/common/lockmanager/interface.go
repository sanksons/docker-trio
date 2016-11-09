package lockmanager

import (
	"time"
)

const LOCK_TYPE_REDLOCK = "redlock"

type MutexConfig struct {
	Retrials     int
	RetrialDelay time.Duration
	ExpireTTL    time.Duration
}

type MutexInterface interface {

	//Define Resource to lock
	SetResource(string)

	//Acquire Lock
	Lock() error

	//Release Lock
	UnLock() error
}

func GetMutex(typ string, config MutexConfig) MutexInterface {
	if typ == LOCK_TYPE_REDLOCK {
		return &RedLock{
			Retrials:     config.Retrials,
			RetrialDelay: config.RetrialDelay,
			ExpireTTL:    config.ExpireTTL,
		}
	}
	//Since we have only one implementation at the moment.
	return &RedLock{
		Retrials:     config.Retrials,
		RetrialDelay: config.RetrialDelay,
		ExpireTTL:    config.ExpireTTL,
	}
}
