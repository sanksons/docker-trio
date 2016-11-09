package supplier

import (
	"common/ResourceFactory"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//function that integrates supplier details,drops existing db if any,
//transforms it and writes it to mongo with counter info
func StartSupplierMigration() error {
	logger.Info("Started Migrating Supplier Table")
	err := initializeCacheObj()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while initialzing cache object: %s", err.Error()))
		return err
	}
	orgInfo, err := getSupplierInfo()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting supplier info: %s", err.Error()))
		return err
	}
	org := TransformSupplier(orgInfo)
	err = checkAndEnsureIndex()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while checking existing index and ensuring indexes for db :%s", err.Error()))
		return err
	}
	err = breakInChunks(org)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in calling breakInChunks :%s", err.Error()))
		return err
	}
	counter, err := getCounter()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting counter :%s", err.Error()))
		return err
	}
	err = updateCounter(counter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while updating counter :%s", err.Error()))
		return err
	}
	logger.Info("Supplier Migration Succesful")
	return nil
}

//function that gets supplier information from catalog_supplier,
//supplier_address and country tables
func getSupplierInfo() ([]OrgMongo, error) {
	logger.Info("Preparing Supplier Info")
	sql := getSupplierSql()
	logger.Debug(fmt.Sprintf("Sql being fired to get supplier info :%s", sql))
	driver, derr := ResourceFactory.GetMySqlDriver(SUPPLIER_MIGRATION)
	if derr != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", derr.Error()))
		return nil, derr
	}
	rows, err := driver.Query(sql)
	if err == nil {
		orgInfo, er := processSupplierRows(rows)
		if er == nil {
			return orgInfo, nil
		}
		return nil, er
	}
	return nil, errors.New(fmt.Sprintf("Error while executing query :%s", err.DeveloperMessage))
}
