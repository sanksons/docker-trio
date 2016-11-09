package post

import (
	proUtil "amenities/products/common"
	proSizeChUtil "amenities/products/common/prodsizechart"
	"common/notification"
	"common/notification/datadog"
	taskPool "common/pool/tasker"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type InsertNode struct {
	id string
}

func (cs *InsertNode) SetID(id string) {
	cs.id = id
}

func (cs InsertNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs InsertNode) Name() string {
	return "InsertNode"
}

func (cs InsertNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Debug("Enter Insert node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.POST_INSERT_NODE)
	defer logger.EndProfile(profiler, proUtil.POST_INSERT_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Insert Node")

	d, _ := io.IOData.Get(proUtil.IODATA)
	iodata, ok := d.([]*ProIOData)
	if !ok {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: "(cs InsertNode)#Execute(): Invalid ProductCreateRequestData Data",
		}
	}
	for i, proIoData := range iodata {
		if proIoData.Status == proUtil.STATUS_FAILURE {
			continue
		}

		//Define a mutex, which we will use to obtain lock.
		mutex := proUtil.ProductMutex{}
		if proIoData.ReqData.ConfigId > 0 {
			mutex.Id = proIoData.ReqData.ConfigId
			mutex.Type = proUtil.LOCK_TYPE_CONFIGID
		} else {
			mutex.Id = proIoData.ReqData.ProductSet
			mutex.Type = proUtil.LOCK_TYPE_PRODUCTSET
		}
		//Try to obtain lock.
		if !mutex.Lock() {
			iodata[i].setFailure(constants.AppError{
				Code:             constants.DbErrorCode,
				Message:          "Unable to Obtain Lock to create product",
				DeveloperMessage: "Unable to Obtain Lock to create product",
			})
			continue
		}

		p, err := cs.GetProduct(proIoData.ReqData.ProductSet, proIoData.ReqData.ConfigId)
		if err != nil {
			//release lock
			mutex.UnLock()
			logger.Error(err)
			iodata[i].setFailure(constants.AppError{
				Code:             constants.DbErrorCode,
				Message:          "Unable to Verify Product Info",
				DeveloperMessage: err.Error(),
			})
			continue
		}
		//Check if we need to create a new product
		//or update existing
		if p == nil {
			//insert
			p, err = cs.CreateNewProduct(proIoData.ReqData)
			if err != nil {
				logger.Error(err)
				iodata[i].setFailure(constants.AppError{
					Code:             constants.DbErrorCode,
					Message:          "Product Insertion failed",
					DeveloperMessage: err.Error(),
				})
				//notify
				notification.SendNotification(
					"Product Create Failed",
					fmt.Sprintf("Message:%s, DevMessage:%s",
						iodata[i].Error.Message, iodata[i].Error.DeveloperMessage,
					),
					[]string{proUtil.TAG_PRODUCT_CREATE, proUtil.TAG_PRODUCT},
					datadog.ERROR,
				)
			} else {
				go func() {
					defer proUtil.RecoverHandler("Insert Product")
					p.Publish("", true)
				}()
				//Sync Product to mysql
				taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT, *p)
				iodata[i].setSuccess(p)
			}
		} else {
			//update
			simpleId, err, isDuplicate := cs.UpdateExistingProduct(proIoData.ReqData, p)
			if err != nil {
				logger.Error(err)
				iodata[i].setFailure(constants.AppError{
					Code:             constants.DbErrorCode,
					Message:          "Product Updation failed",
					DeveloperMessage: err.Error(),
				})
				//notify
				notification.SendNotification(
					"Product Create Failed",
					fmt.Sprintf("Message:%s, DevMessage:%s",
						iodata[i].Error.Message, iodata[i].Error.DeveloperMessage,
					),
					[]string{proUtil.TAG_PRODUCT_CREATE, proUtil.TAG_PRODUCT},
					datadog.ERROR,
				)
			} else {

				go func() {
					defer proUtil.RecoverHandler("Insert/Update Product")
					p.Publish("", true)
				}()
				//Sync Product to mysql
				// Instead of syncing the whole product, lets just sync the simple that got changed.
				productCopy := *p
				for _, sim := range p.Simples {
					if sim.Id == simpleId {
						productCopy.Simples = []*proUtil.ProductSimple{
							sim,
						}
						break
					}
				}
				taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT, productCopy)
				if isDuplicate {
					iodata[i].IsDuplicate = true
				}
				iodata[i].setSuccess(p)
			}
		}
		//release lock
		mutex.UnLock()
	}
	io.IOData.Set(proUtil.IODATA, iodata)
	logger.Debug("Exit Insert node")
	return io, nil
}

