package sizechart

import (
	"amenities/migrations/common/util"
	proUtils "amenities/products/common"
	"amenities/services/sizecharts"
	"common/ResourceFactory"
	"common/appconfig"
	"common/xorm/mysql"
	"fmt"
	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

type SizeChartMig struct {
	Id    int
	State bool
	Msg   string
}

type SizeChartErr struct {
	Id  int    `bson:"id"`
	Msg string `bson:"msg"`
}

func UpdateGivenProductWithSizechart(id int) error {
	var prd proUtils.Product
	mgoSession := ResourceFactory.GetMongoSession("SizeChartService")
	defer mgoSession.Close()
	mongoPrd := mgoSession.SetCollection(util.Products)

	err := mongoPrd.Find(bson.M{"seqId": id}).One(&prd)
	if err != nil {
		return err
	}
	sizeChart := sizecharts.GetSizeChartForProductMigration(prd)
	if sizeChart == nil {
		return fmt.Errorf(
			"UpdateGivenProductWithSizechart() (Migration)[productId:%d]#: SizeChart doesnot exist for this product",
			prd.SeqId,
		)
	}
	var ok bool
	prd.SizeChart, ok = sizeChart.(proUtils.ProSizeChart)
	if !ok {
		logger.Error(fmt.Sprintf(
			"#UpdateGivenProductWithSizechart() Migration: Sizechart assertion failed, looking for proUtil.ProSizeChart but found %v",
			reflect.TypeOf(sizeChart),
		))
	}
	err = prd.InsertOrUpdate("Mongo")
	if err != nil {
		return fmt.Errorf(
			"UpdateGivenProductWithSizechart() Migration: %s",
			err.Error(),
		)
	}
	return nil
}

// This script concurrently updates product with sizechart
func UpdateProductsWithSizeChartConc() error {
	mgoSession := ResourceFactory.GetMongoSession("SizeChartService")
	defer mgoSession.Close()
	conf := config.ApplicationConfig.(*appconfig.AppConfig)

	var prdArr []proUtils.Product
	mongoPrd := mgoSession.SetCollection(util.Products)
	ch := make(chan SizeChartMig)
	x := 0
	// This is No of concurrent products to fetch
	y := conf.ProductConLimit
	count := true

	for count {
		prdCount := 0
		err := mongoPrd.Find(bson.M{"status": "active", "petApproved": 1}).Skip(x).Limit(y).All(&prdArr)
		if err != nil {
			return err
		}
		if len(prdArr) == 0 {
			count = false
			continue
		}
		for _, prd := range prdArr {
			prdCount = prdCount + 1
			go func(prd proUtils.Product) {
				fmt.Println("Started: Adding sizechart for product ", prd.SeqId)
				ch <- AddSizeChartToProd(prd)
			}(prd)
		}
		for i := 0; i < prdCount; i++ {
			s := <-ch
			if !s.State {
				//insert in error collection
				mongodb := mgoSession.SetCollection("sizechartErrors")
				mongodb.Insert(
					SizeChartErr{Id: s.Id, Msg: s.Msg})
			}
			fmt.Println("Finished: Adding sizechart for product ", s.Id)
		}
		fmt.Println(fmt.Sprintf("Processed The CHUNK from %d to %d\n\n", x, x+y))
		x = x + y
	}
	return nil
}

func AddSizeChartToProd(prd proUtils.Product) SizeChartMig {
	sChStatus := SizeChartMig{}

	sizeChart, err := sizecharts.FetchSizeChart(prd)
	if sizeChart == nil {
		sChStatus.Id = prd.SeqId
		sChStatus.Msg = fmt.Sprintf("No SizeChart for: %d. Reason is %s.", prd.SeqId, err.Error())
		sChStatus.State = false
		return sChStatus
	}
	var ok bool
	prd.SizeChart, ok = sizeChart.(proUtils.ProSizeChart)
	if !ok {
		sChStatus.Id = prd.SeqId
		sChStatus.Msg = fmt.Sprintf(
			"Migration: Sizechart assertion failed, looking for proUtil.ProSizeChart but found %v",
			reflect.TypeOf(sizeChart),
		)
		sChStatus.State = false
		return sChStatus
	}
	err = prd.InsertOrUpdate("Mongo")
	if err != nil {
		sChStatus.Id = prd.SeqId
		sChStatus.Msg = fmt.Sprintf("Error in inserting/updating sizechart.%s", err.Error())
		sChStatus.State = false
		return sChStatus
	}
	return SizeChartMig{prd.SeqId, true, "SizeChart Migrated Succesfully"}
}

// This function updates the existing products from mongo with
// applicable sizechart
func UpdateProductsWithSizeChart() error {
	var prd proUtils.Product
	mgoSession := ResourceFactory.GetMongoSession("SizeChartService")
	defer mgoSession.Close()

	mongoPrd := mgoSession.SetCollection(util.Products)
	total, err := mongoPrd.Count()
	if err != nil {
		return err
	}
	// Loops over product and updates the sizechart
	for i := 0; i <= total; i++ {
		err := mongoPrd.Find(bson.M{"status": "active", "petApproved": 1}).Skip(i).Limit(1).One(&prd)
		if err != nil {
			continue
		}
		sizeChart := sizecharts.GetSizeChartForProductMigration(prd)
		if sizeChart == nil {
			logger.Error(fmt.Errorf(
				"UpdateProductsWithSizeChart() (Migration)[productId:%d]#: SizeChart doesnot exist for this product",
				prd.SeqId,
			))
			continue
		}
		var ok bool
		prd.SizeChart, ok = sizeChart.(proUtils.ProSizeChart)
		if !ok {
			logger.Error(fmt.Sprintf(
				"#UpdateProductsWithSizeChart() Migration: Sizechart assertion failed, looking for proUtil.ProSizeChart but found %v",
				reflect.TypeOf(sizeChart),
			))
			continue
		}
		err = prd.InsertOrUpdate("Mongo")
		if err != nil {
			logger.Error(fmt.Errorf(
				"UpdateProductsWithSizeChart() Migration: %s",
				err.Error(),
			))
			continue
		}
	}
	return nil
}

func GetAllSizeChartsMongo() {

	var configIds []int
	var sizeChartCollec []map[string]interface{}
	mgoSession := ResourceFactory.GetMongoSession("AttributeMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.SizeCharts)
	err := mongodb.Find(nil).All(&sizeChartCollec)

	if err != nil {
		logger.Error("Error while getting sizecharts from MongoDB", err.Error())
	}
	for _, v := range sizeChartCollec {
		logger.Info("")
		logger.Info("")
		logger.Info("####Processing the product update for sizechart with ID ", v["seqId"])
		configIds = FindSkusForSC(v["seqId"])
		logger.Info("The array of config for sizechart : ", configIds)
		if len(configIds) != 0 {
			UpdateSizeChartToProducts(configIds, v)
		}

	}
}

func GetPreviousScTypeForSku(configId int) int {
	var res interface{}
	sizeChartType := make(map[string]int)
	var alreadyScTy int
	var preTy string
	sizeChartType["sku"] = 0
	sizeChartType["brand"] = 1
	sizeChartType["brick"] = 2

	criteria := map[string]int{"seqId": configId}
	columns := bson.M{"sizeChart.data.sizeChartType": 1, "_id": 0}
	mgoSession := ResourceFactory.GetMongoSession("SizecharMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.Products)
	err := mongodb.Find(criteria).Select(columns).One(&res)
	if err != nil {
		logger.Info("Could not get previous sizechart type for sku ", configId, "# Error is ", err)
		return -2
	}
	if res == nil || len(res.(bson.M)) == 0 {
		preTy = ""
	} else {
		tempMap := res.(bson.M)
		tempNestedMap := tempMap["sizeChart"].(bson.M)
		tempNestedMap2 := tempNestedMap["data"].(bson.M)
		preTy = tempNestedMap2["sizeChartType"].(string)
	}

	if preTy == "" {
		alreadyScTy = -1
	} else {
		alreadyScTy = sizeChartType[preTy]
	}
	return alreadyScTy
}

func FindSkusForSC(idSizeChart interface{}) []int {
	var skus []int
	var sku int

	sql := getCatalogAdditionalInfoScSku(idSizeChart)
	data, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in running query %v", err.Error()))
		return skus
	}

	rows, _ := data.(*core.Rows)

	for rows.Next() {
		err := rows.Scan(&sku)
		if err != nil {
			logger.Info("Error occured in scanning")
		}
		skus = append(skus, sku)
	}
	rows.Close()
	return skus
}

func FindSkusForSizeChart(categoryId interface{}, brandId interface{}, typeId interface{}) []int {

	var skus []int
	var sku int

	sql := getCatalogConfigSql(categoryId, brandId, typeId)

	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		logger.Error(fmt.Sprintf("Error in running query %v", err.Error()))
		return skus
	}

	rows, e := data.(*core.Rows)
	if !e {

		return skus
	}
	for rows.Next() {
		err := rows.Scan(&sku)
		if err != nil {
			logger.Info("Error occured in scanning")
		}
		skus = append(skus, sku)
	}
	rows.Close()
	return skus
}

func GetSizeChartData(sizeChartId int) ([]SizeChartData, error) {

	sizeData := SizeChartData{}
	var sizeDataArr []SizeChartData
	sql := getSizeChartDataSql(sizeChartId)

	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		logger.Error(fmt.Sprintf("Error in running query %v", err.Error()))
		return sizeDataArr, nil
	}

	rows, ok := data.(*core.Rows)

	if ok {
		for rows.Next() {

			err := rows.ScanStructByIndex(&sizeData)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in getting row for sizechart() %v", err.Error()))
				logger.Info("Error in getting row----->>>", err.Error())
				rows.Close()
				return sizeDataArr, err
			}
			sizeDataArr = append(sizeDataArr, sizeData)
		}
		rows.Close()
	}
	return sizeDataArr, nil

}

