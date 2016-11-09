package common

import (
	factory "common/ResourceFactory"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"time"
)

type Look struct {
	StlCreatedAt *int64   `json:"stl_created_at"`
	StlPrice     *float64 `json:"stl_price"`
}

// Shop Look detials for solr
// ConfigID: Primary key of product
// Returns:
// look Look
func GetLookDetail(pro *Product) (*Look, error) {
	var look Look
	var err error
	look.StlCreatedAt, err = getStlCreatedAt(pro.SeqId)
	if err != nil {
		return nil, err
	}
	if look.StlCreatedAt != nil {
		look.StlPrice, _ = getStlPrice(pro)
	}
	return &look, nil
}

// Shop Look created at by  giving config id of look
// ConfigID: Primary key of product
// Returns:
// createdAt int
func getStlCreatedAt(configId int) (*int64, error) {
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		return nil, fmt.Errorf(
			"getStlCreatedAt(): Cannot get mysql driver: %s",
			err.Error(),
		)
	}
	sql := `SELECT
                csl.created_at AS stl_created_at
            FROM catalog_shop_look_detail AS csld
            INNER JOIN catalog_shop_look AS csl
                ON csld.fk_catalog_shop_look = csl.id_catalog_shop_look
            INNER JOIN catalog_config cc 
	            ON cc.sku = csld.sku
            WHERE (csl.is_active = 1)
            AND (csl.type = 0)
            AND (cc.id_catalog_config = ?)
            AND (is_primary = 1)
            ORDER BY csl.created_at DESC
            LIMIT 1`
	result, sqlerr := driver.Query(sql,
		configId,
	)
	defer result.Close()
	if sqlerr != nil {
		return nil, fmt.Errorf(
			"getStlCreatedAt(): Cannot get mysql driver: %s",
			sqlerr.DeveloperMessage,
		)
	}

	var stlCreatedAt *int64
	for result.Next() {
		var err error
		var CreatedAt *time.Time
		err = result.Scan(&CreatedAt)
		if err != nil {
			return nil, err
		}
		stlCr := CreatedAt.Unix()
		stlCreatedAt = &stlCr
	}
	return stlCreatedAt, nil
}

// Shop Look Price by  giving config id of look
// ConfigID: Primary key of product
// Returns:
// price double
func getStlPrice(pro *Product) (*float64, error) {
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		return nil, fmt.Errorf(
			"getStlCreatedAt(): Cannot get mysql driver: %s",
			err.Error(),
		)
	}
	var StlPrice float64
	//get Price for primary sku
	StlPrice = getSkuPrice(pro)
	sql := `SELECT
	           cc.id_catalog_config AS configId
	           FROM catalog_shop_look_detail AS csld
	           INNER JOIN catalog_shop_look AS csl
	               ON csld.fk_catalog_shop_look = csl.id_catalog_shop_look
	           INNER JOIN catalog_config cc 
	               ON cc.sku = csld.sku
	           WHERE (csl.is_active = 1)
	           AND (csl.type = 0)
	           AND (csld.is_primary = 0)
	           AND (cc.status = 'active')
	           AND (cc.pet_approved = '1')
	           AND csld.fk_catalog_shop_look
	           IN (SELECT
	               csld1.fk_catalog_shop_look
	           FROM catalog_shop_look_detail AS csld1
	           WHERE sku = ?
	           AND csld1.is_primary = 1)`
	result, sqlerr := driver.Query(sql, pro.SKU)
	defer result.Close()
	if sqlerr != nil {
		return nil, fmt.Errorf(
			"getStlCreatedAt(): Cannot get mysql driver: %s",
			sqlerr.DeveloperMessage,
		)
	}
	var configID int
	for result.Next() {
		err := result.Scan(&configID)
		if err != nil {
			logger.Error(fmt.Sprintf(
				"Error while fetching data and secondary sku %s",
				err.Error(),
			))
			return &StlPrice, nil
		}
		pro1, err := GetAdapter(DB_ADAPTER_MONGO).GetById(configID)
		StlPrice = StlPrice + getSkuPrice(&pro1)
	}
	return &StlPrice, nil
}

// return price information to the product by comparsion with special price
// ConfigID: Primary key of product
// Returns:
// price double
func getSkuPrice(pro *Product) float64 {
	var price float64
	var specialPrice float64
	priceInfo := pro.PreparePriceMap()
	price = priceInfo.Price

	if priceInfo.SpecialPrice != nil {
		specialPrice = *priceInfo.SpecialPrice
	}
	if specialPrice > 0 && specialPrice < price {
		price = specialPrice
	}
	return price
}
