package product

import (
	"amenities/migrations/common/util"
	proUtil "amenities/products/common"
	sizeChartServ "amenities/services/sizecharts"
	"common/ResourceFactory"
	"common/xorm/mysql"
	db "database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/core"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
)

type ProductMongo struct {
	Id                   int                   `bson:"seqId"`
	ProductSet           int                   `bson:"productSet"`
	SKU                  string                `bson:"sku"`
	SupplierSKU          *string               `bson:"supplierSkuConfig,omitempty"`
	SupplierName         *string               `bson:"supplierName,omitempty"`
	StatusSupplierConfig string                `bson:"statusSupplierConfig"`
	Name                 string                `bson:"name"`
	Description          string                `bson:"description"`
	UrlKey               string                `bson:"urlKey"`
	SellerId             int                   `bson:"sellerId"`
	DisplayStockedOut    int                   `bson:"displayIfOutOfStock"`
	Status               string                `bson:"status"`
	PetStatus            string                `bson:"petStatus"`
	PetApproved          int                   `bson:"petApproved"`
	CatalogImport        *int                  `bson:"catalogImport,omitempty"`
	AttributeSet         ProAttributeSet       `bson:"attributeSet"`
	BrandId              int                   `bson:"brandId"`
	TY                   int                   `bson:"ty"`
	Group                ProductGroup          `bson:"group"`
	ShipmentType         int                   `bson:"shipmentType"`
	Categories           []int                 `bson:"categories"`
	Leaf                 []int                 `bson:"leaf"`
	PrimaryCategory      int                   `bson:"primaryCategory"`
	Images               []*ProductImage       `bson:"images,omitempty"`
	Videos               []*ProductVideo       `bson:"videos,omitempty"`
	Global               map[string]*Attribute `bson:"global"`
	Attributes           map[string]*Attribute `bson:"attributes"`
	Simples              []*ProductSimple      `bson:"simples"`
	ApprovedAt           *time.Time            `bson:"approvedAt,omitempty"`
	CreatedAt            *time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt            *time.Time            `bson:"updatedAt,omitempty"`
	ActivatedAt          *time.Time            `bson:"activatedAt,omitempty"`
}

func (p *ProductMongo) SetAttributesData(setId int) error {
	attributes, err := processProductAttributes(p.Id, setId, "config")
	if err != nil {
		return errors.New("(p *ProductMongo)#SetAttributesData: " + err.Error())
	}
	//parse
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
	p.Attributes = attrs
	p.Global = global
	return nil
}

func (p *ProductMongo) SetAttributeSetData(attributeSetId int) error {
	set := ProAttributeSet{}
	id := strconv.Itoa(attributeSetId)
	var sql string = `SELECT
			id_catalog_attribute_set,
			name,
			label
		FROM catalog_attribute_set
		WHERE id_catalog_attribute_set=` + id
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		return errors.New("(p *ProductMongo)#SetAttributeSetData: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetAttributeSetData: Assertion failed")
	}
	for rows.Next() {
		err = rows.Scan(&set.Id, &set.Name, &set.Label)
		if err != nil {
			rows.Close()
			return errors.New("(p *ProductMongo)#SetAttributeSetData: " + err.Error())
		}
	}
	rows.Close()
	p.AttributeSet = set
	return nil
}

