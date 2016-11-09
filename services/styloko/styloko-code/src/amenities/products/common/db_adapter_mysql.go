package common

import (
	factory "common/ResourceFactory"
	"common/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jabong/floRest/src/common/sqldb"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type MySqlAdapter struct {
}

func (ma *MySqlAdapter) GetSession() sqldb.SqlDbInterface {
	s, _ := factory.GetMySqlDriver("MYSQL")
	return s
}

func (ma *MySqlAdapter) GetBySkus(skus []string, slice interface{}) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#GetBySkus(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetByIds(ids []int, slice interface{}) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#GetByIds(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetById(id int) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetById(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetBySku(sku string) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetBySku(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetByProductSet(id int) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetByProductSet(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductGroupByName(name string) (ProductGroup, error) {
	//@todo: To be implemented
	return ProductGroup{}, errors.New("(ma *MySqlAdapter)#GetProductGroupByName(): Not Implemented yet")

}

func (ma *MySqlAdapter) GetProductsByGroupId(id int) ([]Product, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductsByGroupId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductIdsBySellerId(int) ([]ProductSmall, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductIdsBySellerId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductsBySellerId(int) ([]Product, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductIdsBySellerId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductIdsByBrandId(id int) ([]ProductSmall, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductIdsBySellerId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductIdBySimpleSku(string) (ProductSmall, error) {
	return ProductSmall{}, errors.New("(ma *MySqlAdapter)#GetProductIdBySimpleSku(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductIdBySimpleId(int) (ProductSmall, error) {
	return ProductSmall{}, errors.New("(ma *MySqlAdapter)#GetProductIdBySimpleSku(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductsByBrandId(id int) ([]Product, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductIdsBySellerId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductsByCategoryId(id int) ([]Product, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductsByCategoryId(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductIdsByCategoryId(id int) ([]ProductSmall, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductIdsByCategoryId(): Not Implemented yet")

}

func (ma *MySqlAdapter) GetProductBySimpleId(id int) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetProductBySimpleId(): Not Implemented yet")

}

func (ma *MySqlAdapter) GetProductByVideoId(id int) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetProductByVideoId(): Not Implemented yet")

}

func (ma *MySqlAdapter) GetProductBySkuAndType(productSku string, productType string) (Product, error) {
	//@todo: To be implemented
	return Product{}, errors.New("(ma *MySqlAdapter)#GetProductBySkuAndType(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetAtrributeByCriteria(attrSrch AttrSearchCondition) (AttributeMongo, error) {
	//@todo: To be implemented
	return AttributeMongo{}, errors.New("(ma *MySqlAdapter)#GetAtrributeByCriteria(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProductsForSeller(sellers []int, limit int, offset int, lastSCId int) ([]Product, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductsForSeller(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetProAttributeSetById(id int) (ProAttributeSet, error) {
	as := ProAttributeSet{}
	//@todo: To be implemented
	return as, errors.New("(ma *MySqlAdapter)#GetProAttributeSetById(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetCategoriesByIds(cats []int) ([]Category, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetCategoriesByIds(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetAttributeMongoById(seqId int) (AttributeMongo, error) {
	tmpAttr := AttributeMongo{}
	//@todo: To be implemented
	return tmpAttr, errors.New("(ma *MySqlAdapter)#GetAttributeMongoById(): Not Implemented yet")
}

func (ma *MySqlAdapter) FindPrimaryCategoryId(catIds []int) (int, error) {
	//@todo: To be implemented
	return 0, errors.New("(ma *MySqlAdapter)#FindPrimaryCategoryId(): Not Implemented yet")
}

func (ma *MySqlAdapter) SaveProduct(pro Product) error {
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#SaveProduct: Could not get connection")
	}
	tx, err := session.GetTxnObj()
	if err != nil {
		return errors.New("(ma *MySqlAdapter)#SaveProduct: " + err.DeveloperMessage)
	}
	var isCommited bool
	//Rollback in defer.
	defer func() {
		if !isCommited && (tx != nil) {
			er := tx.Rollback()
			if er != nil {
				logger.Error(fmt.Sprintf(
					"(ma *MySqlAdapter)#SaveProduct(pro Product)Rollback failed:%s",
					er.Error(),
				))
			} else {
				logger.Error(fmt.Sprintf(
					"(ma *MySqlAdapter)#SaveProduct(pro Product)Rollback success",
				))
			}
		}
	}()
	sync := MySqlSync{
		TxnObj: tx,
	}
	er := sync.SaveProduct(pro)
	if er != nil {
		return errors.New("(ma *MySqlAdapter)#SaveProduct: " + er.Error())
	}
	tx.Commit()
	isCommited = true
	return nil
}

func (ma *MySqlAdapter) AddNode(string, string, interface{}) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#AddNode(): Not Implemented yet")
}

func (ma *MySqlAdapter) DeleteNode(string, string) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#DeleteNode(): Not Implemented yet")
}

func (ma *MySqlAdapter) DeleteImage(imageId int) (int, error) {
	session := ma.GetSession()
	if session == nil {
		return 0, fmt.Errorf("(ma *MySqlAdapter)#DeleteImage: Could not get connection")
	}
	sql := `DELETE FROM catalog_product_image WHERE id_catalog_product_image=?`
	_, err := session.Execute(sql, imageId)
	if err != nil {
		return imageId, errors.New("(ma *MySqlAdapter)#DeleteImage():" + err.DeveloperMessage)
	}
	return imageId, nil
}

func (ma *MySqlAdapter) AddImage(configId int, pi ProductImage) (int, error) {
	session := ma.GetSession()
	if session == nil {
		return 0, fmt.Errorf("(ma *MySqlAdapter)#AddImage: Could not get connection")
	}
	sql := `INSERT INTO catalog_product_image (
	id_catalog_product_image,
	fk_catalog_config,
	image,
	main,
	orientation,
	updated_at,
	original_filename
	) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE fk_catalog_config = VALUES(fk_catalog_config),
	image = VALUES(image), main = VALUES(main), orientation = VALUES(orientation), updated_at =
	VALUES(updated_at), original_filename = VALUES(original_filename)`
	_, err := session.Execute(sql,
		pi.SeqId, configId, pi.ImageNo, pi.Main,
		pi.Orientation, pi.UpdatedAt, pi.OriginalFileName,
	)
	if err != nil {
		return pi.SeqId, errors.New("(ma *MySqlAdapter)#AddImage(): " + err.DeveloperMessage)
	}

	if pi.IfUpdate {
		sql = `UPDATE catalog_config
		SET updated_at = ?,
	  approved_at = ?,
	  activated_at = ?
		WHERE id_catalog_config = ?`
		_, err = session.Execute(sql, pi.UpdatedAt, pi.UpdatedAt, pi.UpdatedAt, configId)
		if err != nil {
			return pi.SeqId, fmt.Errorf("(ma *MySqlAdapter)AddImage:%s", err.DeveloperMessage)
		}
	}
	return pi.SeqId, nil
}

func (ma *MySqlAdapter) SaveVideo(configId int, video ProductVideo) (int, error) {
	session := ma.GetSession()
	if session == nil {
		return 0, fmt.Errorf("(ma *MySqlAdapter)#SaveVideo: Could not get connection")
	}
	txn, txnErr := session.GetTxnObj()
	if txnErr != nil {
		return 0, fmt.Errorf("(ma *MySqlAdapter)#SaveVideo: %s", txnErr.DeveloperMessage)
	}
	query := `Select id_video
		FROM video
		WHERE id_video=?`
	row := txn.QueryRow(query, video.Id)
	var videoId int
	row.Scan(&videoId)
	var err error
	if videoId <= 0 {
		sql1 := `INSERT INTO video (
		id_video, file_name, thumbnail, size, duration,
		status, created_at, updated_at
	) VALUES (
		?,?,?,?,?,?,?,?
	)`
		_, err = txn.Exec(sql1,
			video.Id, video.FileName, video.Thumbnail, video.Size, video.Duration,
			video.Status, video.CreatedAt, video.UpdatedAt,
		)
	} else {
		sql1 := `UPDATE video SET file_name=?, thumbnail=?, size=?,
	    duration=?, status=?, created_at=?, updated_at=? WHERE id_video=?`
		_, err = txn.Exec(sql1,
			video.FileName, video.Thumbnail, video.Size, video.Duration,
			video.Status, video.CreatedAt, video.UpdatedAt, video.Id,
		)
	}
	if err != nil {
		txn.Rollback()
		return 0, fmt.Errorf("(ma *MySqlAdapter)#SaveVideo():%s", err.Error())
	}
	type PV struct {
		IdProductVideo  int
		FkCatalogConfig int
		FkVideo         int
	}
	var pv PV
	searchSql := `Select id_product_video, fk_catalog_config, fk_video
		FROM product_video
		WHERE fk_catalog_config=? AND fk_video=?`
	row = txn.QueryRow(searchSql, configId, video.Id)

	err = row.Scan(&pv.IdProductVideo, &pv.FkCatalogConfig, &pv.FkVideo)
	if err != nil && err != sql.ErrNoRows {
		txn.Rollback()
		return 0, fmt.Errorf("(ma *MySqlAdapter)#SaveVideo():%s", err.Error())
	}
	if pv.IdProductVideo > 0 {
		txn.Commit()
		return video.Id, nil
	}
	insertSql := `INSERT INTO product_video (fk_catalog_config, fk_video) VALUES(?, ?)`
	_, err = txn.Exec(insertSql, configId, video.Id)
	if err != nil {
		txn.Rollback()
		return video.Id, errors.New("(ma *MySqlAdapter)#SaveVideo():" + err.Error())
	}
	txn.Commit()
	return video.Id, nil
}

func (ma *MySqlAdapter) UpdateVideoStatus(videoId int, status string) error {
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateVideoStatus: Could not get connection")
	}
	sql := `UPDATE video SET status=? WHERE id_video=?`
	_, err := session.Execute(sql, status, videoId)
	if err != nil {
		return errors.New("(ma *MySqlAdapter)#UpdateVideoStatus():" + err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) UpdateProductAttribute(prdctUpdtCndn PrdctAttrUpdateCndtn) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#UpdateProductAttribute(): Not Implemented yet")
}

func (ma *MySqlAdapter) UpdatePrice(data PriceUpdate) error {
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdatePrice: Could not get connection")
	}
	tx, txErr := session.GetTxnObj()
	if txErr != nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdatePrice()Cannot get transaction object: %s",
			txErr.DeveloperMessage)
	}

	var (
		isCommited bool
		err        error
		sql        string
	)
	//Rollback in defer
	defer func() {
		if !isCommited {
			er := tx.Rollback()
			if er != nil {
				logger.Error(fmt.Sprintf(
					"(ma *MySqlAdapter)#UpdatePrice()Rollback failed:%s",
					er.Error(),
				))
			}
		}
	}()

	var priceUpdStmt string = ` price=price `
	if data.Price > 0 {
		priceUpdStmt = ` price=` + utils.ToString(data.Price)
	}
	if !data.UpdateSP {
		sql = `UPDATE catalog_simple
	   		  SET
	   	        ` + priceUpdStmt + `
	   		  WHERE
	   		    id_catalog_simple=?`
		_, sqlErr := tx.Exec(sql, data.SimpleId)
		if sqlErr != nil {
			err = fmt.Errorf("%s", sqlErr.Error())
		}
	} else {
		var (
			specialPrice    *string
			specialFromDate *string
			specialToDate   *string
		)
		if data.SpecialPrice == nil {
			specialPrice = nil
			specialFromDate = nil
			specialToDate = nil
		} else {
			sp := utils.ToString(data.SpecialPrice)
			specialPrice = &sp
			sfd := ToMySqlTimeNull(data.SpecialFromDate)
			specialFromDate = sfd
			std := ToMySqlTimeNull(data.SpecialToDate)
			specialToDate = std
		}
		sql = `UPDATE catalog_simple
	            SET
	               ` + priceUpdStmt + ` ,
	              special_price=?,
	              special_from_date=?,
	              special_to_date=?
	            WHERE id_catalog_simple=?`
		_, sqlErr := tx.Exec(sql, specialPrice, specialFromDate,
			specialToDate, data.SimpleId)
		if sqlErr != nil {
			err = fmt.Errorf("%s", sqlErr.Error())
		}
	}
	if err != nil {
		tx.Rollback()
		return errors.New("(ma *MySqlAdapter)#UpdatePrice(): " + err.Error())
	}
	tx.Commit()
	isCommited = true
	return nil
}