func GetSKUSizechartMappingData() ([]SizechartMapping, error) {
	var MappingArr []SizechartMapping
	sql := getSkuSizechartMappingInfo()
	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		return nil, fmt.Errorf("#GetSKUSizechartMappingData():%s", err.Error())
	}
	rows, ok := data.(*core.Rows)

	if ok {
		for rows.Next() {
			mapSizeCh := SizechartMapping{}
			err := rows.ScanStructByIndex(&mapSizeCh)
			if err != nil {
				logger.Error("Error in getting sizechart mapping:%s", err.Error())
				rows.Close()
				return nil, fmt.Errorf("Unable to gets sizechart mapping: %s", err.Error())
			}
			MappingArr = append(MappingArr, mapSizeCh)
		}
		rows.Close()
	}
	return MappingArr, nil
}

func GetDistinctSizeChart() ([]SizeChart, error) {

	var SizeChartArr []SizeChart
	sql := getDistinctSizeChartSql()
	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		logger.Error("Error in fetching data", err.Error())
		logger.Info("Error in fetching data", err.Error())
	}

	rows, ok := data.(*core.Rows)

	if ok {
		for rows.Next() {
			// Get distinct size chart info
			sizeCh := SizeChart{}
			err := rows.ScanStructByIndex(&sizeCh)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in getting distinct sizechart() %v", err.Error()))
				logger.Info("Error in getting row", err.Error())
				rows.Close()
				return nil, err
			}
			SizeChartArr = append(SizeChartArr, sizeCh)
		}
		rows.Close()
	}

	return SizeChartArr, nil
}

