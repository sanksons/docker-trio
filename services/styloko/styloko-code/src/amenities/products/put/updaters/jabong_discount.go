package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	//"github.com/jabong/floRest/src/common/utils/logger"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

type JabongDiscount struct {
	SimpleId               int             `json:"productId" validate:"required"`
	JabongDiscount         float64         `json:"jabongDiscount"`
	FromDate               *string         `json:"jabongDiscountFromDate"`
	ToDate                 *string         `json:"jabongDiscountToDate"`
	JabongDiscountFromDate *time.Time      `json:"-"`
	JabongDiscountToDate   *time.Time      `json:"-"`
	Prd                    proUtil.Product `json:"-"`
}

func (jd *JabongDiscount) Validate() []string {
	errs := put.Validate.Struct(jd)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	if jd.JabongDiscount < 0 || jd.JabongDiscount > 100 {
		return []string{"Jabong Discount can't be less than 0 or more than 100"}
	}
	if jd.JabongDiscount == 0 {
		jd.JabongDiscountFromDate = nil
		jd.JabongDiscountToDate = nil
		return nil
	}
	if jd.FromDate == nil || *jd.FromDate == "" ||
		jd.ToDate == nil || *jd.ToDate == "" {
		return []string{"Jabong Discount Dates cannot be empty if Discount is more than 0"}
	}
	tfrom, err := proUtil.FromMysqlTime(*jd.FromDate, true)
	if err != nil {
		return []string{"Cannot parse Jabong Discount From Date"}
	}
	jd.JabongDiscountFromDate = tfrom
	tto, err := proUtil.FromMysqlTime(*jd.ToDate, true)
	if err != nil {
		return []string{"Cannot parse Jabong Discount To Date"}
	}
	jd.JabongDiscountToDate = tto
	if jd.JabongDiscountToDate.Before(*jd.JabongDiscountFromDate) {
		return []string{"Jabong Discount To Date cannot be before than From Date"}
	}
	_, err = proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySimpleId(jd.SimpleId)
	if err != nil {
		return []string{"Cannot get product by this simple id"}
	}
	return nil
}

func (jd *JabongDiscount) Update() (proUtil.Product, error) {
	jbngDscnt := proUtil.JabongDiscount{
		SimpleId: jd.SimpleId,
		Discount: jd.JabongDiscount,
		FromDate: jd.JabongDiscountFromDate,
		ToDate:   jd.JabongDiscountToDate,
	}
	err := proUtil.GetAdapter(put.DbAdapterName).
		UpdateJabongDiscount(jbngDscnt)
	if err != nil {
		return proUtil.Product{}, err
	}
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySimpleId(jd.SimpleId)
	if err != nil {
		return proUtil.Product{}, err
	}
	jd.Prd = p

	// Add JOB for syncing to mysql
	taskPool.AddProductSyncJob(
		p.SeqId, proUtil.UPDATE_TYPE_JABONG_DISCOUNT, jbngDscnt,
	)
	return p, err
}

func (jd *JabongDiscount) InvalidateCache() error {
	go func() {
		defer proUtil.RecoverHandler("JabongDiscount#Invalidate Cache")
		put.CacheMngr.DeleteById([]int{jd.Prd.SeqId}, true)
	}()
	return nil
}

func (jd *JabongDiscount) Publish() error {
	go func() {
		defer proUtil.RecoverHandler("JabongDiscount#Publish")
		jd.Prd.Publish("", true)
	}()
	return nil
}

func (jd *JabongDiscount) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = jd.SimpleId
	for _, v := range p.Simples {
		if v.Id == jd.SimpleId {
			response["sku"] = v.SKU
			response["sellerSku"] = v.SupplierSKU
			response["updatedAt"] = v.UpdatedAt
		}
	}
	return response
}

//
// Acquire Lock
//
func (jd *JabongDiscount) Lock() bool {
	return true
}

//
// Release Lock
//
func (jd *JabongDiscount) UnLock() bool {
	return true
}
