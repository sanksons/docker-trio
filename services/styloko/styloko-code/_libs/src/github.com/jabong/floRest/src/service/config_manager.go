package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/env"
	utilJson "github.com/jabong/floRest/src/common/utils/json"
)

type ConfigManager struct {
}

func (cm *ConfigManager) InitializeGlobalConfig() {
	if config.GlobalAppConfig == nil {
		config.GlobalAppConfig = new(config.AppConfig)
	}
	fmt.Println("initializing global config ")
	fmt.Println(fmt.Sprintf("globalConfig=%+v", config.GlobalAppConfig))
	cm.Initialize("conf/conf.json", config.GlobalAppConfig)
	fmt.Println(fmt.Sprintf("updated globalConfig=%+v", config.GlobalAppConfig))
}

func (cm *ConfigManager) InitializeAppConfig() {
	fmt.Println("initializing App config")
	if config.ApplicationConfig == nil {
		//TODO - error scenario. It should not be nil
		return
	}
	fmt.Println(fmt.Sprintf("config.ApplicationConfig=%+v", config.ApplicationConfig))
	cm.Initialize("conf/api_conf.json", config.ApplicationConfig)
	fmt.Println(fmt.Sprintf("updated config.ApplicationConfig=%+v", config.ApplicationConfig))
}

func (cm *ConfigManager) Initialize(filePath string, conf interface{}) {

	fmt.Println(fmt.Sprintf("config %+v", conf))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading App Config file %s \n %s", filePath, err))
	}
	err = json.Unmarshal(file, conf)
	if err != nil {
		panic(fmt.Sprintf("Incorrect Json in %s \n %s", filePath, err))
	}
	fmt.Println("Application Config Created")
	fmt.Println(fmt.Sprintf("config %+v", conf))
}

func (cm *ConfigManager) InitializeAppConfigFromJson(confJson string) {
	//Check if the ApplicationConfig is already initialized
	if config.ApplicationConfig == nil {
		fmt.Println("Application Config type is not registered")
		return
	}

	err := json.Unmarshal([]byte(confJson), config.ApplicationConfig)
	if err != nil {
		panic(fmt.Sprintf("Incorrect Json %s \n %s", confJson, err))
	}
	fmt.Println("Application Config Created from Json string")
}

func (cm *ConfigManager) InitializeGlobalConfigFromJson(confJson string) {
	//Check if the ApplicationConfig is already initialized
	if config.GlobalAppConfig != nil {
		return
	}

	config.GlobalAppConfig = new(config.AppConfig)

	err := json.Unmarshal([]byte(confJson), config.GlobalAppConfig)
	if err != nil {
		panic(fmt.Sprintf("Incorrect Json %s \n %s", confJson, err))
	}
	fmt.Println("Global Config Created from Json string")
}

// UpdateConfigFromEnv updates provided config from environment variables
func (cm *ConfigManager) UpdateConfigFromEnv(conf interface{}, ty string) {
	if conf == nil {
		return
	}
	localConfigMap := make(map[string]string)
	if ty == "global" {
		if globalEnvUpdateMap == nil {
			return
		}
		localConfigMap = globalEnvUpdateMap
	} else {
		if configEnvUpdateMap == nil {
			return
		}
		localConfigMap = configEnvUpdateMap
	}

	configEnvUpdateValuesMap := make(map[string]string)
	for k, v := range localConfigMap {
		updatedVal, envValfound := env.GetOsEnviron().Get(v)

		if !envValfound {
			panic(errors.New(fmt.Sprintf("Environment variable %s not found", v)))
		}

		if strings.Trim(updatedVal, " ") == "" {
			panic(errors.New(fmt.Sprintf("Environment variable %s is empty", v)))
		}

		configEnvUpdateValuesMap[k] = updatedVal
	}

	byt, _ := json.Marshal(conf)

	newbyt, juerr := utilJson.UpdateJsonPath(configEnvUpdateValuesMap, byt, ".")
	if juerr != nil {
		panic(juerr)
	}

	if uerr := json.Unmarshal(newbyt, &conf); uerr != nil {
		panic(uerr)
	}
	if ty == "global" {
		fmt.Printf("Updated config from environment variables: %+v\n", config.GlobalAppConfig)
	} else {
		fmt.Printf("Updated config from environment variables: %+v\n", config.ApplicationConfig)
	}
}
