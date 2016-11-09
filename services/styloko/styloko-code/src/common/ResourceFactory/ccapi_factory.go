package ResourceFactory

import (
	"bytes"
	"common/appconfig"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
)

var ccApi *CCApi

func initCCApi() {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	ccApi = &CCApi{}
	ccApi.Url = conf.CCapi.Host + conf.CCapi.Path
	ccApi.UserName = conf.CCapi.Username
	ccApi.Password = conf.CCapi.Password
}

func GetCCApiDriver() *CCApi {
	return ccApi
}

const (
	MAX_TRIAL_LIMIT = 3

	ENTITY_PRODUCTS   = "products"
	ENTITY_BRANDS     = "brands"
	ENTITY_CATEGORIES = "categories"
	ENTITY_SIZECHARTS = "sizechart"
	PRO_MEMCACHE_CALL = "prod_memcache_call"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Minute,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	}}

type CCApi struct {
	Url      string
	UserName string
	Password string
}

//
// Fire a post request with retrials.
//
func (api *CCApi) SendPost(data []byte) error {
	// Prepare reqeust
	var url = api.Url + "/boutique/update"
	logger.Debug("sending request to [%s]", url)
	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		logger.Error(err)
		return fmt.Errorf(
			"(api *CCApi)#SendPost cannot create request: %s",
			err.Error(),
		)
	}
	request.Header.Set("Content-Type", "application/json")
	if api.UserName != "" && api.Password != "" {
		request.SetBasicAuth(api.UserName, api.Password)
	}
	var trials int
	var resp *http.Response
	for trials < MAX_TRIAL_LIMIT {
		resp, err = httpClient.Do(request)
		if err != nil {
			logger.Error(err)
			trials += 1
			time.Sleep(2 * time.Microsecond)
			continue
		}
		defer resp.Body.Close()
		break
	}
	if trials == MAX_TRIAL_LIMIT {
		//request failed, even after trials
		return fmt.Errorf("(api *CCApi)#SendPost: Request Failed.")
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Debug(string(body))
		return fmt.Errorf("status code [%d]", resp.StatusCode)
	}
	return nil
}

//
// Send memcache update for products.
//
func (api *CCApi) UpdateProducts(productIds ...int) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, PRO_MEMCACHE_CALL)
	defer logger.EndProfile(profiler, PRO_MEMCACHE_CALL)

	if len(productIds) == 0 {
		return nil
	}
	//prepare data
	configIds := make([]map[string]int, 0)
	for _, id := range productIds {
		configIds = append(configIds, map[string]int{"config_id": id})
	}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			ENTITY_PRODUCTS: configIds,
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("(api *CCApi)#UpdateProducts cannnot marshal:%s", err.Error())
	}
	err = api.SendPost(dataBytes)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

//
// Send memecache update for brands
//
func (api *CCApi) UpdateBrands(brandIds ...int) error {

	//prepare data
	brands := make([]map[string]int, 0)
	for _, id := range brandIds {
		brands = append(brands, map[string]int{"brandId": id})
	}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			ENTITY_BRANDS: brands,
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("(api *CCApi)#UpdateBrands cannnot marshal:%s", err.Error())
	}
	err = api.SendPost(dataBytes)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("(api *CCApi)#UpdateBrands response:%s", err.Error())
	}
	return nil
}

//
// Send memcache update for categories
//
func (api *CCApi) UpdateCategories(cats ...int) error {
	//prepare data
	catIds := make([]map[string]int, 0)
	for _, id := range cats {
		catIds = append(catIds, map[string]int{"categoryId": id})
	}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			ENTITY_CATEGORIES: catIds,
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("(api *CCApi)#UpdateCategories cannnot marshal:%s", err.Error())
	}
	err = api.SendPost(dataBytes)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

//
// Send memcache update for sizechart
//
func (api *CCApi) UpdateSizecharts(cats ...int) error {
	//prepare data
	catIds := make([]map[string]int, 0)
	for _, id := range cats {
		catIds = append(catIds, map[string]int{"catId": id})
	}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			ENTITY_SIZECHARTS: catIds,
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("(api *CCApi)#UpdateSizecharts cannnot marshal:%s", err.Error())
	}
	err = api.SendPost(dataBytes)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
