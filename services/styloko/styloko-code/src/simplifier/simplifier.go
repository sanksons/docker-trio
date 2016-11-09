package simplifier

import (
	"common/ResourceFactory"
	"common/appconfig"
	"common/pool"
	"common/utils"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// updatePool local pool variable
var updatePool pool.Safe

// StartJudgeDaemon - starts judge daemon
// Initializes a pool of workers and keeps the goroutines in waiting state
// also starts three crons -
// runUpdateJobs - for all newly inserted documents and starts EveryJob()
// runResetJobs - for all incompletely processed chunks,resets is_picked and for complete jobs,set status in redis
// runCleanUpJobs - for all chunks of data that are completely proccessed, it deletes it form database
func StartJudgeDaemon() {
	logger.Info("Starting worker pool for JUDGE DAEMON JOBS")
	updatePool = pool.NewWorkerSafe(JUDGE_DAEMON, JUDGE_DAEMON_POOL_SIZE, JUDGE_DAEMON_QUEUE_SIZE, JUDGE_DAEMON_RETRY_COUNT, JUDGE_DAEMON_WAIT_TIME)
	updatePool.StartWorkers(EveryJob)
	runUpdateJobs()
	runCleanupJob()
	runResetJobs()
}

//reads time from config and fires StartJob in an infinite loop with configurable time delay
func runUpdateJobs() {
	logger.Info("Starting Update Worker")
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	t := config.JudgeDaemon.RunEveryJob
	everyJobRuntime, err := time.ParseDuration(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting RunEvery to time.Duration: %s", err.Error()))
		//set default value
		everyJobRuntime = EVERYJOB_RUNTIME
	}
	go func() {
		utils.RecoverHandler("runUpdateJobs")
		for {
			updatePool.StartJob("Start")
			time.Sleep(everyJobRuntime)
		}
	}()
}

// runCleanupJob reads time from config and fires deleteJobs() in an infinite loop with time delay
func runCleanupJob() {
	logger.Info("Starting Cleanup worker")
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	t := config.JudgeDaemon.RunCleanupJob
	cleanupJobRuntime, err := time.ParseDuration(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting RunCleanupJob to time.Duration: %s", err.Error()))
		//set default value
		cleanupJobRuntime = CLEANUPJOB_RUNTIME
	}
	go func() {
		utils.RecoverHandler("runCleanupJob")
		for {
			deleteJobs()
			time.Sleep(cleanupJobRuntime)
		}
	}()
}

// runResetJobs reads time from config and fires resetIncompleteJobs in an infinite loop with delay
func runResetJobs() {
	logger.Info("Starting reset worker")
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	t := config.JudgeDaemon.RunResetJob
	resetJobRuntime, err := time.ParseDuration(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting RunResetJob to time.Duration: %s", err.Error()))
		//set default value
		resetJobRuntime = RESETJOB_RUNTIME
	}
	go func() {
		utils.RecoverHandler("runResetJobs")
		for {
			resetIncompleteJobs()
			time.Sleep(resetJobRuntime)
		}
	}()
}

//This function gets distinct jobNames from mongo, loops over each name,
//gets all documents from mongo for each name, check if is_completed flag is set true for all,
//if yes, update status to completed in redis and delete all documents from mongo for the jobName
func checkIfJobCompleted(jobName string) {
	logger.Info("Starting checkIfJobCompleted")
	jobInfo := findJobsByJobName(jobName)
	if len(jobInfo) == 0 {
		//return if no documents found
		return
	}
	flag := false
	for _, v := range jobInfo {
		if v.IsCompleted == false {
			//setting flag as true if any document has is_completed =true
			flag = true
		}
	}
	if flag == false {
		//if all documents for the jobName have is_completed=true
		err := changeStatusInRedis(jobName, "completed")
		if err != nil {
			logger.Error(fmt.Sprintf("Error while changing status to completed in redis:%s", err.Error()))
			return
		}
	}
}

//deletes documents from mongo for which is_completed is true for all chunks for a jobName
func deleteJobs() {
	logger.Info("Starting delete jobs")
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	var updateDocs []string
	err := mgoObj.Find(nil).Distinct("job_name", &updateDocs)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting distinct jobNames from mongo:%s", err.Error()))
		//no error returned as cron will retry
		return
	}
	for _, v := range updateDocs {
		jobInfo := findJobsByJobName(v)
		if len(jobInfo) == 0 {
			//return if no documents found
			return
		}
		flag := false
		for _, v := range jobInfo {
			if v.IsCompleted == false {
				//setting flag as true if any document has is_completed =true
				flag = true
			}
		}
		if flag == false {
			err = mgoObj.Remove(bson.M{"job_name": v})
			if err != nil {
				logger.Error(fmt.Sprintf("Error while deleting documents in mongo for jobname(%s) :%s", v, err.Error()))
				//no error returned as cron will retry
			}
		}
	}
}

