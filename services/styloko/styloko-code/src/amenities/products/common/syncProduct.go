package common

import (
	"common/notification"
	"common/utils"
	"database/sql"
	"fmt"
	_ "reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

type MySqlSync struct {
	TxnObj *sql.Tx
}

type CheckAttribute struct {
	Attribute MySqlAttribute
	IsNil     int
}

type MySqlAttribute struct {
	Id            int
	Name          string
	AttributeType string
	ProductType   string
	IsGlobal      int
	SetName       *string
	Mandatory     int
}

//
// Save Full Product data to Mysql.
//
func (sync MySqlSync) SaveProduct(prod Product) error {
	//save system level attributes.
	err := sync.saveSystemTypeAttributes(prod)
	if err != nil {
		logger.Error("(sync MySqlSync)#SaveProduct: System attributes save failed.")
		return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
	}
	//save product categories
	err = sync.SaveCategories(prod.SeqId, prod.Categories)
	if err != nil {
		logger.Error("(sync MySqlSync)#SaveProduct: Categories save failed.")
		return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
	}
	// Update images for product
	if len(prod.Images) > 0 {
		for _, image := range prod.Images {
			if image == nil {
				continue
			}
			err = sync.SaveImage(prod.SeqId, *image)
			if err != nil {
				logger.Error("(sync MySqlSync)#SaveProduct: Image save failed.")
				return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
			}
		}
	}
	// Update videos for the product
	if len(prod.Videos) > 0 {
		for _, video := range prod.Videos {
			if video == nil {
				continue
			}
			err = sync.SaveVideo(prod.SeqId, *video)
			if err != nil {
				logger.Error("(sync MySqlSync)#SaveProduct: Video save failed.")
				return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
			}
		}
	}
	// Save sizechart info for product
	err = sync.SaveSizeChart(prod.SeqId, prod.SizeChart)
	if err != nil {
		logger.Error("(sync MySqlSync)#SaveProduct: SizeChart save failed.")
		return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
	}
	//Save product group information
	if prod.Group != nil && prod.Group.Id != 0 {
		err = sync.SaveProductGroup(prod.SeqId, *prod.Group)
		if err != nil {
			logger.Error("(sync MySqlSync)#SaveProduct: Group save failed.")
			return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
		}
	}
	// Fetch all the attributes for Product and prepare nil attributes list.
	var NilAttributesMap, prodAttributes map[string]MySqlAttribute

	var aErr error
	prodAttributes, aErr = sync.FetchAttributes(prod.AttributeSet.Id)
	if aErr != nil {
		logger.Error(fmt.Sprintf("%s", aErr.Error()))
		return aErr
	}
	NilAttributesMap = make(map[string]MySqlAttribute)

	// Find all nil config level attributes
	for _, prodAttribute := range prodAttributes {
		attr, globalOk := prod.Global[utils.SnakeToCamel(prodAttribute.Name)]
		// if attribute name matches but id doesnot match
		if globalOk && (*attr).Id != prodAttribute.Id {
			globalOk = false
		}
		nonGlobalattr, nonGlobalOk := prod.Attributes[utils.SnakeToCamel(prodAttribute.Name)]
		if nonGlobalOk && (*nonGlobalattr).Id != prodAttribute.Id {
			nonGlobalOk = false
		}
		// If attribute doesnot exist in list of either global/non-global attributes of Product
		if !globalOk && !nonGlobalOk && (prodAttribute.ProductType != PRODUCT_TYPE_SIMPLE) {
			// Key of type : prodId#type#attrId
			mapKey := strconv.Itoa(prod.SeqId) + "#" + prodAttribute.ProductType + "#" + strconv.Itoa(prodAttribute.Id)
			NilAttributesMap[mapKey] = prodAttribute
		}
	}

	// Update global attributes for product
	for _, globalAttr := range prod.Global {
		err = sync.SaveAttributes("global", prod.SeqId, PRODUCT_TYPE_CONFIG, *globalAttr)
		logger.Debug(fmt.Sprintf("%v", *globalAttr))
		if err != nil {
			logger.Error("(sync MySqlSync)#SaveProduct: Global Attribute save failed.")
			return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
		}
	}
	// Insert into catalog_config_attrset table
	qry := `INSERT INTO catalog_config_` + prod.AttributeSet.Name + ` (fk_catalog_config)
			VALUES (` + strconv.Itoa(prod.SeqId) + `) ON DUPLICATE KEY UPDATE
			fk_catalog_config = ?`
	_, err = sync.TxnObj.Exec(qry, prod.SeqId)
	if err != nil {
		logger.Error("(sync MySqlSync)#SaveProduct: Unable to make default entry for attribute.")
		return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
	}
	// Update non-global attributes for product
	for _, nonGlobalAttr := range prod.Attributes {
		err = sync.SaveAttributes(
			prod.AttributeSet.Name, prod.SeqId, PRODUCT_TYPE_CONFIG, *nonGlobalAttr,
		)
		logger.Debug(fmt.Sprintf("%v", *nonGlobalAttr))
		if err != nil {
			logger.Error("(sync MySqlSync)#SaveProduct: Non-Global Attribute save failed.")
			return fmt.Errorf("(sync MySqlSync)#SaveProduct: %s", err.Error())
		}
	}
	// save the product simples
	for _, simple := range prod.Simples {
		//Sync price information only if its a new simple
		//or simple does not belong to retail partner
		exists, err := sync.checkIfSimpleAlreadyExists(simple.Id)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProduct simple 1: %s", err.Error())
		}
		var err1 error
		if !exists || !utils.InArrayInt(RetailPartners, prod.SellerId) {
			err1 = sync.SaveProductSimple(*simple, prod,
				NilAttributesMap, prodAttributes)
		} else {
			err1 = sync.SaveProductSimpleWithoutPriceChange(*simple,
				prod, NilAttributesMap, prodAttributes)
		}
		if err1 != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProduct simple 2: %s", err1.Error())
		}
	}

	// Delete Nil value attributes
	success, err := sync.DeleteNilValuedAttributes(NilAttributesMap)
	if success == false && err != nil {
		logger.Error(fmt.Sprintf("%s", err.Error()))
		return err
	}

	return nil
}

