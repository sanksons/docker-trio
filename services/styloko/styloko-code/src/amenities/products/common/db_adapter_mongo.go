package common

import (
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
	// "common/utils"
	"errors"
	"fmt"
	// "github.com/jabong/floRest/src/common/monitor"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
)

type MongoAdapter struct {
}

func (ma *MongoAdapter) GetSession() *mongodb.MongoDriver {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetSession")
	defer func() {
		logger.EndProfile(profiler, "GetSession")
	}()
	return factory.GetMongoSession(PRODUCT_COLLECTION)
}

func (ma *MongoAdapter) GetBySkus(skus []string, slice interface{}) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetBySkus")
	defer func() {
		logger.EndProfile(profiler, "GetBySkus")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	err := mSession.SetCollection(PRODUCT_COLLECTION).
		Find(M{"sku": M{"$in": skus}}).
		Sort("seqId").
		All(slice)
	return err
}

func (ma *MongoAdapter) GetByIds(ids []int, slice interface{}) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetByIds")
	defer func() {
		logger.EndProfile(profiler, "GetByIds")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	err := mSession.SetCollection(PRODUCT_COLLECTION).
		Find(M{"seqId": M{"$in": ids}}).
		Sort("seqId").
		All(slice)
	return err
}

func (ma *MongoAdapter) GetById(id int) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetById")
	defer func() {
		logger.EndProfile(profiler, "GetById")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"seqId": id,
	}).One(&p)
	if err == mgo.ErrNotFound {
		return p, NotFoundErr
	}
	if err != nil {
		return p, err
	}
	return p, nil
}

func (ma *MongoAdapter) GetBySku(sku string) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetBySku")
	defer func() {
		logger.EndProfile(profiler, "GetBySku")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"sku": sku,
	}).One(&p)
	if err == mgo.ErrNotFound {
		return p, NotFoundErr
	}
	if err != nil {
		return p, err
	}
	return p, nil
}

func (ma *MongoAdapter) GetByProductSet(id int) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetByProductSet")
	defer func() {
		logger.EndProfile(profiler, "GetByProductSet")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"productSet": id,
	}).One(&p)
	if err == mgo.ErrNotFound {
		return p, NotFoundErr
	}
	if err != nil {
		return p, err
	}
	return p, nil
}

