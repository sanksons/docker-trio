package common

import (
	"common/ResourceFactory"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/sqldb"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// This function is worker which updates mysql with sizechart and sku mapping
// ie makes entry in catalog_config_additional_info table
func CreateSizeChartToProductWorker(data interface{}) error {
	dataVal := reflect.ValueOf(data)
	driver, err := ResourceFactory.GetMySqlDriver(SizeChartAPI)

	if err != nil {
		logger.Error(fmt.Sprintf("#CreateSizeChartToProductWorker: Unable to get the mysql driver for MYSQL sync, %s", err.Error()))
		return err
	}

	txnObj, terr := driver.GetTxnObj()

	if terr != nil {
		logger.Error(fmt.Sprintf("#CreateSizeChartToProductWorker: Unable to get mysql transaction object, %s", terr.DeveloperMessage))
		return errors.New(terr.DeveloperMessage)
	}

	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txnObj.Rollback()
		}
	}()

	// iterates over each saved sizechart with corresponding skus
	var categories []int
	for i := 0; i < dataVal.Len(); i++ {
		mysqlSizCh := dataVal.Index(i).Interface().(SizeChartForMysql)
		query, qerr := buildSqlForSkuSizChart(mysqlSizCh, driver, txnObj)
		if qerr != nil {
			logger.Error(fmt.Sprintf("Error while building sql for additional_info table: %s", qerr.Error()))
			continue
		}
		if query == "" { // query already made for this sizechart, so continue
			continue
		}
		logger.Info("QUERY FIRED: " + query)
		err = execOrRollback(txnObj, query)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute query: (%s), Reason is :%s"), query, err.Error())
			continue
		}
		categories = append(categories, mysqlSizCh.CategoryId)
	}
	err = txnObj.Commit()
	if err != nil {
		logger.Error("Commit Failed. Transaction rollback begin.")
		txnObj.Rollback()
		return err
	}
	completeFlag = true
	//memcahe update
	err = ResourceFactory.GetCCApiDriver().UpdateSizecharts(categories...)
	if err != nil {
		logger.Error(err)
	}
	return nil
}

// This function buids sizechart query for skus/ updates the table if sku entry already exists
func buildSqlForSkuSizChart(sizeChart SizeChartForMysql, driver sqldb.SqlDbInterface, txnObj *sql.Tx) (string, error) {
	reward := 0.0
	queryHead := "INSERT into catalog_config_additional_info(" +
		"fk_catalog_config, " +
		"fk_catalog_distinct_sizechart, " +
		"sizechart_type, " +
		"reward_points, " +
		"created_at)Values"
	queryBody := ""

	for _, sku := range sizeChart.Skus {
		config, err := getConfigIdForSku(sku, driver)
		if err != nil {
			continue
		}
		//check if config already exists in table
		ifUpdated, errUp := checkAndUpadateSizChForConfig(config, driver, sizeChart, txnObj)
		if errUp != nil {
			continue
		}
		// 2 means sku existed already and updated succesfully
		// 1 sku existed , but not updated successfully
		// 0 sku doesnot exist in table,entry need to be made
		if ifUpdated == 2 || ifUpdated == 1 {
			continue
		}

		queryBody = queryBody + "( " + strconv.Itoa(config) + ", " +
			strconv.Itoa(sizeChart.SizeChartId) + ", " +
			strconv.Itoa(sizeChart.SizeChartTy) + ", " +
			strconv.FormatFloat(reward, 'f', 1, 64) + ", '" +
			sizeChart.CreatedAt.String() + "'),"
	}
	if queryBody == "" {
		return queryBody, nil
	}
	return queryHead + strings.TrimSuffix(queryBody, ","), nil
}

// This function updates the existing entry of sizechart for sku in db.
func checkAndUpadateSizChForConfig(id int, driver sqldb.SqlDbInterface, sizeChart SizeChartForMysql, txnObj *sql.Tx) (int, error) {
	query := "SELECT * from catalog_config_additional_info where fk_catalog_config = " + strconv.Itoa(id)
	exists := false
	r, err := driver.Query(query)
	if err != nil {
		return 0, fmt.Errorf(err.DeveloperMessage) // O- entry for sku doesnot exists.
	}
	for r.Next() {
		exists = true
	}
	if exists == false {
		return 0, nil
	}
	defer r.Close()

	updateQury := "UPDATE catalog_config_additional_info SET fk_catalog_distinct_sizechart = " +
		strconv.Itoa(sizeChart.SizeChartId) + ", sizechart_type = " + strconv.Itoa(sizeChart.SizeChartTy) +
		" where fk_catalog_config = " + strconv.Itoa(id)
	logger.Info(fmt.Sprintf("Query FIRED : %s", updateQury))
	errQ := execOrRollback(txnObj, updateQury)
	if errQ != nil {
		return 1, errQ
	}
	return 2, nil
}

// This function check if sku exists in system and fetches the config id for same
func getConfigIdForSku(sku string, driver sqldb.SqlDbInterface) (int, error) {
	var id int
	exists := false
	query := "SELECT id_catalog_config from catalog_config where sku = '" + sku + "'"
	result, err := driver.Query(query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occured while getting sku %s for sizechart", sku))
		return -1, errors.New(err.DeveloperMessage)
	}
	for result.Next() {
		exists = true
		err := result.Scan(&id)
		if err != nil {
			logger.Error(fmt.Sprintf("Error occured while getting sku %s for sizechart", sku))
			return -1, err
		}
	}
	result.Close()
	if exists == false {
		return -1, fmt.Errorf(fmt.Sprintf("Sku %s doesnot exists in the system(DB mysql)", sku))
	}
	return id, nil
}

