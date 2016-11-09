package common

import (
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"reflect"
	"strconv"
	"strings"
)

var slrMap map[int]string

var brandPtMap map[int]string

type AttributeMapSet map[string]interface{}

func (attrs AttributeMapSet) SetDispatchLocation(dispatchLoc string) error {
	dispatchLocAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_DISPATCH_LOCATION,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetDispatchLocation: %s", err.Error())
	}
	for _, v := range dispatchLocAttr.Options {

		if strings.ToLower(v.Name) == strings.ToLower(dispatchLoc) {
			attrs[strconv.Itoa(dispatchLocAttr.SeqId)] = strconv.Itoa(v.SeqId)
			break
		}
	}
	return nil
}

func (attrs AttributeMapSet) SetPrePack(oldPackQty int) error {
	// Check if pack_qty attribute is there
	packQtyAtr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_PACKQTY,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetPrePack: %s", err.Error())
	}
	packIdAtr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_PACK_ID,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetPrePack: %s", err.Error())
	}
	packQtyAttrVal, ok := attrs[strconv.Itoa(packQtyAtr.SeqId)]
	if !ok {
		return nil
	}
	if packQtyAttrVal == oldPackQty {
		return nil
	} else {
		nextPackId, err := GetAdapter(DB_ADAPTER_MONGO).GenerateNextSequence(PREPACK_COUNTER)
		if err != nil {
			return err
		}
		attrs[strconv.Itoa(packIdAtr.SeqId)] = nextPackId
	}
	return nil
}

func (attrs AttributeMapSet) SetProcessingTime(pt string, slrSeqId int, brandId int) error {
	ptAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_PROCESSING_TIME,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetProcessingTime: %s", err.Error())
	}
	//if slrSeqId && brandId exist in Map
	//assign value from map else passed value
	if _, ok := slrMap[slrSeqId]; ok {
		if val, exists := brandPtMap[brandId]; exists {
			attrs[strconv.Itoa(ptAttr.SeqId)] = val
			return nil
		}
	}
	attrs[strconv.Itoa(ptAttr.SeqId)] = pt
	return nil
}

//
// Set the values for old attributes, for which we have new option
// type attributes.
// This method syncs values from option type attrs to their corresponding value
// type attrs.
//
func (attrs AttributeMapSet) SetOldAttributes(attrsetId int) error {
	err := attrs.SetOldVariationAttributes(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldFit(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldJeansWashEffect(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldLensType(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldMaterialCode(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldPocket(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldProductWarranty()
	if err != nil {
		return err
	}
	err = attrs.SetOldQualities(attrsetId)
	if err != nil {
		return err
	}
	err = attrs.SetOldSecondaryColor()
	if err != nil {
		return err
	}
	return nil
}

//set fits -> fit
func (attrs AttributeMapSet) SetOldFit(attrsetId int) error {
	//shoes
	if attrsetId != 1 {
		return nil
	}
	//
	// Pick value from fits and put in fit
	//
	fitsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_FITS,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldFit1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(fitsAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldFit2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range fitsAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		fitAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_FIT,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldFit3: %s", err.Error())
		}
		attrs[strconv.Itoa(fitAttr.SeqId)] = optionValue
	}
	return nil
}

func (attrs AttributeMapSet) SetOldQualities(attrsetId int) error {
	//home
	if attrsetId != 18 {
		return nil
	}

	qualitiesAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_QUALITIES,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldQualities1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(qualitiesAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldFit2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range qualitiesAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		qualityAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_QUALITY,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldQualities3: %s", err.Error())
		}
		attrs[strconv.Itoa(qualityAttr.SeqId)] = optionValue
	}
	return nil
}

//materials_code -> material_code
func (attrs AttributeMapSet) SetOldMaterialCode(attrsetId int) error {
	//home
	if attrsetId != 18 {
		return nil
	}
	//
	// Pick value from materials_code and put in material_code
	//
	materialsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_MTRLS_CODE,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldMaterialCode1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(materialsAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldMaterialCode2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range materialsAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		materialAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_MTRL_CODE,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldMaterialCode3: %s", err.Error())
		}
		attrs[strconv.Itoa(materialAttr.SeqId)] = optionValue
	}
	return nil
}

//products_warranty -> product_warranty
func (attrs AttributeMapSet) SetOldProductWarranty() error {

	//
	// Pick value from products_warranty and put in product_warranty
	//
	productswAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_PRDCTS_WRNTY,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldProductWarranty1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(productswAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldProductWarranty2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range productswAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		productwAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_PRDCT_WRNTY,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    true,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldProductWarranty3: %s", err.Error())
		}
		attrs[strconv.Itoa(productwAttr.SeqId)] = optionValue
	}
	return nil
}

// secondary_colors -> secondary_color
func (attrs AttributeMapSet) SetOldSecondaryColor() error {

	//
	// Pick value from fits and put in fit
	//
	colorsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_SCNDRY_CLRS,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldSecondaryColor1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(colorsAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldSecondaryColor2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range colorsAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		colorAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_SCNDRY_CLR,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    true,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldSecondaryColor3: %s", err.Error())
		}
		attrs[strconv.Itoa(colorAttr.SeqId)] = optionValue
	}
	return nil
}

// lens_types -> lens_type
func (attrs AttributeMapSet) SetOldLensType(attrsetId int) error {
	//bags
	if attrsetId != 5 {
		return nil
	}
	//
	// Pick value from lens_types and put in lens_type
	//
	lensTypesAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_LENS_TYPS,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldLensType1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(lensTypesAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldLensType2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range lensTypesAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		lensTypeAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_LENS_TYP,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldLensType3: %s", err.Error())
		}
		attrs[strconv.Itoa(lensTypeAttr.SeqId)] = optionValue
	}
	return nil
}

