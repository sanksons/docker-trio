package standardsize

import (
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

func GetStandardSize(attrStId, brId int,
	lfCtgIds []int, brSz string) (string, bool) {

	mgSes := mongoFactory.GetMongoSession(common.STANDARDSIZE_SEARCH)
	mg := mgSes.SetCollection(common.STANDARDSIZE_COLLECTION)
	mgErr := mgSes.SetCollection(common.STANDARDSIZEERROR_COLLECTION)
	defer mgSes.Close()

	for _, lfCtgId := range lfCtgIds {

		var ip common.StandardSizeStore

		//find this combination with brand Id equal to 0 first
		mg.Find(bson.M{"$and": []bson.M{{"attrbtStId": attrStId},
			{"lfCtgryId": lfCtgId}, {"brndId": 0}}}).One(&ip)
		if ip.SeqId != 0 {
			for _, sizeVal := range ip.Size {
				if sizeVal.StandardSize == brSz {
					return sizeVal.StandardSize, true
				}
			}
		}

		//find this combination with brand Id now
		mg.Find(bson.M{"$and": []bson.M{{"attrbtStId": attrStId},
			{"lfCtgryId": lfCtgId}, {"brndId": brId}}}).One(&ip)
		if ip.SeqId != 0 {
			for _, sizeVal := range ip.Size {
				if sizeVal.BrandSize == brSz {
					return sizeVal.StandardSize, true
				}
			}
		}

		//insert this combination in standard size error collection now
		var ipErr common.StandardSizeError
		mgErr.Find(bson.M{"$and": []bson.M{{"attrbtStId": attrStId},
			{"lfCtgryId": lfCtgId}, {"brndId": brId},
			{"brndSz": brSz}, {"fixed": false}}}).One(&ipErr)

		if ipErr.SeqId != 0 {
			continue
		}
		var ssErr common.StandardSizeError
		nextId := mgSes.GetNextSequence(common.STANDARDSIZEERROR_COLLECTION)
		ssErr.SeqId = nextId
		ssErr.AttributeSetId = attrStId
		ssErr.LeafCategoryId = lfCtgId
		ssErr.BrandId = brId
		ssErr.BrandSize = brSz
		ssErr.CreatedAt = time.Now()
		ssErr.Fixed = false

		err := mgSes.Insert(common.STANDARDSIZEERROR_COLLECTION, ssErr)
		if err != nil {
			logger.Debug(err)
		}

	}
	return "", false
}
