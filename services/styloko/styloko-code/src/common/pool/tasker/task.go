package tasker

import (
	productTasks "amenities/products/common/synctasks"
	factory "common/ResourceFactory"
	"common/constants"
	"common/notification"
	"common/notification/datadog"
	"common/redis"
	"common/utils"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// Stores details of task to be executed.
type Task struct {
	Id           int
	Version      int64
	Type         string
	Resource     int
	ResourceType string
	Status       string
	Trials       int
	Data         []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//
// Insert task into mysql table
//
func (t Task) Add() error {
	driver, err := factory.GetMySqlDriver("Worker")
	if err != nil {
		return fmt.Errorf("(t Task)#Add(): Cannot initiate mysql: %s", err.Error())
	}
	//Set default values
	t.Version = time.Now().Unix()
	t.Status = STATUS_INIT
	t.Trials = 0
	if err != nil {
		return fmt.Errorf("(t Task)#Add(): UnMarshalling Failed: %s", err.Error())
	}
	sql := `INSERT INTO ` + MYSQL_TASK_TABLE + ` (
			version,
    		type,
    		resource,
    		resource_type,
    		status,
    		trials,
    		data,
    		created_at)
    	VALUES
    		(?,?,?,?,?,?,?,NOW())`
	_, sqlerr := driver.Execute(sql,
		t.Version, t.Type, t.Resource, t.ResourceType, t.Status,
		t.Trials, t.Data,
	)
	if sqlerr != nil {
		return fmt.Errorf("(t Task)#Add(): Insertion Failed: %s", sqlerr.DeveloperMessage)
	}
	return nil
}

//
// Set status of task, to one of defined status.
//
func (t Task) SetStatus(status string, errr error) bool {
	var sql string
	sql = `UPDATE ` + MYSQL_TASK_TABLE + `
	         SET ` + "`status`" + ` ='` + status + `'
             WHERE id=?;`
	driver, err := factory.GetMySqlDriver("Worker")
	if err != nil {
		logger.Error(fmt.Errorf("(t Task)#SetStatus(%s):%s", status, err.Error()))
		return false
	}
	_, sqlErr := driver.Execute(sql, t.Id)
	if sqlErr != nil {
		logger.Error(fmt.Errorf("(t Task)#SetStatus(%s):%s", status, sqlErr.DeveloperMessage))
		return false
	}
	if status == STATUS_FAILED && t.Trials == (MAX_TRIAL_LIMIT-1) {
		//notify
		text := fmt.Sprintf("ConfigId:%v, Reason:%v", t.Resource, errr.Error())
		notification.SendNotification(
			fmt.Sprintf("TASKER [%d]: FAILED MAX TRIALS REACHED", t.Id),
			text,
			[]string{"tasker"},
			datadog.ERROR,
		)
	}
	return true
}

//
// Handles processsing of task.
//
func (t Task) Process() error {
	var err error
	var data []byte
	data, err = t.GetData()
	if err != nil {
		return err
	}
	switch t.ResourceType {
	case constants.PRODUCT_RESOURCE_NAME:
		err = productTasks.ProcessTask(t.Type, data, t.Resource, false)
	default:
		err = fmt.Errorf("(t Task) Process(): Undefined Resource")
	}
	return err
}

func (t Task) GetData() ([]byte, error) {
	driver, err := factory.GetMySqlDriver("Worker")
	if err != nil {
		return nil, fmt.Errorf(
			"(t Task)#GetData(): Cannot initiate mysql: %s",
			err.Error(),
		)
	}
	sql := `SELECT data from ` + MYSQL_TASK_TABLE + ` WHERE id=?`
	res, dberr := driver.Query(sql, t.Id)
	if dberr != nil {
		return nil, fmt.Errorf(
			"(t Task)#GetData()1: %s",
			dberr.DeveloperMessage,
		)
	}
	defer res.Close()
	var data []byte
	for res.Next() {
		err := res.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("(sync MySqlSync)#SaveCategories scan error: %s", err.Error())
		}
	}
	return data, nil
}

func (t Task) removeFromProcessingList() bool {

	var trials int = 3
	for trials > 0 {
		trials--
		driver, rerr := redis.GetDriver()
		if rerr != nil {
			continue
		}
		_, err := driver.SREM(PROCESSING_TASKS, utils.ToString(t.Resource))
		if err != nil && !err.IsNotFound() {
			continue
		}
		return true
	}
	return false
}