func (sync MySqlSync) FetchAttributes(attributeSetId int) (map[string]MySqlAttribute, error) {

	ignorableAttributes := []string{
		"beauty_size",
		"ean_code",
		"barcode_ean",
		"is_gift_wrappable",
		//@todo: Need to remove this patch.
		"dispatch_location",
		"processing_time",
	}

	var attributes map[string]MySqlAttribute

	// Fetch all attributes
	attributes = make(map[string]MySqlAttribute)
	query := `SELECT ca.id_catalog_attribute , ca.name AS name, attribute_type,
				product_type, cas.name AS set_name , ca.mandatory as mandatory
			  FROM catalog_attribute AS ca
			  LEFT JOIN catalog_attribute_set AS cas
			  ON  ca.fk_catalog_attribute_set = cas.id_catalog_attribute_set
			  WHERE (id_catalog_attribute_set IS NULL OR id_catalog_attribute_set = ?)
			  AND ca.attribute_type IN ("value", "option", "multi_option")
			  ;
			  `
	rows, err := sync.TxnObj.Query(query, attributeSetId)
	if err != nil {
		return nil, fmt.Errorf("(sync MySqlSync)#FetchAttributes: Unable to fetch attributes. %s ", err.Error())
	}
	for rows.Next() {
		var attr MySqlAttribute
		err := rows.Scan(&attr.Id, &attr.Name, &attr.AttributeType, &attr.ProductType, &attr.SetName, &attr.Mandatory)
		if err != nil {
			return nil, fmt.Errorf("(sync MySqlSync)#FetchAttributes: Unable to fetch attributes. %s ", err.Error())
		}
		if attr.SetName == nil {
			attr.IsGlobal = 1
		} else {
			attr.IsGlobal = 0
		}
		var discard bool
		for _, attrname := range ignorableAttributes {
			if attrname == attr.Name {
				discard = true
			}
		}
		if !discard {
			attributes[attr.Name] = attr
		}
	}
	defer rows.Close()
	return attributes, nil
}

func (sync MySqlSync) DeleteNilValuedAttributes(nilAttributes map[string]MySqlAttribute) (bool, error) {

	for keyAttribute, attribute := range nilAttributes {
		attributeKeys := strings.Split(keyAttribute, "#")
		id, err := utils.GetInt(attributeKeys[0])
		if err != nil {
			return false, fmt.Errorf("(sync MySqlSync)#DeleteNilValuedAttributes: unable to get prodId.", err.Error())
		}
		// check for option type attribute
		if attribute.AttributeType == OPTION_TYPE_SINGLE {
			var err error
			switch attribute.IsGlobal {
			case 1:
				query := `UPDATE catalog_` + attribute.ProductType + ` SET ` +
					`fk_catalog_attribute_option_global_` + attribute.Name + ` = NULL where ` +
					` id_catalog_` + attribute.ProductType + ` = ?`
				_, err = sync.TxnObj.Exec(query, id)
			case 0:
				query := `UPDATE catalog_` + attribute.ProductType + `_` + *attribute.SetName + ` SET ` +
					` fk_catalog_attribute_option_` + *attribute.SetName + `_` + attribute.Name + ` = NULL where ` +
					` fk_catalog_` + attribute.ProductType + ` = ?`
				_, err = sync.TxnObj.Exec(query, id)
			}
			if err != nil {
				logger.Error(attribute)
				return false, fmt.Errorf("#(sync MySqlSync) DeleteNilValuedAttributes: Unable to update option attribute. %s", err.Error())
			}

		} else if attribute.AttributeType == OPTION_TYPE_MULTI {
			var err error
			switch attribute.IsGlobal {
			case 1:
				query := `DELETE from catalog_attribute_link_global_` + attribute.Name + ` where fk_catalog_config = ?`
				_, err = sync.TxnObj.Exec(query, id)

			case 0:
				var idSelect int
				selQry := `SELECT id_catalog_` + attribute.ProductType + `_` + *attribute.SetName + ` from catalog_` + attribute.ProductType + `_` +
					*attribute.SetName + ` where fk_catalog_` + attribute.ProductType + ` = ?`
				res := sync.TxnObj.QueryRow(selQry, id)
				sErr := res.Scan(&idSelect)
				if sErr == sql.ErrNoRows {
					// No entry to be removed
					continue
				} else if sErr != nil {
					logger.Error(attribute)
					return false, fmt.Errorf("#(sync MySqlSync) DeleteNilValuedAttributes: Unable to fetch id_catalog_prodType_attrSet. %s", sErr.Error())
				}
				// delete multi-options for attribute
				delQry := `DELETE from catalog_attribute_link_` + *attribute.SetName + `_` + attribute.Name + ` where ` +
					` fk_catalog_` + attribute.ProductType + `_` + *attribute.SetName + ` = ?`
				_, err = sync.TxnObj.Exec(delQry, idSelect)
			}
			if err != nil {
				logger.Error(attribute)
				return false, fmt.Errorf("#(sync MySqlSync) DeleteNilValuedAttributes: Unable to update multi-option attribute. %s", err.Error())
			}
		} else if attribute.AttributeType == OPTION_TYPE_VALUE {
			var err error
			if attribute.Mandatory == 1 {
				continue
			}
			switch attribute.IsGlobal {
			case 1:
				qry := `UPDATE catalog_` + attribute.ProductType + ` SET ` + attribute.Name + ` = NULL where ` +
					` id_catalog_` + attribute.ProductType + ` = ?`
				_, err = sync.TxnObj.Exec(qry, id)

			case 0:
				qry := `UPDATE catalog_` + attribute.ProductType + `_` + *attribute.SetName + ` SET ` + attribute.Name +
					` = NULL where fk_catalog_` + attribute.ProductType + ` = ?`
				_, err = sync.TxnObj.Exec(qry, id)
			}
			if err != nil {
				logger.Error(attribute)
				// Dont fail even if value type not updated
				return true, nil
			}
		}

	}
	//@TODO: Handle for all types of attributes
	// option type config/simple
	// global -> catalog_attribute_option_global_Atrname// change fk_catalog_attribute_option_global_Atrname
	// non-global ->   change fk_catalog_attribute_option_setName_attrName in catalog_config_attrset

	// multi-option config
	// global -> remove from catalog_attribute_link_global_name
	// non-global-> select from catalog_config_setName and remove from catalog_attribute_link_set_atrname

	//@TODO: Handle for attributes which are different in SC/Styloko

	return true, nil
}

