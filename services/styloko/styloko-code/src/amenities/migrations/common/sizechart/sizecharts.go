package sizechart

import (
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func StartSKUSizechartMapping() error {
	logger.Info("Migarting mapping of skus and sizechart")
	err := MigrateSizechMapping()
	return err
}

func WriteSizeChartToProduct() error {
	fmt.Println("Writing size charts to product")
	UpdateProductsWithSizeChartConc()
	return nil
}

func SingleProductSizeChartUpdate(id int) error {
	err := UpdateGivenProductWithSizechart(id)
	if err != nil {
		return err
	}
	return nil
}

func StartSizeChartMigrationPartial() error {
	logger.Info("Starting SizeChart migration/ Updates from last inserted.")
	fmt.Println("Starting SizeChart migration/ Updates from last inserted.")
	lastId, err := GetLastInsertedSizechartFromMongo()
	if err != nil {
		logger.Error("Error in fetching last inserted sizechart from Mongo ", err.Error())
		return fmt.Errorf("Error in fetching last inserted sizechart from Mongo ", err.Error())
	}
	fmt.Println("The last inserted Id in sizechart collection -> ", lastId)
	if lastId == 0 {
		fmt.Println("Creating sizechart collection from scratch")
	} else {
		fmt.Println("Inserting sizechart after id : ", lastId)
	}
	charts, err := GetSizeChartFromLastInserted(lastId)
	if err != nil {
		logger.Error("Error in fetching the sizechart data from mysql:", err.Error())
		return fmt.Errorf("Error in fetching the sizechart data from mysql:", err.Error())
	}
	if len(charts) == 0 {
		logger.Info("Sizechart collection is already upto date.")
		fmt.Println("Sizechart collection is already upto date.")
	}
	collection := CreateSizeChartCollection(charts)
	err = updateToMongo(collection)
	return err
}

func StartSizeChartMigration() error {
	logger.Info("Starting size chart migration")
	err := checkAndEnsureIndex()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while checking existing index and ensuring indexes for db :%s", err.Error()))
		return err
	}
	val, _ := GetDistinctSizeChart()
	collection := CreateSizeChartCollection(val)
	err = writeToMongo(collection)
	return err
}

func CreateSizeChartCollection(sizeChartArr []SizeChart) []SizeChartMongo {
	var sizeChartCollection []SizeChartMongo

	for _, sizechart := range sizeChartArr {
		var data []SizeChartData
		// gets data for particular size chart
		data, err := GetSizeChartData(sizechart.IdCatalogDistinctSizeChart)

		if err != nil {
			logger.Error(fmt.Sprintf("#CreateSizeChartCollection():Error in getting sizechart data %v", err.Error()))
		}

		doc := createSizeChartMongoDoc(sizechart, data)
		sizeChartCollection = append(sizeChartCollection, doc)
	}

	return sizeChartCollection

}

func createSizeChartMongoDoc(sizeChart SizeChart, sizeData []SizeChartData) SizeChartMongo {
	doc := SizeChartMongo{}
	doc.IdCatalogSizeChart = sizeChart.IdCatalogDistinctSizeChart
	doc.FkCatalogCategory = sizeChart.FkCatalogCategory

	if sizeChart.FkCatalogBrand != nil {
		doc.FkCatalogBrand = *sizeChart.FkCatalogBrand
	}
	if sizeChart.FkCatalogTy != nil {
		doc.FkCatalogTy = *sizeChart.FkCatalogTy
	}
	doc.SizeChartName = sizeChart.SizeChartName
	if sizeChart.SizeChartType != nil {
		doc.SizeChartType = *sizeChart.SizeChartType
	} else {
		doc.SizeChartType = 1
	}
	doc.SizeChartImagePath = sizeChart.SizeChartImage
	if sizeChart.FkAclUser != nil {
		doc.FkAclUser = *sizeChart.FkAclUser
	}
	doc.SizeChartInfo = sizeData
	doc.CreatedAt = sizeChart.CreatedAt
	doc.UpdatedAt = sizeChart.UpdatedAt

	return doc

}
