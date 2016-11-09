//Update-Type: ImageDel
package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	taskPool "common/pool/tasker"
	validator "gopkg.in/go-playground/validator.v8"
	"strconv"
	"strings"
)

type ImageDel struct {
	ImageIds []int `json:"imageIds" validate:"required"`
	Results  map[string]bool
}

func (id *ImageDel) Validate() []string {
	errs := put.Validate.Struct(id)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (id *ImageDel) Update() (proUtil.Product, error) {
	results := make(map[string]bool)
	for _, id := range id.ImageIds {
		pId, err := proUtil.GetAdapter(put.DbAdapterName).DeleteImage(id)
		var hash = strconv.Itoa(pId) + "#" + strconv.Itoa(id)
		if err == proUtil.NotFoundErr {
			results[hash] = true
		} else if err != nil {
			results[hash] = false
		} else {
			results[hash] = true
			// Add JOB for syncing to mysql
			taskPool.AddProductSyncJob(pId, proUtil.UPDATE_TYPE_IMAGEDEL, id)
		}
	}
	id.Results = results
	p := &proUtil.Product{}

	return *p, nil
}

func (id *ImageDel) InvalidateCache() error {
	var proIds []int
	for hash, ok := range id.Results {
		if ok {
			splt := strings.Split(hash, "#")
			proId, _ := strconv.Atoi(splt[0])
			if proId <= 0 {
				continue
			}
			proIds = proUtil.AppendIfMissingInt(proIds, proId)
		}
	}
	go func() {
		defer proUtil.RecoverHandler("ImageDel#Invalidate Cache")
		put.CacheMngr.DeleteById(proIds, true)
	}()
	return nil
}

func (id *ImageDel) Publish() error {
	var proIds []int
	for hash, ok := range id.Results {
		if ok {
			splt := strings.Split(hash, "#")
			proId, _ := strconv.Atoi(splt[0])
			if proId <= 0 {
				continue
			}
			proIds = proUtil.AppendIfMissingInt(proIds, proId)
		}
	}
	pCollection := proUtil.ProductCollection{}
	err := pCollection.LoadByIds(proIds)
	if err != nil {
		return err
	}
	if pCollection.Count > 0 {
		go func(pros []proUtil.Product) {
			defer proUtil.RecoverHandler("ImageDel#Publish")
			for _, p := range pros {
				p.Publish("", true)
			}
		}(pCollection.Products)
	}
	return nil
}

func (id *ImageDel) Response(*proUtil.Product) interface{} {
	response := make(map[string]bool)
	for hash, ok := range id.Results {
		splt := strings.Split(hash, "#")
		key := splt[1]
		response[key] = ok
	}
	return response
}

//
// Acquire Lock
//
func (id *ImageDel) Lock() bool {
	return true
}

//
// Release Lock
//
func (id *ImageDel) UnLock() bool {
	return true
}
