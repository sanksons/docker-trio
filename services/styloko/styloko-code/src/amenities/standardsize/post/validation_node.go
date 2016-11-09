package post

import (
	"amenities/services/attributes"
	"amenities/services/brands"
	"amenities/services/categories"
	"amenities/standardsize/common"
	mongoFactory "common/ResourceFactory"
	"common/appconstant"
	"common/mongodb"
	"common/utils"
	"encoding/json"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// StandardSizeCreateValidation -> struct for node based data
type StandardSizeCreateValidation struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (ss *StandardSizeCreateValidation) SetID(id string) {
	ss.id = id
}

// GetID -> returns current node ID to orchestrator
func (ss StandardSizeCreateValidation) GetID() (id string, err error) {
	return ss.id, nil
}

// Name -> Returns node name to orchestrator
func (ss StandardSizeCreateValidation) Name() string {
	return "StandardSizeCreateValidation"
}

// GetDecision -> Decides which node to run next. Here its a validation node.
func (ss StandardSizeCreateValidation) GetDecision(io workflow.WorkFlowData) (bool, error) {
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.POST_VALIDATE)
	defer func() {
		logger.EndProfile(profiler, common.POST_VALIDATE)
	}()

	data, err := utils.GetPostData(io)
	if err != nil {
		io.IOData.Set(common.STANDARDSIZE_ERROR,
			florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
				Message: common.NO_DATA, DeveloperMessage: "No post data found"})
		return false, nil
	}

	// Path based validation
	pathParams, _ := utils.GetPathParams(io)
	if len(pathParams) != 0 {
		io.IOData.Set(common.STANDARDSIZE_ERROR,
			florest_constants.AppError{Code: appconstant.FunctionalityNotImplementedErrorCode,
				Message: common.INVALID_PATH, DeveloperMessage: "Extra path parameters"})
		return false, nil
	}

	// JSON Unmarshal
	var standardSizeInput []StandardSizeCreateInput
	err = json.Unmarshal(data, &standardSizeInput)
	if err != nil {
		return false, &florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.INVALID_DATA,
			DeveloperMessage: err.Error()}
	}
	//clean standardSizeInput to include only unique values
	uniqueStandardSizeInput := ss.getUnique(standardSizeInput)
	var validArr []StandardSizeCreateInput
	var invalidArr []StandardSizeCreateInput
	for _, x := range uniqueStandardSizeInput {
		apperror, ok := x.validate()
		if ok {
			validArr = append(validArr, x)
			continue
		}
		x.Error = apperror.DeveloperMessage
		invalidArr = append(invalidArr, x)
	}

	dataMap := make(map[string][]StandardSizeCreateInput)
	dataMap["valid"] = validArr
	dataMap["invalid"] = invalidArr
	io.IOData.Set(common.STANDARDSIZE_VALID_DATA, dataMap)
	return true, nil
}

//get respective ids by name and return error of required
func (ssi *StandardSizeCreateInput) getIds(mgoSession *mongodb.MongoDriver) (florest_constants.AppError,
	bool) {

	ssi.AttributeSetID = attributes.GetByName(ssi.AttributeSet, mgoSession).SeqId
	if ssi.AttributeSetID == 0 {
		return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.ERROR_GENERIC,
			DeveloperMessage: common.ERROR_ATTRIBUTE_SET}, false
	}

	lcId, err := strconv.Atoi(ssi.LeafCategory)
	if err != nil {
		return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.ERROR_GENERIC,
			DeveloperMessage: common.ERROR_LEAF_CATEGORY}, false
	}
	catArr := categories.ByIds([]int{lcId})
	if catArr == nil {
		return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.ERROR_GENERIC,
			DeveloperMessage: common.ERROR_LEAF_CATEGORY}, false
	}
	ssi.LeafCategoryID = lcId

	ssi.BrandID = 0
	if ssi.Brand != "" {
		ssi.BrandID = brands.GetByName(ssi.Brand).SeqId
		if ssi.BrandID == 0 {
			return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
				Message:          common.ERROR_GENERIC,
				DeveloperMessage: common.ERROR_BRAND}, false
		}
	}

	return florest_constants.AppError{}, true
}

//execute mongo find contingent on brandId
func (ssi *StandardSizeCreateInput) findMapping(useBrand bool,
	mgoObj *mgo.Collection) common.StandardSizeStore {

	var ip common.StandardSizeStore
	if useBrand {
		mgoObj.Find(bson.M{"$and": []bson.M{{"attrbtStId": ssi.AttributeSetID},
			{"lfCtgryId": ssi.LeafCategoryID}, {"brndId": ssi.BrandID}}}).One(&ip)
	} else {
		mgoObj.Find(bson.M{"$and": []bson.M{{"attrbtStId": ssi.AttributeSetID},
			{"lfCtgryId": ssi.LeafCategoryID}, {"brndId": 0}}}).One(&ip)
	}
	return ip
}

