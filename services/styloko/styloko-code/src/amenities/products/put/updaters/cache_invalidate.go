package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	"common/utils"
	"fmt"
	validator "gopkg.in/go-playground/validator.v8"
	"strconv"
	"strings"
	"time"
)

const (
	CACHE_INV_SELLER    = "seller"
	CACHE_INV_BRAND     = "brand"
	CACHE_INV_CATEGORY  = "category"
	CACHE_INV_SIMPLESKU = "simpleSku"
	CACHE_INV_SIMPLEID  = "simpleId"
	CACHE_INV_CONFIG    = "config"
)

//Update-Type: Cache
type CacheInvalidate struct {
	Id              string `json:"value" validate:"required"`
	Type            string `json:"type" validate:"required"`
	DonotPublish    bool   `json:"doNotPublish" validate:"-"`
	DonotInvalidate bool   `json:"doNotInvalidate" validate:"-"`
}

func (sd *CacheInvalidate) GetAllowedTypes() map[string]string {
	return map[string]string{
		CACHE_INV_SELLER:    "sellerId",
		CACHE_INV_BRAND:     "brandId",
		CACHE_INV_CATEGORY:  "category",
		CACHE_INV_SIMPLESKU: "simpleSku",
		CACHE_INV_SIMPLEID:  "simpleId",
		CACHE_INV_CONFIG:    "config",
	}
}

func (sd *CacheInvalidate) Response(p *proUtil.Product) interface{} {
	return nil
}

func (sd *CacheInvalidate) Update() (proUtil.Product, error) {
	///does not need to update anything here
	p := &proUtil.Product{}
	return *p, nil
}

func (sd *CacheInvalidate) Validate() []string {
	errs := put.Validate.Struct(sd)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	allowedTypes := sd.GetAllowedTypes()
	if _, ok := allowedTypes[sd.Type]; !ok {
		return []string{"Unsupported Type"}
	}
	return nil
}

func (sd *CacheInvalidate) InvalidateCache() error {

	//check if we need to invalidate cache or not.
	if sd.DonotInvalidate {
		return nil
	}

	var result []proUtil.ProductSmall
	var err error
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO)
	switch sd.Type {
	case CACHE_INV_CONFIG:
		configId, _ := utils.GetInt(sd.Id)
		result = []proUtil.ProductSmall{
			proUtil.ProductSmall{
				Id: configId,
			},
		}
	case CACHE_INV_SIMPLEID:
		var proSmall proUtil.ProductSmall
		simpleId, er := strconv.Atoi(sd.Id)
		if er != nil {
			return fmt.Errorf(
				"(sd *CacheInvalidate)#InvalidateCache(): [simpleid:%s]integer conv failed",
				sd.Id,
			)
		}
		proSmall, err = adapter.GetProductIdBySimpleId(simpleId)
		if err != nil {
			return fmt.Errorf(
				"(sd *CacheInvalidate) InvalidateCache()[simpleid:%s] -> %s",
				sd.Id,
				err.Error(),
			)
		}

		result = []proUtil.ProductSmall{proSmall}
	case CACHE_INV_SIMPLESKU:
		var proSmall proUtil.ProductSmall
		proSmall, err = adapter.GetProductIdBySimpleSku(sd.Id)
		if err != nil {
			return err
		}
		result = []proUtil.ProductSmall{proSmall}
	case CACHE_INV_SELLER:
		//get productIds by seller id
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductIdsBySellerId(id)

	case CACHE_INV_BRAND:
		//get productIds by brand id
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductIdsByBrandId(id)

	case CACHE_INV_CATEGORY:
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductIdsByCategoryId(id)
	}

	if err != nil {
		return err
	}
	var pIds []int
	for _, v := range result {
		pIds = append(pIds, v.Id)
	}
	go func() {
		defer proUtil.RecoverHandler("Cache#Invalidate Cache")
		put.CacheMngr.DeleteById(pIds, true)
	}()
	return nil
}

func (sd *CacheInvalidate) Publish() error {

	// check if we need to publish
	if sd.DonotPublish {
		return nil
	}
	var wait bool
	var result []proUtil.Product
	var err error
	adapter := proUtil.GetAdapter(proUtil.DB_ADAPTER_MONGO)
	switch sd.Type {
	case CACHE_INV_CONFIG:
		wait = true
		configId, _ := utils.GetInt(sd.Id)
		var pro proUtil.Product
		pro, err = adapter.GetById(configId)
		result = []proUtil.Product{pro}
	case CACHE_INV_SIMPLEID:
		var pro proUtil.Product
		simpleId, er := strconv.Atoi(sd.Id)
		if er != nil {
			err = fmt.Errorf("(sd *CacheInvalidate)#Publish(): interger conv failed")
			break
		}
		pro, err = adapter.GetProductBySimpleId(simpleId)
		result = []proUtil.Product{pro}
	case CACHE_INV_SIMPLESKU:
		spliArr := strings.Split(sd.Id, "-")
		simpleId, _ := strconv.Atoi(spliArr[1])
		var pro proUtil.Product
		pro, err = adapter.GetProductBySimpleId(simpleId)
		result = []proUtil.Product{pro}

	case CACHE_INV_SELLER:
		//get products by seller id
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductsBySellerId(id)
	case CACHE_INV_BRAND:
		//get products by brand id
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductsByBrandId(id)
	case CACHE_INV_CATEGORY:
		id, err := strconv.Atoi(sd.Id)
		if err != nil {
			return err
		}
		result, err = adapter.GetProductsByCategoryId(id)
	}
	if err != nil {
		return err
	}
	go func(result []proUtil.Product, wait bool) {
		defer proUtil.RecoverHandler("Cache#Publish")
		if wait {
			time.Sleep(time.Second * 10)
		}
		for _, v := range result {
			v.Publish("", true)
			v.PushToMemcache("Cache#Publish")
		}
	}(result, wait)
	return nil
}

//
// Acquire Lock
//
func (sd *CacheInvalidate) Lock() bool {
	return true
}

//
// Release Lock
//
func (sd *CacheInvalidate) UnLock() bool {
	return true
}
