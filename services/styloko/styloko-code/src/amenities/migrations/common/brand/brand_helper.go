package brand

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"fmt"
	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"strings"
)

//funtion takes input *core.Rows, gets data from mysql and
//gives map of Brand and brand ids as output
func processBrandRows(response interface{}) (map[int]*Brand, *string, error) {
	rows, ok := response.(*core.Rows)
	brand := make(map[int]*Brand)
	var bids string
	if ok {
		for rows.Next() {
			b := Brand{}
			err := rows.ScanStructByIndex(&b)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processBrandRows() %v for id %d", err.Error(), b.SeqId))
				rows.Close()
				return nil, nil, err
			}
			brand[b.SeqId] = &b
			bids += fmt.Sprintf("%d,", b.SeqId)
		}
		rows.Close()
	}
	bids = strings.TrimRight(bids, ",")
	return brand, &bids, nil
}

//function takes input *core.Rows, gets data from mysql and
//gives map of relatedBrands as output
func processRelatedBrandRows(response interface{}) (map[int][]RelatedBrand, error) {
	rows, ok := response.(*core.Rows)
	relatedBrand := make(map[int][]RelatedBrand)
	if ok {
		for rows.Next() {
			b := RelatedBrand{}
			err := rows.ScanStructByIndex(&b)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processRelatedBrandRows() %v for id %d", err.Error(), b.SeqId))
				rows.Close()
				return nil, err
			}
			relatedBrand[b.FkCatalogBrand] = append(relatedBrand[b.FkCatalogBrand], b)
		}
	}
	rows.Close()
	return relatedBrand, nil
}

//funtion takes input *core.Rows, gets data from mysql and
//gives map of brandCertificate as output
func processBrandCertificateRows(response interface{}) (map[int][]BrandCertificate, error) {
	rows, ok := response.(*core.Rows)
	brandCertificate := make(map[int][]BrandCertificate)
	if ok {
		for rows.Next() {
			b := BrandCertificate{}
			err := rows.ScanStructByIndex(&b)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processBrandCertificateRows() %v for id %d", err.Error(), b.SeqId))
				rows.Close()
				return nil, err
			}
			brandCertificate[b.FkCatalogBrand] = append(brandCertificate[b.FkCatalogBrand], b)
		}
	}
	rows.Close()
	return brandCertificate, nil
}

//This function writes map[int]BrandMongo into mongo
func writeBrandToMongo(brand map[int]BrandMongo) error {
	logger.Info("Starting to write Brand into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("BrandMigration")
	defer mgoSession.Close()
	for _, v := range brand {
		logger.Debug(fmt.Sprintf("Inserting seqId %d in mongo", v.SeqId))
		upsertVal := true
		updatedVal := map[string]interface{}{"$set": v}
		findCriteria := map[string]interface{}{"seqId": v.SeqId}
		_, err := mgoSession.FindAndModify(util.Brands, updatedVal, findCriteria, upsertVal)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Brand into Mongo %v", err.Error()))
			return err
		}
	}
	return nil
}

//This function takes input map[int]Brand ,map[int][]RelatedBrand and
//map[int][]BrandCertificate as input and places these into to final map[int]BrandMongo
//and returns it
func TransformBrand(brand map[int]*Brand, relatedBrand map[int][]RelatedBrand, brandCertificate map[int][]BrandCertificate) map[int]BrandMongo {
	brandNew := make(map[int]BrandMongo)
	for k, _ := range brand {
		tmp := BrandMongo{}
		tmp.SeqId = brand[k].SeqId
		tmp.Name = brand[k].Name
		tmp.Status = brand[k].Status
		tmp.Position = brand[k].Position
		tmp.UrlKey = brand[k].UrlKey
		tmp.ImageName = brand[k].ImageName
		tmp.BrandClass = brand[k].BrandClass
		tmp.IsExclusive = brand[k].IsExclusive
		tmp.BrandInfo = brand[k].BrandInfo
		tmp.RelatedBrand = relatedBrand[k]
		tmp.Certificate = brandCertificate[k]
		brandNew[k] = tmp

	}
	return brandNew
}

func EnsureIndexInDb() {
	logger.Info("Creating Indexes for new collection to be created")
	mgoSession := ResourceFactory.GetMongoSession("BrandMigration")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(util.Brands)
	var uniqueIndexes = []string{
		"seqId",
		"urlKey",
		"name",
	}
	for _, v := range uniqueIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: true,
			Sparse: true,
		})
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	logger.Info("New indexes created")
}

//This function creates indexes for the collection to be created if the collection does not exist
func checkAndEnsureIndex() error {
	flag := false
	mgoSession := ResourceFactory.GetMongoSession("BrandMigration")
	defer mgoSession.Close()
	logger.Info("Checking if collection already exists")
	colNames, err := mgoSession.CollectionExists()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting collection names from mongo :%s", err.Error()))
		return err
	}
	for _, v := range colNames {
		if v == util.Brands {
			flag = true
		}
	}
	if flag == true {
		logger.Info("Collection already exists so skipping creating indexes")
		return nil
	}
	EnsureIndexInDb()
	return nil
}
