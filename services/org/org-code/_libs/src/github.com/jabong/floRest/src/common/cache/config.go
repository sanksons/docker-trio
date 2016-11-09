package cache

import (
	"time"
)

type Config struct {
	Platform     string
	Host         string
	KeyPrefix    string
	DumpFilePath string
	ExpirySec    int32
	Disabled     bool
	TimeOut      time.Duration
}

type Configs []Config
