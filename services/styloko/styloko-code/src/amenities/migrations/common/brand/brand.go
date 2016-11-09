package brand

import (
	"common/xorm/mysql"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func StartBrandMigration() error {
	//function to store brand table
	err := brandTable()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while migrating brand table :%s", err.Error()))
		return err
	}
	return nil
}

//function that integrates brand details,transforms it and writes it to mongo
func brandTable() error {
	logger.Info("Started Migrating Brand Table")
	err := checkAndEnsureIndex()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while checking existing index and ensuring indexes for db :%s", err.Error()))
		return err
	}
	brand, bids, err := getBrandInfo()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting brand info :%s", err.Error()))
		return err
	}
	relatedBrand, er := getRelatedBrand(*bids)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while getting related brands :%s", err.Error()))
		return er
	}
	brandCertificate, e := getBrandCertificate(*bids)
	if e != nil {
		logger.Error(fmt.Sprintf("Error while getting brand certificates :%s", err.Error()))
		return e
	}
	b := TransformBrand(brand, relatedBrand, brandCertificate)
	err = writeBrandToMongo(b)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while writing brand to mongo :%s", err.Error()))
		return err
	}
	return nil
}

//function that gets brand information from catalog_brand
func getBrandInfo() (map[int]*Brand, *string, error) {
	logger.Info("Preparing Brand Info")
	sql := getBrandSql()
	logger.Debug(fmt.Sprintf("Sql being fired to get brand info :%s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		brand, bids, er := processBrandRows(rows)
		if er == nil {
			return brand, bids, nil
		}
		return nil, nil, er
	}
	return nil, nil, err
}

//function that gets RelatedBrands from catalog_related_brand
func getRelatedBrand(bids string) (map[int][]RelatedBrand, error) {
	logger.Info("Preparing Related Brand")
	sql := getRelatedBrandSql(bids)
	logger.Debug(fmt.Sprintf("Sql being fired to get related brands :%s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		relatedBrand, er := processRelatedBrandRows(rows)
		if er == nil {
			return relatedBrand, nil
		}
		return nil, er
	}
	return nil, err
}

//function that gets BrandCertificates from catalog_brand_certificate
func getBrandCertificate(bids string) (map[int][]BrandCertificate, error) {
	logger.Info("Preparing Brand Certificate")
	sql := getBrandCertificateSql(bids)
	logger.Debug(fmt.Sprintf("Sql being fired to get brand certificate :%s", sql))
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		brandCertificate, er := processBrandCertificateRows(rows)
		if err == nil {
			return brandCertificate, nil
		}
		return nil, er
	}
	return nil, err
}
