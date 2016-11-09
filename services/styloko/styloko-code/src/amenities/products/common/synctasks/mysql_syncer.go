package synctasks

import (
	proUtil "amenities/products/common"
	factory "common/ResourceFactory"
	"common/appconfig"
	"common/notification"
	"common/notification/datadog"
	"common/utils"
	dbsql "database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func ProcessTask(typ string, data []byte, resource int, skipMemcacheUpdate bool) error {
	//Get QC state and set in struct
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	markqcdwn := conf.ProductSyncing.MarkQCDown
	var qcDwn bool
	if strings.ToLower(markqcdwn) == "true" {
		qcDwn = true
	}
	//setup syncer
	syncer := MySqlSyncer{
		Data:       data,
		MarkQCDown: qcDwn,
		ConfigId:   resource,
	}
	var err error
	switch typ {
	// Price update
	case proUtil.UPDATE_TYPE_PRICE:
		err = syncer.PriceSync()
	// Shipment update
	case proUtil.UPDATE_TYPE_SHIPMENT:
		err = syncer.ShipmentSync()
	// Image addition
	case proUtil.UPDATE_TYPE_IMAGEADD:
		err = syncer.AddImage()
	// Image deletion
	case proUtil.UPDATE_TYPE_IMAGEDEL:
		err = syncer.DeleteImage()
	// Video Insert/Update
	case proUtil.UPDATE_TYPE_VIDEO:
		err = syncer.VideoSync()
	// Video status update
	case proUtil.UPDATE_TYPE_VIDEO_STATUS:
		err = syncer.VideoStatusSync()
	// Product Update
	case proUtil.UPDATE_TYPE_PRODUCT:
		err = syncer.ProductSync()
	// Product status update
	case proUtil.UPDATE_TYPE_PRODUCT_STATUS:
		err = syncer.ProductSimpleStatusUpdate()
	// jabong discount update
	case proUtil.UPDATE_TYPE_JABONG_DISCOUNT:
		err = syncer.JabongDiscount()
	// prroduct general attribute update
	case proUtil.SYNC_ATTRIBUTE_GENERAL:
		err = syncer.ProductGeneralAttributeSync()
	// prroduct system attribute update
	case proUtil.SYNC_ATTRIBUTE_SYSTEM:
		err = syncer.ProductSystemAttributeSync()
	default:
		err = fmt.Errorf("ProcessTask()#Not a Valid Task: %s", typ)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("ProcessTask():%s", err.Error()))
		return err
	}
	//check if we need to set QC down
	if syncer.MarkQCDown {
		proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL).SetPetApproved(
			syncer.ConfigId, 0,
		)
	}
	//send call to BOB, Boutique for memcache update
	if !skipMemcacheUpdate {
		go func(configId int) {
			defer proUtil.RecoverHandler("ProcessTask")
			//memcache push
			product, err := proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).GetById(configId)
			if err != nil {
				//notify
				notification.SendNotification(
					"Product Memcache Update Failed",
					fmt.Sprintf("Product:%d, Message1:%s", configId, err.Error()),
					[]string{proUtil.TAG_PRODUCT_SYNC, proUtil.TAG_PRODUCT},
					datadog.ERROR,
				)
			}
			comment := fmt.Sprintf("UpdateType:%s, Product:%s", typ, product.SKU)
			err = product.PushToMemcache(comment)
			if err != nil {
				//notify
				notification.SendNotification(
					"Product Memcache Update Failed",
					fmt.Sprintf("Product:%d, Message2:%s", configId, err.Error()),
					[]string{proUtil.TAG_PRODUCT_SYNC, proUtil.TAG_PRODUCT},
					datadog.ERROR,
				)
			}
			return
		}(syncer.ConfigId)
	}
	return nil
}

type MySqlSyncer struct {
	Data       []byte
	ConfigId   int
	MarkQCDown bool
}

