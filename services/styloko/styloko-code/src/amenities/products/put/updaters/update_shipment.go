package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	"errors"
	validator "gopkg.in/go-playground/validator.v8"
	"strings"
)

//Update-Type: Shipment
type ProductShipmentUpdate struct {
	SkuSimple    string `json:"skuSimple" validate:"required"`
	SimpleId     int    `json:"simpleId" validate:"required"`
	ShipmentType int    `json:"shipmentType" validate:"required"`
}

func (su *ProductShipmentUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = su.SimpleId
	for _, v := range p.Simples {
		if v.Id == su.SimpleId {
			response["sku"] = v.SKU
			response["sellerSku"] = v.SupplierSKU
			response["updatedAt"] = v.UpdatedAt
		}
	}
	return response
}

func (su *ProductShipmentUpdate) getSku() string {
	syrs := strings.Split(su.SkuSimple, "-")
	if len(syrs) == 0 {
		return ""
	}
	return syrs[0]
}

func (su *ProductShipmentUpdate) Publish() error {
	pro, err := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GetProductBySimpleId(su.SimpleId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Shipment#Publish")
		pro.Publish("", true)
	}()
	return nil
}

func (su *ProductShipmentUpdate) Update() (proUtil.Product, error) {
	p := &proUtil.Product{}
	sku := su.getSku()
	err := proUtil.GetAdapter(put.DbAdapterName).UpdateShipmentBySKU(sku, su.ShipmentType)
	if err != nil {
		return *p, errors.New("(su ReqDataShipmentUpdate)#Update(): " + err.Error())
	}
	p.LoadBySku(sku, put.DbAdapterName)
	taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_SHIPMENT, proUtil.M{
		"sku":      sku,
		"shipment": su.ShipmentType,
	})
	return *p, nil
}

func (su *ProductShipmentUpdate) Validate() []string {
	errs := put.Validate.Struct(su)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (su *ProductShipmentUpdate) InvalidateCache() error {
	sku := su.getSku()
	go func() {
		defer proUtil.RecoverHandler("Shipment#Invalidate Cache")
		put.CacheMngr.DeleteBySku([]string{sku}, true)
	}()
	return nil
}

//
// Acquire Lock
//
func (su *ProductShipmentUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (su *ProductShipmentUpdate) UnLock() bool {
	return true
}
