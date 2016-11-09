package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	"common/utils"
	"errors"
	"fmt"

	validator "gopkg.in/go-playground/validator.v8"
)

// Update-Type: Attribute
type ProductAttributeUpdate struct {
	AttributeName string             `json:"attributeName" validate:"required"`
	IsGlobal      bool               `json:"isGlobal"`
	ProductSku    string             `json:"productSku" validate:"required"`
	ProductType   string             `json:"productType" validate:"required,eq=config|eq=simple"`
	Action        int                `json:"action" validate:"required,eq=1|eq=2|eq=3"`
	Value         interface{}        `json:"value" validate:"required"`
	PetApproved   *int               `json:"petApproved" validate:"-"`
	AttrData      *proUtil.Attribute `json:"attrData" validate:"-"`
}

func (pa *ProductAttributeUpdate) Validate() []string {
	errs := put.Validate.Struct(pa)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	if pa.PetApproved != nil && *pa.PetApproved != 0 && *pa.PetApproved != 1 {
		return []string{"Pet Approved can be only 0 or 1"}
	}
	return nil
}

// Handles various cases of attribute update of attribute type 'Value'
// Returns updated map[string]*proUtil.Attribute along with error
func (pa *ProductAttributeUpdate) updateAttrTypeValue(
	pattrMap map[string]*proUtil.Attribute,
	attr proUtil.AttributeMongo,
) (map[string]*proUtil.Attribute, error) {

	attrName := utils.SnakeToCamel(attr.Name)
	var validation string
	if attr.Validation != nil {
		validation = *attr.Validation
	}
	var (
		attrOptionVal interface{}
		err           error
		dtype         string
	)
	switch validation {
	case proUtil.VALIDATION_DECIMAL:
		dtype = "decimal"
		attrOptionVal, err = utils.GetFloat(pa.Value)
	case proUtil.VALIDATION_INTEGER:
		dtype = "integer"
		attrOptionVal, err = utils.GetInt(pa.Value)
	default:
		dtype = "string"
		attrOptionVal = utils.ToString(pa.Value)
	}
	if err != nil {
		return pattrMap, errors.New("Wrong datatype supplied for Attribute, Expecting ->" + dtype)
	}
	//handle ACTION_REPLACE OR ACTION_ADD, Same way
	if pa.Action == proUtil.ACTION_REPLACE || pa.Action == proUtil.ACTION_ADD {
		if _, ok := pattrMap[attrName]; ok {
			pattrMap[attrName].Value = attrOptionVal

		} else {
			attrAdd := attr.Transform2Attribute()
			attrAdd.Value = attrOptionVal
			pattrMap[attrName] = &attrAdd
		}
		pa.AttrData = pattrMap[attrName]
		return pattrMap, nil
	}
	//handle ACTION_REMOVE
	if _, ok := pattrMap[attrName]; !ok {
		//already removed
		return pattrMap, nil
	}
	// check if non mandatory attribute
	if attr.Mandatory == nil || *attr.Mandatory == 0 {
		delete(pattrMap, attrName)
		return pattrMap, nil
	}
	//check if mandatory but default value not provided
	if attr.DefaultValue == nil || *attr.DefaultValue == "" {
		return pattrMap, errors.New(proUtil.MANDATORY_ATTRIBUTE)
	}
	pattrMap[attrName].Value = *attr.DefaultValue
	pa.AttrData = pattrMap[attrName]
	return pattrMap, nil
}

// Handles various cases of attribute update of attribute type 'Option'
// Returns updated map[string]*proUtil.Attribute along with error
func (pa *ProductAttributeUpdate) updateAttrTypeOption(
	pattrMap map[string]*proUtil.Attribute,
	attr proUtil.AttributeMongo,
) (map[string]*proUtil.Attribute, error) {

	attrName := utils.SnakeToCamel(attr.Name)
	//ACTION_REPLACE OR proUtil.ACTION_ADD handle the same way
	if pa.Action == proUtil.ACTION_REPLACE || pa.Action == proUtil.ACTION_ADD {
		attrOptionVal := pa.Value.(string)
		attrOption, err := attr.GetAttrOptionByName(attrOptionVal)
		if err != nil {
			return pattrMap, err
		}
		attrOptionAdd := attrOption.Tranform2AttrOption()
		if _, ok := pattrMap[attrName]; ok {
			pattrMap[attrName].Value = attrOptionAdd
		} else {
			attrAdd := attr.Transform2Attribute()
			attrAdd.Value = attrOptionAdd
			pattrMap[attrName] = &attrAdd
		}
		pa.AttrData = pattrMap[attrName]
		return pattrMap, nil
	}
	//we are here means its ACTION_REMOVE
	if _, ok := pattrMap[attrName]; !ok {
		//its already removed
		return pattrMap, nil
	}
	//if its not mandatory. cool return
	if attr.Mandatory == nil || *attr.Mandatory == 0 {
		delete(pattrMap, attrName)
		return pattrMap, nil
	}
	//mandatory but default value not set its an error
	if attr.DefaultValue == nil || *attr.DefaultValue == "" {
		return pattrMap, errors.New(proUtil.MANDATORY_ATTRIBUTE)
	}
	attrOption, err := attr.GetAttrOptionByName(*attr.DefaultValue)
	if err != nil {
		return pattrMap, err
	}
	attrOptionAdd := attrOption.Tranform2AttrOption()
	pattrMap[attrName].Value = attrOptionAdd
	pa.AttrData = pattrMap[attrName]
	return pattrMap, nil
}

