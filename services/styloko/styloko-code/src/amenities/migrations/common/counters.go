package common

import (
	"amenities/migrations/common/util"
	proUtil "amenities/products/common"
	factory "common/ResourceFactory"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func GetSeqCounterMysql(seqType string) (int, error) {
	mysqlDriver, err := factory.GetDefaultMysqlDriver()
	if err != nil {
		return 0, err
	}
	var index *int
	var result int
	switch seqType {
	case util.PrePack:
		query := `SELECT MAX(pack_id) from catalog_config
			where pack_id >` + strconv.Itoa(util.PrePackStartIndex) + ` and
			pack_id <` + strconv.Itoa(util.PrePackEndIndex)
		row, err := mysqlDriver.Query(query)
		if err != nil {
			return 0, fmt.Errorf(err.DeveloperMessage)
		}
		for row.Next() {
			err := row.Scan(&index)
			if err != nil {
				row.Close()
				return 0, err
			}
		}
		row.Close()
		if index == nil {
			result = util.PrePackStartIndex
		} else {
			result = *index
		}
	}
	return result, nil
}

func GetSeqCounter(seqType string) (int, error) {

	mgoSession := factory.GetMongoSession("CUSTOM")
	defer mgoSession.Close()

	var currentIndex int
	switch seqType {
	case proUtil.DUMMY_IMAGES:
		currentIndex = util.DummyImageStart
	case util.Brands:
		mongodb := mgoSession.SetCollection(util.Brands)
		criteria := bson.M{
			"seqId": bson.M{
				"$gt": util.BrandStartIndex,
				"$lt": util.BrandEndIndex,
			},
		}
		proj := bson.M{"seqId": 1, "_id": 0}
		type Data struct {
			SeqId int `bson:"seqId"`
		}
		var d Data
		err := mongodb.Find(criteria).Select(proj).Sort("-seqId").Limit(1).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.SeqId
		if currentIndex == 0 {
			currentIndex = util.BrandStartIndex
		}

	case util.ProductGroups:
		mongodb := mgoSession.SetCollection(util.ProductGroups)
		criteria := bson.M{
			"seqId": bson.M{
				"$gt": util.PGroupStartIndex,
				"$lt": util.PGroupEndIndex,
			},
		}
		proj := bson.M{"seqId": 1, "_id": 0}
		type Data struct {
			SeqId int `bson:"seqId"`
		}
		var d Data
		err := mongodb.Find(criteria).Select(proj).Sort("-seqId").Limit(1).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.SeqId
		if currentIndex == 0 {
			currentIndex = util.PGroupStartIndex
		}

	case util.SizeCharts:
		mongodb := mgoSession.SetCollection(util.SizeCharts)
		criteria := bson.M{
			"seqId": bson.M{
				"$gt": util.SizeChartStartIndex,
				"$lt": util.SizeChartEndIndex,
			},
		}
		proj := bson.M{"seqId": 1, "_id": 0}
		type Data struct {
			SeqId int `bson:"seqId"`
		}
		var d Data
		err := mongodb.Find(criteria).Select(proj).Sort("-seqId").Limit(1).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.SeqId
		if currentIndex == 0 {
			currentIndex = util.SizeChartStartIndex
		}

	case util.Products:
		mongodb := mgoSession.SetCollection(util.Products)
		criteria := bson.M{
			"seqId": bson.M{
				"$gt": util.ProductStartIndex,
				"$lt": util.ProductEndIndex,
			},
		}
		proj := bson.M{"seqId": 1, "_id": 0}
		type Data struct {
			SeqId int `bson:"seqId"`
		}
		var d Data
		err := mongodb.Find(criteria).Select(proj).Sort("-seqId").Limit(1).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.SeqId
		if currentIndex == 0 {
			currentIndex = util.ProductStartIndex
		}
	case util.Simples:
		mongodb := mgoSession.SetCollection(util.Products)
		pipeline := []bson.M{
			bson.M{
				"$match": bson.M{
					"simples.seqId": bson.M{
						"$gt": util.SimpleStartIndex,
						"$lt": util.SimpleEndIndex,
					},
				},
			},
			bson.M{
				"$unwind": "$simples",
			},
			bson.M{
				"$project": bson.M{
					"simples.seqId": 1,
				},
			},
			bson.M{
				"$match": bson.M{
					"simples.seqId": bson.M{
						"$gt": util.SimpleStartIndex,
						"$lt": util.SimpleEndIndex,
					},
				},
			},
			bson.M{
				"$group": bson.M{
					"_id": "",
					"max": bson.M{
						"$max": "$simples.seqId",
					},
				},
			},
		}

		type Data struct {
			Max int `bson:"max"`
		}
		var d Data
		err := mongodb.Pipe(pipeline).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.Max
		if currentIndex == 0 {
			currentIndex = util.SimpleStartIndex
		}

	case util.ProductImages:
		mongodb := mgoSession.SetCollection(util.Products)
		pipeline := []bson.M{
			bson.M{
				"$match": bson.M{
					"images.seqId": bson.M{
						"$gt": util.ImagesStartIndex,
						"$lt": util.ImagesEndIndex,
					},
				},
			},
			bson.M{
				"$unwind": "$images",
			},
			bson.M{
				"$project": bson.M{
					"images.seqId": 1,
				},
			},
			bson.M{
				"$match": bson.M{
					"images.seqId": bson.M{
						"$gt": util.ImagesStartIndex,
						"$lt": util.ImagesEndIndex,
					},
				},
			},
			bson.M{
				"$group": bson.M{
					"_id": "",
					"max": bson.M{
						"$max": "$images.seqId",
					},
				},
			},
		}

		type Data struct {
			Max int `bson:"max"`
		}
		var d Data
		err := mongodb.Pipe(pipeline).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.Max
		if currentIndex == 0 {
			currentIndex = util.ImagesStartIndex
		}

	case util.ProductVideos:
		mongodb := mgoSession.SetCollection(util.Products)
		pipeline := []bson.M{
			bson.M{
				"$match": bson.M{
					"videos.seqId": bson.M{
						"$gt": util.VideosStartIndex,
						"$lt": util.VideosEndIndex,
					},
				},
			},
			bson.M{
				"$unwind": "$videos",
			},
			bson.M{
				"$project": bson.M{
					"videos.seqId": 1,
				},
			},
			bson.M{
				"$match": bson.M{
					"videos.seqId": bson.M{
						"$gt": util.VideosStartIndex,
						"$lt": util.VideosEndIndex,
					},
				},
			},
			bson.M{
				"$group": bson.M{
					"_id": "",
					"max": bson.M{
						"$max": "$videos.seqId",
					},
				},
			},
		}

		type Data struct {
			Max int `bson:"max"`
		}
		var d Data
		err := mongodb.Pipe(pipeline).One(&d)
		if err != nil && (err != mgo.ErrNotFound) {
			return 0, err
		}
		currentIndex = d.Max
		if currentIndex == 0 {
			currentIndex = util.VideosStartIndex
		}
	}
	return currentIndex, nil
}
