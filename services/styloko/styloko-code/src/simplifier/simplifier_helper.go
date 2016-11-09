package simplifier

import (
	"common/ResourceFactory"
	"common/appconfig"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

var flag bool

//This func reads from mongo a single document whose is_picked is false, returns the document
//after changing its is_picked to true (ensuring atomicity in multi-cluster environment)
func readFromMongo() (UpdateJob, error) {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	var updateDoc UpdateJob
	upsertVal := false
	deleteVal := false
	returnNew := false
	updatedVal := bson.M{"$set": bson.M{"is_picked": true}}
	findCriteria := bson.M{"is_picked": false}
	change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
	_, err := mgoObj.Find(bson.M(findCriteria)).Sort("created_at").Apply(change, &updateDoc)
	if err != nil && err.Error() != "not found" {
		logger.Error(fmt.Sprintf("Error while reading from mongo:%s", err.Error()))
	}
	return updateDoc, err
}

//changes status in redis of jobname passed to the function
func changeStatusInRedis(jobName string, status string) error {
	client, err := ResourceFactory.GetRedisDriver(JUDGE_DAEMON)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting redis driver:%s", err.Error()))
		return err
	}
	client.HSet(jobName, "status", status)
	return nil
}

//This function breaks in configurable chunks and hits api,changes is_completed for the processed chunk
func breakInChunks(prodDoc *UpdateJob) (string, error) {
	var err error
	var stringErr string
	var ErrString string
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	chSize := config.JudgeDaemon.ChunkSize
	apiUrl := config.JudgeDaemon.ApiUrl
	t := config.JudgeDaemon.Timeout

	totalDataSize := len(prodDoc.Data)
	for i := 0; i < totalDataSize; i += chSize {
		var endIndex = totalDataSize
		if (i + chSize) <= totalDataSize {
			endIndex = (i + chSize)
		}
		batch := prodDoc.Data[i:endIndex]
		// Do something with batch
		time.Sleep(100 * time.Millisecond)
		logger.Info(fmt.Sprintf("Prepare chunk data for %d data", len(batch)))
		stringErr, err = hitApiForChunk(batch, getHeaders(prodDoc.Type), apiUrl, t, prodDoc.Type)
		if err != nil {
			if flag {
				ErrString += fmt.Sprintf("%s,%s, ObjectId - %s,%s,%s,\n", "Error while hitting api for Judge Daemon", err.Error(), prodDoc.Id, "", "")
				ErrString += fmt.Sprintf("%s", getFailedSku(batch, getSkuString(prodDoc.Type), flag))
			} else {
				ErrString += fmt.Sprintf("%s,%s, ObjectId - %s\n", "Error while hitting api for Judge Daemon", err.Error(), prodDoc.Id)
				ErrString += fmt.Sprintf("%s", getFailedSku(batch, getSkuString(prodDoc.Type), flag))
			}
			logger.Error(fmt.Sprintf("Error while hitting api for Judge Daemon:%s", err.Error()))
			continue
		}
		ErrString += fmt.Sprintf("%s", stringErr)
	}
	err = changeStateInMongo(prodDoc.JobName, prodDoc.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while changing state in mongo:%s", err.Error()))
	}
	//check if all chunks are processed
	checkIfJobCompleted(prodDoc.JobName)
	return ErrString, nil
}

//get headers according to the update type in the document
func getHeaders(updateType string) map[string]string {
	headers := make(map[string]string, 0)
	if updateType == "attribute" {
		flag = true
		headers["Update-Type"] = "Attribute"
	}
	if updateType == "discount" {
		headers["Update-Type"] = "JabongDiscount"
	}
	if updateType == "price" {
		headers["Update-Type"] = "Price"
	}
	headers["RequestSource"] = "CatalogAdmin"
	return headers
}

//This function hits api for chunk and transforms response and sends back string of errors
func hitApiForChunk(prodChunk []interface{}, headers map[string]string, apiUrl string, t string, updateType string) (string, error) {
	timeout, err := time.ParseDuration(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting timeout to time.Duration: %s", err.Error()))
		return "", err
	}
	prodJson, err := json.Marshal(prodChunk)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting product chunk to json: %s", err.Error()))
		return "", err
	}
	resp, err := http.HttpPut(apiUrl, headers, string(prodJson), timeout*time.Millisecond)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending PUT Request %s", err.Error()))
		return "", err
	}
	stringErr, err := transformResponse(prodChunk, resp, updateType)
	return stringErr, err
}