//
// This functions takes in the transaction object and
// removes the supplied attribute value from the product/simple
// Params:
//  attrName -> name of the attribute on which action needs to be taken
//  isGlobal -> is ita global attribute
//  sku -> sku of config/simple
//  typ -> config/simple
//  tx -> sql transaction object
//
func RemoveAttribute(
	attrName string,
	isGlobal bool,
	sku string,
	typ string,
	tx *dbsql.Tx) error {

	var (
		sql              string
		attributeType    string
		attributeSetName string
		tableName        string
		err              error
		id               int
	)
	//get product id from sku and type
	if typ == proUtil.PRODUCT_TYPE_CONFIG {
		sql = `SELECT id_catalog_config from catalog_config WHERE sku=?`
	} else {
		sql = `SELECT id_catalog_simple from catalog_simple WHERE sku=?`
	}
	row := tx.QueryRow(sql, sku)
	err = row.Scan(&id)
	if err != nil {
		return fmt.Errorf("RemoveAttribute(): Scan failed1, %s", err.Error())
	}

	//handle global and non-global types seperately.
	if isGlobal {
		sql = `SELECT attribute_type FROM catalog_attribute WHERE name=?`
		row := tx.QueryRow(sql, attrName)
		err = row.Scan(&attributeType)
		if err != nil {
			return fmt.Errorf("RemoveAttribute(): Scan failed2, %s", err.Error())
		}
		safeAttrName := utils.SqlSafe(attrName)
		switch attributeType {
		case proUtil.OPTION_TYPE_VALUE:
			tableName = fmt.Sprintf("catalog_%s", typ)
			sql = `UPDATE ` + tableName + ` SET ` + safeAttrName + `= NULL
          			WHERE id_` + tableName + `= ?`

		case proUtil.OPTION_TYPE_SINGLE:
			tableName = fmt.Sprintf("catalog_%s", typ)
			attrNm := fmt.Sprintf("fk_catalog_attribute_option_global_%s", attrName)
			sql = `UPDATE ` + tableName + ` SET ` + attrNm + `= NULL
          			WHERE id_` + tableName + `= ?`

		case proUtil.OPTION_TYPE_MULTI:
			tableName = fmt.Sprintf("catalog_attribute_link_global_%s", attrName)
			sql = `DELETE FROM ` + tableName + ` WHERE fk_catalog_config=?`

		default:
			return fmt.Errorf("RemoveAttribute(): Wrong value for optiontype, %s", attributeType)
		}
	} else {
		//get attribute Set
		sql = `SELECT attribute_type, cas.name
              		FROM catalog_attribute AS ca
              		INNER JOIN catalog_attribute_set AS cas
              		ON ca.fk_catalog_attribute_set = cas.id_catalog_attribute_set
              	WHERE ca.name=?`
		row := tx.QueryRow(sql, attrName)
		err = row.Scan(&attributeType, &attributeSetName)
		if err != nil {
			return fmt.Errorf("RemoveAttribute(): Scan failed3, %s", err.Error())
		}
		switch attributeType {
		case proUtil.OPTION_TYPE_VALUE:
			tableName = fmt.Sprintf("catalog_%s_%s", typ, attributeSetName)
			whereField := fmt.Sprintf("fk_catalog_%s", typ)
			sql = `UPDATE ` + tableName + ` SET ` + utils.SqlSafe(attrName) + `= NULL
          			WHERE ` + whereField + `= ?`

		case proUtil.OPTION_TYPE_SINGLE:
			tableName = fmt.Sprintf("catalog_%s_%s", typ, attributeSetName)
			attrNm := fmt.Sprintf("fk_catalog_attribute_option_%s_%s", attributeSetName, attrName)
			whereField := fmt.Sprintf("fk_catalog_%s", typ)
			sql = `UPDATE ` + tableName + ` SET ` + utils.SqlSafe(attrNm) + `= NULL
          			WHERE ` + whereField + `= ?`

		case proUtil.OPTION_TYPE_MULTI:
			tableName = fmt.Sprintf("catalog_attribute_link_%s_%s", attributeSetName, attrName)
			whereField := fmt.Sprintf("fk_catalog_%s_%s", typ, attributeSetName)
			innerTable := fmt.Sprintf("catalog_%s_%s", typ, attributeSetName)
			innerWhere := fmt.Sprintf("fk_catalog_%s", typ)
			sql = `DELETE FROM ` + tableName + `
				   WHERE ` + whereField + ` = (SELECT
				   	id_catalog_config_app_men
				   	FROM ` + innerTable + `
				   	WHERE ` + innerWhere + `= ?
				   )`
		default:
			return fmt.Errorf("RemoveAttribute(): Wrong type, %s", attributeType)
		}
	}
	_, err = tx.Exec(sql, id)
	if err != nil {
		return fmt.Errorf("Save failed:%s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) ProductGeneralAttributeSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_generalAttribute")
	defer logger.EndProfile(profiler, "product_mysql_sync_generalAttribute")

	var err error
	var msg string = "(syncer MySqlSyncer)#ProductGeneralAttributeSync()"
	type Data struct {
		AttributeName string             `json:"attributeName"`
		IsGlobal      bool               `json:"isGlobal"`
		ProductSku    string             `json:"productSku"`
		ProductType   string             `json:"productType"`
		Action        int                `json:"action"`
		Value         interface{}        `json:"value"`
		PetApproved   int                `json:"petApproved"`
		AttrData      *proUtil.Attribute `json:"attrData"`
	}
	var data Data
	err = json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err.Error())
	}
	// Get transaction object
	mysqlDriver, _ := factory.GetMySqlDriver(proUtil.PRODUCT_COLLECTION)
	tx, txErr := mysqlDriver.GetTxnObj()
	if txErr != nil {
		return fmt.Errorf("%s Cannot get transaction object: %s",
			msg, txErr.DeveloperMessage)
	}
	var isCommited bool
	//Rollback in defer.
	defer func() {
		if !isCommited {
			er := tx.Rollback()
			if er != nil {
				logger.Error(fmt.Sprintf(
					"(ma *MySqlAdapter)#SaveProduct(pro Product)Rollback failed:%s",
					er.Error(),
				))
			}
		}
	}()

	//check if we need to remove attribute
	if data.Action == proUtil.ACTION_REMOVE && data.AttrData == nil {
		err := RemoveAttribute(
			data.AttributeName, data.IsGlobal, data.ProductSku, data.ProductType, tx,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: %s", msg, err.Error())
		}
		tx.Commit()
		isCommited = true
		return nil
	}
	//just sync the attribute to mysql
	var pro proUtil.Product
	var id int
	pro, err = proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO).GetById(syncer.ConfigId)
	if err != nil {
		return fmt.Errorf("%s unable to get product: %s", msg, err.Error())
	}
	id = pro.SeqId
	if data.ProductType == proUtil.PRODUCT_TYPE_SIMPLE {
		for _, s := range pro.Simples {
			if s.SKU == data.ProductSku {
				id = s.Id
				break
			}
		}
	}
	mysqlSync := proUtil.MySqlSync{
		TxnObj: tx,
	}
	var attributeSet string = pro.AttributeSet.Name
	if data.IsGlobal {
		attributeSet = "global"
	}
	err = mysqlSync.SaveAttributes(attributeSet, id, data.ProductType, *data.AttrData)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s unable to save atttribute: %s", msg, err.Error())
	}
	tx.Commit()
	isCommited = true
	return nil
}

