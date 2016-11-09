package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type UpdateOption struct {
	id string
}

func (u *UpdateOption) SetID(id string) {
	u.id = id
}

func (u UpdateOption) GetID() (id string, err error) {
	return u.id, nil
}

func (u UpdateOption) Name() string {
	return "Attribute-option Update api"
}

func (u UpdateOption) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_OPTION_UPDATE)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_OPTION_UPDATE)
	}()
	io.ExecContext.SetDebugMsg(ATTRIBUTE_OPTION_UPDATE, "Attribute-option Update Execute")

	data, _ := io.IOData.Get(PATH_PARAMETERS)
	params := data.(Parameters)
	attrId := params.AttrId
	optionId := params.OptionId

	d, _ := io.IOData.Get(UPDATE_DATA)
	formData, ok := d.(*Option)

	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Form Data is not valid"}
	}

	res, uerr := u.Update(attrId, optionId, formData)
	if uerr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Parameters",
			DeveloperMessage: uerr.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, res)

	syncMap := make(map[string]interface{}, 0)
	syncMap["id"] = attrId
	syncMap["optId"] = optionId
	syncMap["options"] = formData
	attributeUpdatePool.StartJob(syncMap)

	return io, nil
}

func (u UpdateOption) Update(attrId int, optionId int, formData *Option) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_OPTION_UPDATE)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_OPTION_UPDATE)
		mongoDriver.Close()
	}()
	updateQuery := M{
		"$set": M{
			"options.$.value":     formData.Value,
			"options.$.position":  formData.Position,
			"options.$.isDefault": formData.IsDefault,
		},
	}

	query := M{"seqId": attrId,
		"options": M{
			"$elemMatch": M{
				"seqId": optionId,
			},
		},
	}

	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTES_COLLECTION, updateQuery, query, false)
	if err != nil {
		return res, err
	}
	return res, nil
}

// optionUpdateWorker is the worker for option creation
func optionUpdateWorker(data interface{}) error {
	dataMap, _ := data.(map[string]interface{})
	seqId, _ := dataMap["id"].(int)
	optId, _ := dataMap["optId"].(int)
	options, _ := dataMap["options"].(*Option)
	driver, err := factory.GetMySqlDriver(constants.ATTRIBUTEAPI)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", err.Error()))
		return err
	}
	txnObj, serr := driver.GetTxnObj()
	if serr != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire transaction object. Error: %s", serr.DeveloperMessage))
		logger.Debug(serr)
		return errors.New(serr.DeveloperMessage)
	}
	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txnObj.Rollback()
		}
	}()
	attrSet, attrName, err := getAttrSetOrGlobal(seqId)
	if err != nil {
		logger.Error(fmt.Errorf("Error occured in mongo fetch for attribute set. %s", err.Error()))
	}
	tableName := "catalog_attribute_option_" + attrSet + "_" + attrName
	idName := "id_" + tableName
	optName := options.Value
	pos := options.Position
	def := options.IsDefault
	query := "UPDATE " + tableName + " SET name=?, name_en=?, position=?, is_default=? WHERE " + idName + "=?"
	_, err = txnObj.Exec(query, optName, optName, pos, def, optId)
	if err != nil {
		logger.Error("Exec Failed. Transaction rollback begin.")
		txnObj.Rollback()
		return err
	}
	err = txnObj.Commit()
	if err != nil {
		logger.Error("Commit Failed. Transaction rollback begin.")
		txnObj.Rollback()
		return err
	}
	completeFlag = true
	return nil
}

func getAttrSetOrGlobal(id int) (string, string, error) {
	mgoSession := factory.GetMongoSession("AttributeOptionWorker")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(constants.ATTRIBUTES_COLLECTION)
	var attr WAttribute
	err := mgoObj.Find(bson.M{"seqId": id}).One(&attr)
	if err != nil {
		return "", "", err
	}
	if attr.IsGlobal == 0 {
		return *attr.Set.Name, *attr.Name, nil
	}
	return "global", *attr.Name, nil
}
