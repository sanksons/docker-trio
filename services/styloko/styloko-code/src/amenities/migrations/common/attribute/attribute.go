package attribute

import (
	"amenities/migrations/common/util"
	"common/xorm/mysql"
	"fmt"
	"strconv"

	"github.com/jabong/floRest/src/common/utils/logger"
)

func StartAttributeSetMigration() error {
	err := AttributeSetTable()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while migrating attribute set table :%s", err.Error()))
		return err
	}
	return nil
}

func StartAttributeMigration() error {
	err := checkAndEnsureIndex()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while checking existing index and ensuring indexes for db :%s", err.Error()))
		return err
	}
	err = AttributeTable()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while migrating attribute table :%s", err.Error()))
		return err
	}
	return nil
}

func AttributeSetTable() error {
	logger.Info("Started Migrating Attribute Set Table")
	sql := getAttributeSetSql()
	logger.Debug(fmt.Sprintf("Attribute set query being executed is %s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing attribute set query %s", err.Error()))
		return err
	}
	attributeSetRows, err := processAttributeSetRows(rows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while processing attribute set rows %s", err.Error()))
		return err
	}
	err = writeAttributeSetToMongo(attributeSetRows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while writing attribute set to mongo :%s", err.Error()))
		return err
	}
	return nil
}

func AttributeTable() error {
	logger.Info("Started Migrating Attribute Table")
	sql := getAttributeSql()
	logger.Debug(fmt.Sprintf("Attribute query being executed is %s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing attributes query %s", err.Error()))
		return err
	}

	attributeRows, err := processAttributeRows(rows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while processing attribute rows %s", err.Error()))
		return err
	}

	err = writeDataInChunks(attributeRows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while writing data in chunks :%s", err.Error()))
		return err
	}
	return nil
}

func writeDataInChunks(attributeRows []AttributeRow) error {
	limit := util.Chunk
	for i := 0; i < len(attributeRows); {
		if limit > len(attributeRows) {
			limit = len(attributeRows)
		}
		var rows []AttributeRow
		for j := i; j <= limit; j++ {
			if j >= len(attributeRows) {
				break
			}
			rows = append(rows, attributeRows[j])
		}
		logger.Info(fmt.Sprintf("Prepare chunk data for %d attributes", len(rows)))
		err := writeAttributeToMongo(rows)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while writing attribute to mongo :%s", err.Error()))
			return err
		}
		i = limit + 1
		limit = limit + util.Chunk
	}
	return nil
}

func MigrateSingleAttribute(id int) error {
	logger.Info("Started Migrating Attribute by Id")
	idStr := strconv.Itoa(id)
	sql := getAttributeIdSql(idStr)
	logger.Debug(fmt.Sprintf("Attribute by Id query being executed is %s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing attributes by Id query %s", err.Error()))
		return err
	}
	attributeRows, err := processAttributeRows(rows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while processing attribute by Id rows %s", err.Error()))
		return err
	}
	err = writeAttributeToMongo(attributeRows)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while writing attribute by Id to mongo :%s", err.Error()))
		return err
	}
	return nil
}
