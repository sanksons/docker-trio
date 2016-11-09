package servicetest

import (
	"common/appconfig"
	_ "common/appconstant"
	//"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/service"
)

var globalConfigJson = `
 {
    "AppName": "catalog",
    "AppVersion": "1.0.0",
    "ServerPort": "8084",
    "LogConfFile": "conf/logger.json",
    "MonitorConfig": {
        "AppName": "styloko",
        "Platform": "DatadogAgent",
        "AgentServer": "datadog:8125",
        "Verbose": false,
        "Enabled": false,
        "MetricsServer": "datadog:8065"
    },
    "Performance": {
        "UseCorePercentage": 100,
        "GCPercentage": 1000
    }
}`

var appConfigJson = `
	{
    "Cache": {
        "Platform": "centralCache",
        "Host": "http://127.0.0.1:8080",
        "KeyPrefix": "cache/api/v1/buckets/boutique",
        "Disabled": false
    },
     "Org": {
        "Host": "http://org:8083",
        "Path": "/org/v1/sellers",
        "Timeout" : 100000
    },
    "ResponseHeaders": {
        "CacheControl": {
            "ResponseType": "public",
            "NoCache": true,
            "NoStore": false,
            "MaxAgeInSeconds": 300
        }
    },
    "MySqlConfig": {
        "Master": {
            "Username": "MASTER_USERNAME",
            "Password": "MASTER_PASSWORD",
            "Host": "MASTER_HOST",
            "Port": "3306",
            "Dbname": "MASTER_DBNAME",
            "Timezone": "Asia/Kolkata",
            "MaxOpenCon": 0,
            "MaxIdleCon": 0
        },
        "Slave": {
            "Username": "SLAVE_USERNAME",
            "Password": "SLAVE_PASSWORD",
            "Host": "SLAVE_HOST",
            "Port": "3306",
            "Dbname": "SLAVE_DBNAME",
            "Timezone": "Asia/Kolkata",
            "MaxOpenCon": 0,
            "MaxIdleCon": 0
        }
    },
    "MongoDbConfig": {
        "Url": "MONGO_CONNECTION_URL",
        "Dbname": "MONGO_DBNAME"
    },
    "JabongBus": {
        "Url": "http://jabongbus:9807"
    }
}
`

func initTestConfig() {
	service.RegisterConfig(new(appconfig.AppConfig))
	cm := new(service.ConfigManager)
	cm.InitializeGlobalConfigFromJson(globalConfigJson)
	cm.InitializeAppConfigFromJson(appConfigJson)
	service.RegisterConfigEnvUpdateMap(appconfig.MapEnvVariables())
	service.RegisterGlobalEnvUpdateMap(appconfig.MapEnvGlobalVariables())
	cm.UpdateConfigFromEnv(config.ApplicationConfig, "application")
	cm.UpdateConfigFromEnv(config.GlobalAppConfig, "global")
}
