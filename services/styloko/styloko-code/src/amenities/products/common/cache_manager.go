package common

import (
	sizeUtils "amenities/sizechart/common"
	"common/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type CacheManager struct {
	CacheObj cache.CacheInterface
}

const (
	CACHE_KEY_FORMAT_ID      = "product-%s-%s-%d"
	CACHE_KEY_FORMAT_SKU     = "product-%s-%s-%s"
	CACHE_DELETE_RETRY_COUNT = 3
)

//
// Accepts map of [configId]sku
//
func (cm CacheManager) purge(keys map[int]string) {
	if keys == nil || len(keys) <= 0 {
		return
	}
	var keysToBeDel []string
	for id, sku := range keys {
		keyBySku := fmt.Sprintf("%s-%s", sizeUtils.SizeChartCacheKey, sku)
		keyByConfig := fmt.Sprintf("%s-%d", sizeUtils.SizeChartCacheKey, id)
		keysToBeDel = append(keysToBeDel, keyBySku, keyByConfig)
		for _, expanse := range GetAllExpanse() {
			for _, visibility := range GetAllVisibility() {
				keyId := strings.ToLower(fmt.Sprintf(
					CACHE_KEY_FORMAT_ID,
					expanse, visibility, id,
				))
				keySku := strings.ToLower(fmt.Sprintf(
					CACHE_KEY_FORMAT_SKU,
					expanse, visibility, sku,
				))
				keysToBeDel = append(keysToBeDel, keyId, keySku)
			}
		}
		//ebable chunking
		if len(keysToBeDel) > 200 {
			//purge this chunk
			go func(keys []string) {
				utils.RecoverHandler("Cache deletion failed.")
				ok := cm.deleteWithRetry(keys)
				if !ok {
					logger.Error("Multi-delete from cache failed.")
				}
			}(keysToBeDel)
			keysToBeDel = []string{}
		}
	}
	ok := cm.deleteWithRetry(keysToBeDel)
	if !ok {
		logger.Error("Multi-delete from cache failed.")
	}
}

// delete product cache by id
func (cm CacheManager) DeleteById(productIds []int, deleteBySkus bool) {
	if len(productIds) <= 0 {
		return
	}
	keys := make(map[int]string, 0)
	groups := make(map[int]string, 0)
	for _, id := range productIds {
		pro, err := GetAdapter(DB_ADAPTER_MONGO).GetById(id)
		if err != nil {
			logger.Error(err)
			continue
		}
		keys[pro.SeqId] = pro.SKU
		if pro.Group == nil || pro.Group.Id <= 0 {
			continue
		}
		groups[pro.Group.Id] = pro.Group.Name
	}
	for gId, _ := range groups {
		pros, err := GetAdapter(DB_ADAPTER_MONGO).GetProductsByGroupId(gId)
		if err != nil {
			logger.Error(err)
			continue
		}
		for _, p := range pros {
			keys[p.SeqId] = p.SKU
		}
	}
	cm.purge(keys)
}

// delete product cache by sku
func (cm CacheManager) DeleteBySku(skus []string, deleteByIds bool) {
	if len(skus) <= 0 {
		return
	}
	keys := make(map[int]string, 0)
	groups := make(map[int]string, 0)
	for _, sku := range skus {
		pro, err := GetAdapter(DB_ADAPTER_MONGO).GetBySku(sku)
		if err != nil {
			logger.Error(err)
			continue
		}
		keys[pro.SeqId] = pro.SKU
		if pro.Group == nil || pro.Group.Id <= 0 {
			continue
		}
		groups[pro.Group.Id] = pro.Group.Name
	}
	for gId, _ := range groups {
		pros, err := GetAdapter(DB_ADAPTER_MONGO).GetProductsByGroupId(gId)
		if err != nil {
			logger.Error(err)
			continue
		}
		for _, p := range pros {
			keys[p.SeqId] = p.SKU
		}
	}
	cm.purge(keys)
}

// deleteWithRetry deletes keys from cache with provided keys with retries built-in
func (cm CacheManager) deleteWithRetry(keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	for x := 0; x <= CACHE_DELETE_RETRY_COUNT; x++ {
		err := cm.CacheObj.DeleteBatch(keys)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		return true
	}
	return false
}

//
// Sets Product information in cache
//
func (cm CacheManager) Set(
	productId int,
	expanse string,
	visibility string,
	data interface{},
	cacheTTL int,
) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_SET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_SET)
	}()

	e, _ := json.Marshal(data)
	i := cache.Item{
		Key: strings.ToLower(fmt.Sprintf(
			CACHE_KEY_FORMAT_ID,
			expanse, visibility, productId,
		)),
		Value: string(e),
	}
	err := cm.CacheObj.SetWithTimeout(i, false, false, int32(cacheTTL))
	if err != nil {
		logger.Error(err.Error())
	}
}

//
// Get product informationn from cache
//
func (cm CacheManager) GetAllProducts(
	productId []int,
	expanse string,
	visibility string,
) (map[int]*Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_GET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_GET)
	}()
	//prepare keys to be fetched
	var keys []string
	for _, proId := range productId {
		key := strings.ToLower(fmt.Sprintf(
			CACHE_KEY_FORMAT_ID,
			expanse, visibility, proId,
		))
		keys = append(keys, key)
	}
	items, err := cm.CacheObj.GetBatch(keys, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[1]:%s", err.Error()))
		return nil, err
	}
	cdata := make(map[int]*Product, len(productId))
	for k, item := range items {
		kArr := strings.Split(k, "-")
		id, err1 := strconv.Atoi(kArr[len(kArr)-1])
		if err1 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[2]:%s", err1.Error()))
			continue
		}
		if item == nil {
			cdata[id] = nil
			continue
		}
		var v Product
		itemStr, ok := item.Value.(string)
		if !ok || itemStr == "" {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s,%v", "Empty response from Blitz for", item))
			continue
		}
		tmp := json.NewDecoder(strings.NewReader(itemStr))
		tmp.UseNumber()
		err2 := tmp.Decode(&v)
		// err2 := json.Unmarshal([]byte(itemStr), &v)
		if err2 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[4]:%s", err2.Error()))
			continue
		}
		cdata[id] = &v
	}
	return cdata, err
}

