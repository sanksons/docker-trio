package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
)

var cacheImpl cache.CacheInterface
var refreshInterval int
var configKey string
var initialized bool = false

type DynamicConfigManager struct {
	applicationConfig interface{}
}

/**
 * Initialize dynamic config manager and start the refresh timer
 */
func (dcm *DynamicConfigManager) Initialize(applicationConfig interface{}) {
	//Check if the Dynamic config is already initialized
	if initialized {
		return
	}
	dcm.applicationConfig = applicationConfig
	initialized = true
	dynamicConfObj := config.GlobalAppConfig.DynamicConfig
	if dynamicConfObj.Active == false {
		logger.Info("Dynamic configuration is not active. Hence, config auto-refresh will not happen at runtime")
		return
	}
	var err error
	cacheImpl, err = cache.Get(dynamicConfObj.Cache)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize cache to auto-refresh the config \n %s \n %s", dynamicConfObj.Cache, err))
	}

	refreshInterval = dynamicConfObj.RefreshInterval
	configKey = dynamicConfObj.ConfigKey

	go dcm.refreshConfigAtEveryInterval()
	fmt.Println("Dynamic config is initialized")
}

/**
 * Starts the timer to refresh the config at every refresh interval
 */
func (dcm *DynamicConfigManager) refreshConfigAtEveryInterval() {
	refreshNow := time.NewTicker(time.Second * time.Duration(refreshInterval)).C
	for {
		select {
		case <-refreshNow:
			dcm.refreshConfig()
		}
	}
}

/**
 * Gets the updated config from Blitz and merges with the current application config
 */
func (dcm *DynamicConfigManager) refreshConfig() {
	logger.Info("Refreshing the config. Time now : " + time.Now().String())
	data, _ := cacheImpl.Get(configKey, true, true)
	if data != nil && data.Value != nil {
		var raw json.RawMessage
		newAppConfig := dcm.applicationConfig
		if dataValue, ok := data.Value.(string); ok {
			raw = json.RawMessage(dataValue)
		} else {
			logger.Warning(fmt.Sprintf("Error - cannot convert to type string"))
			return
		}
		dataInBytes, err := json.Marshal(&raw)
		if err != nil {
			logger.Warning(fmt.Sprintf("Incorrect Json. Error - %s", err))
			return
		}
		err = json.Unmarshal(dataInBytes, newAppConfig)
		if err != nil {
			logger.Warning(fmt.Sprintf("Incorrect Json. Error - %s", err))
			return
		}
		configCopy := config.ApplicationConfig
		// Mergo : Library to merge structs and maps in Golang.
		err = mergo.MergeWithOverwrite(config.ApplicationConfig, newAppConfig)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to merge the application config. Error - %s", err))
			config.ApplicationConfig = configCopy
		}
	} else {
		logger.Warning("Could not find the dynamic config - key : " + configKey + " in central config cache")
	}
}
