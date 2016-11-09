package updaters

import (
	proUtil "amenities/products/common"
	syncing "amenities/products/common/synctasks"
	put "amenities/products/put"
	"common/notification"
	"common/notification/datadog"
	taskPool "common/pool/tasker"
	"common/redis"
	"common/utils"
	"encoding/json"
	"fmt"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

// Update-Type: Product
type ProductUpdate struct {
	ProductId       int                     `json:"productId" validate:"required"`
	SellerSku       string                  `json:"sellerSku" validate:"required"`
	EanCode         string                  `json:"productIdentifier" validate:"-"`
	SKU             string                  `json:"sku" validate:"required"`
	ApprovalStatus  *int                    `json:"approvalStatus" validate:"-"`
	Status          string                  `json:"status" validate:"required,eq=active|eq=deleted|eq=inactive"`
	Title           string                  `json:"title" validate:"required"`
	Brand           int                     `json:"brand" validate:"required"`
	GroupName       string                  `json:"groupName" validate:"-"`
	AttributeSetId  int                     `json:"attributeSet" validate:"required"`
	ShipmentType    int                     `json:"shipmentType" validate:"required"`
	Description     string                  `json:"description" validate:"required"`
	Price           float64                 `json:"price" validate:"-"`
	SpecialPrice    *float64                `json:"specialPrice" validate:"-"`
	SpecialFromDate *string                 `json:"specialFromDate" validate:"-"`
	SpecialToDate   *string                 `json:"specialToDate" validate:"-"`
	TaxClass        int                     `json:"taxClass" validate:"-"`
	Categories      []int                   `json:"categories" validate:"required"`
	PrimaryCategory int                     `json:"primaryCategory" validate:"-"`
	Attributes      proUtil.AttributeMapSet `json:"attributes" validate:"required"`
	ProductSet      int                     `json:"productSet" validate:"required"`
	ConfigId        int                     `json:"configId" validate:"required"`
	VideoHash       []string                `json:"videoHash" validate:"-"`
	UpdateUID       string                  `json:"-" validate:"-"`
	Mutex           proUtil.ProductMutex    `json:"-" validate:"-"`
}

//
// Convert the struct to json string notation
//
func (pu *ProductUpdate) ToJson() string {
	bytes, _ := json.Marshal(pu)
	return string(bytes)
}

//
// Notification to be sent to datadog on failure.
//
func (pu *ProductUpdate) Notify(err error) {
	var title = "Product Update Failed"
	var text = pu.ToJson()
	if err != nil {
		text = err.Error() + " ------->  " + text
	}

	var tags = []string{proUtil.TAG_PRODUCT_UPDATE, proUtil.TAG_PRODUCT}
	notification.SendNotification(title, text, tags, datadog.ERROR)
}

func (pu *ProductUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = pu.ProductId
	for _, v := range p.Simples {
		if v.Id == pu.ProductId {
			response["sku"] = v.SKU
			response["sellerSku"] = v.SellerSKU
			response["updatedAt"] = v.UpdatedAt.Format(time.RFC3339)
			//prepare video hash
			if len(pu.VideoHash) > 0 {
				videohashMap := make(map[string]bool, 0)
				for _, hash := range pu.VideoHash {
					videohashMap[hash] = false
					for _, video := range p.Videos {
						if hash == video.Hash {
							videohashMap[hash] = true
							break
						}
					}
				}
				response["videoProcessInfo"] = videohashMap
			}
		}
	}
	return response
}

func (pu *ProductUpdate) CheckIfValidSize() error {
	tmpadapter := proUtil.GetAdapter(put.DbAdapterName)
	as, err := tmpadapter.GetProAttributeSetById(pu.AttributeSetId)
	if err != nil {
		return fmt.Errorf("Could not connect to mongo: %s", err.Error())
	}
	sizeattrname, err := as.GetVariationAttributeName()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	criteria := proUtil.AttrSearchCondition{
		Name:        sizeattrname,
		ProductType: proUtil.PRODUCT_TYPE_SIMPLE,
		IsGlobal:    false,
		SetId:       as.Id,
	}
	sizeAttr, err := tmpadapter.GetAtrributeByCriteria(criteria)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	sizeId := sizeAttr.SeqId
	for key, _ := range pu.Attributes {
		if key == utils.ToString(sizeId) {
			return nil
		}
	}
	return fmt.Errorf("Size Not Present, looking for [%s]", sizeattrname)
}

func (pu *ProductUpdate) Validate() []string {
	errs := put.Validate.Struct(pu)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		//send notification
		pu.Notify(errs)
		return msgs
	}
	err := pu.CheckIfValidSize()
	if err != nil {
		pu.Notify(err)
		return []string{err.Error()}
	}
	return nil
}

