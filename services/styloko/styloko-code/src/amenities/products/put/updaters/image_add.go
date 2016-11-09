//Update-Type: ImageAdd
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

type ImageAdd struct {
	OriginalFilename string `json:"originalFileName" validate:"required"`
	IsMain           bool   `json:"isMain" validate:"exists"`
	Orientation      string `json:"orientation" validate:"required,eq=portrait|eq=landscape"`
	ImageNo          int    `json:"imageNo" validate:"required"`
	ConfigId         int    `json:"productId" validate:"required"`
	BGColor          string `json:"bgcolor" validate:"-"`
	ImageId          int    `json:"imageId" validate:"-"`
	UpdateUID        string `json:"-" validate:"-"`
}

func (ia *ImageAdd) Validate() []string {
	errs := put.Validate.Struct(ia)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (ia *ImageAdd) Update() (proUtil.Product, error) {

	p, err := proUtil.GetAdapter(put.DbAdapterName).GetById(ia.ConfigId)
	if err != nil {
		return p, err
	}
	// check if image already exists in system
	imageExists := ia.isImageExists(p)
	if imageExists {
		return p, nil
	}
	image, err := proUtil.NewProductImage(put.DbAdapterName)
	if err != nil {
		return p, err
	}
	image.ImageNo = ia.ImageNo
	image.ImageName = image.PrepareImageName(p)
	image.OriginalFileName = ia.OriginalFilename
	image.Orientation = ia.Orientation
	image.Main = 0
	if ia.IsMain {
		image.Main = 1
	}
	image.IfUpdate = false
	if p.UpdatedAt == nil || p.ActivatedAt == nil || p.ApprovedAt == nil {
		image.IfUpdate = true
	}
	id, err := image.Add(p.SeqId, put.DbAdapterName)
	if err != nil {
		return p, err
	}

	ia.ImageId = id
	//Mark Updated At in Redis.
	ia.UpdateUID, err = p.SetUpdateUID()
	if err != nil {
		//this error can be ignored.
		logger.Error(err)
	}
	pNew := &proUtil.Product{}
	pNew.LoadBySeqId(p.SeqId, put.DbAdapterName)
	if (pNew.ActivatedAt == nil || pNew.ApprovedAt == nil) && pNew.PetApproved == 1 {
		current := time.Now()
		pNew.ActivatedAt = &current
		pNew.ApprovedAt = &current
	}

	// Add JOB for syncing to mysql
	taskPool.AddProductSyncJob(pNew.SeqId, proUtil.UPDATE_TYPE_IMAGEADD, proUtil.M{
		"configId":  p.SeqId,
		"imageData": image,
	})

	//try to set bg color if its main image.
	if ia.IsMain == true {
		err = ia.SetBGColor(pNew.SKU)
		if err != nil {
			logger.Error(fmt.Sprintf("Unable to set BGColor: [%d]", ia.ImageId))
		}
	}
	//check if we need to set petstatus
	if len(p.Images) > 0 {
		return *pNew, nil
	}
	//Update PetStatus for product
	petstats := pNew.GetPetStatus()
	petstats.Image = true
	pNew.SetPetStatus(petstats)
	criteria := proUtil.ProductAttrSystemUpdate{
		ProConfigId: pNew.SeqId,
		AttrName:    proUtil.SYSTEM_PET_STATUS,
		AttrValue:   pNew.PetStatus,
	}
	err = proUtil.GetAdapter(put.DbAdapterName).UpdateProductAttributeSystem(criteria)
	if err != nil {
		return *pNew, err
	}
	// Add Job to update petstatus
	taskPool.AddProductSyncJob(pNew.SeqId, proUtil.SYNC_ATTRIBUTE_SYSTEM, ProductAttributeUpdate{
		AttributeName: proUtil.SYSTEM_PET_STATUS,
		IsGlobal:      true,
		ProductSku:    pNew.SKU,
		ProductType:   proUtil.PRODUCT_TYPE_CONFIG,
		Action:        proUtil.ACTION_REPLACE,
		Value:         pNew.PetStatus,
	})
	return *pNew, err
}

//
// Set BgColor for the Image.
//
func (ia *ImageAdd) SetBGColor(prodSku string) error {
	bgcUpdt := ProductAttributeUpdate{
		AttributeName: proUtil.ATTR_BG_COLOR,
		IsGlobal:      true,
		ProductSku:    prodSku,
		ProductType:   proUtil.PRODUCT_TYPE_CONFIG,
		Action:        proUtil.ACTION_ADD,
		Value:         ia.BGColor,
	}
	_, err := bgcUpdt.Update()
	if err != nil {
		return err
	}
	return nil
}

func (ia *ImageAdd) InvalidateCache() error {
	if !ia.canPublisherInvalidate() {
		return nil
	}
	go func() {
		defer proUtil.RecoverHandler("ImageAdd#Invalidate Cache")
		put.CacheMngr.DeleteById([]int{ia.ConfigId}, true)
	}()
	return nil
}

func (ia *ImageAdd) Publish() error {
	if !ia.canPublisherInvalidate() {
		//skip Publish
		return nil
	}
	pro, err := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GetById(ia.ConfigId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("ImageAdd#Publish")
		pro.Publish("", true)
	}()
	return nil
}

func (ia *ImageAdd) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["productId"] = ia.ConfigId
	if ia.ImageId == 0 || p == nil {
		response["status"] = proUtil.STATUS_FAILURE
		return response
	}
	response["imageId"] = ia.ImageId
	var imagePath string
	var updatedAt string
	for _, v := range p.Images {
		if v.SeqId == ia.ImageId {
			upd := v.UpdatedAt
			updatedAt = proUtil.ToMySqlTime(upd)
			if ia.Orientation == "portrait" {
				imagePath = fmt.Sprintf("/p/%s%d.jpg", v.ImageName, v.ImageNo)
			} else {
				imagePath = fmt.Sprintf("/p/%s%d-ol.jpg", v.ImageName, v.ImageNo)
			}
		}
	}
	response["imagePath"] = imagePath
	response["updatedAt"] = updatedAt
	response["status"] = proUtil.STATUS_SUCCESS
	return response
}

//
// Check if we really need to Invalidate or Publish
//
func (self *ImageAdd) canPublisherInvalidate() bool {
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

//
// Acquire Lock
//
func (ia *ImageAdd) Lock() bool {
	return true
}

//
// Release Lock
//
func (ia *ImageAdd) UnLock() bool {
	return true
}

func (ia *ImageAdd) isImageExists(p proUtil.Product) bool {
	imageExists := false
	for _, prod := range p.Images {
		if prod.ImageNo == ia.ImageNo {
			ia.ImageId = prod.SeqId
			imageExists = true
			break
		}
	}
	return imageExists
}
