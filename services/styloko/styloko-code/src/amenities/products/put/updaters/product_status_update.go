package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
)

// Update-Type: ProductStatus
type ProductStatusUpdate struct {
	SimpleId       int                  `json:"productId" validate:"required"`
	Status         string               `json:"status" validate:"required,eq=active|eq=inactive|eq=deleted"`
	ApprovalStatus *int                 `json:"approvalStatus" validate:"-"`
	ConfigId       int                  `json:"configId" validate:"-"`
	Sku            string               `json:"sku" validate:"-"`
	Mutex          proUtil.ProductMutex `json:"-" validate:"-"`
}

func (ps *ProductStatusUpdate) Validate() []string {
	errs := put.Validate.Struct(ps)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}

	if ps.ConfigId > 0 {
		return nil
	}
	//Incase configId is not supplied in request
	// we will fetch from mongo
	smallP, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductIdBySimpleId(ps.SimpleId)
	if err != nil {
		return []string{"(ps *ProductStatusUpdate) Couls not get Config():" + err.Error()}
	}
	ps.ConfigId = smallP.Id
	return nil
}

func (ps *ProductStatusUpdate) Update() (proUtil.Product, error) {
	//Try to fetch product via simpleId, if not found throw error
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySimpleId(ps.SimpleId)
	if err != nil {
		return proUtil.Product{}, fmt.Errorf(
			"(ps *ProductStatusUpdate)#Update1:%s", err.Error(),
		)
	}

	//prepare list of active, inactive, deleted.
	var inactiveSimples []int
	var activeSimples []int
	var deletedSimples []int
	for _, simple := range p.Simples {
		//exclude simple for which update is fired.
		if simple.Id == ps.SimpleId {
			continue
		}
		if (simple.Status == proUtil.STATUS_INACTIVE) ||
			(simple.Status == proUtil.STATUS_ACTIVE && p.Status == proUtil.STATUS_INACTIVE) {
			inactiveSimples = append(inactiveSimples, simple.Id)
		}
		if simple.Status == proUtil.STATUS_ACTIVE && p.Status == proUtil.STATUS_ACTIVE {
			activeSimples = append(activeSimples, simple.Id)
		}
		if simple.Status == proUtil.STATUS_DELETED {
			deletedSimples = append(deletedSimples, simple.Id)
		}
	}
	for k, simple := range p.Simples {
		if simple.Id == ps.SimpleId {
			p.Simples[k].Status = ps.Status
			if ps.Status == proUtil.STATUS_ACTIVE {
				activeSimples = append(activeSimples, simple.Id)
			}
			if ps.Status == proUtil.STATUS_INACTIVE {
				inactiveSimples = append(inactiveSimples, simple.Id)
			}
			if ps.Status == proUtil.STATUS_DELETED {
				deletedSimples = append(deletedSimples, simple.Id)
			}
			break
		}
	}
	var criteria proUtil.ProUpdateCriteria
	var totalSimplesLen int = len(p.Simples)
	//var inactiveSimplesLen int = len(inactiveSimples)
	var activeSimplesLen int = len(activeSimples)
	var deletedSimplesLen int = len(deletedSimples)
	var configStatus string = p.Status

	if totalSimplesLen == deletedSimplesLen {
		//all are deleted
		configStatus = proUtil.STATUS_DELETED
	} else if activeSimplesLen == 0 {
		configStatus = proUtil.STATUS_INACTIVE
		for _, id := range inactiveSimples {
			for k, simple := range p.Simples {
				if id == simple.Id {
					p.Simples[k].Status = proUtil.STATUS_ACTIVE
				}
			}
		}
	} else if activeSimplesLen == 1 {
		configStatus = proUtil.STATUS_ACTIVE
		for _, id := range inactiveSimples {
			for k, simple := range p.Simples {
				if id == simple.Id {
					p.Simples[k].Status = proUtil.STATUS_INACTIVE
				}
			}
		}
	}

	//update status in mongo
	for _, simple := range p.Simples {
		err := proUtil.GetAdapter(put.DbAdapterName).
			UpdateProductSimpleStatus(simple.Id, simple.Status)
		if err != nil {
			logger.Error(fmt.Sprintf("(ps *ProductStatusUpdate)#Update2: %s", err.Error()))
		}
	}

	//set status
	criteria.Status = proUtil.StringNull{
		Value: &configStatus,
		Isset: true,
	}

	//@PATCH STARTS
	temp := 1
	ps.ApprovalStatus = &temp
	//@PATCH ENDS

	if ps.ApprovalStatus != nil {
		//@patch: set pet approved to be 1 always.
		var petApproved = *ps.ApprovalStatus
		criteria.PetApproved = proUtil.IntNull{
			Value: &petApproved,
			Isset: true,
		}
		if (p.ActivatedAt == nil) && (*ps.ApprovalStatus == 1) {
			//product is activated
			currentTime := time.Now()
			criteria.ActivatedAt = proUtil.TimeNull{
				Value: &(currentTime),
				Isset: true,
			}
		}
	}
	err = proUtil.GetAdapter(put.DbAdapterName).UpdateProduct(p.SeqId, criteria)
	if err != nil {
		logger.Error(fmt.Errorf("(ps *ProductStatusUpdate)#Update3: %s", err.Error()))
	}

	//
	// Add JOB for syncing to mysql
	//
	simpleCount := len(p.Simples)
	var i int
	for _, simple := range p.Simples {
		i = i + 1
		if i == simpleCount {
			taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT_STATUS, proUtil.M{
				"simpleId": simple.Id,
				"status":   simple.Status,
				"criteria": proUtil.ProUpdateCriteria{},
			})
		} else {
			taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT_STATUS, proUtil.M{
				"simpleId": simple.Id,
				"status":   simple.Status,
				"criteria": criteria,
			})
		}
	}
	return p, err
}

func (ps *ProductStatusUpdate) InvalidateCache() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySimpleId(ps.SimpleId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("ProductStatus# Invalidate Cache")
		put.CacheMngr.DeleteById([]int{p.SeqId}, true)
	}()
	return nil
}

func (ps *ProductStatusUpdate) Publish() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductBySimpleId(ps.SimpleId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("ProductStatus#Publish")
		p.Publish("", true)
	}()
	return nil
}

func (ps *ProductStatusUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = ps.SimpleId
	for _, v := range p.Simples {
		if v.Id == ps.SimpleId {
			response["sku"] = v.SKU
			response["sellerSku"] = v.SupplierSKU
			response["updatedAt"] = v.UpdatedAt.Format(time.RFC3339)
		}
	}
	return response
}

//
// Acquire Lock
//
func (self *ProductStatusUpdate) Lock() bool {
	self.Mutex = proUtil.ProductMutex{Id: self.ConfigId}
	return self.Mutex.Lock()
}

//
// Release Lock
//
func (self *ProductStatusUpdate) UnLock() bool {
	return self.Mutex.UnLock()
}