func (ssi *StandardSizeCreateInput) validateInput() (florest_constants.AppError, bool) {

	//check for mandatory inputs
	if ssi.AttributeSet == "" || ssi.LeafCategory == "" || ssi.StandardSize == "" {
		return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.CANNOT_BE_EMPTY,
			DeveloperMessage: common.ERROR_INCOMPLETE_DATA}, false
	}
	return florest_constants.AppError{}, true
}

//validate standard size mapping
func (ssi *StandardSizeCreateInput) validateStndSizeMap(mgoObj *mgo.Collection) (florest_constants.AppError,
	bool) {

	ip := ssi.findMapping(false, mgoObj)
	for _, sv := range ip.Size {
		//check if standard size is already present for given attribute set and leaf category
		if sv.BrandSize == "0" && sv.StandardSize == ssi.StandardSize {
			return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
				Message:          common.ERROR_GENERIC,
				DeveloperMessage: common.VALIDATION_FAILURE_1}, false
		}
	}
	return florest_constants.AppError{}, true
}

//validate brand size mapping
func (ssi *StandardSizeCreateInput) validateBrandSizeMap(mgoObj *mgo.Collection) (common.StandardSizeStore,
	florest_constants.AppError, bool) {

	//check if there is an attribute set and leaf category combination present for standard size
	//if not, then can't map any brand size to this standard size
	ip := ssi.findMapping(false, mgoObj)
	if ip.SeqId == 0 {
		return common.StandardSizeStore{},
			florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
				Message:          common.ERROR_GENERIC,
				DeveloperMessage: common.VALIDATION_FAILURE_2}, false
	}
	return ip, florest_constants.AppError{}, true
}

//validate standard size
func (ssi *StandardSizeCreateInput) validateStndSize(ip common.StandardSizeStore) (florest_constants.AppError,
	bool) {

	ssFound := false
	for _, sv := range ip.Size {
		if sv.StandardSize == ssi.StandardSize {
			ssFound = true
		}
	}
	//standard size needs to be one of the mapping of attribute set and leaf category
	if ssFound == false {
		return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
			Message:          common.ERROR_GENERIC,
			DeveloperMessage: common.VALIDATION_FAILURE_3}, false
	}
	return florest_constants.AppError{}, true
}

//validate standard size and brand size mapping
func (ssi *StandardSizeCreateInput) validateBrandStndSizeMap(mgoObj *mgo.Collection) (florest_constants.AppError,
	bool) {

	ipb := ssi.findMapping(true, mgoObj)
	if ipb.SeqId != 0 {
		for _, sv := range ipb.Size {
			if sv.BrandSize == ssi.BrandSize {
				//check if standard size is already present for this brand and brand size
				//there can be only one standard size mapping to a particular brand and brand size
				if sv.StandardSize == ssi.StandardSize {
					return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
						Message:          common.ERROR_GENERIC,
						DeveloperMessage: common.VALIDATION_FAILURE_4}, false
				} else {
					//check if this brand size is already mapped to another standard size
					//there can be only one brand size to standard size mapping for a particular brand
					return florest_constants.AppError{Code: appconstant.InvalidDataErrorCode,
						Message:          common.ERROR_GENERIC,
						DeveloperMessage: common.VALIDATION_FAILURE_5}, false
				}
			}
		}
	}
	return florest_constants.AppError{}, true
}

//Validate standard size input
func (ssi *StandardSizeCreateInput) validate() (florest_constants.AppError, bool) {

	mgoSession := mongoFactory.GetMongoSession(common.STANDARDSIZE_CREATE)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(common.STANDARDSIZE_COLLECTION)
	err, ok := ssi.validateInput()
	if !ok {
		return err, ok
	}
	err, ok = ssi.getIds(mgoSession)
	if !ok {
		return err, ok
	}
	if ssi.Brand == "" && (ssi.BrandSize == "" || ssi.BrandSize == "0") {
		err, ok := ssi.validateStndSizeMap(mgoObj)
		if !ok {
			return err, ok
		}
	} else {
		ip, err, ok := ssi.validateBrandSizeMap(mgoObj)
		if !ok {
			return err, ok
		} else {
			err, ok := ssi.validateStndSize(ip)
			if !ok {
				return err, ok
			}
			err, ok = ssi.validateBrandStndSizeMap(mgoObj)
			if !ok {
				return err, ok
			}
		}
	}
	return florest_constants.AppError{}, true
}

func (ss *StandardSizeCreateValidation) getUnique(data []StandardSizeCreateInput) map[string]StandardSizeCreateInput {
	uniqueMap := make(map[string]StandardSizeCreateInput, 0)
	for _, v := range data {
		key := fmt.Sprintf("%s%s%s%s%s", v.Brand, v.AttributeSet, v.LeafCategory, v.BrandSize, v.StandardSize)
		uniqueMap[key] = v
	}
	return uniqueMap
}
