//Update-Type: Video
package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
	"time"
)

type VideoUpdate struct {
	VideoId   int    `json:"videoId" validate:"-"`
	FileName  string `json:"fileName" validate:"required"`
	ThumbNail string `json:"thumbNail" validate:"-"`
	Size      int    `json:"size" validate:"required"`
	Duration  int    `json:"duration" validate:"required"`
	Status    string `json:"status" validate:"required,eq=active|eq=approved|eq=deleted"`
	Hash      string `json:"videoHash" validate:"required"`
	ConfigId  int    `json:"productId" validate:"required"`
}

func (vu *VideoUpdate) Validate() []string {
	errs := put.Validate.Struct(vu)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (vu *VideoUpdate) Update() (proUtil.Product, error) {
	var err error
	var video proUtil.ProductVideo
	if vu.VideoId == 0 {
		video, err = proUtil.NewProductVideo(put.DbAdapterName)
		if err != nil {
			return proUtil.Product{}, fmt.Errorf(
				"(vu *VideoUpdate)#Update(): %s", err.Error(),
			)
		}
	} else {
		video.Id = vu.VideoId
		video.CreatedAt = time.Now()
	}
	video.FileName = vu.FileName
	video.Thumbnail = vu.ThumbNail
	video.Size = vu.Size
	video.Duration = vu.Duration
	video.Status = vu.Status
	video.Hash = vu.Hash
	video.UpdatedAt = time.Now()

	id, err := proUtil.GetAdapter(put.DbAdapterName).SaveVideo(
		vu.ConfigId, video,
	)
	if err != nil {
		return proUtil.Product{}, err
	}
	//set video ID
	vu.VideoId = id
	pNew, err := proUtil.GetAdapter(put.DbAdapterName).GetById(vu.ConfigId)
	if err != nil {
		///log error
		logger.Error(err)
	}
	// Add JOB for syncing to mysql
	taskPool.AddProductSyncJob(pNew.SeqId, proUtil.UPDATE_TYPE_VIDEO, proUtil.M{
		"configId":  pNew.SeqId,
		"videoData": video,
	})
	return pNew, err
}

func (vu *VideoUpdate) InvalidateCache() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).GetById(vu.ConfigId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Video#Invalidate Cache")
		put.CacheMngr.DeleteBySku([]string{p.SKU}, true)
	}()
	return nil
}

func (vu *VideoUpdate) Publish() error {
	pro, err := proUtil.GetAdapter(put.DbAdapterName).GetById(vu.ConfigId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Video#Publish")
		pro.Publish("", true)
	}()
	return nil
}

func (vu *VideoUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["id"] = vu.ConfigId
	response["videoId"] = vu.VideoId
	for _, v := range p.Videos {
		if v.Id == vu.VideoId {
			response["updatedAt"] = v.UpdatedAt
		}
	}
	return response
}

//
// Acquire Lock
//
func (vu *VideoUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (vu *VideoUpdate) UnLock() bool {
	return true
}