func (syncer MySqlSyncer) ProductSystemAttributeSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_systemAttribute")
	defer logger.EndProfile(profiler, "product_mysql_sync_systemAttribute")

	var err error
	var msg string = "(syncer MySqlSyncer)#ProductSystemAttributeSync()"
	type Data struct {
		AttributeName string             `json:"attributeName"`
		IsGlobal      bool               `json:"isGlobal"`
		ProductSku    string             `json:"productSku"`
		ProductType   string             `json:"productType"`
		Action        int                `json:"action"`
		Value         interface{}        `json:"value"`
		PetApproved   int                `json:"petApproved"`
		AttrData      *proUtil.Attribute `json:"attrData"`
	}
	var data Data
	err = json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	switch data.AttributeName {
	case "ty":
		tydata := proUtil.ProductAttrSystemUpdate{
			ProConfigId: syncer.ConfigId,
			AttrName:    data.AttributeName,
			AttrValue:   data.Value,
		}
		err = adapter.UpdateProductAttributeSystem(tydata)

	case "pet_status":
		tydata := proUtil.ProductAttrSystemUpdate{
			ProConfigId: syncer.ConfigId,
			AttrName:    data.AttributeName,
			AttrValue:   data.Value,
		}
		err = adapter.UpdateProductAttributeSystem(tydata)

	default:
		err = fmt.Errorf("%s: Unhandeled type [%s]", msg, data.AttributeName)
	}
	return err
}

