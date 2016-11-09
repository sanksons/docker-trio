package product

import (
	factory "common/ResourceFactory"
	dbsql "database/sql"
	"fmt"
)

//
// Delete Product based on the supplied ID.
//
func DeleteProductById(configId int) error {
	//initiate global map.
	AttributesCache = make(map[string][]*CatalogAttribute)
	//get attribute set details
	attrSet, err := getAttributeSet(configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(1):%s", err.Error())
	}
	//get attributes of attribute set
	configAttributes, err := getAttributesOfSet(attrSet.Id, "config")
	if err != nil {
		return fmt.Errorf("DeleteProductById(2):%s", err.Error())
	}
	//get mysql driver
	driver, err := factory.GetMySqlDriver("Default")
	if err != nil {
		return fmt.Errorf("DeleteProductById(3):%s", err.Error())
	}
	//initiate transaction
	txn, terr := driver.GetTxnObj()
	if terr != nil {
		return fmt.Errorf("DeleteProductById(4):%s", terr.DeveloperMessage)
	}
	defer func() {
		txn.Rollback()
	}()

	// Removal Process starts.

	// Remove all multiOPtion type attributes from link table
	for _, attr := range configAttributes {
		if attr.attributeType == "multi_option" {
			err := deleteMultiAttribute(configId, *attr, txn)
			if err != nil {
				return fmt.Errorf("DeleteProductById(5):%s", err.Error())
			}
		}
	}
	//remove link from attributeset table.
	sql := fmt.Sprintf(
		"DELETE FROM catalog_attribute_link_%s WHERE fk_catalog_config=?",
		attrSet.Name,
	)
	_, err = txn.Exec(sql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(6a):%s", err.Error())
	}
	//remove from attributeset table.
	sql = fmt.Sprintf(
		"DELETE FROM catalog_config_%s WHERE fk_catalog_config=?",
		attrSet.Name,
	)
	_, err = txn.Exec(sql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(6b):%s", err.Error())
	}

	//prepare simples deletion
	sql = "SELECT id_catalog_simple from catalog_simple where fk_catalog_config=?"
	result, err := txn.Query(sql, configId)
	if err != nil && err != dbsql.ErrNoRows {
		return fmt.Errorf("DeleteProductById(7):%s", err.Error())
	}
	var simpleIds []int
	for result.Next() {
		var simpleId int
		err := result.Scan(&simpleId)
		if err != nil && err != dbsql.ErrNoRows {
			return fmt.Errorf("DeleteProductById(8):%s", err.Error())
		}
		simpleIds = append(simpleIds, simpleId)
	}
	result.Close()

	for _, simpleId := range simpleIds {
		sql := fmt.Sprintf("DELETE FROM catalog_simple_%s WHERE fk_catalog_simple=?", attrSet.Name)
		_, err = txn.Exec(sql, simpleId)
		if err != nil {
			return fmt.Errorf("DeleteProductById(9):%s", err.Error())
		}
	}
	//delete all simples
	simpleSql := "DELETE FROM catalog_simple where fk_catalog_config = ?;"
	_, err = txn.Exec(simpleSql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(10):%s", err.Error())
	}
	//delete images
	imageSql := "DELETE FROM catalog_product_image where fk_catalog_config= ?;"
	_, err = txn.Exec(imageSql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(11):%s", err.Error())
	}
	//delete category
	catSql := "DELETE FROM catalog_config_has_catalog_category where fk_catalog_config= ?;"
	_, err = txn.Exec(catSql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(12):%s", err.Error())
	}
	//delete config
	configSql := "DELETE FROM catalog_config where id_catalog_config = ?;"
	_, err = txn.Exec(configSql, configId)
	if err != nil {
		return fmt.Errorf("DeleteProductById(13):%s", err.Error())
	}
	txn.Commit()
	return nil
}

func getAttributeSet(configId int) (ProAttributeSet, error) {
	driver, err := factory.GetMySqlDriver("Default")
	if err != nil {
		return ProAttributeSet{}, err
	}
	txn, terr := driver.GetTxnObj()
	if terr != nil {
		return ProAttributeSet{}, fmt.Errorf(terr.DeveloperMessage)
	}
	defer txn.Rollback()
	var attrsetId int
	sql := `SELECT fk_catalog_attribute_set FROM catalog_config 
            where id_catalog_config=?`
	row := txn.QueryRow(sql, configId)
	err = row.Scan(&attrsetId)
	if err != nil {
		return ProAttributeSet{}, fmt.Errorf("getAttributeSet():%s", err.Error())
	}
	var attrSet ProAttributeSet
	sql = `SELECT id_catalog_attribute_set, name, label 
			FROM catalog_attribute_set where id_catalog_attribute_set=?`
	row = txn.QueryRow(sql, attrsetId)
	err = row.Scan(&attrSet.Id, &attrSet.Name, &attrSet.Label)
	if err != nil {
		return attrSet, fmt.Errorf("getAttributeSet():%s", err.Error())
	}
	return attrSet, nil
}

func deleteMultiAttribute(id int, attr CatalogAttribute, txn *dbsql.Tx) error {

	tableName := fmt.Sprintf("catalog_attribute_link_%s_%s", attr.attrSetName, attr.name)
	var where string
	var condition int
	if attr.attrSetName == "global" {
		where = "fk_catalog_config"
		condition = id
	} else {
		selecttableName := fmt.Sprintf("catalog_config_%s", attr.attrSetName)
		filedname := fmt.Sprintf("id_catalog_config_%s", attr.attrSetName)
		sql := "SELECT " + filedname + " FROM " + selecttableName + " WHERE fk_catalog_config=?"
		var reqId int
		row := txn.QueryRow(sql, id)
		err := row.Scan(&reqId)
		if err != nil {
			return fmt.Errorf("deleteMultiAttribute(1):%s", err.Error())
		}
		where = fmt.Sprintf("fk_catalog_config_%s", attr.attrSetName)
		condition = reqId
	}
	delSql := "DELETE FROM " + tableName + " WHERE " + where + "=?"
	_, err := txn.Exec(delSql, condition)
	if err != nil {
		return fmt.Errorf("deleteMultiAttribute(2):%s", err.Error())
	}
	return nil
}
