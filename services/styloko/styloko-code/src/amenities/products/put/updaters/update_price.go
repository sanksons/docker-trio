package updaters

import (
	proUtil "amenities/products/common"
	syncing "amenities/products/common/synctasks"
	put "amenities/products/put"
	"common/notification"
	"common/notification/datadog"
	"common/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

//Update-Type: Price
type PriceUpdate struct {
	SimpleId        int               `json:"simpleId" validate:"required"`
	Price           float64           `json:"price" validate:"-"`
	SpecialPrice    proUtil.FloatNull `json:"specialPrice"`
	SpecialFromDate proUtil.TimeNull  `json:"specialFromDate"`
	SpecialToDate   proUtil.TimeNull  `json:"specialToDate"`
	Prd             proUtil.Product   `json:"-"`
	Updated         bool              `json:"-"`
}

func (pu *PriceUpdate) priceValidator(
	v *validator.Validate,
	structLevel *validator.StructLevel,
) {
	priceupd := structLevel.CurrentStruct.Interface().(PriceUpdate)
	if !(priceupd.SpecialPrice.Isset) && (priceupd.Price < 1) {
		structLevel.ReportError(
			reflect.ValueOf(priceupd.Price),
			"Price",
			"price",
			"required",
		)
	}

	if !((priceupd.SpecialPrice.Isset == priceupd.SpecialFromDate.Isset) ||
		(priceupd.SpecialFromDate.Isset == priceupd.SpecialToDate.Isset)) {
		structLevel.ReportError(
			reflect.ValueOf(priceupd.SpecialPrice),
			"SpecialPrice",
			"specialPrice",
			"either|All",
		)
	}
	if (priceupd.SpecialPrice.Isset && priceupd.SpecialPrice.Value != nil) &&
		(priceupd.SpecialFromDate.Value == nil && priceupd.SpecialToDate.Value == nil) {
		structLevel.ReportError(
			reflect.ValueOf(priceupd.SpecialFromDate),
			"SpecialFromDate",
			"specialFromDate",
			"required",
		)
		structLevel.ReportError(
			reflect.ValueOf(priceupd.SpecialToDate),
			"SpecialToDate",
			"specialToDate",
			"required",
		)
	}

}

func (pu *PriceUpdate) Validate() []string {
	put.Validate.RegisterStructValidation(pu.priceValidator, PriceUpdate{})
	errs := put.Validate.Struct(pu)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (sd *PriceUpdate) InvalidateCache() error {
	if sd.Updated == false {
		return nil
	}
	// Do not Invalidate cache on price update
	// As all systems are picking price directly from memcache.
	/*
		go func() {
			defer proUtil.RecoverHandler("Price#InvalidateCache")
			put.CacheMngr.DeleteById([]int{sd.Prd.SeqId}, true)
		}()
	*/
	return nil
}

func (pricedata *PriceUpdate) Update() (proUtil.Product, error) {
	prd, err := proUtil.GetAdapter(put.DbAdapterName).GetProductBySimpleId(
		pricedata.SimpleId,
	)
	if err != nil {
		return prd, err
	}
	//update data in object
	pricedata.Prd = prd

	dbUpdate := proUtil.PriceUpdate{}

	for index, simple := range prd.Simples {
		if simple.Id == pricedata.SimpleId {
			if pricedata.Price > 0 {
				prd.Simples[index].Price = &pricedata.Price
				dbUpdate.Price = pricedata.Price
			}
			dbUpdate.SimpleId = simple.Id
			//Do not update price, if product belongs to our retail partner.
			if utils.InArrayInt(proUtil.RetailPartners, prd.SellerId) {
				break
			}

			sp := pricedata.SpecialPrice
			if sp.Isset {
				if sp.Value == nil || (*sp.Value < 1) {
					dbUpdate.UpdateSP = true
					prd.Simples[index].SpecialPrice = nil
					prd.Simples[index].SpecialFromDate = nil
					prd.Simples[index].SpecialToDate = nil
				} else {
					dbUpdate.UpdateSP = true
					dbUpdate.SpecialPrice = sp.Value
					dbUpdate.SpecialFromDate = pricedata.SpecialFromDate.Value
					dbUpdate.SpecialToDate = pricedata.SpecialToDate.Value
					prd.Simples[index].SpecialPrice = sp.Value
					prd.Simples[index].SpecialFromDate = pricedata.SpecialFromDate.Value
					prd.Simples[index].SpecialToDate = pricedata.SpecialToDate.Value
				}
			}
			break
		}
	}

	if dbUpdate.SimpleId <= 0 {
		return prd, fmt.Errorf("Could not get simple info")
	}

	// update mysql
	byteData, er := json.Marshal(dbUpdate)
	if er != nil {
		return prd, er
	}
	err = syncing.ProcessTask(proUtil.UPDATE_TYPE_PRICE, byteData, prd.SeqId, true)
	if err != nil {
		return prd, err
	}
	//mark price data as updated
	pricedata.Updated = true
	//set updated product in object
	pricedata.Prd = prd
	//push to memcache
	prd.Pricegetter.UseMaster = true
	err = prd.PushToMemcache("Update Price Call")
	if err != nil {
		//notify
		notification.SendNotification(
			"Product Memcache Update Failed",
			fmt.Sprintf("Product:%d, Message2:%s", prd.SeqId, err.Error()),
			[]string{proUtil.TAG_PRODUCT_SYNC, proUtil.TAG_PRODUCT},
			datadog.ERROR,
		)
	}

	//update in mongo
	proUtil.GetAdapter(put.DbAdapterName).UpdatePrice(dbUpdate)
	return prd, nil
}

func (pricedata *PriceUpdate) Publish() error {
	if pricedata.Updated == false {
		return nil
	}
	go func() {
		defer proUtil.RecoverHandler("Price#Publish")
		//@todo: enable this once tested
		pricedata.Prd.Pricegetter.UseMaster = true
		pricedata.Prd.Publish("", true)
	}()
	return nil
}

func (pu *PriceUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = pu.SimpleId
	for _, v := range p.Simples {
		if v.Id == pu.SimpleId {
			response["sku"] = v.SKU
			response["sellerSku"] = v.SellerSKU
			response["updatedAt"] = v.UpdatedAt.Format(time.RFC3339)
		}
	}
	return response
}

//
// Acquire Lock
//
func (pu *PriceUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (pu *PriceUpdate) UnLock() bool {
	return true
}