//
// This function creates the new group for product and save foreign key in catalog_config
//
func (sync MySqlSync) SaveProductGroup(configId int, group ProductGroup) error {
	// Check and make entry in catalog_config_group for group if doesnot exist
	var id int
	checkQry := `SELECT id_catalog_config_group from catalog_config_group where id_catalog_config_group = ?`
	res := sync.TxnObj.QueryRow(checkQry, group.Id)
	err := res.Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("(sync MySqlSync)#SaveProductGroup: %s", err.Error())
	}
	if err == sql.ErrNoRows {
		//Insert
		insQry := `INSERT into catalog_config_group(id_catalog_config_group, name) VALUES(?,?)`
		_, err = sync.TxnObj.Exec(insQry, group.Id, group.Name)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync) SaveProductGroup -> Insert failed: %s", err.Error())
		}
	}
	// Make entry in catalog_config
	qry := `UPDATE catalog_config
			SET
			fk_catalog_config_group = ` + strconv.Itoa(group.Id) +
		` WHERE id_catalog_config = ` + strconv.Itoa(configId)
	_, err = sync.TxnObj.Exec(qry)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync) SaveProductGroup -> Update failed: %s", err.Error())
	}
	return nil
}

//
// It maps the sizechart info of the product to mysql
//
func (sync MySqlSync) SaveSizeChart(configId int, sizeChart ProSizeChart) error {
	// SizeChart doesnot exist no updates in mysql
	if sizeChart.Id == 0 || sizeChart.Data == nil {
		//logger.Warning("(sync MySqlSync)#SaveSizeChart(): sizechart doesnot exist for the product")
		return nil
	}
	sizeChartId := sizeChart.Id
	var sizeChartTy *int
	//Get the sizechart type
	var getType string = `SELECT sizechart_type
		FROM catalog_distinct_sizechart
		WHERE id_catalog_distinct_sizechart = ?`

	res := sync.TxnObj.QueryRow(getType, sizeChartId)
	err := res.Scan(&sizeChartTy)
	if err != nil {
		logger.Error(fmt.Sprintf("(sync MySqlSync)#SaveSizeChart1: %s", err.Error()))
		title := fmt.Sprintf("Product not synced with sizechart with Id: %d", sizeChart.Id)
		text := fmt.Sprintf("Could not get sizechart in mysql for product %d: %s", configId, err.Error())
		tags := []string{"product-sizechart", "mysql-sync"}
		notification.SendNotification(title, text, tags, "error")
		return nil
	}
	// check for null sizechart ty
	if sizeChartTy == nil {
		var defaultTy int = 1
		sizeChartTy = &defaultTy
	}
	query := `SELECT id_catalog_config_additional_info
		FROM catalog_config_additional_info
		WHERE fk_catalog_config = ?`
	var id int
	currTime := time.Now()
	row := sync.TxnObj.QueryRow(query, configId)
	err = row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(fmt.Sprintf("(sync MySqlSync)#SaveSizeChart2: %s", err.Error()))
		title := fmt.Sprintf("Product not synced with sizechart with Id: %d", sizeChart.Id)
		text := fmt.Sprintf("Error in getting data from sizechart-prod mapping table for productId %d: %s", configId, err.Error())
		tags := []string{"product-sizechart", "mysql-sync"}
		notification.SendNotification(title, text, tags, "error")
		return nil
	}
	if err == sql.ErrNoRows {
		//Insert
		insertQury := `INSERT INTO catalog_config_additional_info (
			fk_catalog_config,
			fk_catalog_distinct_sizechart,
			sizechart_type,
			reward_points,
			created_at
			) VALUES (?,?,?,?,?)`
		_, err = sync.TxnObj.Exec(
			insertQury, configId, sizeChartId, *sizeChartTy, 0.0, utils.ToMySqlTime(&currTime),
		)
		if err != nil {
			logger.Error(fmt.Sprintf("(sync MySqlSync)#SaveSizeChart3: %s", err.Error()))
			title := fmt.Sprintf("Product not synced with sizechart with Id: %d", sizeChart.Id)
			text := fmt.Sprintf("Failed to insert data in sizechart-prod mapping table for productId %d: %s", configId, err.Error())
			tags := []string{"product-sizechart", "mysql-sync"}
			notification.SendNotification(title, text, tags, "error")
			return nil
		}
		return nil
	}
	// Update
	updateQry := `UPDATE catalog_config_additional_info
	SET fk_catalog_distinct_sizechart = ` + strconv.Itoa(sizeChartId) +
		`, sizechart_type = ` + strconv.Itoa(*sizeChartTy) + ` WHERE fk_catalog_config = ` + strconv.Itoa(configId)
	_, err = sync.TxnObj.Exec(updateQry)
	if err != nil {
		logger.Error(fmt.Sprintf("(sync MySqlSync)#SaveSizeChart4: %s", err.Error()))
		title := fmt.Sprintf("Product not synced with sizechart with Id: %d", sizeChart.Id)
		text := fmt.Sprintf("Failed to update data in sizechart-prod mapping table for productId %d: %s", configId, err.Error())
		tags := []string{"product-sizechart", "mysql-sync"}
		notification.SendNotification(title, text, tags, "error")
		return nil
	}
	return nil
}