func (syncer MySqlSyncer) JabongDiscount() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_jabongDiscount")
	defer logger.EndProfile(profiler, "product_mysql_sync_jabongDiscount")

	var jbngDisc proUtil.JabongDiscount
	err := json.Unmarshal(syncer.Data, &jbngDisc)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#JabongDiscount(): %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.UpdateJabongDiscount(jbngDisc)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#JabongDiscount(): %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) ProductSimpleStatusUpdate() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_statusUpdate")
	defer logger.EndProfile(profiler, "product_mysql_sync_statusUpdate")

	type Data struct {
		SimpleId int                       `json:"simpleId"`
		Status   string                    `json:"status"`
		Criteria proUtil.ProUpdateCriteria `json:"criteria"`
	}
	var data Data
	err := json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#ProductSimpleStatusUpdate()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.UpdateProductSimpleStatus(data.SimpleId, data.Status)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#ProductSimpleStatusUpdate()2: %s", err.Error())
	}
	err = adapter.UpdateProduct(syncer.ConfigId, data.Criteria)
	if err != nil {
		logger.Error(fmt.Sprintf("MySqlSyncer#ProductSimpleStatusUpdate()3: %s", err.Error()))
	}
	return nil
}

func (syncer MySqlSyncer) ProductSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_full")
	defer logger.EndProfile(profiler, "product_mysql_sync_full")

	var product proUtil.Product
	err := json.Unmarshal(syncer.Data, &product)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#ProductSync()1 Failed (configId->%d): %s", product.SeqId, err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.SaveProduct(product)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#ProductSync()2 Failed (configId->%d): %s", product.SeqId, err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) VideoStatusSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_videoStatus")
	defer logger.EndProfile(profiler, "product_mysql_sync_videoStatus")

	type Data struct {
		VideoId int    `json:"videoId"`
		Status  string `json:"status"`
	}
	var data Data
	err := json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#VideoStatusSync()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.UpdateVideoStatus(data.VideoId, data.Status)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#VideoStatusSync()2: %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) VideoSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_videoSync")
	defer logger.EndProfile(profiler, "product_mysql_sync_videoSync")

	type Data struct {
		ConfigId  int                  `json:"configId"`
		VideoData proUtil.ProductVideo `json:"videoData"`
	}
	var data Data
	err := json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#VideoSync()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	_, err = adapter.SaveVideo(data.ConfigId, data.VideoData)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#VideoSync()2: %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) DeleteImage() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_deleteImage")
	defer logger.EndProfile(profiler, "product_mysql_sync_deleteImage")

	var imageId int
	err := json.Unmarshal(syncer.Data, &imageId)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#DeleteImage()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	_, err = adapter.DeleteImage(imageId)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#DeleteImage()2: %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) AddImage() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_addImage")
	defer logger.EndProfile(profiler, "product_mysql_sync_addImage")

	type Data struct {
		ConfigId  int                  `json:"configId"`
		ImageData proUtil.ProductImage `json:"imageData"`
	}
	var data Data
	err := json.Unmarshal(syncer.Data, &data)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#AddImage()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	_, err = adapter.AddImage(data.ConfigId, data.ImageData)
	if err != nil {
		return fmt.Errorf("MySqlSyncer#AddImage()2: %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) PriceSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_priceSync")
	defer logger.EndProfile(profiler, "product_mysql_sync_priceSync")

	var priceUpdate proUtil.PriceUpdate
	err := json.Unmarshal(syncer.Data, &priceUpdate)
	if err != nil {
		return fmt.Errorf("(syncer MySqlSyncer)#PriceSync()1: %s", err.Error())
	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.UpdatePrice(priceUpdate)
	if err != nil {
		return fmt.Errorf("(syncer MySqlSyncer)#PriceSync()2: %s", err.Error())
	}
	return nil
}

func (syncer MySqlSyncer) ShipmentSync() error {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "product_mysql_sync_shipmentSync")
	defer logger.EndProfile(profiler, "product_mysql_sync_shipmentSync")

	type Shipment struct {
		Sku      string `json:"sku"`
		Shipment int    `json:"shipment"`
	}
	var shipment Shipment
	err := json.Unmarshal(syncer.Data, &shipment)
	if err != nil {
		return fmt.Errorf("(syncer MySqlSyncer)#ShipmentSync()1: %s", err.Error())

	}
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MYSQL)
	err = adapter.UpdateShipmentBySKU(shipment.Sku, shipment.Shipment)
	if err != nil {
		return fmt.Errorf("(syncer MySqlSyncer)#ShipmentSync()2: %s", err.Error())
	}
	return nil
}