// this fetches all the sizecharts from last inserted sizechart id
func GetSizeChartFromLastInserted(lastInserted int) ([]SizeChart, error) {
	var SizeChartArr []SizeChart
	sql := getDistinctSizeChartAfterId(lastInserted)
	data, err := mysql.GetInstance().Query(sql, false)

	SizeChartArr = []SizeChart{}
	if err != nil {
		logger.Error("#GetSizeChartFromLastInserted():Error in fetching data", err.Error())
		return SizeChartArr, err
	}
	rows, ok := data.(*core.Rows)

	if ok {
		for rows.Next() {
			// Get distinct size chart info
			sizeCh := SizeChart{}
			err := rows.ScanStructByIndex(&sizeCh)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in getting distinct sizechart() %v", err.Error()))
				rows.Close()
				return nil, err
			}
			SizeChartArr = append(SizeChartArr, sizeCh)
		}
		rows.Close()
	}

	return SizeChartArr, nil

}

// It fetches id of last insrted sizechart in mongo
func GetLastInsertedSizechartFromMongo() (int, error) {
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	var id interface{}
	mongoDb := mgoSession.SetCollection(util.SizeCharts)
	total, errCnt := mongoDb.Count()
	if errCnt != nil {
		logger.Error("#GetLastInsertedSizechartFromMongo(): Error in getting fetching count.", errCnt.Error())
		return -1, errCnt
	}
	if total == 0 && errCnt == nil {
		return 0, nil
	}
	err := mongoDb.Find(nil).Sort("-seqId").Select(bson.M{"seqId": 1, "_id": 0}).One(&id)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in gettinglast sizechart inserted", err.Error()))
		return -1, err
	}
	lastId, ok := id.(bson.M)
	if !ok {
		logger.Error("#GetLastInsertedSizechartFromMongo(): Type assertion failed")
		return -1, fmt.Errorf("#GetLastInsertedSizechartFromMongo(): Type assertion failed")
	}
	return lastId["seqId"].(int), nil
}

func updateToMongo(collection []SizeChartMongo) error {
	logger.Info("Starting to write Sizechart into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.SizeCharts)
	for _, doc := range collection {
		err := mongodb.Insert(doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting sizechart into Mongo %v", err.Error()))
			return err
		}
	}
	// Make indexes on mongodb collection sizecharts
	index := mgo.Index{
		Key:        []string{"categoryId", "brandId"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	indErr := mongodb.EnsureIndex(index)
	if indErr != nil {
		logger.Error("#updateToMongo(): Unable to create indexes on sizecharts.", indErr)
		return indErr
	}
	return nil
}

func writeToMongo(collection []SizeChartMongo) error {
	logger.Info("Starting to write Sizechart into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.SizeCharts)
	_ = mongodb.DropCollection()
	for _, doc := range collection {
		err := mongodb.Insert(doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting sizechart into Mongo %v", err.Error()))
		}
	}
	return nil
}

//This function creates indexes for the collection to be created if the collection does not exist
func checkAndEnsureIndex() error {
	flag := false
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	logger.Info("Checking if collection already exists")
	colNames, err := mgoSession.CollectionExists()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting collection names from mongo :%s", err.Error()))
		return err
	}
	for _, v := range colNames {
		if v == util.SizeCharts {
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

func EnsureIndexInDb() {
	logger.Info("Creating Indexes for new collection to be created")
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(util.SizeCharts)
	var normalIndexes = []string{
		"categoryId",
		"brandId",
	}
	for _, v := range normalIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: false,
			Sparse: false,
		})
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	logger.Info("New indexes created")
}