//This function transforms the response from api hit and concats an error string for failures
func transformResponse(prodChunk []interface{}, resp *http.APIResponse, updateType string) (string, error) {
	if resp.HttpStatus/10 != 20 {
		logger.Error(fmt.Sprintf("Error In Job Response : %s", string(resp.Body)))
		return "", errors.New(fmt.Sprintf("Error In Job Response : %s", string(resp.Body)))
	}
	sku := getSkuString(updateType)
	var res FlorestResp
	err := json.Unmarshal(resp.Body, &res)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling json to Job Response : %s", err.Error()))
		return "", err
	}
	var er string
	for k, v := range res.Data {
		v1, _ := v.(map[string]interface{})
		if v1["status"] == "failure" {
			singleProd, ok := prodChunk[k].(bson.M)
			val, ok := v1["error"].(map[string]interface{})
			if !ok {
				if flag {
					er += fmt.Sprintf("%s,%s,%s,%s,%s\n", singleProd[sku], "Invalid Florest Error Json", "Json is incorect", "", "")
				} else {
					er += fmt.Sprintf("%s,%s,%s\n", singleProd[sku], "Invalid Florest Error Json", "Json is incorect")
				}
				logger.Error(fmt.Sprintf("Error while unmarshalling error:%s", err.Error()))
			}
			if flag {
				er += fmt.Sprintf("%s,%s,%s,%s,%s\n", singleProd[sku], singleProd["AttributeName"], singleProd["Value"], val["message"], val["developerMessage"])
			} else {
				er += fmt.Sprintf("%s,%s,%s\n", singleProd[sku], val["message"], val["developerMessage"])
			}
		}
	}
	return er, nil
}

//for different type of updates,different keys are used for sku
func getSkuString(updateType string) string {
	var sku string
	if updateType == "attribute" {
		sku = "ProductSku"
	}
	if updateType == "discount" {
		sku = "productId"
	}
	if updateType == "price" {
		sku = "sku"
	}
	return sku
}

//This function logs progress and errors for the processed chunk
func logProgressAndErrors(stringErr string, jobName string) error {
	err := logError(stringErr, jobName)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while logging errors:%s", err.Error()))
	}
	err = logProgress(jobName)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while logging progress:%s", err.Error()))
	}
	return err
}

//This function inserts all the compiled errors into mongo
func logError(stringErr string, jobName string) error {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON_ERRORS)
	defer mgoSession.Close()
	var errStruct ErrorStruct
	errStruct.JobName = jobName
	errStruct.ErrMsg = stringErr
	err := mgoObj.Insert(errStruct)
	return err
}

//This function changes state in mongo for the passed id
func changeStateInMongo(jobname string, id bson.ObjectId) error {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	updatedVal := bson.M{"$set": bson.M{"is_completed": true}}
	findCriteria := bson.M{"_id": id}
	err := mgoObj.Update(findCriteria, updatedVal)
	return err
}

//This function logs progress in redis by checking no of is_completed/total no of chunks for each jobanme
func logProgress(jobName string) error {
	mgoSession := ResourceFactory.GetMongoSessionWithDb(JUDGE_DAEMON, JUDGE)
	mgoObj := mgoSession.SetCollection(JUDGE_DAEMON)
	defer mgoSession.Close()
	var updateDocs []UpdateJob
	err := mgoObj.Find(bson.M{"job_name": jobName}).Select(bson.M{"is_completed": 1, "job_name": 1}).All(&updateDocs)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting count for completed jobs:%s", err.Error()))
		return err
	}
	var count int
	for _, v := range updateDocs {
		if v.IsCompleted == true {
			count++
		}
	}
	progress := int(count * 100 / len(updateDocs))
	progString := strconv.Itoa(progress)
	client, err := ResourceFactory.GetRedisDriver(JUDGE_DAEMON)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting redis client:%s", err.Error()))
		return err
	}
	client.HSet(jobName, "progress", progString)
	return nil
}

func getFailedSku(batch []interface{}, sku string, flag bool) string {
	var errString string
	for k, _ := range batch {
		vNew := batch[k].(bson.M)
		if flag {
			errString += fmt.Sprintf("%s,%s,%s,%s,%s\n", "Failed Sku", vNew[sku], "", "", "")
		} else {
			errString += fmt.Sprintf("%s,%s,%s\n", "Failed Sku", vNew[sku], "")
		}
	}
	return errString
}
