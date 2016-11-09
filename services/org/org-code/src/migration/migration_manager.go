package migration

import (
	"common/ResourceFactory"
	"common/appconstant"
	"common/notification"
	"common/pool"
	"common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"migration/common"
	"migration/erp"
	"migration/supplier"
)

type Manager struct {
	id string
}

func (m *Manager) SetID(id string) {
	m.id = id
}

func (m Manager) GetID() (id string, err error) {
	return m.id, nil
}

func (m Manager) Name() string {
	return "Migration Manager"
}

func (m Manager) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	defer pool.RecoverHandler("MigrationAPI")
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.MIGRATION_MANAGER)
	defer func() {
		logger.EndProfile(profiler, common.MIGRATION_MANAGER)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+m.Name(), rc)
	io.ExecContext.SetDebugMsg(common.MIGRATION_MANAGER, "Migration Manager execution started")
	exists := checkFlagInMongo()
	if exists {
		unsetFlagInMongo()
		//getting headers
		header, err := utils.GetRequestHeader(io)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while getting headers %s:", err.Error()))
			setFlagInMongo()
			notification.SendNotification("Error while getting Headers for Migration", err.Error(), nil, "error")
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting Headers", DeveloperMessage: err.Error()}
		}
		//checking if header for supplier is set
		if _, ok := header["Supplier"]; ok {
			err := supplier.StartSupplierMigration()
			if err != nil {
				logger.Error(fmt.Sprintf("Error while migrating supplier :%s", err.Error()))
				setFlagInMongo()
				notification.SendNotification("Supplier Migration Failed", err.Error(), nil, "error")
				return io, &florest_constants.AppError{Code: appconstant.MigrationError, Message: "Supplier Migration Failed", DeveloperMessage: err.Error()}
			}
			//checking if header for erp is set
		} else if _, ok := header["Erp"]; ok {
			//getting post data
			data, err := utils.GetPostData(io)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while getting post data :%s", err.Error()))
				setFlagInMongo()
				notification.SendNotification("Error while getting post data for Erp Migration", err.Error(), nil, "error")
				return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting post data", DeveloperMessage: err.Error()}
			}
			//starting erp migration with the post data
			resp, err := erp.StartErpMigration(data)
			//if err other than inserting in mongo
			if err != nil && len(resp) == 0 {
				logger.Error(fmt.Sprintf("Error while migrating erp data :%s", err.Error()))
				setFlagInMongo()
				notification.SendNotification("Erp Migration Failed", err.Error(), nil, "error")
				return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Erp Migration Failed", DeveloperMessage: err.Error()}
			} else {
				//if error while inserting in mongo
				setFlagInMongo()
				io.IOData.Set(florest_constants.RESULT, resp)
				return io, nil
			}
			//if no header was set,return error
		} else {
			logger.Error("No header found")
			setFlagInMongo()
			notification.SendNotification("Error while Migration", "Required Header not found", nil, "error")
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while migration", DeveloperMessage: "Required Header not found"}
		}
		//if flag was not set in mongo,migration already running
	} else {
		logger.Error("Error while migrating supplier : Already running")
		setFlagInMongo()
		notification.SendNotification("Migration Failed", "Migration Already Running", nil, "error")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Migration Failed", DeveloperMessage: "Migration Already Running"}
	}
	setFlagInMongo()
	io.IOData.Set(florest_constants.RESULT, "Migration Finished")
	return io, nil
}

//This function sets flag as true in mongo for _id:seller
func setFlagInMongo() {
	mgoSession := ResourceFactory.GetMongoSession(common.MIGRATION)
	mgoObj := mgoSession.SetCollection(common.FLAG_COLLECTION)
	defer mgoSession.Close()
	findCriteria := map[string]interface{}{"_id": common.SELLER}
	updateCriteria := map[string]interface{}{"flag": true}
	_, err := mgoObj.Upsert(findCriteria, updateCriteria)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while setting flag in mongo :%s", err.Error()))
		return
	}
	logger.Info("Key successfully set in mongo")
}

//This function checks if the flag value for _id:sellers is true,
//if flag value is true returns true else false
func checkFlagInMongo() bool {
	mgoSession := ResourceFactory.GetMongoSession(common.MIGRATION)
	mgoObj := mgoSession.SetCollection(common.FLAG_COLLECTION)
	defer mgoSession.Close()
	findCriteria := map[string]interface{}{"_id": common.SELLER}
	m := make(map[string]interface{})
	err := mgoObj.Find(findCriteria).One(&m)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting flag value from mongo :%s", err.Error()))
		return false
	}
	if _, ok := m["flag"]; !ok {
		logger.Error("Error while reading flag value from mongo : Does not exist")
		return false
	}
	return m["flag"].(bool)
}

//This function sets flag as false in mongo for _id:seller
func unsetFlagInMongo() {
	mgoSession := ResourceFactory.GetMongoSession(common.MIGRATION)
	mgoObj := mgoSession.SetCollection(common.FLAG_COLLECTION)
	defer mgoSession.Close()
	findCriteria := map[string]interface{}{"_id": common.SELLER}
	updateCriteria := map[string]interface{}{"flag": false}
	_, err := mgoObj.Upsert(findCriteria, updateCriteria)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unsetting flag in mongo :%s", err.Error()))
		return
	}
	logger.Info("Key successfully unset in mongo")
}