// Makes an array of proUtil.AttrOption from 'Value' field of the input
// by getting details of proUtil.AttributeMongoOption by name and converting them to a map
func (pa *ProductAttributeUpdate) makeAttrOptionArr(attr proUtil.AttributeMongo) (
	[]proUtil.AttrOption, error,
) {
	retArr := make([]proUtil.AttrOption, 0)
	var paValIntrfc []interface{}
	var ok bool
	paValIntrfc, ok = pa.Value.([]interface{})
	if !ok {
		paValIntrfc = []interface{}{pa.Value}
	}
	for _, paVal := range paValIntrfc {
		paValue, ok := paVal.(string)
		if !ok {
			return retArr, errors.New("(pa *ProductAttributeUpdate)#makeAttrOptionArr: Assertion failed")
		}
		attrOption, err := attr.GetAttrOptionByName(paValue)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%v: %v", err.Error(), paValue))
		}
		attrOptionAdd := attrOption.Tranform2AttrOption()
		retArr = append(retArr, attrOptionAdd)
	}
	return retArr, nil
}

// Handles various cases of attribute update of attribute type 'multi_option'
// Returns updated map[string]*proUtil.Attribute along with error
func (pa *ProductAttributeUpdate) updateAttrTypeMultiOption(
	pattrMap map[string]*proUtil.Attribute,
	attr proUtil.AttributeMongo,
) (map[string]*proUtil.Attribute, error) {

	attrName := utils.SnakeToCamel(attr.Name)
	attrOptionArr, err := pa.makeAttrOptionArr(attr)
	if err != nil {
		return pattrMap, err
	}
	//handle ACTION_REPLACE
	if pa.Action == proUtil.ACTION_REPLACE {
		if _, ok := pattrMap[attrName]; ok {
			pattrMap[attrName].Value = attrOptionArr
		} else {
			attrAdd := attr.Transform2Attribute()
			attrAdd.Value = attrOptionArr
			pattrMap[attrName] = &attrAdd
		}
		pa.AttrData = pattrMap[attrName]
		return pattrMap, nil
	}
	//handle ACTION_ADD
	if pa.Action == proUtil.ACTION_ADD {
		if _, ok := pattrMap[attrName]; ok {
			pattrMapAttr := (*pattrMap[attrName]).AddMulOptionAttr(attrOptionArr)
			pattrMap[attrName] = &pattrMapAttr
		} else {
			attrAdd := attr.Transform2Attribute()
			attrAdd.Value = attrOptionArr
			pattrMap[attrName] = &attrAdd
		}
		pa.AttrData = pattrMap[attrName]
		return pattrMap, nil
	}
	//handle ACTION_REMOVE
	if pa.Action == proUtil.ACTION_REMOVE {
		//if attr does not have any associated options.
		if _, ok := pattrMap[attrName]; !ok {
			return pattrMap, nil
		}
		attrRemove := (*pattrMap[attrName]).RemoveMulOptionAttr(attrOptionArr)
		pattrMap[attrName] = &attrRemove
		valArr := attrRemove.Value.([]proUtil.AttrOption)
		if len(valArr) > 0 {
			pa.AttrData = pattrMap[attrName]
			return pattrMap, nil
		}
		//atrribute options are empty.
		if attr.Mandatory == nil || *attr.Mandatory == 0 {
			//remove the atttribute from pattrMap and return
			delete(pattrMap, attrName)
			return pattrMap, nil
		}
		//if default value is not provided for a mandatory attribute send error.
		if attr.DefaultValue == nil || *attr.DefaultValue == "" {
			return pattrMap, errors.New(proUtil.MANDATORY_ATTRIBUTE)
		}
		attrOption, err := attr.GetAttrOptionByName(*attr.DefaultValue)
		if err != nil {
			return pattrMap, err
		}
		attrOptionAdd := attrOption.Tranform2AttrOption()
		pattrMap[attrName].Value = []proUtil.AttrOption{attrOptionAdd}
		pa.AttrData = pattrMap[attrName]
		return pattrMap, nil
	}
	return pattrMap, errors.New(
		"(pa *ProductAttributeUpdate) updateAttrTypeMultiOption: Undefined condition to handle",
	)
}

