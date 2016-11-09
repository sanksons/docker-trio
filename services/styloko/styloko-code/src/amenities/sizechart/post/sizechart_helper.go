package post

import (
	sizeUtils "amenities/sizechart/common"
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
)

func createSizeChartName(brandId int, categoryId int, catalogType string) string {
	var sizeChartKey string
	categoryName := strings.Replace(getCategoryNameById(categoryId), " ", "-", -1)
	sizeChartKey = "category:" + categoryName + "_" + strconv.Itoa(categoryId)

	if brandId != 0 {
		brandName := strings.Replace(getBrandNameById(brandId), " ", "-", -1)
		sizeChartKey = sizeChartKey + "__brand:" + brandName
	}

	if catalogType != "all" && catalogType != "" {
		sizeChartKey = sizeChartKey + "__ty:" + catalogType
	}

	return sizeChartKey

}

// Get category Name  by ID
func getCategoryNameById(id int) string {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.GetCategoryNameById)
	type category struct {
		Name string
	}
	cat := category{}
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.GetCategoryNameById)
		mongoDriver.Close()
	}()
	mgoObj := mongoDriver.SetCollection(sizeUtils.CategoryCollection)
	err := mgoObj.Find(bson.M{"seqId": id}).Select(bson.M{"name": 1, "_id": 0}).One(&cat)
	if err != nil {
		return ""
	}
	return cat.Name
}

// Get brand name by Id
func getBrandNameById(id int) string {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.GetBrandNameById)
	type brand struct {
		Name string
	}
	brd := brand{}
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.GetBrandNameById)
		mongoDriver.Close()
	}()
	mgoObj := mongoDriver.SetCollection(sizeUtils.BrandCollection)
	err := mgoObj.Find(bson.M{"seqId": id}).Select(bson.M{"name": 1, "_id": 0}).One(&brd)
	if err != nil {
		return ""
	}
	return brd.Name
}

// Get typeId by url-key
func getTypeIdByUrlKey(ty string) (int, string) {
	var id int
	rowExist := false
	query := "SELECT id_catalog_ty from catalog_ty where url_key = '" + ty + "'"

	if ty == "all" {
		return 0, ""
	}
	driver, er := factory.GetMySqlDriver("sizechart")
	if er != nil {
		return 0, er.Error()
	}
	res, err := driver.Query(query)

	if err != nil {
		return 0, err.DeveloperMessage
	}
	for res.Next() {
		rowExist = true
		err := res.Scan(&id)
		if err != nil {
			res.Close()
			return 0, err.Error()
		}
	}
	res.Close()
	if rowExist == false {
		return 0, "Category Type doesnot exist for the given url-key"
	}
	return id, ""
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// save sizechart_active and column on mousehover for sizechart
// It is default for every sizechart and stored with categories
func SaveDefaultSettingForSizeChart(categoryId int) error {
	// category update API to update the category.
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.UpdateDefaultSetting)
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.UpdateDefaultSetting)
		mongoDriver.Close()
	}()
	mgoObj := mongoDriver.SetCollection(sizeUtils.CategoryCollection)
	criteria := bson.M{"seqId": categoryId}
	upadteVal := bson.M{"$set": bson.M{"szchrtActv": 1, "dispSzConv": sizeUtils.ColumnsOnMouseHover}}
	err := mgoObj.Update(criteria, upadteVal)

	if err != nil {
		logger.Error("#SaveDefaultSettingForSizeChart(): unable to store default sizechart setting", err.Error())
		return err
	}
	return nil
}
