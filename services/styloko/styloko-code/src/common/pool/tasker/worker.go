package tasker

import (
	factory "common/ResourceFactory"
	"common/notification"
	"common/notification/datadog"
	"common/redis"
	"common/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

//
// Inorder to start a new pool call Initialize() method
// of this struct.
//
type Tasker struct {
	UUID         string
	ResourceType string
	TaskQueue    []string
	SleepTime    time.Duration
	FetchLimit   int
}

//
// This method runs a infinite loop to process the tasks.
//
func (t Tasker) Initiate() {

	go func() {
		defer RecoverHandler("PRODUCT WORKER CRASHED")

		//release TASKER lock as soon as service boots.
		t.releaseLock(t.GetKeyName(MAINJOB_FLAG_PREFIX))

		//Wait for 30 seconds for all servers can boot.
		time.Sleep(time.Second * 30)

		//start infinite running loop
		for true {
			//wait for a while
			time.Sleep(time.Second * t.SleepTime)
			//Add a Random delay so that workers can be seperated
			time.Sleep(time.Microsecond * time.Duration(utils.Random(12, 90)))
			profiler := logger.NewProfiler()

			//acquire lock to do work
			if t.getLock(t.GetKeyName(MAINJOB_FLAG_PREFIX)) {
				logger.StartProfile(profiler, "product_mysql_sync_lockTime")
				//fmt.Printf("[%s] LOCKED  [%v]\n", t.UUID, time.Now().UnixNano())
				err := t.doWork()
				if err != nil {
					logger.Error(err.Error())
				}
				t.releaseLock(t.GetKeyName(MAINJOB_FLAG_PREFIX))
				logger.EndProfile(profiler, "product_mysql_sync_lockTime")
				//fmt.Printf("[%s] RELEASED[%v]\n", t.UUID, time.Now().UnixNano())
			}
		}
	}()

	//Restart hanged tasks
	go func() {
		defer RecoverHandler("Hanged Tasks WORKER CRASHED")
		//start infinite running loop
		for true {
			//wait for a while
			time.Sleep(time.Minute * time.Duration(utils.Random(5, 15)))
			if t.canRunJob(HANGEDJOB_LAST_RUN_FLAG_PREFIX, (HANGEDJOB_INTERVAL * time.Minute)) {
				err := t.restartHangedTasks()
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}()

	//CleanUp old jobs
	go func() {
		defer RecoverHandler("Cleanup Tasks WORKER CRASHED")
		//start infinite running loop
		for true {
			//wait for a while
			time.Sleep(time.Hour * 1)
			if t.canRunJob(CLEANUPJOB_LAST_RUN_FLAG_PREFIX, (CLEANUPJOB_INTERVAL * time.Hour)) {
				err := t.cleanupOldTasks()
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}()
}

//
// This method performs the work to fetch records
// from DB and send for processing
//
func (t Tasker) doWork() error {
	//Specify recover handler for this func
	defer RecoverHandler("TASKER")
	//Re-Initialize queue with currently processing tasks.
	t.reInitializeQueue()
	tx, err := t.startTransaction()
	if err != nil {
		return fmt.Errorf("(t Tasker)#doWork(): %s", err.Error())
	}

	//fetch tasks to be processed
	//this function has checks for filtering the tasks
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_sql")
	tasks := t.fetchTasks(tx)
	logger.EndProfile(profiler, "product_mysql_sync_sql")

	if tasks == nil || len(tasks) <= 0 {
		//no task found
		tx.Tx.Rollback()
		return err
	}
	var processingTasks []Task
	for _, task := range tasks {
		//check if we already have a task with same resource id.
		if !t.isInQueue(task.Resource) {
			//add resource id to queue.
			(&t).addToQueue(task.Resource)
			processingTasks = append(processingTasks, task)
		}
	}
	//mark tasks to be in processing
	err = t.markProcessingTasks(processingTasks, tx)
	if err != nil {
		tx.Tx.Rollback()
		return err
	}
	//Release the connection and lock
	tx.Tx.Commit()
	//Send tasks for processing
	for _, task := range processingTasks {
		//start a routine with recovery handler
		go func(task Task) {
			defer TaskerRecoverHandler(task)
			//fmt.Println(fmt.Sprintf("process task:%d", task.Id))
			time.Sleep(time.Microsecond * 10)
			ProcessTask(task)
		}(task)
	}
	return nil
}

//
// Mark tasks to processing state
//
func (t Tasker) markProcessingTasks(tasks []Task, SqlTx *SqlTxn) error {

	var ids []string
	var resources []string
	for _, task := range tasks {
		ids = append(ids, strconv.Itoa(task.Id))
		resources = append(resources, strconv.Itoa(task.Resource))
	}
	sql := `UPDATE ` + MYSQL_TASK_TABLE + `
			SET ` + "`status`" + `='` + STATUS_PROCESSING + `',
				trials=trials+1
	        WHERE
	            id IN (` + strings.Join(ids, ",") + `);`

	_, err := SqlTx.Tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("(t Tasker)#MarkProcessingTasks1: %s ", err.Error())
	}
	//Add processing list in redis as well.
	driver, rerr := redis.GetDriver()
	if rerr != nil {
		return fmt.Errorf("(t Tasker)#MarkProcessingTasks2: %s ", rerr.Error())
	}
	//@todo: Should we not ignore this error?
	driver.SADD(PROCESSING_TASKS, resources...)
	return nil
}

//
// Start a New Transaction and return transaction object.
//
func (t Tasker) startTransaction() (*SqlTxn, error) {
	driver, err := factory.GetMySqlDriver("Worker")
	if err != nil {
		return nil, fmt.Errorf(
			"(t Task)#FetchTasks(): Cannot initiate mysql: %s",
			err.Error(),
		)
	}
	txnObj, sqlErr := driver.GetTxnObj()
	if err != nil {
		return nil, fmt.Errorf(
			"(t Task)#FetchTasks(): Cannot initiate transaction: %s",
			sqlErr.DeveloperMessage,
		)
	}
	return &SqlTxn{
		Tx: txnObj,
	}, nil
}

//
// Checks if the supplied resource id already exists in our queue.
//
func (t Tasker) isInQueue(resourceId int) bool {
	for _, id := range t.TaskQueue {
		if id == utils.ToString(resourceId) {
			return true
		}
	}
	return false
}

//
// Adds the supplied resourceId to queue, so it can be checked.
//
func (t *Tasker) addToQueue(resourceId int) bool {
	t.TaskQueue = append(t.TaskQueue, utils.ToString(resourceId))
	return true
}

//
// Flush all the data in queue, need to do this before each new fetch.
//
func (t Tasker) reInitializeQueue() bool {
	t.TaskQueue = []string{}
	driver, rerr := redis.GetDriver()
	if rerr != nil {
		logger.Error(fmt.Errorf("(t Tasker) reInitializeQueue(): %s", rerr.Error()))
		return true
	}
	members, err := driver.SMembers(PROCESSING_TASKS)
	if err != nil {
		logger.Error(fmt.Errorf("(t Tasker) reInitializeQueue(2): %s", err.Error()))
		return true
	}
	t.TaskQueue = members
	return true
}

func (t Tasker) getFetchSql() string {
	sql := `SELECT
	id, version, type, resource, resource_type,
	status, trials, created_at, updated_at
	FROM ` + MYSQL_TASK_TABLE + ` AS s
	WHERE
		resource_type="` + t.ResourceType + `"
	AND
	    (
	      (status = '` + STATUS_INIT + `')
	        OR
	       (
	         status = '` + STATUS_FAILED + `'
	           AND
	         trials < ` + strconv.Itoa(MAX_TRIAL_LIMIT) + `
	       )
	    )
	ORDER BY id
	LIMIT %d;`
	sql = fmt.Sprintf(sql, t.FetchLimit)
	return sql
}

//
// Check if a particular job can be run at the moment or not.
//
func (t Tasker) canRunJob(flag string, allowedInterval time.Duration) bool {
	key := t.GetKeyName(flag)
	driver, err := redis.GetDriver()
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(t Tasker)#canRunJob()Cannot Get Redis Driver: %s", err.Error(),
		))
		return false
	}
	lastRun, rerr := driver.GetSetTime(key, time.Now().Unix())
	if rerr != nil {
		logger.Error(fmt.Sprintf(
			"(t Tasker)#canRunJob()1: %s", rerr.Error(),
		))
		return false
	}
	if lastRun == nil {
		//its first run
		return false
	}
	nextRun := lastRun.Add(allowedInterval).Unix()
	currTime := time.Now().Unix()
	if currTime >= nextRun {
		return true
	} else {
		err := driver.Set(key, lastRun.Unix())
		if err != nil {
			logger.Error(fmt.Sprintf(
				"(t Tasker)#canRunJob()2: %s", err.Error(),
			))
		}
	}
	return false
}

//
// If a task is in processing state for a long period of time
// it is said to be hanged.
//
func (t Tasker) restartHangedTasks() error {
	sql := `UPDATE ` + MYSQL_TASK_TABLE + `
	SET ` + "`status`" + ` = "` + STATUS_FAILED + `"
    WHERE (DATE_ADD(updated_at, INTERVAL 10 MINUTE) < NOW())
    AND ` + "`status`" + `='` + STATUS_PROCESSING + `';`

	driver, err := factory.GetMySqlDriver("Worker")
	if err != nil {
		return fmt.Errorf(
			"(t Task)#restartHangedTasks(): Cannot initiate mysql: %s",
			err.Error(),
		)
	}
	_, res := driver.Execute(sql)
	if res != nil {
		return fmt.Errorf(
			"(t Task)#restartHangedTasks(): %s",
			res.DeveloperMessage,
		)
	}
	return nil
}

//
// Delete Old tasks.
//
func (t Tasker) cleanupOldTasks() error {
	sql := `DELETE FROM ` + MYSQL_TASK_TABLE + `
    WHERE  (
    	` + "`status`" + ` = "` + STATUS_SUCCESS + `"
    		 OR
    	` + "`status`" + ` = "` + STATUS_CANCELED + `"
    	     OR
    	(
    		` + "`status`" + ` = "` + STATUS_FAILED + `"
    		       AND
    		trials >= ` + utils.ToString(MAX_TRIAL_LIMIT) + `
    	)
    )
    AND
    (DATE_ADD(updated_at, INTERVAL 48 HOUR) < NOW())
    `
	driver, err := factory.GetMySqlDriver(MYSQL_POOL_NAME)
	if err != nil {
		return fmt.Errorf(
			"(t Task)#cleanupOldTasks(): Cannot initiate mysql: %s",
			err.Error(),
		)
	}
	_, res := driver.Execute(sql)
	if res != nil {
		return fmt.Errorf(
			"(t Task)#cleanupOldTasks(): %s",
			res.DeveloperMessage,
		)
	}
	return nil
}

//
// Fetch the list of tasks.
//
func (t Tasker) fetchTasks(db *SqlTxn) []Task {
	sql := t.getFetchSql()
	result, err := db.Tx.Query(sql)
	if err != nil {
		fmt.Println("(t Tasker) FetchTasks: " + err.Error())
		return nil
	}
	//@todo: check if this needs to be done or not
	defer result.Close()
	var tasks []Task
	for result.Next() {
		var task Task
		result.Scan(
			&task.Id, &task.Version, &task.Type, &task.Resource,
			&task.ResourceType, &task.Status, &task.Trials,
			&task.CreatedAt, &task.UpdatedAt,
		)
		tasks = append(tasks, task)
	}
	//fmt.Println("Tasks count:", len(tasks))
	return tasks
}

//
// Prepare Redis key name based on the prefix supplied.
//
func (t Tasker) GetKeyName(prefix string) string {
	return (prefix + strings.ToLower(t.ResourceType))
}

//
// Release acquired lock.
//
func (t Tasker) releaseLock(lock string) bool {
	var trials int
	var maxTrials int = 5
	var status bool
	for !status && (trials < maxTrials) {
		if trials != 0 {
			time.Sleep(time.Millisecond * 3)
		}
		trials = trials + 1
		driver, err := redis.GetDriver()
		if err != nil {
			logger.Error(fmt.Sprintf(
				"(t Tasker)#releaseLock()Cannot Get Redis Driver: %s", err.Error(),
			))
			continue
		}
		rerr := driver.Set(lock, AVAILABLE)
		if rerr != nil {
			logger.Error(fmt.Sprintf(
				"(t Tasker)#releaseLock()1: %s", rerr.Error(),
			))
			continue
		}
		status = true
	}
	//If release lock failed notify.
	if !status {
		//notify
		notification.SendNotification(
			"TASKER: RELEASE LOCK FAILED",
			"FLAG: TASKER_BUSY",
			[]string{"tasker"},
			datadog.ERROR,
		)
	}
	return status
}

//
// Acquire a Lock.
//
func (t Tasker) getLock(lock string) bool {
	driver, err := redis.GetDriver()
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(t Tasker)#getLock()Cannot Get Redis Driver: %s", err.Error(),
		))
		return false
	}
	wasBusy, rerr := driver.GetSetBool(lock, BUSY)
	if rerr != nil {
		logger.Error(fmt.Sprintf(
			"(t Tasker)#getLock()1: %s", rerr.Error(),
		))
		return false
	}
	if wasBusy {
		return false
	}
	return true
}

//
// Send supplied task for processing
//
func ProcessTask(task Task) {
	err := task.Process()
	if err != nil {
		//set task failed
		task.SetStatus(STATUS_FAILED, err)
		logger.Error(fmt.Errorf("Tasker#ProcessTask() : %s", err.Error()))
	} else {
		//set task success
		task.SetStatus(STATUS_SUCCESS, fmt.Errorf("No Error"))
	}
	task.removeFromProcessingList()
}
