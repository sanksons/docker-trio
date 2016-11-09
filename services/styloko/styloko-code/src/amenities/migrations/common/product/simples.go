package product

import (
	"amenities/migrations/common/util"
	"common/xorm/mysql"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/go-xorm/core"
)

type ProductSimple struct {
	Id                   int                   `bson:"seqId" mapstructure:"id_catalog_simple"`
	SKU                  string                `bson:"sku" mapstructure:"sku"`
	SupplierSKU          *string               `bson:"supplierSku" mapstructure:"sku_supplier_simple"`
	SellerSKU            *string               `bson:"sellerSku"`
	BarcodeEan           *string               `bson:"barcodeEan"`
	EanCode              *string               `bson:"eanCode"`
	StatusSupplierSimple string                `bson:"statusSupplierSimple" mapstructure:"status_supplier_simple"`
	Quantity             int                   `bson:"quantity"`
	Price                *float64              `bson:"price" mapstructure:"price"`
	OriginalPrice        *float64              `bson:"originalPrice" mapstructure:"original_price"`
	SpecialPrice         *float64              `bson:"specialPrice" mapstructure:"special_price"`
	SpecialFromDate      *time.Time            `bson:"specialFromDate" mapstructure:"special_from_date"`
	SpecialToDate        *time.Time            `bson:"specialToDate" mapstructure:"special_to_date"`
	TaxClass             *int                  `bson:"taxClass" mapstructure:"fk_catalog_tax_class"`
	Attributes           map[string]*Attribute `bson:"attributes"`
	Global               map[string]*Attribute `bson:"global"`
	Status               string                `bson:"status" mapstructure:"status"`
	CreationSource       *string               `bson:"creationSource" mapstructure:"creation_source_simple"`
	CreatedAt            time.Time             `bson:"createdAt" mapstructure:"created_at"`
	UpdatedAt            time.Time             `bson:"updatedAt" mapstructure:"updated_at"`
}

func (s *ProductSimple) SetAttributes(attributeSetId int) error {
	attributes, err := processProductAttributes(s.Id, attributeSetId, "simple")
	if err != nil {
		return errors.New("(s *ProductSimple)#SetAttributes: " + err.Error())
	}
	global := make(map[string]*Attribute)
	attrs := make(map[string]*Attribute)
	for _, v := range attributes {
		cName := util.SnakeToCamel(v.Name)
		if v.IsGlobal {
			global[cName] = v
		} else {
			attrs[cName] = v
		}
	}
	s.Attributes = attrs
	s.Global = global
	return nil
}

func (s *ProductSimple) SetQuantity() error {
	var quantity int
	var simpleId = strconv.Itoa(s.Id)
	var sqlQ string = `SELECT quantity FROM catalog_stock WHERE fk_catalog_simple=` + simpleId
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		return errors.New("(s *ProductSimple)#SetQuantity: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(s *ProductSimple)#SetQuantity: Assertion failed")
	}
	for rows.Next() {
		err := rows.Scan(&quantity)
		if (err != nil) && (err != sql.ErrNoRows) {
			rows.Close()
			return errors.New("(s *ProductSimple)#SetQuantity: Unable to get Quantity")
		}
	}
	rows.Close()
	s.Quantity = quantity
	return nil
}

func GetVideosSeqCounter() (int, error) {
	var count int
	var sqlQ string = `SELECT MAX(id_video)
    FROM video;`
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		return count, errors.New("GetVideosSeqCounter():" + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return count, errors.New("GetVideosSeqCounter(): Assertion failed")
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if (err != nil) && (err != sql.ErrNoRows) {
			rows.Close()
			return count, errors.New("GetVideosSeqCounter(): Unable to get Quantity")
		}
	}
	rows.Close()
	return count, nil
}

func GetImageSeqCounter() (int, error) {
	var count int
	var sqlQ string = `SELECT MAX(id_catalog_product_image)
    FROM catalog_product_image;`
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		return count, errors.New("GetImageSeqCounter():" + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return count, errors.New("GetImageSeqCounter(): Assertion failed")
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if (err != nil) && (err != sql.ErrNoRows) {
			rows.Close()
			return count, errors.New("GetImageSeqCounter(): Unable to get Quantity")
		}
	}
	rows.Close()
	return count, nil
}

func GetSimpleSeqCounter() (int, error) {
	var count int
	var sqlQ string = `SELECT MAX(id_catalog_simple)
    FROM catalog_simple;`
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		return count, errors.New("getSimpleSeqCounter():" + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return count, errors.New("getSimpleSeqCounter(): Assertion failed")
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if (err != nil) && (err != sql.ErrNoRows) {
			rows.Close()
			return count, errors.New("getSimpleSeqCounter(): Unable to get Quantity")
		}
	}
	rows.Close()
	return count, nil
}