//
// This is a jugaad as asked by @apoorva.
//
func (sync MySqlSync) SaveProductSimpleWithoutPriceChange(simple ProductSimple, prod Product,
	NilAttributesMap map[string]MySqlAttribute,
	prodAttributes map[string]MySqlAttribute) error {

	//Get Supplier Info
	var statusSupplierSimple string
	statusSupplierSimple = "active"
	// Save system level attrbutes of simple
	query := `INSERT into catalog_simple (id_catalog_simple,
	fk_catalog_config,
	fk_catalog_import,
	sku,
	status,
	status_supplier_simple,
	created_at_external,
	updated_at_external,
	created_at,
	updated_at,
	sku_supplier_simple,
	barcode_ean,
	seller_sku,
	ean_code,
	price,
	original_price,
	creation_source_simple,
	fk_catalog_tax_class,
	jabong_discount_from,
	jabong_discount_to,
	jabong_discount)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	 ON DUPLICATE KEY UPDATE
		 fk_catalog_config = VALUES(fk_catalog_config),
		 fk_catalog_import = VALUES(fk_catalog_import),
		 sku = VALUES(sku),
		 status= VALUES(status),
		 status_supplier_simple = VALUES(status_supplier_simple),
		 created_at_external = VALUES(created_at_external),
		 updated_at_external= VALUES(updated_at_external),
		 created_at = VALUES(created_at),
		 updated_at= VALUES(updated_at),
		 sku_supplier_simple = VALUES(sku_supplier_simple),
		 barcode_ean= VALUES(barcode_ean),
		 seller_sku= VALUES(seller_sku),
		 ean_code= VALUES(ean_code),
		 price = VALUES(price),
		 original_price = VALUES(original_price),
		 creation_source_simple = VALUES(creation_source_simple),
		 fk_catalog_tax_class = VALUES(fk_catalog_tax_class),
		 jabong_discount_from = VALUES(jabong_discount_from),
		 jabong_discount_to = VALUES(jabong_discount_to),
		 jabong_discount = VALUES(jabong_discount)`
	_, err := sync.TxnObj.Exec(query,
		simple.Id,
		prod.SeqId,
		prod.CatalogImport,
		simple.SKU,
		simple.Status,
		statusSupplierSimple,
		simple.CreatedAt,
		simple.UpdatedAt,
		simple.CreatedAt,
		simple.UpdatedAt,
		simple.SupplierSKU,
		simple.BarcodeEan,
		simple.SellerSKU,
		simple.EanCode,
		simple.Price,
		simple.OriginalPrice,
		simple.CreationSource,
		simple.TaxClass,
		simple.JabongDiscountFromDate,
		simple.JabongDiscountToDate,
		simple.JabongDiscount)
	if err != nil {
		logger.Debug("#SaveProductSimple() : Save Simple failed", err.Error())
		return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
	}
	// For particular simple check if attribute exists for product and prepare nil map
	for _, prodAttribute := range prodAttributes {
		attr, globalOk := simple.Global[utils.SnakeToCamel(prodAttribute.Name)]
		// if attribute name matches but id doesnot match
		if globalOk && (*attr).Id != prodAttribute.Id {
			globalOk = false
		}
		nonGlobalattr, nonGlobalOk := simple.Attributes[utils.SnakeToCamel(prodAttribute.Name)]
		if nonGlobalOk && (*nonGlobalattr).Id != prodAttribute.Id {
			nonGlobalOk = false
		}

		// If attribute doesnot exist in list of either global/non-global attributes of Product
		if !globalOk && !nonGlobalOk && prodAttribute.ProductType != PRODUCT_TYPE_CONFIG {
			// Key of type : prodId#type#attrId
			mapKey := strconv.Itoa(simple.Id) + "#" + prodAttribute.ProductType + "#" + strconv.Itoa(prodAttribute.Id)
			(NilAttributesMap)[mapKey] = prodAttribute
		}

	}

	// Save the attributes of simples
	for _, attribute := range simple.Attributes {
		err := sync.SaveAttributes(prod.AttributeSet.Name, simple.Id, PRODUCT_TYPE_SIMPLE, *attribute)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
		}
	}
	//Save global atttributes of simple
	for _, attribute := range simple.Global {
		err := sync.SaveAttributes("global", simple.Id, PRODUCT_TYPE_SIMPLE, *attribute)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
		}
	}
	return nil
}

