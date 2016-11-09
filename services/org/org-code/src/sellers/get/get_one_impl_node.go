package get

import (
	"common/appconfig"
	"common/appconstant"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"io/ioutil"
	"net/http"
	"reflect"
	"sellers/common"
	"strconv"
)

type GetSeller struct {
	id string
}

func (s *GetSeller) SetID(id string) {
	s.id = id
}

func (s GetSeller) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetSeller) Name() string {
	return "GET seller by id"
}

func (s GetSeller) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_SELLER_GET_ONE)
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_GET, "Single seller get execution started")
	data, _ := io.IOData.Get(GET_ONE)
	//converting id from string to int
	id, err := strconv.Atoi(data.(string))
	if err != nil {
		logger.Error(fmt.Sprintf("Error while converting id %v to int :%s", id, err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid Id", DeveloperMessage: "Seller Id should be int"}
	}
	// check from cache
	v, err := s.FromCache(data.(string))
	if err == nil {
		io.IOData.Set(florest_constants.RESULT, v)
		return io, nil
	}
	//Gets data from mongo for the id passed and returns
	//error if no data was found for the id.
	org, err := common.GetById(id, "Get_By_Id")
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.ResourceNotFoundCode, Message: "Invalid Id", DeveloperMessage: "Seller Id does not exist"}
	}
	//setting phn as "" to not show in json
	org.Phone = ""
	//setting data in RESULT
	io.IOData.Set(florest_constants.RESULT, org)
	go s.SetCache(data.(string), org)
	return io, nil
}

// set seller information to cache
func (s GetSeller) SetCache(key string, data interface{}) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_SET_CACHE)
	defer func() {
		logger.EndProfile(profiler, SELLER_SET_CACHE)
	}()
	e, _ := json.Marshal(data)
	i := cache.Item{
		Key:   fmt.Sprintf("%s-%s", common.SELLERS, key),
		Value: string(e),
	}
	err := cacheObj.Set(i, false, false)
	if err != nil {
		logger.Error(err.Error())
	}
}

// get seller information from cache
func (s GetSeller) FromCache(key string) (interface{}, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET_FROM_CACHE)
	logger.StartProfile(profiler, SELLER_GET_FROM_CACHE_ACTUAL)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET_FROM_CACHE)
	}()
	item := useDefaultHttp(fmt.Sprintf("%s-%s", common.SELLERS, key))
	logger.EndProfile(profiler, SELLER_GET_FROM_CACHE_ACTUAL)
	if item == nil {
		return nil, errors.New(fmt.Sprintf("Error while getting key from blitz:%s", key))
	}
	return item, nil
}

func useDefaultHttp(key string) (out interface{}) {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	url := fmt.Sprintf("%s/%s/entity/%s", config.Cache.Host, config.Cache.KeyPrefix, key)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while forming new request:%s", err.Error()))
		return nil
	}
	response, err := client.Do(request)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending request:%s", err.Error()))
		return nil
	}
	defer response.Body.Close()
	request.Close = true
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while reading body of response:%s", err.Error()))
		return nil
	}
	var output Output
	err = json.Unmarshal(contents, &output)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while unmarshalling blitz response:%s", err.Error()))
		return nil
	}
	if output.Status.Success == true {
		t := reflect.ValueOf(output.Data.Value).Kind()
		if t == reflect.String {
			v, c := output.Data.Value.(string)
			if c {
				json.Unmarshal([]byte(v), &out)
			}
		} else {
			out = output.Data.Value
		}
		return out
	}
	return nil
}

type Output struct {
	Status status `json:"status"`
	Data   data   `json:"data"`
}

type status struct {
	Success bool `json:"success"`
}

type data struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