//
// get product based on the supplied productSet ID
//
func (cs InsertNode) GetProduct(productSet int, configId int) (*proUtil.Product, error) {
	pr := &proUtil.Product{}
	var err error
	if configId > 0 {
		err = pr.LoadBySeqId(configId, storageAdapter)
	} else {
		err = pr.LoadByProductSet(productSet)
	}
	if err == nil {
		return pr, nil
	}
	if err == proUtil.NotFoundErr {
		return nil, nil
	}
	return nil, fmt.Errorf("(cs InsertNode)#GetProduct: %s", err.Error())
}

//
// SAVE new product in database
// CASE: Config and Simple both are New.
//
func (cs InsertNode) CreateNewProduct(data *ProductCreateRequestData) (
	*proUtil.Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, "CreateNewProduct")
	defer func() {
		logger.EndProfile(profiler, "CreateNewProduct")
	}()
	p, err := proUtil.CreateNewProduct(storageAdapter)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct1: %s", err.Error())
	}
	p.BrandId = data.Brand
	p.Name = data.Title
	p.Description = data.Description
	p.ProductSet = data.ProductSet
	p.TaxClass = &data.TaxClass
	p.SellerId = data.Seller
	p.SetProductGroup(data.GroupName)
	p.SetUrlKey()
	p.SetAttributeSet(data.AttributeSet, proUtil.DB_READ_ADAPTER)
	p.PrimaryCategory = data.PrimaryCategory
	p.Categories = data.Categories
	p.PrepareLeafCategories(proUtil.DB_READ_ADAPTER)
	p.PrepareSKU(proUtil.DB_READ_ADAPTER)
	p.ShipmentType = data.ShipmentType
	p.Status = data.Status
	currentTime := time.Now()
	p.SetPetStatus(proUtil.PetStatus{Created: true, Edited: true})

	//@PATCH STARTS
	temp := 1
	data.ApprovalStatus = &temp
	//@PATCH ENDS

	//Mark pet approved as always 1.
	if data.ApprovalStatus == nil {
		p.PetApproved = 1
	} else {
		p.PetApproved = *data.ApprovalStatus
	}
	p.ActivatedAt = &currentTime
	p.ApprovedAt = &currentTime

	//Set label to identify that product is created through SC.
	err = data.Attributes.SetSCProduct(proUtil.MRK_AS_SC_PRODUCT, false)
	if err != nil {
		return nil, fmt.Errorf(
			"(cs InsertNode)#CreateNewProduct2, unable to set SCProduct: %s",
			err.Error(),
		)
	}
	// [starts] set PT and dispatch location
	seller, err := p.GetSellerInfoWithRetries(conf.Org, false)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct3, cannot get seller: %s", err.Error())
	}
	err = data.Attributes.SetDispatchLocation(seller.SellerCustInfo.DispatchLocation)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct4, dispatch location: %s", err.Error())
	}
	err = data.Attributes.SetProcessingTime(seller.SellerCustInfo.ProcessingTime, seller.Id, p.BrandId)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct5, processing time: %s", err.Error())
	}
	err = data.Attributes.SetOldAttributes(p.AttributeSet.Id)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct6: %s", err.Error())
	}
	err = data.Attributes.SetPrePack(0)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct7: %s", err.Error())
	}
	err = p.SetAttributes(data.Attributes, true, proUtil.DB_READ_ADAPTER, false)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct8: %s", err.Error())
	}
	_, err = p.SetCatalogType(data.Attributes, proUtil.DB_READ_ADAPTER)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct9: %s", err.Error())
	}
	simple, err := proUtil.CreateNewSimple(storageAdapter)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct10: %s", err.Error())
	}
	err = simple.SetSKU(p.SKU)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct11: %s", err.Error())
	}
	simple.Price = &data.Price
	simple.OriginalPrice = &data.Price
	simple.SpecialPrice = data.SpecialPrice
	if data.SpecialFromDate != nil {
		t, err := proUtil.FromMysqlTime(*data.SpecialFromDate, true)
		if err != nil {
			logger.Error(err)
		}
		simple.SpecialFromDate = t
	}
	if data.SpecialToDate != nil {
		t, err := proUtil.FromMysqlTime(*data.SpecialToDate, true)
		if err != nil {
			logger.Error(err)
		}
		*t = (*t).Add(time.Duration(proUtil.TO_DATE_DIFF) * time.Second)
		simple.SpecialToDate = t
	}

	simple.CreationSource = &data.PlatformIdentifier
	simple.TaxClass = &data.TaxClass
	simple.Status = data.Status
	simple.SellerSKU = data.SellerSKU
	simple.BarcodeEan = data.SellerSKU
	simple.EanCode = data.EanCode
	err = simple.SetAttributes(p.AttributeSet.Id, data.Attributes,
		true, proUtil.DB_READ_ADAPTER, false)
	if err != nil {
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct12: %s", err.Error())
	}
	simple.SetStandardSize(&p, storageAdapter)
	p.AppendSimple(simple)
	// Update the product with sizechart data
	p, err = proSizeChUtil.UpdateProductWithSizeChart(p)
	if err != nil {
		logger.Error(err)
	}
	err = p.InsertOrUpdate(storageAdapter)
	if err != nil {
		logger.Error("Product Insertion Failed:" + err.Error())
		return nil, fmt.Errorf("(cs InsertNode)#CreateNewProduct13: %s", err.Error())
	}
	return &p, nil
}

