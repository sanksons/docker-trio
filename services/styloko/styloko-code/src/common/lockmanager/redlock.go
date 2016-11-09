package lockmanager

import (
	"common/redis"
	"common/utils"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	UNLOCK_SCRIPT = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        `
)

type RedLock struct {
	Retrials     int
	RetrialDelay time.Duration
	ExpireTTL    time.Duration
	Resource     string
	Value        string
}

func (self *RedLock) SetResource(resource string) {
	self.Resource = resource
}

func (self *RedLock) Lock() error {
	var trials = self.Retrials
	var gErr error
	for trials > 0 {
		err := self.lock()
		if err == nil {
			break
		}
		gErr = err
		trials--
		time.Sleep(self.RetrialDelay)
		// randomize time wait
		rtime := time.Duration(utils.Random(1, 15))
		time.Sleep(rtime * time.Millisecond)
	}
	if trials == 0 {
		return gErr
	}
	return nil
}

func (self *RedLock) UnLock() error {
	var trials = self.Retrials
	var gErr error
	for trials > 0 {
		shouldRetry, err := self.unlock()
		if err == nil {
			break
		}
		if err != nil && !shouldRetry {
			trials = 0
			break
		}
		gErr = err
		trials--
		time.Sleep(self.RetrialDelay)
		rtime := time.Duration(utils.Random(1, 15))
		time.Sleep(rtime * time.Millisecond)
	}
	if trials == 0 {
		return gErr
	}
	return nil
}

func getRandStr() string {
	b := make([]byte, 16)
	crand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (self *RedLock) lock() error {
	driver, err := redis.GetDriver()
	if err != nil {
		return fmt.Errorf(
			"(self *RedLock)#lock():1:Cannot Get Redis Driver: %s", err.Error(),
		)
	}
	self.Value = getRandStr()
	locked, rerr := driver.SetNX(self.Resource, self.Value, self.ExpireTTL)
	if rerr != nil {
		return fmt.Errorf(
			"(self *RedLock)#lock():2: %s", rerr.Error(),
		)
	}
	if !locked {
		return fmt.Errorf(
			"(self *RedLock)#lock():3: Lock Busy",
		)
	}
	return nil
}

// Release lock
//
// Returns:
// bool -> should retry?
// error -> error found in releasing lock
//
func (self *RedLock) unlock() (bool, error) {
	driver, err := redis.GetDriver()
	if err != nil {
		return true, fmt.Errorf(
			"(self *RedLock)#lock():1:Cannot Get Redis Driver: %s", err.Error(),
		)
	}
	val, err := driver.EVAL(UNLOCK_SCRIPT, []string{self.Resource}, []string{self.Value})
	if err != nil {
		return true, fmt.Errorf("(self *RedLock)#unlock():%s", err.Error())
	}
	if utils.ToString(val) == "0" {
		return false, fmt.Errorf("(self *RedLock)#unlock():Key value didnt matched")
	}
	return false, nil
}