//
// This saves the product simple to mysql.
//
func (sync MySqlSync) SaveProductSimple(simple ProductSimple, prod Product,
	NilAttributesMap map[string]MySqlAttribute,
	prodAttributes map[string]MySqlAttribute) error {

	//Get Supplier Info
	var statusSupplierSimple string
	statusSupplierSimple = "active"
	// Save system level attrbutes of simple
	query := `INSERT into catalog_simple (id_catalog_simple,
	fk_catalog_config,
	fk_catalog_import,
	sku,
	status,
	status_supplier_simple,
	created_at_external,
	updated_at_external,
	created_at,
	updated_at,
	sku_supplier_simple,
	barcode_ean,
	seller_sku,
	ean_code,
	price,
	original_price,
	special_price,
	special_from_date,
	special_to_date,
	creation_source_simple,
	fk_catalog_tax_class,
	jabong_discount_from,
	jabong_discount_to,
	jabong_discount)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	 ON DUPLICATE KEY UPDATE
		 fk_catalog_config = VALUES(fk_catalog_config),
		 fk_catalog_import = VALUES(fk_catalog_import),
		 sku = VALUES(sku),
		 status= VALUES(status),
		 status_supplier_simple = VALUES(status_supplier_simple),
		 created_at_external = VALUES(created_at_external),
		 updated_at_external= VALUES(updated_at_external),
		 created_at = VALUES(created_at),
		 updated_at= VALUES(updated_at),
		 sku_supplier_simple = VALUES(sku_supplier_simple),
		 barcode_ean= VALUES(barcode_ean),
		 seller_sku= VALUES(seller_sku),
		 ean_code= VALUES(ean_code),
		 price = VALUES(price),
		 original_price = VALUES(original_price),
		 special_price = VALUES(special_price),
		 special_from_date = VALUES(special_from_date),
		 special_to_date = VALUES(special_to_date),
		 creation_source_simple = VALUES(creation_source_simple),
		 fk_catalog_tax_class = VALUES(fk_catalog_tax_class),
		 jabong_discount_from = VALUES(jabong_discount_from),jabong_discount_to = VALUES(jabong_discount_to),jabong_discount = VALUES(jabong_discount)`
	_, err := sync.TxnObj.Exec(query,
		simple.Id,
		prod.SeqId,
		prod.CatalogImport,
		simple.SKU,
		simple.Status,
		statusSupplierSimple,
		simple.CreatedAt,
		simple.UpdatedAt,
		simple.CreatedAt,
		simple.UpdatedAt,
		simple.SupplierSKU,
		simple.BarcodeEan,
		simple.SellerSKU,
		simple.EanCode,
		simple.Price,
		simple.OriginalPrice,
		simple.SpecialPrice,
		simple.SpecialFromDate,
		simple.SpecialToDate,
		simple.CreationSource,
		simple.TaxClass,
		simple.JabongDiscountFromDate,
		simple.JabongDiscountToDate,
		simple.JabongDiscount)
	if err != nil {
		logger.Debug("#SaveProductSimple() : Save Simple failed", err.Error())
		return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
	}
	// For particular simple check if attribute exists for product and prepare nil map
	for _, prodAttribute := range prodAttributes {
		attr, globalOk := simple.Global[utils.SnakeToCamel(prodAttribute.Name)]
		// if attribute name matches but id doesnot match
		if globalOk && (*attr).Id != prodAttribute.Id {
			globalOk = false
		}
		nonGlobalattr, nonGlobalOk := simple.Attributes[utils.SnakeToCamel(prodAttribute.Name)]
		if nonGlobalOk && (*nonGlobalattr).Id != prodAttribute.Id {
			nonGlobalOk = false
		}

		// If attribute doesnot exist in list of either global/non-global attributes of Product
		if !globalOk && !nonGlobalOk && prodAttribute.ProductType != PRODUCT_TYPE_CONFIG {
			// Key of type : prodId#type#attrId
			mapKey := strconv.Itoa(simple.Id) + "#" + prodAttribute.ProductType + "#" + strconv.Itoa(prodAttribute.Id)
			(NilAttributesMap)[mapKey] = prodAttribute
		}

	}

	// Save the attributes of simples
	for _, attribute := range simple.Attributes {
		err := sync.SaveAttributes(prod.AttributeSet.Name, simple.Id, PRODUCT_TYPE_SIMPLE, *attribute)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
		}
	}
	//Save global atttributes of simple
	for _, attribute := range simple.Global {
		err := sync.SaveAttributes("global", simple.Id, PRODUCT_TYPE_SIMPLE, *attribute)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveProductSimple:%s", err.Error())
		}
	}
	return nil
}

//
// It saves the video info for product
//
func (sync MySqlSync) SaveVideo(configId int, prodVideo ProductVideo) error {
	query := `INSERT INTO video (
		id_video,
		file_name,
		thumbnail,
		size,
		duration,
		status,
		created_at,
		updated_at
		) VALUES (?,?,?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE
		file_name = VALUES(file_name),
		thumbnail = VALUES(thumbnail),
		size = VALUES(size),
		duration = VALUES(duration),
		status = VALUES(status),
		created_at = VALUES(created_at),
		updated_at = VALUES(updated_at)`

	_, err := sync.TxnObj.Exec(
		query, prodVideo.Id, prodVideo.FileName, prodVideo.Thumbnail,
		prodVideo.Size, prodVideo.Duration, prodVideo.Status,
		prodVideo.CreatedAt, prodVideo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#SaveVideo: %s", err.Error())
	}
	// Save mapping of product and video
	type PV struct {
		IdProductVideo  int
		FkCatalogConfig int
		FkVideo         int
	}
	var pv PV
	searchSql := `Select id_product_video, fk_catalog_config, fk_video
		FROM product_video
		WHERE fk_catalog_config=? AND fk_video=?`

	row := sync.TxnObj.QueryRow(searchSql, configId, prodVideo.Id)
	err = row.Scan(&pv.IdProductVideo, &pv.FkCatalogConfig, &pv.FkVideo)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("(sync MySqlSync)#SaveVideo:%s", err.Error())
	}
	//check if nothing to do.
	if pv.IdProductVideo > 0 {
		return nil
	}
	//else insert
	insertSql := `INSERT INTO product_video (fk_catalog_config, fk_video) VALUES(?, ?)`
	_, err = sync.TxnObj.Exec(insertSql, configId, prodVideo.Id)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#SaveVideo():" + err.Error())
	}
	return nil
}

//
// It saves the image mapping to mysql.
//
func (sync MySqlSync) SaveImage(confidId int, prodImg ProductImage) error {
	query := `INSERT INTO catalog_product_image (
		id_catalog_product_image,
		fk_catalog_config,
		image,
		main,
		orientation,
		updated_at,
		original_filename
		) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE fk_catalog_config = VALUES(fk_catalog_config),
		image = VALUES(image), main = VALUES(main), orientation = VALUES(orientation), updated_at = VALUES(updated_at),
		original_filename =  VALUES(original_filename)`

	_, err := sync.TxnObj.Exec(query,
		prodImg.SeqId, confidId, prodImg.ImageNo, prodImg.Main,
		prodImg.Orientation, prodImg.UpdatedAt, prodImg.OriginalFileName,
	)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync) SaveImage: %s", err.Error())
	}
	return nil
}

// Determines the difference between two sets containing int values
func (sync MySqlSync) FindArrayDifference(setA []int, setB []int) []int {
	// Calculates A-B ie elements of A not having elements  of B
	if len(setA) == 0 {
		return setB
	}
	res := []int{}
	for i := 0; i < len(setA); i++ {
		ind := sync.FindInSlice(setB, setA[i])
		if ind < 0 {
			res = append(res, setA[i])
		}
	}
	return res
}

func (sync MySqlSync) FindInSlice(setA []int, ele int) int {
	for i := 0; i < len(setA); i++ {
		if setA[i] == ele {
			return i
		}
	}
	return -1
}

// This function converts the int array to delimiter separated string
func (sync MySqlSync) IntArrayToString(a []int, delim string) string {
	res := ""
	for i := 0; i < len(a); i++ {
		res = res + strconv.Itoa(a[i]) + delim
	}
	return strings.TrimSuffix(res, delim)
}