// jeans_wash_effects -> jeans_wash_effect
func (attrs AttributeMapSet) SetOldJeansWashEffect(attrsetId int) error {
	//app_men, app_women
	if attrsetId != 3 && attrsetId != 4 {
		return nil
	}
	//
	// Pick value from jeans_wash_effects and put in jeans_wash_effect
	//
	effectsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_JEANS_WES,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldJeansWashEffect1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(effectsAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldJeansWashEffect2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range effectsAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		effectAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_JEANS_WE,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldJeansWashEffect3: %s", err.Error())
		}
		attrs[strconv.Itoa(effectAttr.SeqId)] = optionValue
	}
	return nil
}

// pocket -> pockets
func (attrs AttributeMapSet) SetOldPocket(attrsetId int) error {
	//app_men, app_women
	if attrsetId != 3 && attrsetId != 4 {
		return nil
	}
	//
	// Pick value from pocket and put in pockets
	//
	pocketAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_POCKET,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return fmt.Errorf("SetOldPocket1: %s", err.Error())
	}
	if val, ok := attrs[strconv.Itoa(pocketAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldPocket2: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range pocketAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		pocketsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_POCKETS,
				ProductType: PRODUCT_TYPE_CONFIG,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return fmt.Errorf("SetOldPocket3: %s", err.Error())
		}
		attrs[strconv.Itoa(pocketsAttr.SeqId)] = optionValue
	}
	return nil
}

func (attrs AttributeMapSet) SetOldVariationAttributes(attrsetId int) error {

	//bags, beauty, fragrance, home, toys
	var allowedAttributeSets = []int{5, 21, 20, 18, 22}
	var goahead bool
	for _, v := range allowedAttributeSets {
		if v == attrsetId {
			goahead = true
		}
	}
	if !goahead {
		return nil
	}

	//
	// Pick value from variations and put in variation
	//
	variationsAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_VARIATIONS,
			ProductType: PRODUCT_TYPE_SIMPLE,
			IsGlobal:    false,
			SetId:       attrsetId,
		},
	)
	if err == NotFoundErr {
		return nil
	}
	if err != nil {
		return err
	}
	if val, ok := attrs[strconv.Itoa(variationsAttr.SeqId)]; ok {
		if val == nil {
			return nil
		}
		optionId, isAsserted := val.(string)
		if !isAsserted {
			return fmt.Errorf(
				"(attrs AttributeMapSet)#SetOldVariationAttributes1: expected string, found %v",
				reflect.TypeOf(val),
			)
		}
		var optionValue string
		for _, v := range variationsAttr.Options {
			if strconv.Itoa(v.SeqId) == optionId {
				optionValue = v.Name
				break
			}
		}
		if optionValue == "" {
			logger.Error(fmt.Sprintf("Cannot find any option for option id:%v", optionId))
			return nil
		}
		variationAttr, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
			AttrSearchCondition{
				Name:        ATTR_VARIATION,
				ProductType: PRODUCT_TYPE_SIMPLE,
				IsGlobal:    false,
				SetId:       attrsetId,
			},
		)
		if err != nil {
			return err
		}
		attrs[strconv.Itoa(variationAttr.SeqId)] = optionValue
	}
	return nil
}

//
// Set sc_product attribute value
//
func (attrs AttributeMapSet) SetSCProduct(value string, override bool) error {

	scProduct, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(
		AttrSearchCondition{
			Name:        ATTR_SC_PRODUCT,
			ProductType: PRODUCT_TYPE_CONFIG,
			IsGlobal:    true,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"(attrs AttributeMapSet) SetSCProduct: Unable to get sc_product: %s", err.Error(),
		)
	}
	_, ok := attrs[strconv.Itoa(scProduct.SeqId)]
	if override || !ok {
		attrs[strconv.Itoa(scProduct.SeqId)] = value
	}
	return nil
}