func (ma *MongoAdapter) GetProductGroupByName(name string) (ProductGroup, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductGroupByName")
	defer func() {
		logger.EndProfile(profiler, "GetProductGroupByName")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	pGroup := ProductGroup{}
	col := mSession.SetCollection(PGROUP_COLLECTION)
	err := col.Find(M{
		"name": name,
	}).One(&pGroup)

	if err == nil {
		return pGroup, nil
	}
	if err == mgo.ErrNotFound {
		//insert new
		pGroup, er := ma.CreateNewPGroup(name)
		if er != nil {
			return pGroup, er
		}
		return pGroup, nil
	}
	return pGroup, err
}

func (ma *MongoAdapter) GetProductsByGroupId(id int) ([]Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductsByGroupId")
	defer func() {
		logger.EndProfile(profiler, "GetProductsByGroupId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	if id <= 0 {
		return nil, nil
	}
	products := []Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(
		M{"group.seqId": id}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductIdsBySellerId(sellerId int) ([]ProductSmall, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductIdsBySellerId")
	defer func() {
		logger.EndProfile(profiler, "GetProductIdsBySellerId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []ProductSmall{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"sellerId": sellerId,
	}).
		Select(M{
			"seqId": 1, "sku": 1, "_id": 0,
		}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductsBySellerId(sellerId int) ([]Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductsBySellerId")
	defer func() {
		logger.EndProfile(profiler, "GetProductsBySellerId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"sellerId": sellerId,
	}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductIdsByBrandId(brandId int) ([]ProductSmall, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductIdsByBrandId")
	defer func() {
		logger.EndProfile(profiler, "GetProductIdsByBrandId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []ProductSmall{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"brandId": brandId,
	}).
		Select(M{
			"seqId": 1, "sku": 1, "_id": 0,
		}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductsByCategoryId(catId int) ([]Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductsByCategoryId")
	defer func() {
		logger.EndProfile(profiler, "GetProductsByCategoryId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"categories": catId,
	}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductsByBrandId(brandId int) ([]Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductsByBrandId")
	defer func() {
		logger.EndProfile(profiler, "GetProductsByBrandId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"brandId": brandId,
	}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductIdsByCategoryId(catId int) ([]ProductSmall, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductIdsByCategoryId")
	defer func() {
		logger.EndProfile(profiler, "GetProductIdsByCategoryId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	products := []ProductSmall{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"categories": catId}).
		Select(M{
			"seqId": 1, "sku": 1, "_id": 0,
		}).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ma *MongoAdapter) GetProductIdBySimpleSku(sku string) (ProductSmall, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductIdBySimpleSku")
	defer func() {
		logger.EndProfile(profiler, "GetProductIdBySimpleSku")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	proSmall := ProductSmall{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"simples.sku": sku}).
		Select(M{
			"seqId": 1, "sku": 1, "_id": 0,
		}).One(&proSmall)
	if err != nil {
		return proSmall, err
	}
	return proSmall, nil
}

func (ma *MongoAdapter) GetProductIdBySimpleId(id int) (ProductSmall, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductIdBySimpleId")
	defer func() {
		logger.EndProfile(profiler, "GetProductIdBySimpleId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	proSmall := ProductSmall{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(M{
		"simples.seqId": id}).
		Select(M{
			"seqId": 1, "sku": 1, "_id": 0,
		}).One(&proSmall)
	if err != nil {
		return proSmall, err
	}
	return proSmall, nil
}

func (ma *MongoAdapter) GetProductBySimpleId(id int) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductBySimpleId")
	defer func() {
		logger.EndProfile(profiler, "GetProductBySimpleId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(
		M{"simples.seqId": id}).One(&p)
	if err != nil {
		return p, errors.New("(*MongoAdapter)#GetProductBySimpleId:" + err.Error())
	}
	return p, nil
}

func (ma *MongoAdapter) GetProductByVideoId(id int) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductByVideoId")
	defer func() {
		logger.EndProfile(profiler, "GetProductByVideoId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	coll := mSession.SetCollection(PRODUCT_COLLECTION)
	err := coll.Find(M{"videos.seqId": id}).One(&p)
	if err != nil {
		return p, errors.New("(*MongoAdapter)#GetProductByVideoId:" + err.Error())
	}
	return p, nil
}

func (ma *MongoAdapter) GetProductBySkuAndType(productSku string, productType string) (Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductBySkuAndType")
	defer func() {
		logger.EndProfile(profiler, "GetProductBySkuAndType")
	}()
	p := Product{}
	var err error
	mSession := ma.GetSession()
	defer mSession.Close()
	coll := mSession.SetCollection(PRODUCT_COLLECTION)
	if productType == PRODUCT_TYPE_CONFIG {
		err = coll.Find(M{"sku": productSku}).One(&p)
	} else {
		err = coll.Find(M{"simples.sku": productSku}).One(&p)
	}
	if err != nil {
		return p, errors.New("(*MongoAdapter)#GetProductBySkuAndType:" + err.Error())
	}
	return p, nil
}

func (ma *MongoAdapter) GetAtrributeByCriteria(attrSrch AttrSearchCondition) (AttributeMongo, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetAtrributeByCriteria")
	defer func() {
		logger.EndProfile(profiler, "GetAtrributeByCriteria")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	//convert our bool value to int
	var isGlobal int
	if attrSrch.IsGlobal {
		isGlobal = 1
	}
	//find all the attributes with supplied name
	attributes := []AttributeMongo{}
	err := mSession.SetCollection(ATTRIBUTE_COLLECTION).Find(
		M{"name": attrSrch.Name},
	).All(&attributes)

	if err != nil {
		logger.Error(err)
		return AttributeMongo{}, errors.New("(*MongoAdapter)#GetAtrributeByCriteria:" + err.Error())
	}

	for _, attr := range attributes {
		if (attr.ProductType == attrSrch.ProductType) &&
			(attr.IsGlobal == isGlobal) &&
			(isGlobal == 1 || (attr.Set.SeqId == attrSrch.SetId)) {
			//go WarmAttributeCache(attr)
			return attr, nil
		}
	}
	return AttributeMongo{}, errors.New("(*MongoAdapter)#GetAtrributeByCriteria:Invalid Attribute Details")
}

func (ma *MongoAdapter) GetProAttributeSetById(id int) (ProAttributeSet, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProAttributeSetById")
	defer func() {
		logger.EndProfile(profiler, "GetProAttributeSetById")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	as := ProAttributeSet{}
	err := mSession.SetCollection(ATTRIBUTESET_COLLECTION).
		Find(M{"seqId": id}).One(&as)
	if err != nil {
		return as, err
	}
	return as, nil
}

func (ma *MongoAdapter) GetAttributeMongoById(seqId int) (AttributeMongo, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetAttributeMongoById")
	defer func() {
		logger.EndProfile(profiler, "GetAttributeMongoById")
	}()
	mSession := ma.GetSession()

	defer mSession.Close()
	tmpAttr := AttributeMongo{}
	err := mSession.SetCollection(ATTRIBUTE_COLLECTION).Find(M{
		"seqId": seqId,
	}).One(&tmpAttr)
	if err != nil {
		logger.Error(err)
		return tmpAttr, err
	}

	//go WarmAttributeCache(tmpAttr)
	return tmpAttr, nil
}

func (ma *MongoAdapter) GetAttributeMongoByName(name string) (AttributeMongo, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetAttributeMongoByName")
	defer func() {
		logger.EndProfile(profiler, "GetAttributeMongoByName")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	tmpAttr := AttributeMongo{}
	err := mSession.SetCollection(ATTRIBUTE_COLLECTION).Find(M{
		"name": name,
	}).One(&tmpAttr)
	if err != nil {
		return tmpAttr, err
	}
	return tmpAttr, nil
}

func (ma *MongoAdapter) GetCategoriesByIds(cats []int) ([]Category, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetCategoriesByIds")
	defer func() {
		logger.EndProfile(profiler, "GetCategoriesByIds")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{
		"seqId": M{
			"$in": cats,
		},
	}
	c := []Category{}
	err := mSession.SetCollection(CATEGORY_COLLECTION).Find(criteria).All(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ma *MongoAdapter) FindPrimaryCategoryId(catIds []int) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "FindPrimaryCategoryId")
	defer func() {
		logger.EndProfile(profiler, "FindPrimaryCategoryId")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{
		"seqId": M{
			"$in": catIds,
		},
	}
	c := []Category{}
	err := mSession.SetCollection(CATEGORY_COLLECTION).Find(criteria).All(&c)
	if err != nil {
		return 0, err
	}
	for _, v := range c {
		if v.Parent == 1 {
			return v.Id, nil
		}
	}
	return 0, err

}

func (ma *MongoAdapter) GetProductsForSeller(
	sellers []int, limit int, offset int, lastScId int,
) ([]Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductsForSeller")
	defer func() {
		logger.EndProfile(profiler, "GetProductsForSeller")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	//get last id
	type LastId struct {
		Name  string `bson:"_id"`
		Value int    `bson:"value"`
	}
	var lastId LastId
	mSession.SetCollection(COUNTER_COLLECTION).Find(
		M{"_id": SSR_COUNTER_NAME},
	).One(&lastId)

	if lastScId > 0 {
		lastId.Value = lastScId
	}
	var products []Product
	criteria := M{
		"sellerId":      M{"$gt": 0},
		"status":        M{"$exists": true},
		"brandId":       M{"$gt": 0},
		"attributeSet":  M{"$exists": true},
		"shipmentType":  M{"$exists": true},
		"simples.seqId": M{"$gt": 0},
	}
	if len(sellers) > 0 {
		criteria["sellerId"] = M{"$in": sellers}
	}
	criteria["seqId"] = M{"$gt": lastId.Value}

	coll := mSession.SetCollection(PRODUCT_COLLECTION)
	err := coll.Find(criteria).Limit(limit).Sort("seqId").
		All(&products)

	//set new id
	lproducts := len(products)
	var newLastId int
	if lproducts > 0 {
		newLastId = products[len(products)-1].SeqId
	} else {
		newLastId = limit + lastId.Value
	}

	mSession.SetCollection(COUNTER_COLLECTION).Upsert(
		M{"_id": SSR_COUNTER_NAME},
		M{"$set": M{"value": newLastId}},
	)
	return products, err
}

func (ma *MongoAdapter) CreateNewPGroup(name string) (ProductGroup, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "CreateNewPGroup")
	defer func() {
		logger.EndProfile(profiler, "CreateNewPGroup")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := ProductGroup{}
	seqId := mSession.GetNextSequence(PGROUP_COLLECTION)
	if seqId <= 0 {
		return p, errors.New("(*MongoAdapter)#CreateNewPGroup: Unable to Generate Sequence")
	}
	//insert in product group

	p.Id = seqId
	p.Name = name
	mSession.SetCollection(PGROUP_COLLECTION).Insert(p)
	return p, nil
}

func (ma *MongoAdapter) SaveProduct(p Product) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "SaveProduct")
	defer func() {
		logger.EndProfile(profiler, "SaveProduct")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{"seqId": p.SeqId}
	t := time.Now()
	p.UpdatedAt = &t
	_, err := mSession.SetCollection(PRODUCT_COLLECTION).Upsert(criteria, p)
	if err != nil {
		return err
	}
	return nil
}

func (ma *MongoAdapter) AddNode(productSku string, nodeName string, data interface{}) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "AddNode")
	defer func() {
		logger.EndProfile(profiler, "AddNode")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{
		"sku": productSku,
	}
	set := M{
		nodeName:    data,
		"updatedAt": time.Now(),
	}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set},
	)
	return err
}

func (ma *MongoAdapter) DeleteNode(sku string, nodeName string) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "DeleteNode")
	defer func() {
		logger.EndProfile(profiler, "DeleteNode")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{
		"sku": sku,
	}
	set := M{
		nodeName:    nil,
		"updatedAt": time.Now(),
	}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set},
	)
	return err
}

func (ma *MongoAdapter) AddImage(productId int, pi ProductImage) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "AddImage")
	defer func() {
		logger.EndProfile(profiler, "AddImage")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{"seqId": productId}
	data := M{
		"$push": M{
			"images": M{
				"seqId":            pi.SeqId,
				"imageNo":          pi.ImageNo,
				"main":             pi.Main,
				"orientation":      pi.Orientation,
				"originalFilename": pi.OriginalFileName,
				"imageName":        pi.ImageName,
				"updatedAt":        pi.UpdatedAt,
			},
		},
	}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(criteria, data)
	if err != nil {
		return pi.SeqId, err
	}

	if pi.IfUpdate {
		set := M{
			"approvedAt":  pi.UpdatedAt,
			"activatedAt": pi.UpdatedAt,
			"updatedAt":   pi.UpdatedAt,
		}
		err = mSession.SetCollection(PRODUCT_COLLECTION).Update(criteria, M{"$set": set})
		return pi.SeqId, err
	}
	return pi.SeqId, nil
}

func (ma *MongoAdapter) DeleteImage(imageId int) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "DeleteImage")
	defer func() {
		logger.EndProfile(profiler, "DeleteImage")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	p := Product{}
	coll := mSession.SetCollection(PRODUCT_COLLECTION)

	err := coll.Find(M{"images.seqId": imageId}).One(&p)
	if err == mgo.ErrNotFound {
		return 0, NotFoundErr
	}
	if err != nil {
		return 0, errors.New("(ma *MongoAdapter) DeleteImage: " + err.Error())
	}
	err = coll.Update(
		M{"images.seqId": imageId},
		M{
			"$pull": M{
				"images": M{
					"seqId": imageId,
				},
			},
		},
	)
	if err != nil {
		return p.SeqId, errors.New("(ma *MongoAdapter) DeleteImage: " + err.Error())
	}
	return p.SeqId, nil
}

func (ma *MongoAdapter) SaveVideo(configId int, video ProductVideo) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "SaveVideo")
	defer func() {
		logger.EndProfile(profiler, "SaveVideo")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	coll := mSession.SetCollection(PRODUCT_COLLECTION)

	//first pull the old one then push the new one
	err := coll.Update(M{"seqId": configId}, M{"$pull": M{
		"videos": M{
			"seqId": video.Id,
		},
	}})
	if err != nil {
		//@todo:: check what happens when video does not exists
		return 0, errors.New("(*MongoAdapter)#SaveVideo: Pull " + err.Error())
	}
	err = coll.Update(M{"seqId": configId}, M{"$push": M{
		"videos": video,
	}})
	if err != nil {
		return 0, errors.New("(*MongoAdapter)#SaveVideo: Push " + err.Error())
	}
	return video.Id, nil
}

func (ma *MongoAdapter) UpdateVideoStatus(videoId int, status string) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateVideoStatus")
	defer func() {
		logger.EndProfile(profiler, "UpdateVideoStatus")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{"videos.seqId": videoId}
	set := M{"videos.$.status": status}
	set["videos.$.updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) UpdateProductAttribute(prdctUpdtCndn PrdctAttrUpdateCndtn) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateProductAttribute")
	defer func() {
		logger.EndProfile(profiler, "UpdateProductAttribute")
	}()
	var err error
	mSession := ma.GetSession()
	defer mSession.Close()
	coll := mSession.SetCollection(PRODUCT_COLLECTION)
	var set M
	if prdctUpdtCndn.ProductType == PRODUCT_TYPE_CONFIG {
		if prdctUpdtCndn.IsGlobal == true {
			set = M{"global": prdctUpdtCndn.PattrMap}
		} else {
			set = M{"attributes": prdctUpdtCndn.PattrMap}
		}
		set["updatedAt"] = time.Now()
		err = coll.Update(M{"sku": prdctUpdtCndn.ProductSku}, M{"$set": set})
	} else {
		if prdctUpdtCndn.IsGlobal == true {
			set = M{"simples.$.global": prdctUpdtCndn.PattrMap}
		} else {
			set = M{"simples.$.attributes": prdctUpdtCndn.PattrMap}
		}
		set["simples.$.updatedAt"] = time.Now()
		set["updatedAt"] = time.Now()
		err = coll.Update(M{"simples.sku": prdctUpdtCndn.ProductSku}, M{"$set": set})
	}
	return err
}

func (ma *MongoAdapter) UpdatePrice(pricedata PriceUpdate) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdatePrice")
	defer func() {
		logger.EndProfile(profiler, "UpdatePrice")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{
		"simples.seqId": pricedata.SimpleId,
	}
	set := M{}
	if pricedata.Price > 0 {
		set = M{"simples.$.price": pricedata.Price}
	}

	if pricedata.UpdateSP {
		set["simples.$.specialPrice"] = pricedata.SpecialPrice
		set["simples.$.specialFromDate"] = pricedata.SpecialFromDate
		set["simples.$.specialToDate"] = pricedata.SpecialToDate
	}
	set["simples.$.updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) UpdateShipmentBySKU(sku string, shipmentType int) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateShipmentBySKU")
	defer func() {
		logger.EndProfile(profiler, "UpdateShipmentBySKU")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{"sku": sku}
	set := M{
		"shipmentType": shipmentType,
		"updatedAt":    time.Now(),
	}
	err := mSession.SetCollection(PRODUCT_COLLECTION).
		Update(criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) GenerateNextSequence(collectionName string) (int, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GenerateNextSequence")
	defer func() {
		logger.EndProfile(profiler, "GenerateNextSequence")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	var counterName string
	switch collectionName {
	case PRODUCT_COLLECTION:
		counterName = PRODUCT_COLLECTION
	case SIMPLE_COLLECTION:
		counterName = SIMPLE_COLLECTION
	case PIMAGE_COLLECTION:
		counterName = PIMAGE_COLLECTION
	case PVIDEO_COLLECTION:
		counterName = PVIDEO_COLLECTION
	case PREPACK_COUNTER:
		counterName = PREPACK_COUNTER
	case DUMMY_IMAGES:
		counterName = DUMMY_IMAGES
	default:
		return 0, errors.New(
			"(ma *MongoAdapter) GenerateNextSequence: Wrong collection Type supplied",
		)
	}
	counter := mSession.GetNextSequence(counterName)
	if counter > 0 {
		return counter, nil
	}
	return 0, errors.New(
		"(ma *MongoAdapter) GenerateNextSequence: Cannot generate sequence",
	)
}

func (ma *MongoAdapter) UpdateProductStatus(id int, status string) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateProductStatus")
	defer func() {
		logger.EndProfile(profiler, "UpdateProductStatus")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{"seqId": id}
	set := M{"status": status}
	set["updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) UpdateProductSimpleStatus(id int, status string) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateProductSimpleStatus")
	defer func() {
		logger.EndProfile(profiler, "UpdateProductSimpleStatus")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	criteria := M{"simples.seqId": id}
	set := M{"simples.$.status": status}
	set["simples.$.updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) UpdateJabongDiscount(jd JabongDiscount) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateJabongDiscount")
	defer func() {
		logger.EndProfile(profiler, "UpdateJabongDiscount")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{
		"simples.seqId": jd.SimpleId,
	}
	set := M{"simples.$.jabongDiscount": jd.Discount}
	set["simples.$.jabongDiscountFromDate"] = jd.FromDate
	set["simples.$.jabongDiscountToDate"] = jd.ToDate
	set["simples.$.updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(
		criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) SetPetApproved(configId int, petApproved int) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "SetPetApproved")
	defer func() {
		logger.EndProfile(profiler, "SetPetApproved")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{"seqId": configId}
	set := M{"petApproved": petApproved}
	set["updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).
		Update(criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) UpdateProductAttributeSystem(input ProductAttrSystemUpdate) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateProductAttributeSystem")
	defer func() {
		logger.EndProfile(profiler, "UpdateProductAttributeSystem")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	criteria := M{
		"seqId": input.ProConfigId,
	}
	set := M{}
	switch input.AttrName {
	case SYSTEM_TY:
		set["ty"] = input.AttrValue.(int)
	case SYSTEM_PET_STATUS:
		set["petStatus"] = input.AttrValue.(string)
	default:
		//handle later
	}
	set["updatedAt"] = time.Now()
	err := mSession.SetCollection(PRODUCT_COLLECTION).
		Update(criteria, M{"$set": set})
	return err
}

func (ma *MongoAdapter) GetAttributeMapping(name string) (AttrMapping, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetAttributeMapping")
	defer func() {
		logger.EndProfile(profiler, "GetAttributeMapping")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	am := AttrMapping{}
	err := mSession.SetCollection(ATTRIBUTE_MAP_COLLECTION).
		Find(M{"from": name}).One(&am)
	return am, err
}

func (ma *MongoAdapter) UpdateProduct(configId int, criteria ProUpdateCriteria) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "UpdateProduct")
	defer func() {
		logger.EndProfile(profiler, "UpdateProduct")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()
	upd := M{
		"seqId": configId,
	}
	set := M{}
	if criteria.ActivatedAt.Isset {
		set["activatedAt"] = criteria.ActivatedAt.Value
	}
	if criteria.PetApproved.Isset {
		set["petApproved"] = criteria.PetApproved.Value
	}
	if criteria.Status.Isset {
		set["status"] = criteria.Status.Value
	}
	err := mSession.SetCollection(PRODUCT_COLLECTION).Update(upd, M{"$set": set})
	if err != nil {
		return fmt.Errorf("(ma *MongoAdapter) UpdateProduct: %s", err.Error())
	}
	return nil
}

func (ma *MongoAdapter) ResetSSRCounter() error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "ResetSSRCounter")
	defer func() {
		logger.EndProfile(profiler, "ResetSSRCounter")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	_, err := mSession.SetCollection(COUNTER_COLLECTION).Upsert(
		M{"_id": SSR_COUNTER_NAME},
		M{"$set": M{"value": 0}},
	)
	if err != nil {
		return fmt.Errorf(
			"(ma *MongoAdapter)#ResetSSRCounter(): %s",
			err.Error(),
		)
	}
	return errors.New("(ma *MongoAdapter)#ResetSSRCounter(): Not Implemented yet")
}

func (ma *MongoAdapter) GetProductBySellerIdSku(sellerId int,
	sellerSkuArr []string) ([]ProductSmallSimples, error) {

	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "GetProductSellerIdSku")
	defer func() {
		logger.EndProfile(profiler, "GetProductSellerIdSku")
	}()
	mSession := ma.GetSession()
	defer mSession.Close()

	skuQuery := M{"$in": sellerSkuArr}
	query := M{"sellerId": sellerId, "simples.sellerSku": skuQuery}
	fields := M{"simples.sku": 1, "simples.sellerSku": 1, "_id": 0}
	result := []ProductSmallSimples{}

	err := mSession.SetCollection(PRODUCT_COLLECTION).Find(query).Select(fields).All(&result)
	if err != nil {
		return nil, fmt.Errorf("(ma *MongoAdapter)#GetProductBySellerIdSku(): %v", err.Error())
	}
	return result, nil
}
