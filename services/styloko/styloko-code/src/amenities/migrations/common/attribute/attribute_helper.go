package attribute

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"common/xorm/mysql"
	"errors"
	"fmt"
	"time"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
)

func processAttributeSetRows(response interface{}) ([]AttributeSet, error) {
	rows, ok := response.(*core.Rows)
	if ok {
		var AttributeSetArr []AttributeSet
		for rows.Next() {
			s := AttributeSet{}
			err := rows.Scan(
				&s.SeqId,
				&s.Name,
				&s.Label,
				&s.Identifier,
			)
			logger.Debug(fmt.Sprintf("Started processing rows for attribute set %d", s.SeqId))

			if err != nil {
				logger.Error(fmt.Sprintf("Error in processAttributeSetRows() %v", err.Error()))
				rows.Close()
				return nil, err
			}
			s.CreatedAt = time.Now()
			s.UpdatedAt = time.Now()
			if s.Name != "mobile" {
				AttributeSetArr = append(AttributeSetArr, s)
			}

		}
		rows.Close()
		return AttributeSetArr, nil
	}
	logger.Error("Error while assering into *core.Rows")
	return nil, errors.New("Error while assering into *core.Rows")
}

func writeAttributeSetToMongo(AttributeSet []AttributeSet) error {
	logger.Info("Starting to write Attribute Set into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("AttributeMigration")
	defer mgoSession.Close()
	for _, v := range AttributeSet {
		logger.Debug(fmt.Sprintf("Inserting seqId %d in mongo", v.SeqId))
		upsertVal := true
		updatedVal := map[string]interface{}{"$set": v}
		findCriteria := map[string]interface{}{"seqId": v.SeqId}
		_, err := mgoSession.FindAndModify(util.AttributeSets, updatedVal, findCriteria, upsertVal)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Attribute Set into Mongo %v", err.Error()))
			return err
		}
	}
	return nil
}

func processAttributeRows(response interface{}) ([]AttributeRow, error) {
	rows, ok := response.(*core.Rows)
	if ok {
		var AttributeArr []AttributeRow
		var setId int
		var setName string
		var mandatoryField *int
		for rows.Next() {
			s := AttributeRow{}
			err := rows.Scan(
				&s.SeqId,
				&setName,
				&setId,
				&s.Name,
				&s.IsGlobal,
				&s.ProductType,
				&s.AttributeType,
				&s.MaxLength,
				&s.DecimalPlaces,
				&s.DefaultValue,
				&s.UniqueValue,
				&s.AliceExport,
				&s.DwhFieldName,
				&s.SolrSearchable,
				&s.SolrFilter,
				&s.SolrSuggestions,
				&s.ExportPosition,
				&s.ExportName,
				&s.ImportName,
				&s.ImportConfigIdentifier,
				&s.Label,
				&s.Description,
				&s.PetMode,
				&s.PetType,
				&s.PetOverview,
				&s.PetOverviewFilter,
				&s.PetQc,
				&s.Validation,
				&mandatoryField,
				&s.MandatoryImport,
				&s.Extra,
				&s.Visible,
				&s.CreatedAt,
				&s.UpdatedAt,
				&s.DataType,
			)
			logger.Debug(fmt.Sprintf("Started processing rows for attribute %d", s.SeqId))

			if s.IsGlobal == 0 {
				s.Set.SeqId = setId
				s.Set.Name = setName
			}

			if err != nil {
				logger.Error(fmt.Sprintf("Error in processAttributeRows() %v for id %d", err.Error(), s.SeqId))
				rows.Close()
				return nil, err
			}

			if mandatoryField == nil || *mandatoryField == 0 {
				s.Mandatory = false
			} else {
				s.Mandatory = true
			}

			// Tax Rate will be mandatory
			if s.Name == "tax_rate" {
				s.Mandatory = true
			}

			if *s.PetType == "checkbox" {
				str := "integer"
				s.Validation = &str
				defVal := "0"
				s.DefaultValue = &defVal
			}

			if *s.PetType == "numberfield" {
				if s.Validation == nil {
					logger.Error(fmt.Sprintf("Attribute id %d having name %s does not have a valid validation, validation is null", s.SeqId, s.Name))
					str := "decimal"
					s.Validation = &str
				}
			}

			if s.Validation != nil && *s.Validation == "email" {
				str := "letters"
				s.Validation = &str
			}

			if s.CreatedAt == nil {
				time := time.Now()
				s.CreatedAt = &time
			}

			if s.UpdatedAt == nil {
				time := time.Now()
				s.UpdatedAt = &time
			}

			s.IsActive = 1
			if s.AttributeType == `multi_option` || s.AttributeType == `option` {
				attrOpts := getOptions(setName, s.Name)
				s.Options = attrOpts
			}

			AttributeArr = append(AttributeArr, s)
		}
		rows.Close()
		return AttributeArr, nil
	}
	logger.Error("Error while assering into *core.Rows")
	return nil, errors.New("Error while assering into *core.Rows")
}

