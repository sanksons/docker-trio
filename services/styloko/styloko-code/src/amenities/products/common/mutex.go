package common

import (
	"common/lockmanager"
	"common/notification"
	"common/notification/datadog"
	"fmt"
	_ "github.com/jabong/floRest/src/common/utils/logger"
	"time"
)

const (
	LOCK_TYPE_CONFIGID   = "config"
	LOCK_TYPE_PRODUCTSET = "set"
	LOCK_PREFIX          = "mutex"

	LOCK_RETRIAL_LIMIT = 10
	LOCK_RETRIAL_DELAY = 150 * time.Millisecond
	LOCK_EXPIRE_TIME   = 2 * time.Minute
)

type ProductMutex struct {
	Id    int
	Type  string
	Mutex lockmanager.MutexInterface
}

func (mutex *ProductMutex) GetKeyName() string {
	var keyName string
	if mutex.Type == LOCK_TYPE_PRODUCTSET {
		keyName = fmt.Sprintf("styloko_%s_%s_%d", LOCK_PREFIX, LOCK_TYPE_PRODUCTSET, mutex.Id)
	} else {
		keyName = fmt.Sprintf("styloko_%s_%s_%d", LOCK_PREFIX, LOCK_TYPE_CONFIGID, mutex.Id)
	}
	return keyName
}

//
// Try to get a lock
//
func (self *ProductMutex) Lock() bool {
	config := lockmanager.MutexConfig{
		Retrials:     LOCK_RETRIAL_LIMIT,
		RetrialDelay: LOCK_RETRIAL_DELAY,
		ExpireTTL:    LOCK_EXPIRE_TIME,
	}
	self.Mutex = lockmanager.GetMutex(lockmanager.LOCK_TYPE_REDLOCK, config)
	self.Mutex.SetResource(self.GetKeyName())
	err := self.Mutex.Lock()
	if err != nil {
		return false
	}
	return true
}

//
// Release lock
//
func (self *ProductMutex) UnLock() bool {
	err := self.Mutex.UnLock()
	if err != nil {
		notification.SendNotification(
			"ProductMutex UnLock Failed",
			fmt.Sprintf("lockkey:%s", self.GetKeyName()),
			[]string{"ProductMutex"},
			datadog.ERROR,
		)
		return false
	}
	return true
}
