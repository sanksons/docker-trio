package common

import (
	"common/ResourceFactory"
	"gopkg.in/mgo.v2/bson"
)

// This function returns applicable sizechart for given sku, and sizechart type
func CheckGivenSizeChartForProduct(catId int, brandId int, ty int, sizeChartTy int) *SizeChartMongo {
	// ty = 0 means ty doesnot exist for the product
	var criteria map[string]interface{}
	var res SizeChartMongo
	mongoDriver := ResourceFactory.GetMongoSession(SizeChartSer)
	defer mongoDriver.Close()
	mgoObj := mongoDriver.SetCollection(SizeChartCollec)
	if ty == 0 {
		criteria = bson.M{"categoryId": catId, "brandId": brandId, "sizeChartType": sizeChartTy}
	} else {
		criteria = bson.M{"categoryId": catId, "brandId": brandId, "tyId": ty, "sizeChartType": sizeChartTy}
	}
	err := mgoObj.Find(criteria).Sort("-updatedAt").Limit(1).One(&res)

	if err != nil {
		return nil
	}
	return &res
}

// This function returns the sizechart by id
func GetSizeChartById(sizeChartId int) *SizeChartMongo {
	var res SizeChartMongo
	mongoDriver := ResourceFactory.GetMongoSession(SizeChartSer)
	defer mongoDriver.Close()
	mgoObj := mongoDriver.SetCollection(SizeChartCollec)

	err := mgoObj.Find(bson.M{"seqId": sizeChartId}).One(&res)
	if err != nil {
		return nil
	}
	return &res
}

// It check and get the SKU level sizechart for product
func GetSkuLevelSizeChartProd(sku string) *SizeChartMongo {
	mongoDriver := ResourceFactory.GetMongoSession(SizeChartSer)
	defer mongoDriver.Close()
	type Data struct {
		SeqId int `bson:"sizechartId"`
	}
	var d Data

	mgoObj := mongoDriver.SetCollection(SizeChartMappingCollec)
	err := mgoObj.Find(bson.M{"sku": sku}).Select(bson.M{"sizechartId": 1, "_id": 0}).Sort("-sizechartId").Limit(1).One(&d)
	if err == nil {
		// sku level sizechart exists for product
		res := GetSizeChartById(d.SeqId)
		return res
	}
	// either error occcured or no sizechart exists
	return nil
}
