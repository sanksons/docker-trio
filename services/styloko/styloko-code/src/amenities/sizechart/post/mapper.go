package post

import (
	ProUtils "amenities/products/common"
	"amenities/services/products"
	sizeUtils "amenities/sizechart/common"
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func getHeader(data []sizeUtils.SizeChartData) []string {
	var headers []string
	checkDup := make(map[string]int)

	for _, sizeChartData := range data {
		header := sizeChartData.ColumnHeader

		if _, ok := checkDup[header]; !ok {
			headers = append(headers, header)
			checkDup[header] = 1
		}
	}
	return headers
}

/*
This function updates the single product with sizechart data
*/
func UpdateProduct(sku string, preTy string, doc sizeUtils.SizeChartMongo) bool {
	sizeChartType := make(map[string]int)
	var alreadyScTy int
	sizeChartType["sku"] = 0
	sizeChartType["brand"] = 1
	sizeChartType["brick"] = 2

	if preTy == "" {
		alreadyScTy = -1
	} else {
		alreadyScTy = sizeChartType[preTy]
	}

	sChartPdp := CreateDisplayableSizeChart(sku, alreadyScTy, doc.SizeChartType, doc)
	if sChartPdp == nil {
		return false
	}
	sizeChProd := sizeUtils.ProdSChart{Id: doc.IdCatalogSizeChart, SizeChart: *sChartPdp}

	err := products.AddNode("SizeChart", sku, sizeChProd)

	if err != nil {
		logger.Error("mapper.go: #updateProduct() -> Error while updating sizechart to product " + sku + " : " + err.Error())
		return false
	}
	return true
}

// create a displayable PDP page sizechart for particular SKU
func CreateDisplayableSizeChart(sku string, preTy int, currTy int, doc sizeUtils.SizeChartMongo) *sizeUtils.SChart {
	k := 0
	mismatch := true
	upload := false
	chart := sizeUtils.SChart{}
	sizeIndexMapping := make(map[string]int)
	chart.Sizes = make(map[string][]string)

	// Type Mapping
	sizeChartTypeMap := make(map[int]string)
	sizeChartTypeMap = map[int]string{0: "sku", 1: "brand", 2: "brick"}

	// Rules to create sizechart for product,
	// depending upon the previously uploaded sizechart type and currently being uploaded
	switch {
	case preTy == -1, preTy == currTy, currTy < preTy:
		upload = true
	}

	if upload == false {
		logger.Info("mapper.go: #createDisplayableSizeChart() ->Sizechart of higher priority already exists for sku " + sku)
		return nil
	}

	sizeSku, err := products.GetSimpleSizesForSKU(sku, "mongo")
	if err != nil {
		logger.Error("mapper.go: #createDisplayableSizeChart() ->Sizes not found for the sku " + sku)
		return nil
	}

	chart.Headers = getHeader(doc.SizeChartInfo)
	chart.ImageName = doc.SizeChartImagePath
	chart.ScType = sizeChartTypeMap[currTy]

	for _, size := range sizeSku {
		chart.Sizes[strconv.Itoa(k)] = []string{size}
		sizeIndexMapping[size] = k
		k++
	}

	for _, data := range doc.SizeChartInfo {
		if data.RowHeaderType != "" {
			if index, ok := sizeIndexMapping[data.RowHeaderType]; ok {
				chart.Sizes[strconv.Itoa(index)] = append(chart.Sizes[strconv.Itoa(index)], data.Value)
				mismatch = false
			}
		}
	}
	if mismatch {
		logger.Info("mapper.go: #createDisplayableSizeChart() ->Mismatch of sizes for sku " + sku + ".")
		return nil
	}
	return &chart
}

func storeSizechartSkuMapping(sku string, sizechartId int) error {
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer mongoDriver.Close()

	mongoObj := mongoDriver.SetCollection(sizeUtils.SizeChartMappingCollec)
	doc := sizeUtils.SizeChartSkuMapping{sku, sizechartId}
	err := mongoObj.Insert(doc)
	if err != nil {
		return err
	}
	return nil
}

// This function finds the skus for sizechart and sizechart type previously uploadded to it
func getSkusForSizeChart(doc sizeUtils.SizeChartMongo) map[string]string {
	category := doc.FkCatalogCategory
	brand := doc.FkCatalogBrand
	ty := doc.FkCatalogTy

	mapSkuScType := make(map[string]string)
	var p []ProUtils.Product
	var criteria map[string]interface{}
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer mongoDriver.Close()
	mgoObj := mongoDriver.SetCollection(sizeUtils.ProductsCollection)

	if brand == 0 && ty == 0 {
		criteria = bson.M{"leaf": category}
	} else if brand != 0 && ty == 0 {
		criteria = bson.M{"leaf": category, "brandId": brand}
	} else {
		criteria = bson.M{"leaf": category, "brandId": brand, "ty": ty}
	}
	err := mgoObj.Find(criteria).Select(bson.M{"sku": 1, "_id": 0}).All(&p)

	if err != nil {
		logger.Error("mapper.go: #getSkusForSizeChart: Unable to find skus for sizechart")
		return mapSkuScType
	}
	// Check if product doesnt belong to sku level sizechart
	var prd ProUtils.Product
	mgoObjChk := mongoDriver.SetCollection(sizeUtils.SizeChartMappingCollec)

	for _, product := range p {
		errChk := mgoObjChk.Find(bson.M{"sku": product.SKU}).One(&prd)
		if errChk == nil {
			continue
		}
		if product.SizeChart.Data == nil {
			mapSkuScType[product.SKU] = ""
		} else {
			chart := product.SizeChart.Data.(sizeUtils.SChart)
			mapSkuScType[product.SKU] = chart.ScType
		}
	}
	return mapSkuScType
}
