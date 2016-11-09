package common

import (
	"common/ResourceFactory"
	"common/appconfig"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
	"time"
)

//This function accepts a []map[string]interface{} and depending on the command attached,
//calls the InsertInMysql or UpdateInMysql,
func SyncWithMysql(dataMap interface{}) error {
	dataMapArr := dataMap.([]map[string]interface{})
	data, err := utils.ConvertStructArrToMapArr(dataMapArr[0]["data"])
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting data []struct to []map[string]interface{} :%s", err.Error()))
		return err
	}
	command := dataMapArr[0]["command"]
	for _, v := range data {
		if command == SELLER_INSERT {
			err := InsertInMysql(v)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while inserting in MySql :%s", err.Error()))
				return err
			}
		}
		if command == SELLER_UPDATE {
			err := UpdateInMysql(v)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while updating in MySql :%s", err.Error()))
				return err
			}
		}
	}
	logger.Info("SYNC DONE.")
	return nil
}

//This function inserts the correctly inserted seller details into respective columns in mysql
func InsertInMysql(dataMap map[string]interface{}) error {
	driver, errs := ResourceFactory.GetMySqlDriver(SELLER_INSERT)
	if errs != nil || driver == nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", errs.Error()))
		return errs
	}
	txObj, err := driver.GetTxnObj()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting transaction object: %s", err.DeveloperMessage))
		return errors.New("Error while getting transaction object")
	}

	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txObj.Rollback()
		}
	}()

	catSupSql, status := getInsertCatalogSupplierQuery(dataMap)
	logger.Info(fmt.Sprintf("catalog_supplier sql being fired:%s", catSupSql))
	r, er := txObj.Exec(catSupSql, int(dataMap["seqId"].(float64)), dataMap["ordrEml"], dataMap["orgName"], dataMap["slrName"], dataMap["phn"], status, `NOW()`, `NOW()`)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while executing catalog supplier query %s", er.Error()))
		return er
	}
	//execute and get last inserted id
	supInsId, _ := r.LastInsertId()
	countrySql := getInsertCountryQuery(dataMap)
	logger.Info(fmt.Sprintf("country sql being fired:%s", countrySql))
	res, e := txObj.Exec(countrySql, dataMap["cntryCode"])
	if e != nil {
		logger.Error(fmt.Sprintf("Error while executing country query %s", e.Error()))
		return e
	}
	//execute and get last inserted id
	countryInsId, _ := res.LastInsertId()
	supAddrSql := getInsertSupplierAddressQuery(dataMap)
	logger.Info(fmt.Sprintf("supplier_address sql being fired:%s", supAddrSql))
	_, serr := txObj.Exec(supAddrSql, dataMap["addr1"], dataMap["city"], int(dataMap["pstcode"].(float64)), `NOW()`, `NOW()`, supInsId, countryInsId)
	if serr != nil {
		logger.Error(fmt.Sprintf("Error while executing supplier address query %s", serr.Error()))
		return serr
	}
	cerr := txObj.Commit()
	if cerr != nil {
		txObj.Rollback()
	}
	completeFlag = true
	return nil
}

//This function updates the correctly updated seller details into respective columns in mysql
func UpdateInMysql(dataMap map[string]interface{}) error {
	driver, derr := ResourceFactory.GetMySqlDriver(SELLER_UPDATE)
	if derr != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", derr.Error()))
		return derr
	}
	txObj, err := driver.GetTxnObj()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting transaction object: %s", err.DeveloperMessage))
		return errors.New("Error while getting transaction object")
	}

	completeFlag := false
	defer func() {
		if !completeFlag {
			logger.Error("Transaction has failed probably due to panic. Rollback begins.")
			txObj.Rollback()
		}
	}()

	catSupSql := getUpdateCatalogSupplierQuery(int(dataMap["seqId"].(float64)))
	logger.Info(fmt.Sprintf("catalog_supplier sql being fired:%s", catSupSql))
	_, er := txObj.Exec(catSupSql, dataMap["ordrEml"], dataMap["orgName"], dataMap["slrName"], dataMap["phn"], dataMap["status"], `NOW()`)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while executing catalog supplier query %s", er.Error()))
		return er
	}
	serr := txObj.Commit()
	if serr != nil {
		txObj.Rollback()
	}

	completeFlag = true
	return nil
}

//generates insert query for catalog_supplier table
func getInsertCatalogSupplierQuery(dataMap map[string]interface{}) (string, string) {
	var status string
	if dataMap["status"] == nil {
		status = "inactive"
	}
	suppCol := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s", "id_catalog_supplier", "order_email", "name_en", "name", "phone", "status", "updated_at", "created_at")
	return fmt.Sprintf("%s (%s) VALUES (?,?,?,?,?,?,?,?)", "INSERT INTO catalog_supplier", suppCol), status
}

//generates insert query for country table
func getInsertCountryQuery(dataMap map[string]interface{}) string {
	return fmt.Sprintf("%s (%s) VALUES (?)", "INSERT INTO country ", "iso2_code")
}

//generates insert query for supplier_address table
func getInsertSupplierAddressQuery(dataMap map[string]interface{}) string {
	//append both the inserted ids as relevant fks
	addrCol := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s", "street", "city", "postcode", "updated_at", "created_at", "fk_id_catalog_supplier", "fk_country")
	return fmt.Sprintf("%s (%s) VALUES (?,?,?,?,?,?,?)", "INSERT INTO supplier_address", addrCol)
}

//generates update query for catalog_supplier table
func getUpdateCatalogSupplierQuery(suppId int) string {
	catSupSql := `UPDATE catalog_supplier SET order_email=?,name_en=?,name=?,phone=?,status=?,updated_at=?`
	return fmt.Sprintf("%s WHERE id_catalog_supplier = %d", catSupSql, suppId)
}

//This function invalidates products cache for updated sellers by
// preparing data and calling styloko products api
func InvalidateProductsForUpdatedSellers(data interface{}) error {
	res := data.([]Schema)
	slrData := prepareSellerData(res)
	err := sendDataToProdApi(slrData)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending data to product API :%s", err.Error()))
		return err
	}
	return nil
}

//This function prepeares seller data for consumption by styloko products api
func prepareSellerData(res []Schema) []SellerData {
	slrData := make([]SellerData, 0)
	for k, _ := range res {
		tmp := SellerData{}
		tmp.Value = strconv.Itoa(res[k].SeqId)
		tmp.Type = "seller"
		slrData = append(slrData, tmp)
	}
	return slrData
}

//This function sends passed data to styloko product api
func sendDataToProdApi(slrData []SellerData) error {
	jsonData, err := json.Marshal(slrData)
	if err != nil {
		logger.Error("Error while marshalling request", err)
		return err
	}
	logger.Info(fmt.Sprintf("Json Data sent to Orchestrator:%s", string(jsonData)))
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	headers := make(map[string]string)
	headers["Update-Type"] = "Cache"
	//read time from conf
	t, err := time.ParseDuration(config.Styloko.Timeout)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting timeout to time.Duration: %s", err.Error()))
		return err
	}
	resp, err := http.HttpPut(config.Styloko.Url, headers, string(jsonData), t*time.Millisecond)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending POST Request %s", err.Error()))
		return err
	}
	logger.Info(fmt.Sprintf("Response from Styloko Product API %s", string(resp.Body)))
	return nil
}
