package put

import (
	"common/ResourceFactory"
	"common/appconstant"
	"common/notification"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"sellers/common"
	"strconv"
	"time"
)

type Update struct {
	id string
}

func (s *Update) SetID(id string) {
	s.id = id
}

func (s Update) GetID() (id string, err error) {
	return s.id, nil
}

func (s Update) Name() string {
	return "Update seller by id"
}

func (s Update) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_UPDATE)
	defer func() {
		logger.EndProfile(profiler, SELLER_UPDATE)
	}()
	mgoSession := ResourceFactory.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_UPDATE, "Seller update execution started")
	logger.Info("Update Seller Started")

	p, _ := io.IOData.Get(ORG_UPDATE_DATA)
	//parsing data into []common.Schema from interface{}
	data, pOk := p.([]common.Schema)
	if !pOk || data == nil {
		logger.Error("OrgUpdate.invalid type of columns")
		notification.SendNotification("Seller Update Json Incorrect", "Data tpe mismatch", nil, "error")
		return io, &florest_constants.AppError{Code: appconstant.DataTypeMismatch, Message: "Json Incorrect"}
	}
	//calling update function to update data into mongo
	resp := s.Update(mgoObj, data)
	//setting value recieved from update() into RESULT
	io.IOData.Set(florest_constants.RESULT, resp)
	logger.Info("Update seller finished")
	//getting correctly inserted ids in resp
	sellerData, err := common.GetDataForInsertedIds(resp)
	if err != nil {
		notification.SendNotification("Error while getting correctly inserted ids after Seller Update", err.Error(), nil, "error")
		logger.Error(fmt.Sprintf("Error while getting correctly inserted ids :%s", err.Error()))
	}
	//start workers only if any id was updated
	if len(sellerData) != 0 {
		//starting job to sending updated ids to product API
		prodUpdatePool.StartJob(sellerData)
		//start job to sync with mysql
		s.SyncUpdatedIdsInMysql(sellerData)
		//starting job to update on erp
		//commented only for perf
		updateErpPool.StartJob(sellerData)
	}
	common.PushDataToMemcache(sellerData, "update", "updated seller")
	return io, nil
}

//function takes []common.Schema and returns []Response with
// seqId as the id of the updated seller and result as true or false
//if the seller update was success or failure respectively.
//It also deletes key from cache on updation.
func (s Update) Update(mgoObj *mgo.Collection, data []common.Schema) []Response {
	var arrUpdate []Response
	for _, v := range data {
		var temp Response
		//setting current time as updated time
		time := time.Now()
		v.UpdatedAt = &time
		orgInfo, err := common.Insert(mgoObj, v, false)
		if err != nil {
			temp.SeqId = v.SeqId
			temp.Result = false
			arrUpdate = append(arrUpdate, temp)
			notification.SendNotification("Insert in Mongo failed while Updating Seller", fmt.Sprintf("SeqId : %d , Error : %s", v.SeqId, err.Error()), nil, "error")
			logger.Error(fmt.Sprintf("Error while updating into Mongo.%s", err))
		} else {
			temp.SeqId = orgInfo.SeqId
			temp.Result = true
			arrUpdate = append(arrUpdate, temp)
		}
		go s.Delete(temp.SeqId)
	}
	return arrUpdate
}

//Deleting key from cache
func (s Update) Delete(key int) {
	err := cacheObj.Delete(fmt.Sprintf("%s-%s", common.SELLERS, strconv.Itoa(key)))
	if err != nil {
		logger.Error(fmt.Sprintf("Error while deleting key from cache: %v", err.Error()))
	}
}

//preparing data to call worker to sync with mysql
func (s Update) SyncUpdatedIdsInMysql(data []common.Schema) {
	dataMapArr := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	temp["data"] = data
	temp["command"] = SELLER_UPDATE
	dataMapArr = append(dataMapArr, temp)
	sellerUpdatePool.StartJob(dataMapArr)
}
