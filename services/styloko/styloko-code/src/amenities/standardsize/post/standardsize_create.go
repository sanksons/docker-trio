package post

import (
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"common/appconstant"
	"common/mongodb"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

// StandardSizeCreate -> struct for node based data
type StandardSizeCreate struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (sl *StandardSizeCreate) SetID(id string) {
	sl.id = id
}

// GetID -> returns current node ID to orchestrator
func (sl StandardSizeCreate) GetID() (id string, err error) {
	return sl.id, nil
}

// Name -> Returns node name to orchestrator
func (sl StandardSizeCreate) Name() string {
	return "StandardSizeCreate"
}

// Execute -> Starts node execution
func (sl StandardSizeCreate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.POST_CREATE)
	defer func() {
		logger.EndProfile(profiler, common.POST_CREATE)
	}()
	tmp, err := io.IOData.Get(common.STANDARDSIZE_VALID_DATA)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.InconsistantDataStateErrorCode,
			Message: "Failure in getting input", DeveloperMessage: "Cannot retrieve post data"}
	}

	// Assertions
	tmp1, _ := tmp.(map[string][]StandardSizeCreateInput)
	validArr := tmp1["valid"]
	invalidArr := tmp1["invalid"]

	for _, x := range validArr {
		errMessage, _, ok := x.createStandardSize()
		if !ok {
			x.Error = errMessage
			invalidArr = append(invalidArr, x)
		}
	}
	result := make(map[string]interface{})
	result["errors"] = invalidArr
	io.IOData.Set(florest_constants.RESULT, result)
	return io, nil
}

//transform standard size create input to mongo schema object
func (ssi *StandardSizeCreateInput) transformStandardSizeCreate(attributeSetId,
	leafCategoryId, brandId int, mgoSession *mongodb.MongoDriver) common.StandardSizeStore {

	nextId := mgoSession.GetNextSequence(common.STANDARDSIZE_COLLECTION)
	var final common.StandardSizeStore
	final.SeqId = nextId
	final.AttributeSetId = attributeSetId
	final.LeafCategoryId = leafCategoryId
	final.BrandId = brandId
	final.Size = make([]common.BrandStandardSize, 1)
	if brandId == 0 {
		final.Size[0] = common.BrandStandardSize{BrandSize: "0",
			StandardSize: ssi.StandardSize}
	} else {
		final.Size[0] = common.BrandStandardSize{BrandSize: ssi.BrandSize,
			StandardSize: ssi.StandardSize}
	}
	return final
}

//check if given inputs are present in standard size error collection; if yes, then set fixed to true
func (ssi *StandardSizeCreateInput) clearStndSizeError() {
	mgoSessionError := mongoFactory.GetMongoSession(common.STANDARDSIZEERROR_CREATE)
	defer mgoSessionError.Close()
	mgoObjError := mgoSessionError.SetCollection(common.STANDARDSIZEERROR_COLLECTION)
	var ipe common.StandardSizeError
	mgoObjError.Find(bson.M{"$and": []bson.M{{"attrbtStId": ssi.AttributeSetID},
		{"lfCtgryId": ssi.LeafCategoryID}, {"brndId": ssi.BrandID},
		{"brndSz": ssi.BrandSize}, {"fixed": false}}}).One(&ipe)
	if ipe.SeqId != 0 {
		ipe.Fixed = true
		_, err := mgoSessionError.FindAndModify(common.STANDARDSIZEERROR_COLLECTION,
			bson.M{"$set": ipe}, bson.M{"seqId": ipe.SeqId}, false)
		if err != nil {
			logger.Debug(err)
		}
	}
}

//create or update standard size only or standard brand size
func (ssi *StandardSizeCreateInput) createOrUpdateBrandStndSize(useBrand bool) (string,
	int, bool) {

	mgoSession := mongoFactory.GetMongoSession(common.STANDARDSIZE_CREATE)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(common.STANDARDSIZE_COLLECTION)
	var (
		ip  common.StandardSizeStore
		bid int
		bsz string
	)
	if useBrand {
		bid = ssi.BrandID
		bsz = ssi.BrandSize
	} else {
		bid = 0
		bsz = "0"
	}

	mgoObj.Find(bson.M{"$and": []bson.M{{"attrbtStId": ssi.AttributeSetID},
		{"lfCtgryId": ssi.LeafCategoryID}, {"brndId": bid}}}).One(&ip)

	if ip.SeqId == 0 {
		//create a new standard size mapping to attribute set and leaf category
		final := ssi.transformStandardSizeCreate(ssi.AttributeSetID,
			ssi.LeafCategoryID, bid, mgoSession)

		err := mgoSession.Insert(common.STANDARDSIZE_COLLECTION, final)

		if err != nil {
			logger.Debug(err)
			return err.DeveloperMessage, 0, false
		}
		return "", final.SeqId, true
	} else {
		//append a new standard size mapping to existing attribute set and leaf category
		ip.Size = append(ip.Size, common.BrandStandardSize{BrandSize: bsz,
			StandardSize: ssi.StandardSize})

		_, err := mgoSession.FindAndModify(common.STANDARDSIZE_COLLECTION,
			bson.M{"$set": ip}, bson.M{"seqId": ip.SeqId}, false)

		if err != nil {
			logger.Debug(err)
			return err.Error(), 0, false
		}
		return "", ip.SeqId, true
	}
}

func (ssi *StandardSizeCreateInput) createStandardSize() (string, int, bool) {
	//MONGO PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.MONGO_CREATE)
	defer func() {
		logger.EndProfile(profiler, common.MONGO_CREATE)
	}()

	ssi.clearStndSizeError()

	if ssi.Brand == "" && (ssi.BrandSize == "" || ssi.BrandSize == "0") {
		err, id, ok := ssi.createOrUpdateBrandStndSize(false)
		return err, id, ok
	} else {
		err, id, ok := ssi.createOrUpdateBrandStndSize(true)
		return err, id, ok
	}
}