//
// It saves category-config mapping in mysql
//
func (sync MySqlSync) SaveCategories(configId int, categories []int) error {
	// Mapping which is existing in DB but not in new data is deleted.
	// Mapping which is not in DB but is there in new data is inserted
	var queryBody string
	var cat int
	catArr := []int{}
	checkQuery := `SELECT fk_catalog_category
		FROM catalog_config_has_catalog_category
		WHERE fk_catalog_config = ?`
	res, err := sync.TxnObj.Query(checkQuery, configId)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#SaveCategories: %s", err.Error())
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&cat)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveCategories scan error: %s", err.Error())
		}
		catArr = append(catArr, cat)
	}
	res.Close()
	//get categories to be inserted.
	if len(categories) > 0 {
		toBeInsert := sync.FindArrayDifference(categories, catArr)
		if len(toBeInsert) > 0 {
			queryHead := `INSERT into catalog_config_has_catalog_category
					(fk_catalog_config, fk_catalog_category) VALUES`
			for _, category := range toBeInsert {
				queryBody = queryBody + `(` + strconv.Itoa(configId) + `, ` + strconv.Itoa(category) + `),`
			}
			_, err = sync.TxnObj.Exec(queryHead + strings.TrimSuffix(queryBody, ","))
			if err != nil {
				return fmt.Errorf("(sync MySqlSync)#SaveCategories insert error: %s", err.Error())
			}
		}
	}

	// get categories to be deleted
	if len(catArr) > 0 {
		tobeDeleted := sync.FindArrayDifference(catArr, categories)
		if len(tobeDeleted) == 0 {
			//nothing to delete.
			return nil
		}
		delQuery := `DELETE from catalog_config_has_catalog_category WHERE fk_catalog_config = ` +
			strconv.Itoa(configId) + ` and fk_catalog_category IN(` + sync.IntArrayToString(tobeDeleted, ",") + `)`

		_, err = sync.TxnObj.Exec(delQuery)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#SaveCategories delete error: %s", err.Error())
		}
	}
	return nil
}

// this function obtains seller info using seller ID
func (sync MySqlSync) getSupplierInfoFromSellerId(id int) (*string, string, error) {
	var supplierName string
	var statusSupplierConfig string
	query := `SELECT name, status
			  from catalog_supplier
			  where id_catalog_supplier = ?`
	res := sync.TxnObj.QueryRow(query, id)

	err := res.Scan(&supplierName, &statusSupplierConfig)
	if err == sql.ErrNoRows {
		// if seller doesnot exist, send default values
		return nil, `active`, nil
	}
	if err != nil {
		return nil, "", err
	}
	return &supplierName, statusSupplierConfig, nil
}

// It updates system type attribute of the product
func (sync MySqlSync) saveSystemTypeAttributes(prod Product) error {
	var catalogTy *int
	// Get the supplier data
	supplierName, statusSupplierConfig, err := sync.getSupplierInfoFromSellerId(
		prod.SellerId,
	)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync) saveSystemTypeAttributes: %s", err.Error())
	}
	// By default seller status is active
	statusSupplierConfig = "active"

	// handle case if ty is zero, store NULL in mysql
	if prod.TY == 0 {
		catalogTy = nil
	} else {
		catalogTy = &prod.TY
	}
	query := `INSERT into catalog_config (
				id_catalog_config, sku, status, status_supplier_config,
				name, fk_catalog_brand, fk_catalog_import, display_if_out_of_stock,
				fk_catalog_attribute_set, pet_status,
				pet_approved, created_at, updated_at, activated_at, supplier_name,
				sku_supplier_config, description, fk_catalog_supplier, fk_catalog_ty,
				fk_catalog_shipment_type, product_set
				)
				VALUES (?,?,?,?,?,?,?,?,
						?,?,?,?,?,?,?,
						?,?,?,?,?,?)
					ON DUPLICATE KEY UPDATE
					sku = VALUES(sku), status = VALUES(status), status_supplier_config = VALUES(status_supplier_config),
					name = VALUES(name), fk_catalog_brand = VALUES(fk_catalog_brand), fk_catalog_import = VALUES(fk_catalog_import),
					display_if_out_of_stock = VALUES(display_if_out_of_stock),
					fk_catalog_attribute_set = VALUES(fk_catalog_attribute_set), pet_status = VALUES(pet_status),
					pet_approved = VALUES(pet_approved), created_at = VALUES(created_at), updated_at = VALUES(updated_at),
					activated_at = VALUES(activated_at), supplier_name = VALUES(supplier_name), sku_supplier_config = VALUES(sku_supplier_config),
					description = VALUES(description), fk_catalog_supplier = VALUES(fk_catalog_supplier), fk_catalog_ty = VALUES(fk_catalog_ty),
					fk_catalog_shipment_type = VALUES(fk_catalog_shipment_type), product_set = VALUES(product_set)`

	_, err = sync.TxnObj.Exec(
		query, prod.SeqId, prod.SKU, prod.Status, statusSupplierConfig,
		prod.Name, prod.BrandId, prod.CatalogImport, prod.DisplayStockedOut,
		prod.AttributeSet.Id, prod.PetStatus,
		prod.PetApproved, prod.CreatedAt, prod.UpdatedAt, prod.ActivatedAt,
		supplierName, prod.SupplierSKU, prod.Description, prod.SellerId,
		catalogTy, prod.ShipmentType, prod.ProductSet)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#saveSystemTypeAttributes:%s", err.Error())
	}
	return nil
}

