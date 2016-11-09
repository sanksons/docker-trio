package common

import (
	"common/ResourceFactory"
	"common/constants"
	"errors"
	"fmt"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// CategoryUpdateWorker -> Worker that updates MySQL based on changes in Mongo
func CategoryUpdateWorker(args interface{}) error {
	driver, err := ResourceFactory.GetMySqlDriver(constants.CATEGORY_UPDATE)
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

	queryPrimary, id, segIds, err := categoryUpdateQuery(args)
	segDelete, segUpdate := getSegments(segIds, id)

	// Primary Update category Query
	logger.Info("QUERY FIRED: " + queryPrimary)
	err = execOrRollback(txnObj, queryPrimary)
	if err != nil {
		return err
	}

	// Delete query for previous relations
	logger.Info("QUERY FIRED: " + segDelete)
	err = execOrRollback(txnObj, segDelete)
	if err != nil {
		return err
	}

	// Insert new relations
	for x := range segUpdate {
		logger.Info("QUERY FIRED: " + segUpdate[x])
		err = execOrRollback(txnObj, segUpdate[x])
		if err != nil {
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
	//send memcache update
	err = ResourceFactory.GetCCApiDriver().UpdateCategories(id)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(fmt.Sprintf("Category with ID %d has been updated successfully", id))

	return nil
}

// CategoryCreateWorker -> Worker that updates MySQL based on changes in Mongo
// TODO This function is untested. Only pushed to be a reference point to others.
func CategoryCreateWorker(args interface{}) error {
	driver, err := ResourceFactory.GetMySqlDriver(constants.CATEGORY_CREATE)
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
	insertCategory, segmentQueries, rightQuery, leftQuery, id := categoryInsertQuery(args)

	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txnObj.Rollback()
		}
	}()

	// Right value modification
	logger.Info("QUERY FIRED: " + rightQuery)
	err = execOrRollback(txnObj, rightQuery)
	if err != nil {
		return err
	}

	// Left value modification
	logger.Info("QUERY FIRED: " + leftQuery)
	err = execOrRollback(txnObj, leftQuery)
	if err != nil {
		return err
	}

	// Insert Category begin
	logger.Info("QUERY FIRED: " + insertCategory)
	err = execOrRollback(txnObj, insertCategory)
	if err != nil {
		return err
	}

	// Insert segments begin
	for x := range segmentQueries {
		logger.Info("QUERY FIRED: " + segmentQueries[x])
		err = execOrRollback(txnObj, segmentQueries[x])
		if err != nil {
			return err
		}
	}

	err = txnObj.Commit()
	if err != nil {
		logger.Error("Category create failed.")
		logger.Error("Commit Failed. Transaction rollback begin.")
		txnObj.Rollback()

		return err
	}
	completeFlag = true
	err = ResourceFactory.GetCCApiDriver().UpdateCategories(id)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(fmt.Sprintf("Category with ID %d has been created successfully", id))
	return nil
}