func (pu *ProductUpdate) Update() (proUtil.Product, error) {
	p := &proUtil.Product{}
	//TODO: maybe we can get this data from blitz
	err := p.LoadBySeqId(pu.ConfigId, put.DbAdapterName)
	if err != nil {
		if err == proUtil.NotFoundErr {
			driver, rerr := redis.GetDriver()
			if rerr == nil {
				driver.SADD(proUtil.MIGRATE_PRODUCTS_KEY, utils.ToString(pu.ConfigId))
			}
		}
		//send notification
		pu.Notify(err)
		return *p, err
	}
	petstats := p.GetPetStatus()
	petstats.Edited = true
	petstats.Created = true
	p.SetPetStatus(petstats)

	//@PATCH STARTS
	temp := 1
	pu.ApprovalStatus = &temp
	//@PATCH ENDS

	if pu.ApprovalStatus == nil {
		p.PetApproved = 1
	} else {
		if (p.ActivatedAt == nil) && (*pu.ApprovalStatus == 1) {
			//product is activated
			currentTime := time.Now()
			p.ActivatedAt = &currentTime
		}
		p.PetApproved = *pu.ApprovalStatus
	}
	p.Name = pu.Title
	p.BrandId = pu.Brand
	p.SetAttributeSet(pu.AttributeSetId, proUtil.DB_READ_ADAPTER)
	p.ShipmentType = pu.ShipmentType
	p.Description = pu.Description
	p.TaxClass = &pu.TaxClass
	p.PrimaryCategory = pu.PrimaryCategory
	p.Categories = pu.Categories
	p.PrepareLeafCategories(proUtil.DB_READ_ADAPTER)
	if p.Group != nil && p.Group.Name != pu.GroupName {
		p.SetProductGroup(pu.GroupName)
	}
	scvalue := p.GetScProductValue()
	pu.Attributes.SetSCProduct(scvalue, true)

	// [starts] set PT and dispatch location
	if put.Conf.ProductSellerUpdate {
		seller, err := p.GetSellerInfoWithRetries(put.Conf.Org, false)
		if err != nil {
			return *p, fmt.Errorf("(pu *ProductUpdate)#Update()cannot get seller: %s", err.Error())
		}
		err = pu.Attributes.SetDispatchLocation(seller.SellerCustInfo.DispatchLocation)
		if err != nil {
			return *p, fmt.Errorf("(pu *ProductUpdate)#Update()dispatch location: %s", err.Error())
		}
		err = pu.Attributes.SetProcessingTime(seller.SellerCustInfo.ProcessingTime, seller.Id, p.BrandId)
		if err != nil {
			return *p, fmt.Errorf("(pu *ProductUpdate)#Update() processing time: %s", err.Error())
		}
	}
	// [ends] set PT and dispatch location

	packQtyAttr, ok := p.Global[utils.SnakeToCamel(proUtil.ATTR_PACKQTY)]
	var pckQty int
	if ok {
		packQty, err := packQtyAttr.GetValue("value")
		if err != nil {
			return *p, fmt.Errorf("(pu *ProductUpdate)#Update()packQty: %s", err.Error())
		}
		pckQty, err = utils.GetInt(packQty)
		if err != nil {
			return *p, fmt.Errorf("(pu *ProductUpdate)#Update()packQty: %s", err.Error())
		}
	}
	err = pu.Attributes.SetPrePack(pckQty)
	if err != nil {
		return *p, fmt.Errorf("(pu *ProductUpdate)#Update()prepack: %s", err.Error())
	}
	err = pu.Attributes.SetOldAttributes(p.AttributeSet.Id)
	if err != nil {
		return *p, fmt.Errorf("(pu *ProductUpdate)#Update()oldattrs: %s", err.Error())
	}
	p.SetAttributes(pu.Attributes, true, proUtil.DB_READ_ADAPTER, true)

	//TODO : why not do this call in few minutes or seconds or only once ??
	_, err = p.SetCatalogType(pu.Attributes, proUtil.DB_READ_ADAPTER)
	if err != nil {
		return *p, fmt.Errorf("(pu *ProductUpdate)#Update(): %s", err.Error())
	}

	var simpleFoundIndex int = -1
	for i, simple := range p.Simples {
		if simple.Id != pu.ProductId {
			continue
		}
		simpleFoundIndex = i
		p.Simples[i].SetAttributes(
			p.AttributeSet.Id, pu.Attributes, true, proUtil.DB_READ_ADAPTER, true,
		)
		p.Simples[i].Status = pu.Status
		p.Simples[i].SetStandardSize(p, put.DbAdapterName)
		p.Simples[i].SellerSKU = pu.SellerSku
		p.Simples[i].BarcodeEan = pu.SellerSku
		if pu.EanCode != "" {
			p.Simples[i].EanCode = pu.EanCode
		}

		if pu.Price > 0 {
			p.Simples[i].Price = &pu.Price

			var upd proUtil.PriceUpdate
			upd.SimpleId = p.Simples[i].Id
			upd.Price = pu.Price

			if !utils.InArrayInt(proUtil.RetailPartners, p.SellerId) {

				upd.UpdateSP = true
				if pu.SpecialPrice != nil {
					p.Simples[i].AddSpecialPrice(
						pu.SpecialPrice,
						pu.SpecialFromDate,
						pu.SpecialToDate,
					)
					upd.SpecialPrice = p.Simples[i].SpecialPrice
					upd.SpecialFromDate = p.Simples[i].SpecialFromDate
					upd.SpecialToDate = p.Simples[i].SpecialToDate
				} else {
					p.Simples[i].RemoveSpecialPrice()
				}

			}

			//update special price in mysql
			go func() {
				defer proUtil.RecoverHandler("Price synchronous mysql sync")
				byteData, _ := json.Marshal(upd)
				syncing.ProcessTask(proUtil.UPDATE_TYPE_PRICE, byteData, p.SeqId, false)

			}()
		}
		p.Simples[i].SetUpdatedAt()
	}
	if simpleFoundIndex < 0 {
		//send notification
		pu.Notify(fmt.Errorf("ConfigId and SimpleId do not match"))
		return *p, fmt.Errorf("ConfigId and SimpleId do not match")
	}
	// update sizechart to product
	//*p, _ = proSizeChUtil.UpdateProductWithSizeChart(*p)
	err = p.InsertOrUpdate(put.DbAdapterName)
	if err != nil {
		//send notification
		pu.Notify(err)
		return *p, err
	}
	//Set update UID
	pu.UpdateUID, _ = p.SetUpdateUID()

	// Add JOB for syncing to mysql
	// Instead of syncing the whole product, lets just sync the simple that got changed.
	productCopy := *p
	productCopy.Simples = []*proUtil.ProductSimple{
		productCopy.Simples[simpleFoundIndex],
	}
	taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT, productCopy)
	return *p, err
}

