package post

import (
	sizeUtils "amenities/sizechart/common"
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
	"github.com/jabong/floRest/src/common/utils/logger"
	"reflect"
	"time"
)

//This func create a collection data for sizechart
func createSizeChartCollec(sizeChart sizeUtils.SizeChart, brandWiseAllSc interface{}, ty string) ([]sizeUtils.SizeChartMongo, string) {
	var collection []sizeUtils.SizeChartMongo
	var mongoDoc sizeUtils.SizeChartMongo
	var chartTy int
	var mongoDriver *mongodb.MongoDriver
	allSC := reflect.ValueOf(brandWiseAllSc)
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer mongoDriver.Close()

	for i := 0; i < allSC.Len(); i++ {

		brandWiseData := allSC.Index(i)

		brandWise := brandWiseData.Interface().(sizeUtils.BrandWiseScData)
		mongoDoc.IdCatalogSizeChart = mongoDriver.GetNextSequence(sizeUtils.SizeChartCollec)
		mongoDoc.SizeChartName = createSizeChartName(brandWise.BrandId, sizeChart.CategoryId, sizeChart.CatalogType)
		mongoDoc.FkCatalogBrand = brandWise.BrandId
		mongoDoc.FkCatalogCategory = sizeChart.CategoryId

		// check if brand is empty, then sizechart is brick level ie generic
		// otherwise sizechart is specific ie brand level
		if brandWise.BrandId == 0 {
			chartTy = sizeUtils.SizeChartLevelGeneric
		} else {
			chartTy = sizeUtils.SizeChartLevelSpecific
		}
		// if header SizeChart-Type is SKU
		if ty == sizeUtils.SKU {
			chartTy = sizeUtils.SizeChartLevelSku
		}
		mongoDoc.SizeChartType = chartTy

		if ty != sizeUtils.SKU {
			categoryTy, err := getTypeIdByUrlKey(sizeChart.CatalogType)
			if err != "" {
				return collection, err
			}
			mongoDoc.FkCatalogTy = categoryTy
		}
		mongoDoc.SizeChartImagePath = sizeChart.ImageName
		mongoDoc.FkAclUser = sizeChart.AclUserId
		mongoDoc.UpdatedAt = time.Now()
		mongoDoc.CreatedAt = time.Now()
		mongoDoc.SizeChartInfo = brandWise.ScData

		collection = append(collection, mongoDoc)
	}

	return collection, ""
}

func storeSizeChartMongo(collection []sizeUtils.SizeChartMongo) bool {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.SizeChartCreationMongo)
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.SizeChartCreationMongo)
		mongoDriver.Close()
	}()
	for _, doc := range collection {
		err := mongoDriver.Insert(sizeUtils.SizeChartCollec, doc)

		if err != nil {
			return false
		}
	}
	return true
}