// Handles various cases of attribute update of all attribute types
func (pa *ProductAttributeUpdate) updateAttribute(pattrMap map[string]*proUtil.Attribute,
	attr proUtil.AttributeMongo) (map[string]*proUtil.Attribute, error) {

	switch attr.AttributeType {
	case "value":
		pattrMap, err := pa.updateAttrTypeValue(pattrMap, attr)
		if err != nil {
			return pattrMap, err
		}
		break
	case "option":
		pattrMap, err := pa.updateAttrTypeOption(pattrMap, attr)
		if err != nil {
			return pattrMap, err
		}
		break
	case "multi_option":
		pattrMap, err := pa.updateAttrTypeMultiOption(pattrMap, attr)
		if err != nil {
			return pattrMap, err
		}
		break
	default:
		return pattrMap, errors.New("Only value, option, multi option attribute type handled.")
	}
	return pattrMap, nil
}

// Handles all attributes except system type
func (pa *ProductAttributeUpdate) updateAttributeAll(prdct proUtil.Product,
	attr proUtil.AttributeMongo) (proUtil.Product, error) {

	var (
		pattr map[string]*proUtil.Attribute
		err   error
	)
	if pa.ProductType == proUtil.PRODUCT_TYPE_CONFIG {
		if pa.IsGlobal {
			pattr, err = pa.updateAttribute(prdct.Global, attr)
		} else {
			pattr, err = pa.updateAttribute(prdct.Attributes, attr)
		}
		if err != nil {
			return proUtil.Product{}, err
		}
		if pa.IsGlobal {
			prdct.Global = pattr
		} else {
			prdct.Attributes = pattr
		}
	} else {
		// find simple product by sku from the array of simples
		var sp *proUtil.ProductSimple
		for index := 0; index < len(prdct.Simples); index++ {
			if prdct.Simples[index].SKU == pa.ProductSku {
				sp = prdct.Simples[index]
				break
			}
		}
		if pa.IsGlobal == true {
			pattr, err = pa.updateAttribute(sp.Global, attr)
		} else {
			pattr, err = pa.updateAttribute(sp.Attributes, attr)
		}
		if err != nil {
			return proUtil.Product{}, err
		}
		if pa.IsGlobal == true {
			sp.Global = pattr
		} else {
			sp.Attributes = pattr
		}
	}
	// update in mongo
	prdctUpdtCndtn := proUtil.PrdctAttrUpdateCndtn{pa.ProductSku,
		pa.ProductType, pa.IsGlobal, pattr}
	err = proUtil.GetAdapter(put.DbAdapterName).
		UpdateProductAttribute(prdctUpdtCndtn)
	if err != nil {
		return proUtil.Product{}, err
	}
	return prdct, nil
}

func (pa *ProductAttributeUpdate) updateAttributeSystem(
	prdct proUtil.Product,
	attr proUtil.AttributeMongo,
) (proUtil.Product, error) {

	var prdctAttrSysUpdt proUtil.ProductAttrSystemUpdate
	var err error

	switch pa.AttributeName {
	case "ty":
		name, ok := pa.Value.(string)
		if !ok {
			return prdct, errors.New("Cannot Assert [Need string value]")
		}
		idTy, err := proUtil.GetTyByName(name)
		if err != nil {
			return prdct, fmt.Errorf("cannot get tyByName: %s", err.Error())
		}
		tyArr, err := proUtil.GetTyByCategory(prdct.Categories)
		if err != nil {
			return prdct, errors.New("Cannot get ty for product categories")
		}
		var validTy int
		for _, ty := range tyArr {
			if ty != idTy {
				continue
			}
			validTy = idTy
			break
		}
		if validTy <= 0 {
			return prdct, errors.New("Cannnot match TY with any of the categories of product")
		}
		prdct.TY = validTy
		prdctAttrSysUpdt = proUtil.ProductAttrSystemUpdate{
			prdct.SeqId, "ty", validTy,
		}

	default:
		return prdct, fmt.Errorf("Update of %v not allowed with this API", pa.AttributeName)
	}
	//check if we need to update db
	if prdctAttrSysUpdt.ProConfigId > 0 {
		err = proUtil.GetAdapter(put.DbAdapterName).
			UpdateProductAttributeSystem(prdctAttrSysUpdt)
	}
	return prdct, err
}