//finds all documents in mongo for the jobname passed
func findJobsByJobName(jobName string) []UpdateJob {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	var updateDoc []UpdateJob
	err := mgoObj.Find(bson.M{"job_name": jobName}).Select(bson.M{"is_picked": 1, "is_completed": 1, "job_name": 1}).All(&updateDoc)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting documents by jobname in mongo:%s", err.Error()))
		//no error returned as cron will retry and len(updateDoc) checked at func return
	}
	return updateDoc
}

//This function reads documents from mongo whose is_picked is false,changes status in redis from pending
//to running for the picked job, breaks the data in the document to configurable chunks and hits the api
//,tranfsorms the response,logs progress in redis after completion of a chunk and logs error in mongo
func EveryJob(data interface{}) error {
	logger.Info("Starting EveryJob")
	prodDoc, err := readFromMongo()
	if err != nil {
		//no error returned as cron will try again
		return nil
	}
	err = changeStatusInRedis(prodDoc.JobName, "running")
	if err != nil {
		logger.Error(fmt.Sprintf("Error while changing status in Redis:%s", err.Error()))
		//no error returned as next chunk can change progress or eventually background job will change status
	}
	stringErr, err := breakInChunks(&prodDoc)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while breaking into chunks:%s", err.Error()))
		logError(stringErr, prodDoc.JobName)
		//no error returned as error already logged in mongo for failed chunk
	}
	//updating UpdatedAt of the processed chunk
	prodDoc.UpdatedAt = time.Now()
	err = logProgressAndErrors(stringErr, prodDoc.JobName)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while logging progress and errors:%s", err.Error()))
	}
	return err
}

// //This function finds all incompletely run documents whose os_picked is true and is_completed is false
// //and sets is_picked as false to be picked up by the updateJob cron later
func resetIncompleteJobs() {
	logger.Info("Starting to reset incomplete Jobs")
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	t := config.JudgeDaemon.RunResetJob
	resetJobRuntime, err := time.ParseDuration(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting resetJobRuntime to time.Duration: %s", err.Error()))
	}
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	var updateDocs []UpdateJob
	err = mgoObj.Find(bson.M{"is_picked": true, "is_completed": false}).Select(bson.M{"_id": 1}).All(&updateDocs)
	if err != nil && err.Error() != "not found" {
		logger.Error(fmt.Sprintf("Error while finding incomplete jobs from mongo:%s", err.Error()))
		return
	}
	//if no documents were found return
	if len(updateDocs) == 0 {
		return
	}
	//range over found documents and check if updated at > resetJobRuntime then elete documents from db
	for _, v := range updateDocs {
		if time.Now().Sub(v.UpdatedAt) > resetJobRuntime {
			err := mgoObj.Update(bson.M{"_id": v.Id}, bson.M{"$set": bson.M{"is_picked": false}})
			if err != nil {
				logger.Error(fmt.Sprintf("Error while updating is_picked in mongo:%s", err.Error()))
			}
		}
	}
}
