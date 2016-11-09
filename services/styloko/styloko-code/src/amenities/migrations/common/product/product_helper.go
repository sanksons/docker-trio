package product

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"common/xorm/mysql"
	"database/sql"
	"errors"
	"github.com/go-xorm/core"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func GetProductsFromMysql(
	offset string,
	limit string,
	criteria ProductFetchCriteria,
) (interface{}, error) {

	var sqlQ string
	switch criteria.GetCriteriaType() {
	case CRITERIA_TYPE_SELLER:
		sqlQ = `SELECT * FROM catalog_config
	    WHERE id_catalog_config IS NOT NULL
	    AND fk_catalog_supplier=` + strconv.Itoa(criteria.SellerId) + `
	     LIMIT ` + offset + `,` + limit
	case CRITERIA_TYPE_STATUS:
		if criteria.Status == "active" {
			sqlQ = `SELECT * FROM catalog_config
	    WHERE id_catalog_config IS NOT NULL
	    AND status = 'active' AND pet_approved=1
	     LIMIT ` + offset + `,` + limit
		} else {
			sqlQ = `SELECT * FROM catalog_config
	    WHERE id_catalog_config IS NOT NULL
	    AND (status != 'active') OR (status = 'active' AND pet_approved=0)
	     LIMIT ` + offset + `,` + limit
		}
	case CRITERIA_TYPE_PARTIAL:
		sqlQ = `SELECT * FROM catalog_config 
			WHERE id_catalog_config >= ` + strconv.Itoa(criteria.MinId) + ` AND 
			id_catalog_config <= ` + strconv.Itoa(criteria.MaxId) + `
	     LIMIT ` + offset + `,` + limit
	case CRITERIA_TYPE_BRAND:
		sqlQ = `SELECT * FROM catalog_config
	    WHERE id_catalog_config IS NOT NULL
	    AND fk_catalog_brand=` + strconv.Itoa(criteria.BrandId) + `
	     LIMIT ` + offset + `,` + limit
	case CRITERIA_TYPE_PROMOTION:
		sqlQ = `SELECT * FROM catalog_config cc JOIN 
		catalog_attribute_link_global_promotion cp ON
		cc.id_catalog_config = cp.fk_catalog_config
	    WHERE cc.id_catalog_config IS NOT NULL
	    AND cp.fk_catalog_attribute_option_global_promotion =` + strconv.Itoa(criteria.PromotionId) + `
	     LIMIT ` + offset + `,` + limit
	default:

	}
	printInfo(sqlQ)
	rows, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExhausted
		} else {
			return nil, errors.New("GetProductsFromMysql(): " + err.Error())
		}
	}
	return rows, nil
}

