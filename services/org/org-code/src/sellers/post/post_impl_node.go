package post

import (
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/mongodb"
	"common/notification"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"sellers/common"
	"time"
)

type Insert struct {
	id string
}

func (s *Insert) SetID(id string) {
	s.id = id
}

func (s Insert) GetID() (id string, err error) {
	return s.id, nil
}

func (s Insert) Name() string {
	return "INSERT seller by id"
}

func (s Insert) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_INSERT)
	defer func() {
		logger.EndProfile(profiler, SELLER_INSERT)
	}()
	mgoSession := mongo.GetMongoSession(common.SELLERS)
	mgoObj := mgoSession.SetCollection(common.SELLERS_COLLECTION)
	defer mgoSession.Close()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_INSERT, "Seller insert execution started")
	logger.Info("Insert Seller Started")

	p, _ := io.IOData.Get(ORG_DATA)
	//parsing data into []common.Schema from interface{}
	data, pOk := p.([]common.Schema)
	if !pOk || data == nil {
		logger.Error("OrgCreate.invalid type of columns")
		notification.SendNotification("Seller Create Json Incorrect", "Data type mismatch", nil, "error")
		return io, &florest_constants.AppError{Code: appconstant.DataTypeMismatch, Message: "Json Incorrect"}
	}
	//calling insert function to insert data into mongo
	logger.Info(fmt.Sprintf("Inserting data in mongo :%v", data))
	resp := s.CheckAndInsert(mgoSession, mgoObj, data)
	//setting value recieved from update() into RESULT
	io.IOData.Set(florest_constants.RESULT, resp)
	logger.Info(fmt.Sprintf("Getting data for inserted seller :%v", resp))
	data, errs := common.GetDataForInsertedIds(resp)
	if errs != nil {
		notification.SendNotification("Error while getting data for inserted sellers after Seller Create", errs.Error(), nil, "error")
		logger.Error(fmt.Sprintf("Error while getting data for inserted sellers : %s", errs.Error()))
		return io, errs
	}
	logger.Debug(fmt.Sprintf("If len(data) not 0, call worker to sync with mysql :%v", data))
	if len(data) != 0 {
		s.SyncInsertedIdsInMysql(data)
	}
	common.PushDataToMemcache(data, "insert", "inserted seller")
	logger.Info("Insert seller finished")
	return io, nil
}

//This function return a newly generated seqId by incrementing one
//from the seqId founf in counters collection and checks if it already exists
//sends 0 and err if id exists, id and nil otherwise
func (s Insert) GetNewSeqId(mgoSession *mongodb.MongoDriver) (int, error) {
	bsonMap := make(map[string]interface{})
	id := mgoSession.GetNextSequence(common.SELLERS_COLLECTION)
	bsonMap["seqId"] = id
	ok, _, err := common.CheckIfKeyExists(bsonMap)
	if ok {
		return 0, err
	}
	return id, nil
}

//Checks if sequence Id already exists, if not data is inserted in mongo
// and and []Response is returned with seller name and newly created sequence Id
func (s Insert) CheckAndInsert(mgoSession *mongodb.MongoDriver, mgoObj *mgo.Collection, data []common.Schema) []Response {
	var arrInsert []Response
	for _, v := range data {
		var temp Response
		id, err := s.GetNewSeqId(mgoSession)
		if err == nil {
			v.SeqId = id
			time := time.Now()
			v.CreatedAt = &time
			v.UpdatedAt = &time
			orgInfo, err := common.Insert(mgoObj, v, true)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while inserting into Mongo  %s", err))
				notification.SendNotification("Insert in Mongo failed while Creating Seller", fmt.Sprintf("Org Name : %s , Error : %s", v.OrgName, err.Error()), nil, "error")
				temp.Name = v.OrgName
				temp.SeqId = 0
				temp.Error = err.Error()
				arrInsert = append(arrInsert, temp)
				continue
			}
			temp.Name = orgInfo.OrgName
			temp.SeqId = orgInfo.SeqId
			arrInsert = append(arrInsert, temp)
			continue
		}
		notification.SendNotification("Unable to create seqId while Creating Seller", fmt.Sprintf("Org Name : %s , Error : %s", v.OrgName, err.Error()), nil, "error")
		temp.Name = v.OrgName
		temp.SeqId = id
		arrInsert = append(arrInsert, temp)
	}
	return arrInsert
}

func (s Insert) SyncInsertedIdsInMysql(data []common.Schema) {
	//preparing data to call worker to sync with mysql
	dataMapArr := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	temp["data"] = data
	temp["command"] = SELLER_INSERT
	dataMapArr = append(dataMapArr, temp)
	logger.Info("Starting job,worker called")
	sellerInsertPool.StartJob(dataMapArr)
}