// This worker dumps the sizechart data to mysql tables like
// catalog_distinct_sizechart, catalog_sizechart, catalog_category_has_sizechart_image
func CreateSizeChartWorker(data interface{}) error {
	dataVal := reflect.ValueOf(data)
	driver, err := ResourceFactory.GetMySqlDriver(SizeChartAPI)

	if err != nil {
		logger.Error(fmt.Sprintf("Unable to get the mysql driver for MYSQL sync, %s", err.Error()))
		return err
	}

	txnObj, terr := driver.GetTxnObj()

	if terr != nil {
		logger.Error(fmt.Sprintf("Unable to get mysql transaction object, %s", terr.DeveloperMessage))
		return errors.New(terr.DeveloperMessage)
	}

	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txnObj.Rollback()
		}
	}()

	// Query for one sizechart at a time
	for i := 0; i < dataVal.Len(); i++ {
		sizeMongo := dataVal.Index(i).Interface()
		sizechart := sizeMongo.(SizeChartMongo)
		distinctSizeCh := buildSqlQuerySizechart(sizechart)
		sizeChData := buildSizeChartDataQuery(sizechart)
		sizeChImg := buildSizeChartImageQuery(sizechart)

		logger.Info("QUERY FIRED: " + distinctSizeCh)
		err = execOrRollback(txnObj, distinctSizeCh)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute query: (%s), Reason is :%s"), distinctSizeCh, err.Error())
			return err
		}

		logger.Info("QUERY FIRED: " + sizeChData)
		err = execOrRollback(txnObj, sizeChData)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute query: (%s), Reason is :%s"), sizeChData, err.Error())
			return err
		}

		logger.Info("QUERY FIRED: " + sizeChImg)
		err = execOrRollback(txnObj, sizeChImg)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute query: (%s), Reason is :%s"), sizeChImg, err.Error())
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

func execOrRollback(txnObj *sql.Tx, query string) (err error) {
	_, err = txnObj.Exec(query)
	if err != nil {
		logger.Error("Exec Failed. Transaction rollback begin. Error :", err.Error())
		txnObj.Rollback()
	}
	return err
}

func buildSqlQuerySizechart(data SizeChartMongo) string {
	typeQueryPart := "fk_catalog_ty,"
	typeQueryVal := strconv.Itoa(data.FkCatalogTy) + ", "
	if data.FkCatalogTy == 0 {
		typeQueryPart = ""
		typeQueryVal = ""
	}

	query := "INSERT into catalog_distinct_sizechart(" +
		" id_catalog_distinct_sizechart," +
		" fk_catalog_category," +
		" fk_catalog_brand," + typeQueryPart +
		" sizechart_name," +
		" sizechart_type," +
		" fk_acl_user," +
		" created_at," +
		" updated_at)" +
		" Values(" +
		strconv.Itoa(data.IdCatalogSizeChart) + ", " +
		strconv.Itoa(data.FkCatalogCategory) + ", " +
		strconv.Itoa(data.FkCatalogBrand) + ", " +
		typeQueryVal +
		"'" + strings.Replace(data.SizeChartName, "'", "\\'", -1) + "', " +
		strconv.Itoa(data.SizeChartType) + ", " +
		strconv.Itoa(data.FkAclUser) + ", '" +
		data.CreatedAt.String() + "', '" +
		data.UpdatedAt.String() + "')"
	return query
}

func buildSizeChartDataQuery(data SizeChartMongo) string {
	var queryBody string
	queryHead := "INSERT Into catalog_sizechart(" +
		" fk_catalog_distinct_sizechart," +
		" fk_catalog_category," +
		" fk_catalog_brand," +
		" brand," +
		" column_header," +
		" row_header_name," +
		" row_header_type," +
		" value," +
		" created_at) VALUES"

	for _, sizeData := range data.SizeChartInfo {

		queryBody = queryBody + "( " +
			strconv.Itoa(data.IdCatalogSizeChart) + ", " +
			strconv.Itoa(data.FkCatalogCategory) + ", " +
			strconv.Itoa(data.FkCatalogBrand) + ", '" +
			strings.Replace(sizeData.Brand, "'", "\\'", -1) + "', '" +
			sizeData.ColumnHeader + "', '" +
			sizeData.RowHeaderName + "', '" +
			sizeData.RowHeaderType + "', '" +
			sizeData.Value + "', '" +
			data.CreatedAt.String() + "'),"
	}
	return queryHead + strings.TrimSuffix(queryBody, ",")
}

func buildSizeChartImageQuery(data SizeChartMongo) string {
	var query string

	query = "INSERT into catalog_category_has_sizechart_image(" +
		" fk_catalog_distinct_sizechart, " +
		" fk_catalog_category, " +
		" fk_catalog_brand, " +
		" brand, " +
		" image_path) VALUES(" +
		strconv.Itoa(data.IdCatalogSizeChart) + ", " +
		strconv.Itoa(data.FkCatalogCategory) + ", " +
		strconv.Itoa(data.FkCatalogBrand) + ", '" +
		strings.Replace(data.SizeChartInfo[0].Brand, "'", "\\'", -1) + "', '" +
		strings.Replace(data.SizeChartImagePath, "'", "\\'", -1) + "')"
	return query
}