// Attribute Update main function, calls updateAttributeAll or updateAttributeSystem further
func (pa *ProductAttributeUpdate) Update() (proUtil.Product, error) {
	prdct, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySkuAndType(pa.ProductSku, pa.ProductType)
	if err != nil {
		return proUtil.Product{}, err
	}
	attrSrchCndtn := proUtil.AttrSearchCondition{
		pa.AttributeName,
		pa.ProductType,
		pa.IsGlobal,
		prdct.AttributeSet.Id}
	attr, err := proUtil.GetAdapter(put.DbAdapterName).
		GetAtrributeByCriteria(attrSrchCndtn)
	if err != nil {
		return proUtil.Product{}, err
	}
	//check if this attribute needs to be ignored
	attrIgnore := proUtil.AttributesInfo[proUtil.ATTRIBUTES_IGNORE].([]interface{})
	for _, attrIg := range attrIgnore {
		if attr.Name == attrIg.(string) {
			return proUtil.Product{}, errors.New("Update of this attribute is not allowed")
		}
	}
	// check if its a system type attribute
	if attr.AttributeType == proUtil.OPTION_TYPE_SYSTEM {
		prdct, err = pa.updateAttributeSystem(prdct, attr)
		if err == nil {
			// Add sync Job
			taskPool.AddProductSyncJob(prdct.SeqId, proUtil.SYNC_ATTRIBUTE_SYSTEM, *pa)
		}
	} else {
		//value, option, multiOption attrs
		prdct, err = pa.updateAttributeAll(prdct, attr)
		if err != nil {
			// Add sync Job
			return proUtil.Product{}, err
		}
		//mysql sync
		taskPool.AddProductSyncJob(prdct.SeqId, proUtil.SYNC_ATTRIBUTE_GENERAL, *pa)
		if attr.Name != proUtil.ATTR_PACK_QTY {
			return prdct, nil
		}

		attrVal, err := utils.GetInt(pa.Value)
		if err != nil || attrVal <= 1 {
			return prdct, nil
		}
		nextPrepckId, err := proUtil.GetAdapter(put.DbAdapterName).GenerateNextSequence("prepack")
		if err != nil {
			return proUtil.Product{}, err
		}

		pa1 := ProductAttributeUpdate{
			AttributeName: proUtil.ATTR_PACKID,
			IsGlobal:      true,
			ProductSku:    pa.ProductSku,
			ProductType:   proUtil.PRODUCT_TYPE_CONFIG,
			Action:        proUtil.ACTION_REPLACE,
			Value:         nextPrepckId,
		}
		attrSrchCndtn := proUtil.AttrSearchCondition{
			pa1.AttributeName,
			pa1.ProductType,
			pa1.IsGlobal,
			prdct.AttributeSet.Id,
		}
		attr, err := proUtil.GetAdapter(put.DbAdapterName).
			GetAtrributeByCriteria(attrSrchCndtn)
		prdct, err = pa1.updateAttributeAll(prdct, attr)
		if err == nil {
			// Add sync Job
			taskPool.AddProductSyncJob(prdct.SeqId, proUtil.SYNC_ATTRIBUTE_GENERAL, pa1)
		}
	}
	if err == nil && pa.PetApproved != nil {
		//update petApproved
		prdct.PetApproved = *pa.PetApproved
		err = proUtil.GetAdapter(put.DbAdapterName).
			SetPetApproved(prdct.SeqId, *pa.PetApproved)
	}
	return prdct, err
}

func (pa *ProductAttributeUpdate) InvalidateCache() error {
	prdct, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySkuAndType(pa.ProductSku, pa.ProductType)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Attribute#Invalidate Cache")
		put.CacheMngr.DeleteBySku([]string{prdct.SKU}, true)
	}()
	return nil
}

func (pa *ProductAttributeUpdate) Publish() error {
	prdct, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySkuAndType(pa.ProductSku, pa.ProductType)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Attribute#Publish")
		prdct.Publish("", true)
	}()
	return nil
}

func (pa *ProductAttributeUpdate) Response(p *proUtil.Product) interface{} {
	return p.SKU
}

//
// Acquire Lock
//
func (pa *ProductAttributeUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (pa *ProductAttributeUpdate) UnLock() bool {
	return true
}