func (pu *ProductUpdate) InvalidateCache() error {
	if !pu.canPublishorInvalidate() {
		return nil
	}
	go func() {
		defer proUtil.RecoverHandler("ProductUpdate#Invalidate Cache")
		put.CacheMngr.DeleteById([]int{pu.ConfigId}, true)
	}()
	return nil
}

func (pu *ProductUpdate) Publish() error {
	if !pu.canPublishorInvalidate() {
		return nil
	}
	pro, err := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GetBySku(
		pu.SKU,
	)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("ProductUpdate#Publish")
		pro.Publish("", true)
	}()
	return nil
}

//
// Acquire Lock
//
func (self *ProductUpdate) Lock() bool {
	self.Mutex = proUtil.ProductMutex{Id: self.ConfigId}
	return self.Mutex.Lock()
}

//
// Release Lock
//
func (self *ProductUpdate) UnLock() bool {
	return self.Mutex.UnLock()
}

//
// Check if we really need to Invalidate or Publish
//
func (self *ProductUpdate) canPublishorInvalidate() bool {
	if self.UpdateUID == "" {
		return true
	}
	pro := &proUtil.Product{SeqId: self.ConfigId}
	uid, err := pro.GetUpdateUID()
	if err != nil {
		return true
	}
	if uid == self.UpdateUID {
		return true
	}
	return false
}