func (p *ProductMongo) SetVideos() error {
	var id string = strconv.Itoa(p.Id)
	var sql string = `SELECT
    		id_video,
    		file_name,
    		thumbnail,
    		size,
    		duration,
    		status,
    		created_at,
    		updated_at
    	FROM product_video
        INNER JOIN video
        ON product_video.fk_video = video.id_video
        WHERE  fk_catalog_config =  ` + id + `;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return nil
		}
		return errors.New("(p *ProductMongo)#SetVideos: Assertion failed")
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetVideos: Assertion failed")
	}
	var videos []*ProductVideo
	for rows.Next() {
		video := &ProductVideo{}
		err := rows.ScanStructByIndex(video)
		if err != nil {
			continue
		}
		videos = append(videos, video)
	}
	rows.Close()
	p.Videos = videos
	return nil
}

func (p *ProductMongo) SetImages() error {
	var id string = strconv.Itoa(p.Id)
	var sql string = `SELECT
			catalog_product_image.id_catalog_product_image,
               catalog_product_image.image,
               catalog_product_image.main,
               catalog_product_image.updated_at,
               catalog_product_image.original_filename,
               (SELECT
                  REPLACE(CONCAT(IF(catalog_brand.name IS NOT NULL,
                  catalog_brand.name,
                  catalog_config.id_catalog_config),
                  "-",
                  catalog_config.name,
                  "-",

               RIGHT(UNIX_TIMESTAMP(catalog_product_image.updated_at),
               4) ,
               "-",
               IF(catalog_product_image.fk_catalog_config < 100,
               RPAD(REVERSE(CAST(catalog_product_image.fk_catalog_config AS CHAR)),
               3,
               0),
               REVERSE(CAST(catalog_product_image.fk_catalog_config AS CHAR))),
               "-"),
               " ",
               "-")) AS image_name, orientation FROM
               catalog_product_image
            INNER JOIN
               catalog_config
                  ON catalog_config.id_catalog_config = catalog_product_image.fk_catalog_config
            LEFT JOIN
               catalog_brand
                  ON catalog_brand.id_catalog_brand = catalog_config.fk_catalog_brand
            WHERE catalog_product_image.fk_catalog_config = ` + id + ` ORDER BY
            catalog_product_image.fk_catalog_config ASC;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return nil
		}
		return errors.New("(p *ProductMongo)#SetImages: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetImages: Assertion failed")
	}
	images := []*ProductImage{}
	for rows.Next() {
		var tmp ProductImage
		err := rows.Scan(&tmp.Id,
			&tmp.ImageNo,
			&tmp.Main,
			&tmp.UpdatedAt,
			&tmp.OriginalFileName,
			&tmp.ImageName,
			&tmp.Orientation)
		if err != nil {
			continue
		}
		tmp.ImageName = SanitizeImageName(tmp.ImageName)
		images = append(images, &tmp)
	}
	rows.Close()
	p.Images = images
	return nil
}

func (p *ProductMongo) SetCategoriesAndLeaf() error {
	var id string = strconv.Itoa(p.Id)
	var categories, leaf []int
	var sql string = `SELECT
		cc.id_catalog_category,
		IF (rgt-lft = 1, 1,0) AS is_leaf
	FROM catalog_category AS cc
	INNER JOIN catalog_config_has_catalog_category AS cchcc
	ON cc.id_catalog_category = cchcc.fk_catalog_category
	WHERE cchcc.fk_catalog_config = ` + id + ` ORDER BY cc.lft ASC`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return nil
		}
		return errors.New("(p *ProductMongo)#SetCategoriesAndLeaf: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetCategoriesAndLeaf: Assertion failed")
	}

	for rows.Next() {
		var categoryId, isLeaf int
		err := rows.Scan(&categoryId, &isLeaf)
		if err != nil {
			continue
		}
		categories = append(categories, categoryId)
		if isLeaf == 1 {
			leaf = append(leaf, categoryId)
		}
	}
	rows.Close()
	p.Categories = categories
	p.Leaf = leaf
	if len(p.Leaf) > 0 {
		p.PrimaryCategory = p.Leaf[0]
	} else if len(p.Categories) > 0 {
		p.PrimaryCategory = p.Categories[len(p.Categories)-1]
	} else {
		p.PrimaryCategory = 1
	}
	return nil
}

