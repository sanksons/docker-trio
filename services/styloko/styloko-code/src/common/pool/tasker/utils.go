package tasker

import (
	factory "common/ResourceFactory"
	"common/constants"
	"common/notification"
	"common/notification/datadog"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

const (
	MYSQL_TASK_TABLE = "styloko_mysql_syncing"
	MAX_TRIAL_LIMIT  = 3
	MAX_ADD_LIMIT    = 3

	MYSQL_POOL_NAME = "Worker"

	//Task status
	STATUS_INIT       = "init"
	STATUS_PROCESSING = "processing"
	STATUS_FAILED     = "failed"
	STATUS_SUCCESS    = "success"
	STATUS_CANCELED   = "canceled"

	//flags
	PROCESSING_TASKS                = "{styloko}_tasker_processing"
	MAINJOB_FLAG_PREFIX             = "{styloko}_tasker_main_"
	BUSY                            = "1"
	AVAILABLE                       = "0"
	HANGEDJOB_LAST_RUN_FLAG_PREFIX  = "{styloko}_tasker_hanged_"
	HANGEDJOB_INTERVAL              = 10 //minutes
	CLEANUPJOB_LAST_RUN_FLAG_PREFIX = "{styloko}_tasker_clean_"
	CLEANUPJOB_INTERVAL             = 5 //hour
)

//Mysql Transaction
type SqlTxn struct {
	Tx *sql.Tx
}

//
// Cancel All non processed tasks for a Resource.
//
func CancelTasks(resourceType string, resourceId int) error {
	var sql string
	sql = `UPDATE ` + MYSQL_TASK_TABLE + `
	         SET ` + "`status`" + ` ='` + STATUS_CANCELED + `'
             WHERE
               resource_type='` + resourceType + `'
             AND (
               (status = '` + STATUS_PROCESSING + `')
                 OR
               (status = '` + STATUS_INIT + `')
                 OR
               (status = '` + STATUS_FAILED + `' AND trials < ` + strconv.Itoa(MAX_TRIAL_LIMIT) + `)
             )
             AND
               resource=?`
	driver, err := factory.GetMySqlDriver(MYSQL_POOL_NAME)
	if err != nil {
		logger.Error(fmt.Errorf("CancelTasks(%s, %d) failed:", resourceType, resourceId, err.Error()))
		return err
	}
	_, sqlErr := driver.Execute(sql, resourceId)
	if sqlErr != nil {
		logger.Error(fmt.Errorf("CancelTasks(%s, %d) failed:", resourceType, resourceId, sqlErr.DeveloperMessage))
		return err
	}
	return nil
}

//
// Add sync task to mysql, with re-trials
//
func AddProductSyncJob(resourceId int, typ string, data interface{}) {
	go func(resourceId int, typ string, data interface{}) {
		//recover here
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error(fmt.Sprintf("[PANIC] occured with ADD TASK"))
			}
		}()
		byteData, _ := json.Marshal(data)
		task := Task{
			ResourceType: constants.PRODUCT_RESOURCE_NAME,
			Type:         typ,
			Resource:     resourceId,
			Data:         byteData,
		}
		var trials int
		for trials < MAX_ADD_LIMIT {
			err := task.Add()
			if err != nil {
				trials += 1
				time.Sleep(2 * time.Microsecond)
				continue
			}
			break
		}
		if trials == MAX_ADD_LIMIT {
			logger.Error(fmt.Errorf("Task Add Failed: %d-%s", resourceId, typ))
		}
	}(resourceId, typ, data)
}

//
// Handles Recovery when a go routine, executing a task syncing fails
//
func TaskerRecoverHandler(task Task) {
	if rec := recover(); rec != nil {
		task.SetStatus(STATUS_FAILED, fmt.Errorf("[PANIC] occured with TASK:%d", task.Id))
		logger.Error(fmt.Sprintf("[PANIC] occured with TASK:%d", task.Id))
		trace := make([]byte, 4096)
		count := runtime.Stack(trace, true)
		logger.Error(fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace))
		logger.Error(fmt.Sprintf("Reason for panic: %s", rec))

		//notify
		trace = make([]byte, 1024)
		count = runtime.Stack(trace, true)
		stackTrace := fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace)
		title := fmt.Sprintf("Panic occured, Task failed[%d]", task.Id)
		text := fmt.Sprintf("Panic reason %s\n\nStack Trace: %s", rec, stackTrace)
		tags := []string{"tasker-error", "panic"}
		notification.SendNotification(title, text, tags, datadog.ERROR)
	}
}

//
// General recovery handler for tasker
//
func RecoverHandler(data string) {
	if rec := recover(); rec != nil {
		logger.Error(fmt.Sprintf("[PANIC] occured:%s", data))
		trace := make([]byte, 4096)
		count := runtime.Stack(trace, true)
		logger.Error(fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace))
		logger.Error(fmt.Sprintf("Reason for panic: %s", rec))

		//notify
		trace = make([]byte, 1024)
		count = runtime.Stack(trace, true)
		stackTrace := fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace)
		title := fmt.Sprintf("Panic occured")
		text := fmt.Sprintf("Panic reason %s\n\nStack Trace: %s", rec, stackTrace)
		tags := []string{"tasker-error", "panic"}
		notification.SendNotification(title, text, tags, datadog.ERROR)
	}
}
