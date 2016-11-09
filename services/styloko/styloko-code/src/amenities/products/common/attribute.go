package common

import (
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//List of attributes required for servicibility check
var ServicibilityAttributes []string = []string{
	"192", "193", "194", "195", "196", "230",
}

//Attribute as in Attribute collection
type AttributeMongo struct {
	SeqId                  int                    `bson:"seqId"`
	Set                    AttributeSet           `bson:"set"`
	GlobalIdentifier       *string                `bson:"-"`
	IsGlobal               int                    `bson:"isGlobal"`
	Name                   string                 `bson:"name"`
	Label                  string                 `bson:"label"`
	Description            *string                `bson:"description"`
	ProductType            string                 `bson:"productType"`
	AttributeType          string                 `bson:"attributeType"`
	MaxLength              *int                   `bson:"maxLength"`
	DecimalPlaces          *int                   `bson:"decimalPlaces"`
	DefaultValue           *string                `bson:"defaultValue"`
	UniqueValue            *string                `bson:"uniqueValue"`
	PetType                *string                `bson:"petType"`
	PetMode                *string                `bson:"petMode"`
	Validation             *string                `bson:"validation"`
	Mandatory              *int                   `bson:"mandatory"`
	MandatoryImport        *int                   `bson:"mandatoryImport"`
	AliceExport            string                 `bson:"aliceExport"`
	PetQc                  *int                   `bson:"petQc"`
	ImportConfigIdentifier *int                   `bson:"importConfigIdentifier"`
	SolrSearchable         int                    `bson:"solrSearchable"`
	SolrFilter             int                    `bson:"solrFilter"`
	SolrSuggestions        int                    `bson:"solrSuggestions"`
	Visible                int                    `bson:"visible"`
	CreatedAt              time.Time              `bson:"creatdAt"`
	UpdatedAt              time.Time              `bson:"updatedAt"`
	IsActive               int                    `bson:"isActive"`
	Options                []AttributeMongoOption `bson:"options"`
}

type AttrMapping struct {
	AttrFrom string            `bson:"from"`
	AttrTo   string            `bson:"to"`
	Mapping  map[string]string `bson:"mapping"`
}

var attrTy *AttributeMongo

//
// Transforms attribute from AttributeMongo type to Attribute type
//
func (attr AttributeMongo) Transform2Attribute() Attribute {
	attribute := Attribute{}
	attribute.Id = attr.SeqId
	if attr.IsGlobal == 1 {
		attribute.IsGlobal = true
	}
	attribute.Label = attr.Label
	attribute.Name = attr.Name
	attribute.AliceExport = attr.AliceExport
	attribute.OptionType = attr.AttributeType
	attribute.SolrSearchable = attr.SolrSearchable
	attribute.SolrFilter = attr.SolrFilter
	attribute.SolrSuggestions = attr.SolrSuggestions
	return attribute
}

// Gets AttributeMongoOption type from AttributeMongo whose name matches optionName
func (attr AttributeMongo) GetAttrOptionByName(optionName string) (AttributeMongoOption, error) {
	options := attr.Options
	if len(options) == 0 {
		return AttributeMongoOption{}, errors.New(NO_ATTRIBUTE_OPTIONS)
	} else {
		for _, option := range options {
			if strings.ToLower(option.Name) == strings.ToLower(optionName) {
				return option, nil
			}
		}
	}
	return AttributeMongoOption{}, errors.New(NO_ATTRIBUTE_OPTIONS)
}

//AttributeOptionn as in Attribute collection.
type AttributeMongoOption struct {
	SeqId            int     `bson:"seqId"`
	GlobalIdentifier *string `bson:"-"`
	Name             string  `bson:"value"`
	Position         int     `bson:"position"`
	IsDefault        int     `bson:"isDefault"`
}

// Tranforms AttributeMongoOption to AttrOption
func (mOption AttributeMongoOption) Tranform2AttrOption() AttrOption {
	attrOpt := AttrOption{}
	attrOpt.Id = mOption.SeqId
	attrOpt.Value = mOption.Name
	return attrOpt
}

//AttributeOption as we embed in Product
type AttrOption struct {
	Id    int    `bson:"seqId" json:"seqId"`
	Value string `bson:"value" json:"value"`
}

//Attribute as we embed in Product
type Attribute struct {
	Id              int         `bson:"seqId" json:"seqId"`
	IsGlobal        bool        `bson:"isGlobal" json:"isGlobal"`
	Label           string      `bson:"label" json:"label"`
	Name            string      `bson:"name" json:"name"`
	Value           interface{} `bson:"value" json:"value"`
	AliceExport     string      `bson:"aliceExport" json:"aliceExport"`
	OptionType      string      `bson:"optionType" json:"optionType"`
	SolrSearchable  int         `bson:"solrSearchable" json:"solrSearchable"`
	SolrFilter      int         `bson:"solrFilter" json:"solrFilter"`
	SolrSuggestions int         `bson:"solrSuggestions" json:"solrSuggestions"`
}

func (attr Attribute) ToString() string {
	bytesdata, err := json.Marshal(attr)
	if err != nil {
		return ""
	}
	return string(bytesdata)
}

//
// This method returns the value of a particular attribute.
// Param:
//  valtype as id => gives id of option, multiOption type attrs
//  valtype as value => gives value of option, multiOption type attrs
//
// Returns:
//  string => option and value type attrs
//  []string => multiOption attrs
//
func (a Attribute) GetValue(valtype string) (interface{}, error) {
	switch a.OptionType {
	case OPTION_TYPE_MULTI:
		is, ok := a.Value.([]interface{})
		if !ok {
			return nil, errors.New("(a Attribute)#GetValue: OPTION_TYPE_MULTI Assertion failed")
		}
		var result []interface{}
		for _, m := range is {
			switch v := m.(type) {
			case bson.M:
				if valtype == "value" {
					result = append(result, v["value"])
				} else {
					if v["seqId"] == nil {
						continue
					}
					result = append(result, v["seqId"])
				}
			case AttrOption:
				if valtype == "value" {
					result = append(result, v.Value)
				} else {
					if v.Id == 0 {
						continue
					}
					result = append(result, v.Id)
				}
			case map[string]interface{}:
				if valtype == "value" {
					result = append(result, v["value"])
				} else {
					if v["seqId"] == nil {
						continue
					}
					result = append(result, v["seqId"])
				}
			default:
				continue
			}
		}
		return result, nil
	case OPTION_TYPE_SINGLE:
		switch val := a.Value.(type) {
		case bson.M:
			if valtype == "value" {
				return val["value"], nil
			} else {
				return val["seqId"], nil
			}
		case AttrOption:
			if valtype == "value" {
				return val.Value, nil
			} else {
				return val.Id, nil
			}
		case map[string]interface{}:
			if valtype == "value" {
				return val["value"], nil
			} else {
				return val["seqId"], nil
			}
		default:
			return nil, fmt.Errorf("(a Attribute)#GetValue: Could not match type [%v]", reflect.TypeOf(a.Value))
		}
	default:
		return a.Value, nil
	}
	return "", nil
}

func (a Attribute) GetValueForPromotion() (interface{}, error) {
	switch a.OptionType {
	case OPTION_TYPE_MULTI:
		is, ok := a.Value.([]interface{})
		if !ok {
			return nil, errors.New("(a Attribute)#GetValue: OPTION_TYPE_MULTI Assertion failed")
		}
		var result []interface{}
		for _, m := range is {
			result = append(result, m)
		}
		return result, nil
	case OPTION_TYPE_SINGLE:
		switch val := a.Value.(type) {
		case bson.M:
			return val["value"], nil
		case AttrOption:
			return val.Value, nil
		case map[string]interface{}:
			return val["value"], nil
		default:
			return nil, fmt.Errorf("(a Attribute)#GetValue: Could not match type [%v]", reflect.TypeOf(a.Value))
		}
	default:
		return a.Value, nil
	}
	return "", nil
}

// Removes attrOptionArr values from value of Attribute
func (pattrMapAttr Attribute) RemoveMulOptionAttr(attrOptionArr []AttrOption) Attribute {

	var inValArr []AttrOption
	inValIntrfcArr := pattrMapAttr.Value.([]interface{})
	for _, inValIntrfc := range inValIntrfcArr {
		inVal := inValIntrfc.(bson.M)
		remove := false
		for _, attrOpt := range attrOptionArr {
			if attrOpt.Id == inVal["seqId"] {
				remove = true
				break
			}
		}
		if !remove {
			inValArr = append(inValArr, AttrOption{inVal["seqId"].(int), inVal["value"].(string)})
		}
	}
	pattrMapAttr.Value = inValArr
	return pattrMapAttr
}

// Adds attrOptionArr values to value of Attribute
func (pattrMapAttr Attribute) AddMulOptionAttr(attrOptionArr []AttrOption) Attribute {

	var inValArr []AttrOption
	inValIntrfcArr := pattrMapAttr.Value.([]interface{})
	for _, inValIntrfc := range inValIntrfcArr {
		inVal := inValIntrfc.(bson.M)
		inValArr = append(inValArr, AttrOption{inVal["seqId"].(int), inVal["value"].(string)})
	}

	for _, attrOpt := range attrOptionArr {
		already := false
		inValIntrfcArr := pattrMapAttr.Value.([]interface{})
		for _, inValIntrfc := range inValIntrfcArr {
			inVal := inValIntrfc.(bson.M)
			if attrOpt.Id == inVal["seqId"] {
				already = true
				break
			}
		}
		if !already {
			inValArr = append(inValArr, attrOpt)
		}
	}
	pattrMapAttr.Value = inValArr
	return pattrMapAttr
}

func getMappedAttributes(attrId int, attrOptId, adapter string,
	attrSrch AttrSearchCondition) (string, string, error) {

	attrFrom, err := GetAdapter(adapter).GetAttributeMongoById(attrId)
	if err != nil {
		return "", "", err
	}
	attrMap, err := GetAdapter(adapter).GetAttributeMapping(attrFrom.Name)
	if err != nil {
		return "", "", err
	}
	attrSrch.Name = attrMap.AttrTo
	attrTo, err := GetAdapter(adapter).GetAtrributeByCriteria(attrSrch)
	if err != nil {
		return "", "", fmt.Errorf("No filter attribute found")
	}
	attrOptIdInt, err := strconv.Atoi(attrOptId)
	if err != nil || attrOptIdInt == 0 {
		return "", "", fmt.Errorf("Invalid attribute option")
	}
	attrOptName := ""
	for _, val := range attrFrom.Options {
		if val.SeqId == attrOptIdInt {
			attrOptName = val.Name
			break
		}
	}
	if attrOptName == "" {
		return "", "", fmt.Errorf("Invalid attribute option")
	}
	foundMapping := ""
	for key, val := range attrMap.Mapping {
		if key == attrOptName {
			foundMapping = val
		}
	}
	if foundMapping == "" {
		return "", "", fmt.Errorf("No mapping found for this attribute option")
	}
	mapId := 0
	for _, val := range attrTo.Options {
		if val.Name == foundMapping {
			mapId = val.SeqId
		}
	}
	if mapId == 0 {
		return "", "", fmt.Errorf("No mapping found for this attribute option")
	}
	Id := strconv.Itoa(attrTo.SeqId)
	mapIdStr := strconv.Itoa(mapId)
	return Id, mapIdStr, nil
}

func addInMap(attrIp, attrOp map[string]interface{}, key, adapter string,
	attrSrch AttrSearchCondition) map[string]interface{} {

	if val, ok := attrIp[key]; ok {
		if val == nil {
			logger.Warning(fmt.Errorf("addInMap: KEY:%s, no value for this key", key))
			return attrOp
		}
		keyInt, _ := strconv.Atoi(key)
		var valFinal string
		if tmpArr, ok := val.([]interface{}); ok {
			for _, tmpVal := range tmpArr {
				valFinal = tmpVal.(string)
			}
		} else {
			valFinal = val.(string)
		}
		if valFinal == "" {
			logger.Error(fmt.Errorf("addInMap: KEY:%s, no value for this key", key))
			return attrOp
		}
		keyOp, valOp, err := getMappedAttributes(keyInt, valFinal, adapter, attrSrch)
		if err != nil {
			logger.Error(fmt.Errorf("addInMap: KEY:%s, %s", key, err.Error()))
			return attrOp
		}
		attrOp[keyOp] = valOp
	}
	logger.Warning(fmt.Errorf("addInMap: KEY:%s, no such key found", key))
	return attrOp
}

func PrepareDefaultValueAttrsConfig() map[string]interface{} {
	attrList := make(map[string]interface{}, 0)
	//set cancelable
	cancelable, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_IS_CANCELABLE,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_IS_CANCELABLE, err.Error()))
	}
	attrList[utils.ToString(cancelable.SeqId)] = DEF_IS_CANCELABLE

	//set cod
	cod, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_IS_COD,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_IS_COD, err.Error()))
	}
	attrList[utils.ToString(cod.SeqId)] = DEF_IS_COD

	//set fragile
	fragile, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_IS_FRAGILE,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_IS_FRAGILE, err.Error()))
	}
	attrList[utils.ToString(fragile.SeqId)] = DEF_IS_FRAGILE

	//set surface
	surface, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_IS_SURFACE,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_IS_SURFACE, err.Error()))
	}
	attrList[utils.ToString(surface.SeqId)] = DEF_IS_SURFACE

	//set processing
	pt, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_PROCESSING_TIME,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_PROCESSING_TIME, err.Error()))
	}
	attrList[utils.ToString(pt.SeqId)] = DEF_PROCESSING_TIME

	//set shippingAmount
	shippingAmount, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_SHIPPING_AMOUNT,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_SHIPPING_AMOUNT, err.Error()))
	}
	attrList[utils.ToString(shippingAmount.SeqId)] = DEF_SHIPPING_AMOUNT

	//set blockcatalog
	blockcatalog, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_BLOCK_CATALOG,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_BLOCK_CATALOG, err.Error()))
	}
	attrList[utils.ToString(blockcatalog.SeqId)] = DEF_BLOCK_CATALOG

	//set packId
	packId, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_PACK_ID,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_PACK_ID, err.Error()))
	}
	attrList[utils.ToString(packId.SeqId)] = DEF_PACK_ID

	//set packQty
	packQty, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_PACK_QTY,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_PACK_QTY, err.Error()))
	}
	attrList[utils.ToString(packQty.SeqId)] = DEF_PACK_QTY

	//set vatCmpnyContri
	vatCmpnyContri, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_VAT_CMPNY_CONTRI,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_VAT_CMPNY_CONTRI, err.Error()))
	}
	attrList[utils.ToString(vatCmpnyContri.SeqId)] = DEF_VAT_CMPNY_CONTRI

	//set vatCustContri
	vatCustContri, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_VAT_CUST_CONTRI,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_VAT_CUST_CONTRI, err.Error()))
	}
	attrList[utils.ToString(vatCustContri.SeqId)] = DEF_VAT_CUST_CONTRI

	//set imageOrientation
	imageOrientation, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_IMAGE_ORIENTATION,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf("PrepareDefaultValueAttrs(%s): %s", ATTR_IMAGE_ORIENTATION, err.Error()))
	}
	attrList[utils.ToString(imageOrientation.SeqId)] = DEF_IMAGE_ORIENTATION
	return attrList
}

