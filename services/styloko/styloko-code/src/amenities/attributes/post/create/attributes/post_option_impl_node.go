package attributes

import (
	factory "common/ResourceFactory"
	"common/appconstant"
	constants "common/constants"
	"common/mongodb"
	"errors"
	"fmt"
	"sort"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

type InsertOption struct {
	id string
}

func (u *InsertOption) SetID(id string) {
	u.id = id
}

func (u InsertOption) GetID() (id string, err error) {
	return u.id, nil
}

func (u InsertOption) Name() string {
	return "Attribute-Option Insertion api"
}

func (u InsertOption) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_OPTION_INSERT)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_OPTION_INSERT)
	}()
	io.ExecContext.SetDebugMsg(ATTRIBUTE_OPTION_INSERT, "Attribute-option insertion Execute")

	data, _ := io.IOData.Get(PATH_PARAMETERS)
	params := data.(Parameters)
	attrId := params.AttrId

	d, _ := io.IOData.Get(INSERT_DATA)
	formData, ok := d.([]Option)
	if !ok {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Flag",
			DeveloperMessage: "Form Data is not valid"}
	}
	res, uerr := u.Insert(attrId, formData)
	if uerr != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Invalid Parameters",
			DeveloperMessage: uerr.Error()}
	}
	io.IOData.Set(florest_constants.RESULT, res)

	u.addFilterAttributes(attrId, formData)

	syncMap := make(map[string]interface{}, 0)
	syncMap["id"] = attrId
	syncMap["options"] = formData
	attributeCreatePool.StartJob(syncMap)

	return io, nil
}

func (u InsertOption) Insert(attrId int, formData []Option) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_OPTION_INSERT)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_OPTION_INSERT)
		mongoDriver.Close()
	}()
	var optionArray []Option
	baseSeq := u.GetOptionSeq(attrId)
	for x := range formData {
		formData[x].SeqId = x + baseSeq
		optionArray = append(optionArray, formData[x])
	}
	updateQuery := M{
		"$push": M{
			"options": M{
				"$each": optionArray,
			},
		},
	}
	query := M{"seqId": attrId}
	res, err := mongoDriver.FindAndModify(constants.ATTRIBUTES_COLLECTION, updateQuery, query, true)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (u InsertOption) GetOptionSeq(attrId int) int {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, POST_GET_OPTION_SEQ)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, POST_GET_OPTION_SEQ)
		mongoDriver.Close()
	}()
	var options CheckOptions
	attributeObj := mongoDriver.SetCollection(constants.ATTRIBUTES_COLLECTION)
	attributeObj.
		Find(bson.M{"seqId": attrId}).
		Select(
			bson.M{
				"_id":           0,
				"options.seqId": 1,
			}).
		One(&options)

	seqId := 1
	if len(options.Options) > 0 {
		var optionsIds []int
		for _, option := range options.Options {
			optionsIds = append(optionsIds, option.SeqId)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(optionsIds)))
		seqId = seqId + optionsIds[0]
	}
	return seqId
}

func (u InsertOption) addFilterAttributes(attrId int, formData []Option) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, ATTRIBUTE_OPTION_INSERT)
	mongoDriver := factory.GetMongoSession(constants.ATTRIBUTEAPI)
	defer func() {
		logger.EndProfile(profiler, ATTRIBUTE_OPTION_INSERT)
		mongoDriver.Close()
	}()

	var attr Attribute
	mongoDriver.FindOne(constants.ATTRIBUTES_COLLECTION, M{"seqId": attrId}, &attr)
	for _, opt := range formData {
		if (attr.AttributeType == "option" || attr.AttributeType == "multi_option") &&
			opt.FilterAttrName != "" {
			err := u.addFilterAttribute(*attr.Name, opt.FilterAttrName,
				opt.Value, opt.FilterAttrVal, mongoDriver)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (u InsertOption) addFilterAttribute(attrFrom, attrTo,
	attrFromOpt, attrToOpt string, mongoDriver *mongodb.MongoDriver) error {
	query := M{"from": attrFrom,
		"to": attrTo,
	}
	key := fmt.Sprintf("mapping.%s", attrFromOpt)
	updateQuery := M{"$set": M{key: attrToOpt}}
	_, err := mongoDriver.FindAndModify(constants.ATTRIBUTEMAPPING_COLLECTION, updateQuery, query, false)
	if err != nil {
		return fmt.Errorf("No mapping found from %s to %s in %s collection",
			attrFrom, attrTo, constants.ATTRIBUTEMAPPING_COLLECTION)
	}
	return err
}

// ptionCreateWorker is the worker for option creation
func optionCreateWorker(data interface{}) error {
	dataMap, _ := data.(map[string]interface{})
	seqId, _ := dataMap["id"].(int)
	options, _ := dataMap["options"].([]Option)
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
	for x := range options {
		optName := options[x].Value
		pos := options[x].Position
		def := options[x].IsDefault
		optId := options[x].SeqId
		query := "INSERT INTO " + tableName + " (" + idName + ", name, name_en, position, is_default) VALUES (?,?,?,?,?)"
		_, err = txnObj.Exec(query, optId, optName, optName, pos, def)
		if err != nil {
			logger.Error("Exec Failed. Transaction rollback begin.")
			txnObj.Rollback()
			return err
		}
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