// get product information from cache
func (cm CacheManager) GetAllProductsBySku(
	skus []string,
	expanse string,
	visibility string,
) (map[string]*Product, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_GET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_GET)
	}()
	//prepare keys to be fetched
	var keys []string
	for _, sku := range skus {
		key := strings.ToLower(fmt.Sprintf(
			CACHE_KEY_FORMAT_SKU,
			expanse, visibility, sku,
		))
		keys = append(keys, key)
	}
	items, err := cm.CacheObj.GetBatch(keys, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[1]:%s", err.Error()))
		return nil, err
	}
	cdata := make(map[string]*Product, len(skus))
	for k, item := range items {
		kArr := strings.Split(k, "-")
		sku := kArr[len(kArr)-1]
		if item == nil {
			cdata[sku] = nil
			continue
		}
		var v Product
		itemStr, ok := item.Value.(string)
		if !ok || itemStr == "" {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s,%v", "Empty response from Blitz", item))
			continue
		}
		tmp := json.NewDecoder(strings.NewReader(itemStr))
		tmp.UseNumber()
		err1 := tmp.Decode(&v)
		// err1 := json.Unmarshal([]byte(itemStr), &v)
		if err1 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s", err1.Error()))
			continue
		}
		cdata[sku] = &v
	}
	return cdata, err
}

//
// Get product informationn from cache
//
func (cm CacheManager) GetAll(
	productId []int,
	expanse string,
	visibility string,
) (map[int]*ProductCache, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_GET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_GET)
	}()
	//prepare keys to be fetched
	var keys []string
	for _, proId := range productId {
		key := strings.ToLower(fmt.Sprintf(
			CACHE_KEY_FORMAT_ID,
			expanse, visibility, proId,
		))
		keys = append(keys, key)
	}
	items, err := cm.CacheObj.GetBatch(keys, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[1]:%s", err.Error()))
		return nil, err
	}
	cdata := make(map[int]*ProductCache, len(productId))
	for k, item := range items {
		kArr := strings.Split(k, "-")
		id, err1 := strconv.Atoi(kArr[len(kArr)-1])
		if err1 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[2]:%s", err1.Error()))
			continue
		}
		if item == nil {
			cdata[id] = nil
			continue
		}
		var v ProductCache
		itemStr, ok := item.Value.(string)
		if !ok || itemStr == "" {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s:%v", "Empty response from Blitz for key", item.Key))
			continue
		}
		tmp := json.NewDecoder(strings.NewReader(itemStr))
		tmp.UseNumber()
		err2 := tmp.Decode(&v)
		// err2 := json.Unmarshal([]byte(itemStr), &v)
		if err2 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAll[4]:%s", err2.Error()))
			continue
		}
		cdata[id] = &v
	}
	return cdata, err
}

// get product information from cache
func (cm CacheManager) GetAllBySku(
	skus []string,
	expanse string,
	visibility string,
) (map[string]*ProductCache, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_GET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_GET)
	}()
	//prepare keys to be fetched
	var keys []string
	for _, sku := range skus {
		key := strings.ToLower(fmt.Sprintf(
			CACHE_KEY_FORMAT_SKU,
			expanse, visibility, sku,
		))
		keys = append(keys, key)
	}
	items, err := cm.CacheObj.GetBatch(keys, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[1]:%s", err.Error()))
		return nil, err
	}
	cdata := make(map[string]*ProductCache, len(skus))
	for k, item := range items {
		kArr := strings.Split(k, "-")
		sku := kArr[len(kArr)-1]
		if item == nil {
			cdata[sku] = nil
			continue
		}
		var v ProductCache
		itemStr, ok := item.Value.(string)
		if !ok || itemStr == "" {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s,%v", "Empty response from Blitz", item))
			continue
		}
		tmp := json.NewDecoder(strings.NewReader(itemStr))
		tmp.UseNumber()
		err1 := tmp.Decode(&v)
		// err1 := json.Unmarshal([]byte(itemStr), &v)
		if err1 != nil {
			logger.Error(fmt.Sprintf("(cm CacheManager)#GetAllSku[3]:%s", err1.Error()))
			continue
		}
		cdata[sku] = &v
	}
	return cdata, err
}

// set product information by sku
func (cm CacheManager) SetBySku(
	sku string,
	expanse string,
	visibility string,
	data interface{},
	cacheTTL int,
) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRODUCT_CACHE_SET)
	defer func() {
		logger.EndProfile(profiler, PRODUCT_CACHE_SET)
	}()
	e, _ := json.Marshal(data)
	i := cache.Item{
		Key: strings.ToLower(fmt.Sprintf(CACHE_KEY_FORMAT_SKU,
			expanse, visibility, sku,
		)),
		Value: string(e),
	}
	err := cm.CacheObj.SetWithTimeout(i, false, false, int32(cacheTTL))
	if err != nil {
		logger.Error(err.Error())
	}
}
