package simplifier

import (
	"time"
)

const (
	JUDGE                    = "judge"
	JUDGE_DAEMON             = "judge_jobs"
	JUDGE_DAEMON_POOL_SIZE   = 2
	JUDGE_DAEMON_QUEUE_SIZE  = 5
	JUDGE_DAEMON_RETRY_COUNT = 5
	JUDGE_DAEMON_WAIT_TIME   = 1000
	JUDGE_DAEMON_ERRORS      = "judge_jobs_errors"
	ERRORS                   = "ERRORS"
	CSV_ERRORS               = "CSV_ERRORS"
	EVERYJOB_RUNTIME         = time.Duration(2) * time.Minute
	CLEANUPJOB_RUNTIME       = time.Duration(168) * time.Hour
	RESETJOB_RUNTIME         = time.Duration(10) * time.Minute
)