func ProcessProduct(product generalMap) ProductMigrationStatus {
	p := &ProductMongo{}
	st := ProductMigrationStatus{}
	var ok bool
	var err error
	err = func() error {
		p.Id, ok = dbToInt(product["id_catalog_config"])
		if !ok {
			return errors.New("ProcessProduct(): Cannot set productID")
		}
		//Do not reset product_set if its SC product.
		if p.Id > util.ProductStartIndex && p.Id < util.ProductEndIndex {
			p.ProductSet = 0
		}
		printInfo("Started: Processing product (" + strconv.Itoa(p.Id) + ")")
		attributeSetId, ok := dbToInt(product["fk_catalog_attribute_set"])
		if !ok {
			return errors.New("ProcessProduct(): Cannot set AttributeSetID")
		}
		p.Name, ok = dbToString(product["name"])
		if !ok {
			return errors.New("ProcessProduct(): Cannot set Name")
		}
		p.SKU, ok = dbToString(product["sku"])
		if !ok {
			return errors.New("ProcessProduct(): Cannot set SKU")
		}
		p.BrandId, ok = dbToInt(product["fk_catalog_brand"])
		if !ok {
			return errors.New("ProcessProduct(): Cannot set Brand")
		}

		//p.ShipmentType, _ = dbToInt(product["fk_catalog_shipment_type"])
		//Hardcode shipment type
		p.ShipmentType = 3
		p.TY, _ = dbToInt(product["fk_catalog_ty"])
		p.DisplayStockedOut, _ = dbToInt(product["display_if_out_of_stock"])
		p.PetStatus, _ = dbToString(product["pet_status"])
		p.PetApproved, _ = dbToInt(product["pet_approved"])
		p.Description, _ = dbToString(product["description"])
		if p.Description == "" {
			p.Description = "Example description"
		}
		p.Status, _ = dbToString(product["status"])
		p.StatusSupplierConfig, _ = dbToString(product["status_supplier_config"])
		suppliername, _ := dbToString(product["supplier_name"])
		p.SupplierName = &suppliername

		t1, ok1 := product["approved_at"].(time.Time)
		if ok1 {
			p.ApprovedAt = &t1
		}
		t2, ok2 := product["activated_at"].(time.Time)
		if ok2 {
			p.ActivatedAt = &t2
		}
		t3, ok3 := product["created_at"].(time.Time)
		if ok3 {
			p.CreatedAt = &t3
		}
		t4, ok4 := product["updated_at"].(time.Time)
		if ok4 {
			p.UpdatedAt = &t4
		}
		err = p.SetAttributesData(attributeSetId)
		if err != nil {
			return err
		}
		err = p.SetAttributeSetData(attributeSetId)
		if err != nil {
			return err
		}
		err = p.SetVideos()
		if err != nil {
			return err
		}
		err = p.SetImages()
		if err != nil {
			return err
		}
		err = p.SetCategoriesAndLeaf()
		if err != nil {
			return err
		}
		err = p.SetSimples()
		if err != nil {
			return err
		}
		err = p.SetSupplierInfo(product)
		if err != nil {
			return err
		}
		err = p.SetUrlKey()
		if err != nil {
			return err
		}
		err = p.SetProductGroup(product)
		if err != nil {
			return err
		}
		err = p.Write2Mongo()
		if err != nil {
			return err
		}
		return nil
	}()
	printInfo("Finished: Processing product (" + strconv.Itoa(p.Id) + ")")
	st.Id = p.Id
	st.State = true
	if err != nil {
		st.State = false
		st.Msg = err.Error()
	}
	return st
}

func ProcessProductRows(rows *core.Rows) int {
	defer recoverHandler("ProcessProductRows")
	var count = 0
	ch := make(chan ProductMigrationStatus)
	for rows.Next() {
		product := make(generalMap)
		err := rows.ScanMap(&product)
		if err != nil {
			printErr(errors.New("ProcessProductRows(): Unable to Scan Product"))
			continue
		}
		count = count + 1
		go func(product generalMap) {
			defer recoverHandler("ProcessProductRows")
			ch <- ProcessProduct(product)
		}(product)
	}
	for i := 0; i < count; i++ {
		s := <-ch
		if !s.State {
			//insert in error collection
			mgoSession := ResourceFactory.GetMongoSession("Products")
			defer mgoSession.Close()
			mongodb := mgoSession.SetCollection(PRO_ERR_COLL)
			mongodb.Insert(
				ProductErr{Id: s.Id, Msg: s.Msg})
			printErr(errors.New(" Failed Insertion(" + strconv.Itoa(s.Id) + "):" + s.Msg))
		}
	}
	return count
}

func getLastMongoId() (int, error) {
	mgoSession := ResourceFactory.GetMongoSession(util.Products)
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.Products)

	type Data struct {
		SeqId int `bson:"seqId"`
	}
	var d Data
	err := mongodb.Find(nil).Sort("-seqId").Select(bson.M{"seqId": 1, "_id": 0}).One(&d)
	if err != nil && err == mgo.ErrNotFound {
		return 0, nil
	}
	if err != nil {
		return -1, err
	}
	return d.SeqId, nil
}

func getLastMysqlId() (int, error) {
	var max int
	var sql string = `SELECT MAX(id_catalog_config) FROM catalog_config;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		return 0, errors.New("getLastMysqlId():" + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return 0, errors.New("getLastMysqlId(): Assertion failed")
	}
	for rows.Next() {
		err = rows.Scan(&max)
		if err != nil {
			rows.Close()
			return max, errors.New("getLastMysqlId(): " + err.Error())
		}
	}
	rows.Close()
	return max, nil
}
