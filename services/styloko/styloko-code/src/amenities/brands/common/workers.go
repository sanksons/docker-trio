package common

import (
	"common/ResourceFactory"
	"common/constants"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//Worker that updates MySQL based on changes in Mongo
func BrandUpdateWorker(args interface{}) error {
	driver, err := ResourceFactory.GetMySqlDriver(constants.BRAND_UPDATE)
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

	dataMap, _ := args.(map[string]interface{})
	brand, _ := dataMap["brandInfo"].(Brand)

	err = brandUpdateQuery(brand, txnObj)
	if err != nil {
		logger.Error("#BrandUpdateWorker(): Brand Update failed")
		return err
	}
	deleteRelatedBrandEntriesQuery := deletedRelatedBrands(brand.SeqId)
	relatedBrandIdsUpdateQuery := relatedBrandIdsUpdate(brand.SeqId, brand.RelatedBrand)

	// Delete previous relations in catalog_related_brand
	logger.Info("QUERY FIRED: " + deleteRelatedBrandEntriesQuery)
	err = execOrRollback(txnObj, deleteRelatedBrandEntriesQuery)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing delete brand entry query :%s", err.Error()))
		return err
	}

	// Insert new relations in catalog_related_brand
	for x := range relatedBrandIdsUpdateQuery {
		logger.Info("QUERY FIRED: " + relatedBrandIdsUpdateQuery[x])
		err = execOrRollback(txnObj, relatedBrandIdsUpdateQuery[x])
		if err != nil {
			logger.Error(fmt.Sprintf("Error while related brand update query :%s", err.Error()))
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
	err = ResourceFactory.GetCCApiDriver().UpdateBrands(brand.SeqId)
	if err != nil {
		logger.Error(err)
	}

	logger.Info(fmt.Sprintf("Brand with ID %d has been updated successfully", brand.SeqId))
	return nil
}

//Worker that updates MySQL based on changes in Mongo
func BrandCreateWorker(args interface{}) error {
	driver, err := ResourceFactory.GetMySqlDriver(constants.BRAND_CREATE)
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

	dataMap, _ := args.(map[string]interface{})
	brand, _ := dataMap["brandInfo"].(Brand)

	err = brandInsertQuery(brand, txnObj)

	if err != nil {
		logger.Error("#BrandCreateWorker(): brand creation failed")
		return err
	}
	// Related Brand queries
	relatedBrandQueries := getRelatedBrandsCreate(brand.SeqId, brand.RelatedBrand)

	// Insert related brands begin
	for x := range relatedBrandQueries {
		logger.Info("QUERY FIRED: " + relatedBrandQueries[x])
		err = execOrRollback(txnObj, relatedBrandQueries[x])
		if err != nil {
			return err
		}
	}

	err = txnObj.Commit()
	if err != nil {
		logger.Error("Brand create failed.")
		logger.Error("Commit Failed. Transaction rollback begin.")
		txnObj.Rollback()
		return err
	}
	completeFlag = true
	//send memcache update
	err = ResourceFactory.GetCCApiDriver().UpdateBrands()
	if err != nil {
		logger.Error(err)
	}

	logger.Info(fmt.Sprintf("Brand with ID %d has been created successfully", brand.SeqId))
	return nil
}