//
// Save attributes
//
func (sync MySqlSync) SaveAttributes(attributeSet string, id int,
	productType string,
	attribute Attribute,
) error {

	//patch to avoid wrong attributes
	if attribute.Id == 0 {
		return nil
	}

	var err error
	switch attribute.OptionType {
	case OPTION_TYPE_VALUE:
		err = sync.saveValueTypeAttribute(
			id, attributeSet, attribute, productType,
		)
	case OPTION_TYPE_SINGLE:
		err = sync.saveOptionTypeAttribute(
			id, attributeSet, attribute, productType,
		)
	case OPTION_TYPE_MULTI:
		if productType == PRODUCT_TYPE_CONFIG {
			err = sync.updateMultiOptionTypeAttribute(
				id, attribute, attributeSet,
			)
		}
	default:
		logger.Error(fmt.Errorf(
			"Undefined Option Type:[%s]", attribute.OptionType,
		))
		err = nil
	}
	if err != nil {
		logger.Error(attribute.ToString())
		return fmt.Errorf("(sync MySqlSync)#SaveAttributes: %s", err.Error())
	}
	return nil
}

// This function saves the option type attributes
func (sync MySqlSync) saveOptionTypeAttribute(id int, attrSet string,
	attribute Attribute, proType string,
) error {
	tableName := "catalog_" + proType + "_" + attrSet
	fieldName := `fk_catalog_attribute_option_` + attrSet + `_` + attribute.Name
	//Get Attribute Value
	attrVal, gErr := attribute.GetValue("id")
	if gErr != nil {
		logger.Error(attribute.ToString())
		return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: %s", gErr.Error())
	}
	//check If its Global attribute
	if attrSet == "global" {
		// For isReturnable attribute we save different value from mongo.
		if attribute.Name == "is_returnable" {
			isReturnableVal, err := utils.GetInt(attrVal)
			if err != nil {
				return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: Could not get attribute Val", err.Error())
			}
			if isReturnableVal == 4 {
				attrVal = 0
			}
		}
		qury := `UPDATE catalog_` + proType + ` SET ` + fieldName + ` = ? where id_catalog_` + proType + ` = ?`
		_, err := sync.TxnObj.Exec(qury, attrVal, id)
		if err != nil {
			logger.Error(attribute.ToString())
			logger.Error(qury)
			return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: %s", err.Error())
		}
		return nil
	}
	// non global
	var tabId int
	checkQuery := `SELECT id_` + tableName + `
		FROM ` + tableName + `
		WHERE fk_catalog_` + proType + ` = ` + strconv.Itoa(id)
	row := sync.TxnObj.QueryRow(checkQuery)
	err := row.Scan(&tabId)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(checkQuery)
		logger.Error(attribute.ToString())
		return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: %s", err.Error())
	}
	if err == sql.ErrNoRows {
		// Row doesnot matches
		query := `INSERT into ` + tableName + `(fk_catalog_` + proType + `, ` + fieldName + `) VALUES (?,?)`
		_, err := sync.TxnObj.Exec(query, id, attrVal)
		if err != nil {
			logger.Error(query)
			logger.Error(attribute.ToString())
			return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: %s", err.Error())
		}
		return nil
	}
	if tabId <= 0 {
		return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute : tabId cannot be 0")
	}
	updateQury := `UPDATE ` + tableName + ` SET ` + fieldName + ` = ? WHERE fk_catalog_` + proType + ` = ?`
	_, err = sync.TxnObj.Exec(updateQury, attrVal, id)
	if err != nil {
		logger.Error(attribute.ToString())
		logger.Error(updateQury)
		return fmt.Errorf("(sync MySqlSync) saveOptionTypeAttribute: %s", err.Error())
	}
	return nil
}

// this function saves attribute of type value
func (sync MySqlSync) saveValueTypeAttribute(Id int, attrSet string,
	attribute Attribute, proType string,
) error {
	tableName := "catalog_" + proType + "_" + attrSet
	if attrSet == "global" {
		attrVal := attribute.Value
		if attrVal == nil {
			return nil
		}
		tableName = "catalog_" + proType
		query := `UPDATE ` + tableName + ` SET ` + attribute.Name + ` = ? where id_catalog_` + proType + ` = ?`
		_, err := sync.TxnObj.Exec(query, attrVal, Id)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#saveValueTypeAttribute: %s", err.Error())
		}
		return nil
	}
	//Update variation and beauty_size for beauty attribute set for catalog_simple.
	if attrSet == "beauty" && proType == PRODUCT_TYPE_SIMPLE {
		err := sync.UpdateVariationBeautySizeForSimple(Id, attribute)
		if err != nil {
			return err
		}
		return nil
	}
	// check if value exists in table
	var id int
	checkQuery := `SELECT id_` + tableName + ` from ` + tableName + ` where fk_catalog_` + proType + ` = ` + strconv.Itoa(Id)
	row := sync.TxnObj.QueryRow(checkQuery)
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("(sync MySqlSync)#saveValueTypeAttribute: %s", err.Error())
	}
	if err == sql.ErrNoRows {
		query := `INSERT into ` + tableName + `(fk_catalog_` + proType + `, ` + attribute.Name + `)VALUES(?,?)`
		_, err := sync.TxnObj.Exec(query, Id, attribute.Value)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#saveValueTypeAttribute: %s", err.Error())
		}
		return nil
	}
	if id < 0 {
		return fmt.Errorf("(sync MySqlSync)#saveValueTypeAttribute: id cannot be 0")
	}
	updateQury := `UPDATE ` + tableName + ` SET ` + attribute.Name + ` = ? where fk_catalog_` + proType + ` = ?`
	_, err = sync.TxnObj.Exec(updateQury, attribute.Value, Id)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#saveValueTypeAttribute: %s", err.Error())
	}
	return nil
}