func (p *ProductMongo) SetUrlKey() error {
	var id string = strconv.Itoa(p.Id)
	var sql string = `SELECT url_key
    	FROM catalog_config
		INNER JOIN  catalog_brand
		ON catalog_config.fk_catalog_brand = catalog_brand.id_catalog_brand
		WHERE id_catalog_config=` + id + `;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		return errors.New("(p *ProductMongo)#SetUrlKey: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetUrlKey: Assertion failed")
	}
	var brandurlKey string
	for rows.Next() {
		err := rows.Scan(&brandurlKey)
		if err != nil {
			continue
		}
	}
	rows.Close()
	p.UrlKey = fmt.Sprintf("%s-%s-%d",
		strings.Replace(brandurlKey, "/", "-", -1),
		sanitize.BaseName(strings.Replace(p.Name, "/", "-", -1)),
		p.Id)
	return nil
}

func (p *ProductMongo) SetSimples() error {
	var simples []*ProductSimple
	var id string = strconv.Itoa(p.Id)
	var sql string = `SELECT id_catalog_simple,
		sku,
		sku_supplier_simple,
		seller_sku,
		barcode_ean,
		ean_code,
		price,
		original_price,
		special_price,
		special_from_date,
		special_to_date,
		fk_catalog_tax_class,
		status,
		creation_source_simple,
		created_at,
		updated_at,
		status_supplier_simple
	 FROM catalog_simple WHERE fk_catalog_config=` + id + `;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return nil
		}
		return errors.New("(p *ProductMongo)#SetSimples: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetSimples: Assertion failed")
	}

	for rows.Next() {
		s := &ProductSimple{}
		err := rows.Scan(&s.Id, &s.SKU, &s.SupplierSKU, &s.SellerSKU, &s.BarcodeEan, &s.EanCode,
			&s.Price, &s.OriginalPrice, &s.SpecialPrice, &s.SpecialFromDate, &s.SpecialToDate,
			&s.TaxClass, &s.Status, &s.CreationSource,
			&s.CreatedAt, &s.UpdatedAt, &s.StatusSupplierSimple)
		if err != nil {
			rows.Close()
			return errors.New("(p *ProductMongo)#SetSimples: " + err.Error())
		}
		var simplesku = s.SKU
		if s.SellerSKU == nil && s.BarcodeEan == nil {
			s.SellerSKU = &simplesku
			s.BarcodeEan = &simplesku
		}
		if s.SellerSKU == nil && s.BarcodeEan != nil {
			s.SellerSKU = s.BarcodeEan
		}
		if s.SellerSKU != nil && s.BarcodeEan == nil {
			s.BarcodeEan = s.SellerSKU
		}
		if s.EanCode == nil {
			s.EanCode = s.BarcodeEan
		}
		s.SetQuantity()
		s.SetAttributes(p.AttributeSet.Id)
		simples = append(simples, s)
	}
	rows.Close()
	p.Simples = simples
	return nil
}

func (p *ProductMongo) SetProductGroup(rawProduct generalMap) error {

	groupId, _ := dbToInt(rawProduct["fk_catalog_config_group"])
	if groupId <= 0 {
		return nil
	}
	var groupIdStr string = strconv.Itoa(groupId)
	var sql string = `SELECT name
    	FROM catalog_config_group
    	WHERE id_catalog_config_group=` + groupIdStr
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return nil
		}
		return errors.New("(p *ProductMongo)#SetProductGroup: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return errors.New("(p *ProductMongo)#SetProductGroup: Assertion failed")
	}
	var g string
	for rows.Next() {
		rows.Scan(&g)
	}
	rows.Close()
	grp := ProductGroup{}
	grp.Id = groupId
	grp.Name = g
	p.Group = grp
	return nil
}

func (p *ProductMongo) SetSupplierInfo(rawProduct generalMap) error {
	supplierSKU, _ := dbToString(rawProduct["sku_supplier_config"])
	p.SellerId, _ = dbToInt(rawProduct["fk_catalog_supplier"])
	if supplierSKU != "" {
		p.SupplierSKU = &supplierSKU
	}
	return nil
}

func (p *ProductMongo) Write2Mongo() error {
	mgoSession := ResourceFactory.GetMongoSession("Products")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.Products)
	if len(p.Global) == 0 {
		mgoSession.Close()
		return errors.New("Global key is empty")
	}
	/*
		if p.PrimaryCategory == 0 {
			mgoSession.Close()
			return errors.New("Primary category cannot be 0")
		}
	*/
	if p.AttributeSet.Id == 17 {
		mgoSession.Close()
		return errors.New("Mobile set ignored")
	}
	_, err := mongodb.Upsert(bson.M{"seqId": p.Id}, p)
	mgoSession.Close()
	if err != nil {
		return errors.New("(p *ProductMongo)#Write2Mongo():Unable to write" + err.Error())
	}
	// size chart insertion
	pro, err := proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).GetById(p.Id)
	if err != nil {
		return errors.New("SizeChart Append Failed: " + err.Error())
	}
	data := sizeChartServ.GetSizeChartForProductMigration(pro)
	err = proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).AddNode(p.SKU, "sizeChart", data)
	if err != nil {
		return errors.New("SizeChart Append Failed: " + err.Error())
	}
	return nil
}