func (ma *MySqlAdapter) UpdateShipmentBySKU(sku string, shipmentType int) error {
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateShipmentBySKU: Could not get connection")
	}
	sql := `UPDATE catalog_config SET fk_catalog_shipment_type=? WHERE sku=?`
	_, sqlErr := session.Execute(sql, shipmentType, sku)
	if sqlErr != nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateShipmentBySKU():%s", sqlErr.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) GenerateNextSequence(string) (int, error) {
	//@todo: To be implemented
	return 0, errors.New("(ma *MySqlAdapter)#GenerateNextSequence(): Not Implemented yet")
}

func (ma *MySqlAdapter) UpdateProductStatus(configid int, status string) error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#UpdateProductStatus(): Not Implemented yet")
}

func (ma *MySqlAdapter) UpdateProductSimpleStatus(id int, status string) error {
	//@todo: To be implemented
	///@todo: write impl
	sql := `UPDATE catalog_simple SET status = ? WHERE id_catalog_simple = ?`
	session := ma.GetSession()
	_, err := session.Execute(sql, status, id)
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateProductSimpleStatus: Could not get connection")
	}
	if err != nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateProductSimpleStatus:%s", err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) UpdateJabongDiscount(disc JabongDiscount) error {
	sql := `UPDATE catalog_simple SET jabong_discount=?, jabong_discount_from=?,
		jabong_discount_to=? WHERE id_catalog_simple=?`

	var discount *float64
	var from *string
	var to *string
	if disc.Discount > 0 {
		discount = &disc.Discount
		fr := ToMySqlTime(disc.FromDate)
		from = &fr
		t := ToMySqlTime(disc.ToDate)
		to = &t
	}
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateJabongDiscount: Could not get connection")
	}
	_, err := session.Execute(sql, discount, from, to, disc.SimpleId)
	if err != nil {
		return fmt.Errorf(
			"(ma *MySqlAdapter)#UpdateJabongDiscount:%s", err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) SetPetApproved(configId int, petApproved int) error {
	sql := `UPDATE catalog_config SET pet_approved=? WHERE id_catalog_config=?`
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#SetPetApproved: Could not get connection")
	}
	_, err := session.Execute(sql, petApproved, configId)
	if err != nil {
		return fmt.Errorf("(ma *MySqlAdapter) SetPetApproved:%s", err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) UpdateProductAttributeSystem(input ProductAttrSystemUpdate) error {
	var sql string
	var msg string = "(ma *MySqlAdapter)#UpdateProductAttributeSystem()"
	switch input.AttrName {
	case SYSTEM_TY:
		strVal, ok := input.AttrValue.(string)
		if !ok {
			return fmt.Errorf("%s : %s", msg, "[ty] string assrt failed")
		}
		sql = `UPDATE catalog_config
					SET fk_catalog_ty= (
						SELECT id_catalog_ty FROM catalog_ty
						WHERE name = "` + strVal + `"
						)
					WHERE id_catalog_config=?;`
	case SYSTEM_PET_STATUS:
		strVal, ok := input.AttrValue.(string)
		if !ok {
			return fmt.Errorf("%s : %s", msg, "[pet_status] string assrt failed")
		}
		sql = `UPDATE catalog_config
					SET pet_status= "` + strVal + `"
					WHERE id_catalog_config=?;`
	default:
	}
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateProductAttributeSystem: Could not get connection")
	}
	_, err := session.Execute(sql, input.ProConfigId)
	if err != nil {
		logger.Error(sql)
		return fmt.Errorf("%s:%s", msg, err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) UpdateProduct(configId int, criteria ProUpdateCriteria) error {
	var sql string
	sql = `UPDATE catalog_config SET `
	pieces := []string{}
	if criteria.ActivatedAt.Isset {
		var tmp string
		if criteria.ActivatedAt.Value == nil {
			tmp = fmt.Sprintf("activated_at = NULL")
		} else {
			tmp = fmt.Sprintf("activated_at = '%s'", criteria.ActivatedAt.GetStringValue())
		}
		pieces = append(pieces, tmp)
	}

	if criteria.PetApproved.Isset {
		var tmp string
		if criteria.PetApproved.Value == nil {
			tmp = fmt.Sprintf("`pet_approved` = NULL")
		} else {
			tmp = fmt.Sprintf("`pet_approved` =%d",
				*criteria.PetApproved.Value,
			)
		}
		pieces = append(pieces, tmp)
	}
	if criteria.Status.Isset {
		var tmp string
		if criteria.Status.Value == nil {
			tmp = fmt.Sprintf("`status` = NULL")
		} else {
			tmp = fmt.Sprintf("`status` = '%s'", criteria.Status.GetStringValue())
		}
		pieces = append(pieces, tmp)
	}
	if len(pieces) <= 0 {
		return nil
	}
	sql = fmt.Sprintf("%s %s WHERE id_catalog_config=%d",
		sql, strings.Join(pieces, ","), configId,
	)
	session := ma.GetSession()
	if session == nil {
		return fmt.Errorf("(ma *MySqlAdapter)#UpdateProduct: Could not get connection")
	}
	_, err := session.Execute(sql)
	if err != nil {
		return fmt.Errorf("%s:%s", "(ma *MySqlAdapter)#UpdateProduct", err.DeveloperMessage)
	}
	return nil
}

func (ma *MySqlAdapter) GetAttributeMapping(name string) (AttrMapping, error) {
	//@todo: To be implemented
	return AttrMapping{}, errors.New("(ma *MySqlAdapter)#GetAttributeMapping(): Not Implemented yet")
}

func (ma *MySqlAdapter) GetAttributeMongoByName(name string) (AttributeMongo, error) {
	//@todo: To be implemented
	return AttributeMongo{}, errors.New("(ma *MySqlAdapter)#GetAttributeMongoByName(): Not Implemented yet")
}

func (ma *MySqlAdapter) ResetSSRCounter() error {
	//@todo: To be implemented
	return errors.New("(ma *MySqlAdapter)#ResetSSRCounter(): Not Implemented yet")

}

func (ma *MySqlAdapter) GetProductBySellerIdSku(sellerId int,
	sellerSkuArr []string) ([]ProductSmallSimples, error) {
	//@todo: To be implemented
	return nil, errors.New("(ma *MySqlAdapter)#GetProductSellerIdSku(): Not Implemented yet")
}