//
// Process AttributeMapSet annd returns list of Attributes
//
func (attrs AttributeMapSet) ProcessAtributes(productType string,
	skipSystem bool, adapter string,
) ([]*Attribute, error) {

	attributeColl := []*Attribute{}
	for key, attr := range attrs {
		// NOTE: if this flag is set means, we need to ignore attribute
		//       and not append it in attribute list.
		var ignoreAttribute bool

		if attr == nil {
			continue
		}
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			logger.Error(err)
			continue
		}
		tmpAttr, err := GetAdapter(adapter).GetAttributeMongoById(keyInt)
		if err != nil {
			logger.Error(err)
			continue
		}
		if tmpAttr.ProductType != productType {
			continue
		}
		if skipSystem && (tmpAttr.AttributeType == OPTION_TYPE_SYSTEM) {
			continue
		}
		a := &Attribute{}
		a.Id = tmpAttr.SeqId

		a.IsGlobal = false
		if tmpAttr.IsGlobal == 1 {
			a.IsGlobal = true
		}
		a.Label = tmpAttr.Label
		a.Name = tmpAttr.Name
		a.OptionType = tmpAttr.AttributeType
		a.AliceExport = tmpAttr.AliceExport
		a.SolrFilter = tmpAttr.SolrFilter
		a.SolrSearchable = tmpAttr.SolrSearchable
		a.SolrSuggestions = tmpAttr.SolrSuggestions

		switch tmpAttr.AttributeType {
		case OPTION_TYPE_SINGLE:
			val, ok := attr.(string)
			if !ok {
				logger.Error(tmpAttr.Name + "Option Assertion failed.")
				continue
			}
			//search attribute options data for correct option
			valInt, err := strconv.Atoi(val)
			if err != nil || (valInt < 0) {
				logger.Error(err)
				continue
			}
			var attrOption AttrOption
			for _, opt := range tmpAttr.Options {
				if opt.SeqId == valInt {
					attrOption = AttrOption{
						Id:    valInt,
						Value: opt.Name,
					}
				}
			}
			a.Value = attrOption
			//check if we dont have a valid value for this attribute.
			if attrOption.Id <= 0 {
				ignoreAttribute = true
			}

		case OPTION_TYPE_MULTI:
			var valSlice []interface{}
			var ok bool
			valSlice, ok = attr.([]interface{})
			if !ok {
				//if its not an array, it may be string
				valSlice = []interface{}{attr}
			}
			var aoSlice []AttrOption
			for _, valI := range valSlice {
				val, ok := valI.(string)
				if !ok {
					logger.Error(tmpAttr.Name + ":Option Assertion failed.")
					continue
				}
				valInt, err := strconv.Atoi(val)
				if err != nil || (valInt <= 0) {
					logger.Error(err)
					continue
				}
				for _, opt := range tmpAttr.Options {
					if opt.SeqId == valInt {
						aoSlice = append(aoSlice, AttrOption{
							Id:    valInt,
							Value: opt.Name,
						})
					}
				}
			}
			if len(aoSlice) > 0 {
				a.Value = aoSlice
			} else {
				ignoreAttribute = true
			}

		case OPTION_TYPE_VALUE:

			switch v := attr.(type) {
			case string:
				var validation string
				if tmpAttr.Validation != nil {
					validation = *tmpAttr.Validation
				}
				switch validation {
				case VALIDATION_DECIMAL:
					a.Value, _ = strconv.ParseFloat(v, 64)
				case VALIDATION_INTEGER:
					a.Value, _ = strconv.Atoi(v)
				default:
					a.Value = v
				}
			case bool:
				//we store tinyint inplace of bool
				if v {
					a.Value = 1
				} else {
					a.Value = 0
				}
			case int:
				a.Value = v
			case int64:
				a.Value = int(v)
			case int32:
				a.Value = int(v)
			case float32:
				a.Value = float64(v)
			case float64:
				a.Value = v
			default:
				ignoreAttribute = true
			}
		default:
			continue
			//leave this
		}
		if !ignoreAttribute {
			attributeColl = append(attributeColl, a)
		}
	}
	return attributeColl, nil
}

func IntializeGlobalVariables(slrTmpMap, brandPtTmpMap map[string]string) {
	slrMap = make(map[int]string)
	brandPtMap = make(map[int]string)
	for k, v := range slrTmpMap {
		key, _ := strconv.Atoi(k)
		slrMap[key] = v
	}
	for k, v := range brandPtTmpMap {
		key, _ := strconv.Atoi(k)
		brandPtMap[key] = v
	}
}