func GetMappedAttributesConfig(attributeSetId int,
	attrs map[string]interface{}, adapter string) map[string]interface{} {

	attrMapping := PrepareDefaultValueAttrsConfig()

	// map color to color family
	color, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        ATTR_COLOR,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err == nil {
		tmpMap := addInMap(attrs, attrMapping, utils.ToString(color.SeqId),
			adapter, AttrSearchCondition{
				"", PRODUCT_TYPE_CONFIG, true, attributeSetId},
		)
		attrMapping = tmpMap
	} else {
		logger.Error(fmt.Errorf("GetMappedAttributesConfig(%s): %s", ATTR_COLOR, err.Error()))
	}

	attrSrch := AttrSearchCondition{"", PRODUCT_TYPE_CONFIG, false, attributeSetId}
	switch attributeSetId {
	case 1:
		//map upper_material_details.
		upperMaterialDetails, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
			Name:        ATTR_UPPR_MTRL_DTLS,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attributeSetId,
		})
		if err != nil {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_UPPR_MTRL_DTLS, err.Error(),
			))
			break
		}
		tmpMap := addInMap(attrs, attrMapping, utils.ToString(upperMaterialDetails.SeqId),
			adapter, attrSrch,
		)
		attrMapping = tmpMap

	case 3:
		//map fabric_details.
		fabricDetails, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_FABRIC_DETAILS,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err != nil {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_FABRIC_DETAILS, err.Error(),
			))
			break
		}
		tmpMap := addInMap(
			attrs, attrMapping, utils.ToString(fabricDetails.SeqId),
			adapter, attrSrch,
		)
		attrMapping = tmpMap

	case 5:
		//map frame_material_detail
		frameMaterialDetail, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_FRAME_MTRL_DTLS,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err != nil {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_FRAME_MTRL_DTLS, err.Error(),
			))
			break
		}
		tmpMap := addInMap(
			attrs, attrMapping, utils.ToString(frameMaterialDetail.SeqId), adapter, attrSrch,
		)
		attrMapping = tmpMap

		//map frame_color
		frameColor, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_FRAME_COLOR,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err == nil {
			tmpMap := addInMap(
				attrs, attrMapping, utils.ToString(frameColor.SeqId), adapter, attrSrch,
			)
			attrMapping = tmpMap
		} else {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_FRAME_COLOR, err.Error(),
			))
		}

	case 6:
		//map strap_color
		strapColor, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_STRAP_COLOR,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err == nil {
			tmpMap := addInMap(
				attrs, attrMapping, utils.ToString(strapColor.SeqId), adapter, attrSrch,
			)
			attrMapping = tmpMap
		} else {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_STRAP_COLOR, err.Error(),
			))
		}

		//map strap_material_detail
		strapMaterialDetail, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_STRAP_MTRL_DTL,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err == nil {
			tmpMap := addInMap(
				attrs, attrMapping, utils.ToString(strapMaterialDetail.SeqId), adapter, attrSrch,
			)
			attrMapping = tmpMap
		} else {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_STRAP_MTRL_DTL, err.Error(),
			))
		}

		//map dial_color
		dialColor, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_DIAL_COLOR,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attributeSetId,
			},
		)
		if err == nil {
			tmpMap := addInMap(
				attrs, attrMapping, utils.ToString(dialColor.SeqId), adapter, attrSrch,
			)
			attrMapping = tmpMap
		} else {
			logger.Error(fmt.Errorf(
				"GetMappedAttributesConfig(%s): %s", ATTR_DIAL_COLOR, err.Error(),
			))
		}
	}
	return attrMapping
}

func GetMappedAttributesSimple(attributeSetId int) map[string]interface{} {
	global := make(map[string]interface{}, 0)
	switch attributeSetId {
	//we dont have any such case till now
	default:
		return global
	}
}

func GetCatalogTyAttribute(adapter string) (AttributeMongo, error) {
	if attrTy == nil {
		tmpAttr, err := GetAdapter(adapter).GetAttributeMongoByName("ty")
		attrTy = &tmpAttr
		return *attrTy, err
	}
	return *attrTy, nil
}
