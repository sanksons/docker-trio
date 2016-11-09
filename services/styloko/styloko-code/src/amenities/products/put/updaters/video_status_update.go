// Update-Type: Video Status
package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	validator "gopkg.in/go-playground/validator.v8"
)

type VideoStatusUpdate struct {
	VideoId int    `json:"videoId" validate:"required"`
	Status  string `json:"status" validate:"required,eq=active|eq=approved|eq=deleted"`
}

func (vs *VideoStatusUpdate) Validate() []string {
	errs := put.Validate.Struct(vs)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (vs *VideoStatusUpdate) Update() (proUtil.Product, error) {
	err := proUtil.GetAdapter(put.DbAdapterName).
		UpdateVideoStatus(vs.VideoId, vs.Status)
	if err != nil {
		return proUtil.Product{}, err
	}

	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductByVideoId(vs.VideoId)

	// Add JOB for syncing to mysql
	taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_VIDEO_STATUS, proUtil.M{
		"videoId": vs.VideoId,
		"status":  vs.Status,
	})
	return p, err
}

func (vs *VideoStatusUpdate) InvalidateCache() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductByVideoId(vs.VideoId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("VideoStatus#Invalidate Cache")
		put.CacheMngr.DeleteBySku([]string{p.SKU}, true)
	}()
	return nil
}

func (vs *VideoStatusUpdate) Publish() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).
		GetProductByVideoId(vs.VideoId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("VideoStatus#Invalidate Cache")
		p.Publish("", true)
	}()
	return nil
}

func (vs *VideoStatusUpdate) Response(p *proUtil.Product) interface{} {
	return vs.VideoId
}

//
// Acquire Lock
//
func (vs *VideoStatusUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (vs *VideoStatusUpdate) UnLock() bool {
	return true
}
