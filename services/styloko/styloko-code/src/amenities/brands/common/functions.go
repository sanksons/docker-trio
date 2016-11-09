package common

import (
	"common/ResourceFactory"
	"common/constants"
	"database/sql"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
	s "strings"
)

//creates update query for the Brand
func brandUpdateQuery(brand Brand, txnObj *sql.Tx) error {
	updateQuery := `UPDATE catalog_brand SET status = ?, name = ?, position = ?,
	url_key = ?, image_name = ?, brand_class = ?, is_exclusive = ?, brand_info = ?
	WHERE id_catalog_brand = ?`
	_, err := txnObj.Exec(updateQuery, brand.Status, brand.Name, brand.Position,
		brand.UrlKey, brand.ImageName, brand.BrandClass, brand.IsExclusive,
		brand.BrandInfo, brand.SeqId)
	if err != nil {
		return fmt.Errorf("#brandUpdateQuery(): %s", err.Error())
	}
	return nil
}

//deletes existing related brands against corresponding id
func deletedRelatedBrands(id int) string {
	deleteRelatedBrandEntriesQuery := "DELETE FROM catalog_related_brand WHERE fk_catalog_brand=" + string(strconv.Itoa(id))
	return deleteRelatedBrandEntriesQuery
}

//inserts related brands added/created against the corresponding id
func relatedBrandIdsUpdate(id int, relatedBrandIds []RelatedBrand) (updateRelatedBrandQuery []string) {
	initialPart := "INSERT INTO catalog_related_brand (fk_catalog_brand,id_related_brand) VALUES"
	for x, _ := range relatedBrandIds {
		tmp := initialPart + "(" + strconv.Itoa(id) + "," + strconv.Itoa(relatedBrandIds[x].IdRelatedBrand) + ")"
		updateRelatedBrandQuery = append(updateRelatedBrandQuery, tmp)
	}
	return updateRelatedBrandQuery
}

//returns array of related brands queries.
// No delete queries in this scenario
func getRelatedBrandsCreate(id int, relatedBrandIds []RelatedBrand) (updateQuery []string) {
	insertQuery := "INSERT INTO catalog_related_brand (fk_catalog_brand, id_related_brand) VALUES ("
	for x, _ := range relatedBrandIds {
		tmp := insertQuery + strconv.Itoa(id) + "," + strconv.Itoa(relatedBrandIds[x].IdRelatedBrand) + ")"
		updateQuery = append(updateQuery, tmp)
	}

	return updateQuery
}

//returns everything i.e brand create query, related brand queries, brand certificate queries.
//One must fire them all and then commit the transaction.
func brandInsertQuery(brand Brand, txnObj *sql.Tx) error {
	//insert Brand basic info
	brandInfoQuery := `INSERT into catalog_brand(id_catalog_brand, status, name, name_en,
	position, url_key, image_name, brand_class, is_exclusive, brand_info) VALUES
	(?,?,?,?,?,?,?,?,?,?)`
	_, err := txnObj.Exec(brandInfoQuery, brand.SeqId, brand.Status, brand.Name, s.ToLower(brand.Name),
		brand.Position, brand.UrlKey, brand.ImageName, brand.BrandClass, brand.IsExclusive,
		brand.BrandInfo)
	if err != nil {
		return fmt.Errorf("#brandInsertQuery(): %s", err.Error())
	}
	return nil
}

//Executes the query, if it fails, then db rollback is called.
func execOrRollback(txnObj *sql.Tx, query string) (err error) {
	_, err = txnObj.Exec(query)
	if err != nil {
		logger.Error("Exec Failed. Transaction rollback begin.")
		txnObj.Rollback()
	}
	return err
}

//takes map[string]interface{} and checks whether the string key exists
//in mongo or not,returns true if it exists, false otherwise along with the found object.
func CheckIfKeyExists(bsonMap map[string]interface{}) (bool, []Brand, error) {
	mgoSession := ResourceFactory.GetMongoSession(BRAND_OPERATION)
	mgoObj := mgoSession.SetCollection(constants.BRAND_COLLECTION)
	defer mgoSession.Close()
	brandData := []Brand{}
	if val, ok := bsonMap["seqId"]; ok {
		val = val.(int)
	}
	err := mgoObj.Find(bsonMap).All(&brandData)
	if len(brandData) != 0 {
		return true, brandData, err
	}
	return false, brandData, nil
}