//
// SAVE existing product
// CASE: config is old but simple is new.
// bool => Add a flag to check if its duplicate product.
//
func (cs InsertNode) UpdateExistingProduct(
	data *ProductCreateRequestData, p *proUtil.Product) (int, error, bool) {

	//check if we already have a product [sellerId + sellerSku]
	for _, sim := range p.Simples {
		if sim.SellerSKU == data.SellerSKU {
			//we already have product with this seller sku.
			return sim.Id, nil, true
		}
	}
	p.BrandId = data.Brand
	p.Name = data.Title
	p.Description = data.Description
	p.ProductSet = data.ProductSet
	p.TaxClass = &data.TaxClass
	p.SellerId = data.Seller
	p.SetProductGroup(data.GroupName)
	p.SetUrlKey()
	p.SetAttributeSet(data.AttributeSet, proUtil.DB_READ_ADAPTER)
	p.PrimaryCategory = data.PrimaryCategory
	p.Categories = data.Categories
	p.PrepareLeafCategories(proUtil.DB_READ_ADAPTER)
	p.ShipmentType = data.ShipmentType
	p.Status = data.Status
	petstats := p.GetPetStatus()
	petstats.Edited = true
	p.SetPetStatus(petstats)
	scvalue := p.GetScProductValue()
	data.Attributes.SetSCProduct(scvalue, true)
	err := data.Attributes.SetOldAttributes(p.AttributeSet.Id)
	if err != nil {
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct1: %s", err.Error()), false
	}
	err = p.SetAttributes(data.Attributes, true, proUtil.DB_READ_ADAPTER, true)
	if err != nil {
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct2: %s", err.Error()), false
	}
	simple, err := proUtil.CreateNewSimple(storageAdapter)
	if err != nil {
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct3: %s", err.Error()), false
	}
	err = simple.SetSKU(p.SKU)
	if err != nil {
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct4: %s", err.Error()), false
	}
	simple.Price = &data.Price
	simple.OriginalPrice = &data.Price
	simple.SpecialPrice = data.SpecialPrice
	if data.SpecialFromDate != nil {
		t, err := proUtil.FromMysqlTime(*data.SpecialFromDate, true)
		if err != nil {
			logger.Error(err)
		}
		simple.SpecialFromDate = t
	}
	if data.SpecialToDate != nil {
		t, err := proUtil.FromMysqlTime(*data.SpecialToDate, true)
		if err != nil {
			logger.Error(err)
		}
		*t = (*t).Add(time.Duration(proUtil.TO_DATE_DIFF) * time.Second)
		simple.SpecialToDate = t
	}
	simple.SellerSKU = data.SellerSKU
	simple.CreationSource = &data.PlatformIdentifier
	simple.TaxClass = &data.TaxClass
	simple.Status = data.Status
	simple.BarcodeEan = data.SellerSKU
	simple.EanCode = data.EanCode
	err = simple.SetAttributes(p.AttributeSet.Id, data.Attributes,
		true, proUtil.DB_READ_ADAPTER, false)
	if err != nil {
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct5: %s", err.Error()), false
	}
	simple.SetStandardSize(p, storageAdapter)
	p.AppendSimple(simple)
	// Update the product with sizechart data
	*p, _ = proSizeChUtil.UpdateProductWithSizeChart(*p)
	err = p.InsertOrUpdate(storageAdapter)
	if err != nil {
		logger.Error("Product Updation failed:" + err.Error())
		return 0, fmt.Errorf("(cs InsertNode)#UpdateExistingProduct6: %s", err.Error()), false
	}
	return simple.Id, nil, false
}