// this function updates multi_option type attribute value for Global/Non global attributes
func (sync MySqlSync) updateMultiOptionTypeAttribute(
	configId int, attribute Attribute, attributeSet string,
) error {
	var msg string = "(sync MySqlSync)#updateMultiOptionTypeAttribute()"
	tableName := `catalog_attribute_link_` + attributeSet + `_` + attribute.Name
	attrValues, err := attribute.GetValue("id")
	if err != nil {
		logger.Error(attribute.ToString())
		return fmt.Errorf("%s: %s", msg, err.Error())
	}
	attrValuesAr, ok := attrValues.([]interface{})
	if !ok {
		logger.Error(attribute.ToString())
		return fmt.Errorf("%s: Failed to type assert multioption values.", msg)
	}
	var query string
	var fieldVal1 int
	fieldName1 := "fk_catalog_config"
	fieldName2 := `fk_catalog_attribute_option_` + attributeSet + `_` + attribute.Name
	fieldVal1 = configId
	if attributeSet != "global" {
		fieldName1 = "fk_catalog_config_" + attributeSet
		//Get fieldval1
		tblName := `catalog_config_` + attributeSet
		query := `SELECT id_` + tblName + ` from ` + tblName + ` where fk_catalog_config = ?`
		res := sync.TxnObj.QueryRow(query, configId)
		err := res.Scan(&fieldVal1)
		// If no rows for config in catalog_config_attrset add empty row for sku
		if err == sql.ErrNoRows {
			insQury := `INSERT into ` + tblName + ` (fk_catalog_config)VALUES(?)`
			resIns, err := sync.TxnObj.Exec(insQury, configId)
			if err != nil {
				logger.Error(insQury)
				return fmt.Errorf("%s: Unable to insert fk_catalog_config to catalog_config_attrset.", msg, err.Error())
			}
			fieldValInt64, err := resIns.LastInsertId()
			if err != nil {
				return fmt.Errorf("%s: Unable to get last inserted id for catalog_config_attrset", msg, err.Error())
			}
			fieldVal1, _ = utils.GetInt(fieldValInt64)
		} else if err != nil {
			logger.Error(attribute.ToString())
			logger.Error(query)
			return fmt.Errorf("%s: Unable to get id_catalog_config_attrset for config.", msg, err.Error())
		}
	}
	// Get existing multi-options in mysql DB for attribute
	existingMOptions := []int{}
	selQury := `SELECT ` + fieldName2 + ` from ` + tableName + ` where ` + fieldName1 + ` = ?`
	res, sErr := sync.TxnObj.Query(selQury, fieldVal1)
	if sErr != nil {
		logger.Error(attribute.ToString())
		logger.Error(selQury)
		return fmt.Errorf("%s: Unable to get existing multioptions from mysql: %s.", msg, sErr.Error())
	}
	defer res.Close()
	for res.Next() {
		var id int
		err = res.Scan(&id)
		if err != nil {
			logger.Error(attribute.ToString())
			return fmt.Errorf("%s: Unable to get existing multioptions from mysql: %s.", msg, err.Error())
		}
		existingMOptions = append(existingMOptions, id)
	}
	res.Close()
	var newMOptions []int
	for _, attrValIntrface := range attrValuesAr {
		val, err := utils.GetInt(attrValIntrface)
		if err != nil {
			logger.Error(attribute.ToString())
			return fmt.Errorf("%s: Couldnot convert newMultioption to int: %s.", msg, err.Error())
		}
		newMOptions = append(newMOptions, val)
	}

	// Delete the one which are in old not in new
	if len(existingMOptions) > 0 {
		for _, delMOptionVal := range sync.FindArrayDifference(existingMOptions, newMOptions) {
			delQuery := `DELETE from ` + tableName + ` where ` + fieldName1 + ` = ? and ` +
				fieldName2 + ` = ?`
			_, err := sync.TxnObj.Exec(delQuery, fieldVal1, delMOptionVal)
			if err != nil {
				logger.Error(attribute.ToString())
				logger.Error(delQuery)
				return fmt.Errorf("%s: Unable to delete old Multioptionval.: %s", msg, err.Error())
			}
		}
	}

	// Insert the one which are not in old but are in new
	if len(newMOptions) > 0 {
		for _, attrValId := range sync.FindArrayDifference(newMOptions, existingMOptions) {
			query = `INSERT into ` + tableName + ` (` + fieldName1 + `, ` + fieldName2 + `) VALUES(?,?)
			ON DUPLICATE KEY UPDATE ` + fieldName2 + ` = VALUES(` + fieldName2 + `)`
			_, err = sync.TxnObj.Exec(query, fieldVal1, attrValId)
			if err != nil {
				logger.Error(attribute.ToString())
				logger.Error(query)
				return fmt.Errorf("%s : unable to update multioption attribute: %s", msg, err.Error())
			}
		}
	}

	return nil
}

func (sync MySqlSync) UpdateVariationBeautySizeForSimple(Id int, attribute Attribute) error {
	var id int
	checkQuery := `SELECT id_catalog_simple_beauty from catalog_simple_beauty where fk_catalog_simple = ` + strconv.Itoa(Id)
	//logger.Error(checkQuery)
	row := sync.TxnObj.QueryRow(checkQuery)
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("(sync MySqlSync)#UpdateVariationBeautySizeForSimple: %s", err.Error())
	}
	if err == sql.ErrNoRows {
		query := `INSERT into catalog_simple_beauty (fk_catalog_simple , ` + attribute.Name + `, beauty_size)VALUES(?,?,?)`
		_, err := sync.TxnObj.Exec(query, Id, attribute.Value, attribute.Value)
		if err != nil {
			return fmt.Errorf("(sync MySqlSync)#UpdateVariationBeautySizeForSimple(insert): %s", err.Error())
		}
		return nil
	}
	if id < 0 {
		return fmt.Errorf("(sync MySqlSync)#UpdateVariationBeautySizeForSimple: id cannot be 0")
	}
	updateQury := `UPDATE catalog_simple_beauty SET ` + attribute.Name + ` = ? , beauty_size = ? where fk_catalog_simple = ?`
	_, err = sync.TxnObj.Exec(updateQury, attribute.Value, attribute.Value, Id)
	if err != nil {
		return fmt.Errorf("(sync MySqlSync)#UpdateVariationBeautySizeForSimple(update): %s", err.Error())
	}
	return nil
}

func (sync MySqlSync) checkIfSimpleAlreadyExists(simpleId int) (bool, error) {
	var id int
	sql1 := `SELECT id_catalog_simple FROM catalog_simple WHERE id_catalog_simple = ?`
	row := sync.TxnObj.QueryRow(sql1, simpleId)
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("(sync MySqlSync)#checkIfSimpleAlreadyExists: %s", err.Error())
	}
	if id > 0 {
		return true, nil
	}
	return false, nil
}