func writeAttributeToMongo(Attribute []AttributeRow) error {
	logger.Info("Starting to write Attribute into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("AttributeMigration")
	defer mgoSession.Close()
	for _, v := range Attribute {
		logger.Debug(fmt.Sprintf("Inserting seqId %d in mongo", v.SeqId))
		upsertVal := true
		updatedVal := map[string]interface{}{"$set": v}
		findCriteria := map[string]interface{}{"seqId": v.SeqId}
		_, err := mgoSession.FindAndModify(util.Attributes, updatedVal, findCriteria, upsertVal)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Attribute into Mongo %v", err.Error()))
			return err
		}
	}
	return nil
}

func getOptions(attrSetName string, attrName string) []Option {
	var attrOpts []Option
	zeroSeqExist := false
	count := 0
	optionNo := 0
	var sql string
	sql = getAttributeOptionSql(attrSetName, attrName)
	logger.Info(fmt.Sprintf("Attribute option query being executed is %s", sql))
	response, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in running query %v", err.Error()))
		return nil
	}
	rows, ok := response.(*core.Rows)
	if ok {
		for rows.Next() {
			var attrOpt Option
			err = rows.Scan(
				&attrOpt.SeqId,
				&attrOpt.Value,
				&attrOpt.Position,
				&attrOpt.IsDefault,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in reading row to struct for Attribute Options %v for id %d", err.Error(), attrOpt.SeqId))
				rows.Close()
				return nil
			}
			if attrOpt.SeqId == 0 {
				optionNo = count
				zeroSeqExist = true
			}
			attrOpts = append(attrOpts, attrOpt)
			count++
		}
		if zeroSeqExist {
			attrOpts[optionNo].SeqId = count
		}
		rows.Close()
	}
	if len(attrOpts) == 0 {
		defVal := "default"
		zero := 0
		opt := Option{1, &defVal, &zero, &zero}
		attrOpts = append(attrOpts, opt)
		sql := updateAttributeOption(attrSetName, attrName, opt)
		logger.Info(fmt.Sprintf("Attribute option query being executed is %s", sql))
		r, err := mysql.GetInstance().Query(sql, true)
		rows, ok := r.(*core.Rows)
		if ok {
			rows.Close()
		}
		if err != nil {
			logger.Error(fmt.Sprintf("Error in running query %v", err.Error()))
			return nil
		}
	}
	return attrOpts
}

//This function creates indexes for the collection to be created if the collection does not exist
func checkAndEnsureIndex() error {
	flag := false
	mgoSession := ResourceFactory.GetMongoSession("AttributeMigration")
	defer mgoSession.Close()
	logger.Info("Checking if collection already exists")
	colNames, err := mgoSession.CollectionExists()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting collection names from mongo :%s", err.Error()))
		return err
	}
	for _, v := range colNames {
		if v == util.Attributes {
			flag = true
		}
	}
	if flag == true {
		logger.Info("Collection already exists so skipping creating indexes")
		return nil
	}
	EnsureIndexInDb()
	return nil
}

func EnsureIndexInDb() {
	logger.Info("Creating Indexes for new collection to be created")
	mgoSession := ResourceFactory.GetMongoSession("AttributeMigration")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(util.Attributes)
	var uniqueIndexes = []string{
		"seqId",
	}
	var normalIndexes = []string{
		"name",
	}
	for _, v := range normalIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: false,
			Sparse: false,
		})
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	for _, v := range uniqueIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: true,
			Sparse: true,
		})
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	logger.Info("New indexes created")
}
