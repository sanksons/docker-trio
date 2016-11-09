package servicetest

import (
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/service"
)

var globalConfigJson = `
    {
	   "AppName":"org",
	   "AppVersion":"1.0.0",
	   "ServerPort":"8083",
	   "LogConfFile":"conf/logger.json",
	   "MonitorConfig":{  
	      "AppName":"org",
	      "Platform":"DatadogAgent",
	      "AgentServer":"datadog:8125",
	      "Verbose":false,
	      "Enabled":true,
	      "MetricsServer":"datadog:8065"
	   },
	   "Performance":{  
	      "UseCorePercentage":100,
	      "GCPercentage":1000
	   }
   }`

var appConfigJson = `
  {
  "Cache": {
    "Platform": "centralCache",
    "Host": "http://blitz:8080",
    "KeyPrefix": "cache/api/v1/buckets/boutique",
    "Disabled": false
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
      "Username":"MASTER_USERNAME",
      "Password" : "MASTER_PASSWORD",
      "Host" : "MASTER_HOST",
      "Port" : "3306",
      "Dbname" : "MASTER_DBNAME",
      "Timezone" : "Asia/Kolkata",
      "MaxOpenCon" : 0,
      "MaxIdleCon" : 0
    },
    "Slave": {
      "Username":"SLAVE_USERNAME",
      "Password" : "SLAVE_PASSWORD",
      "Host" : "SLAVE_HOST",
      "Port" : "3306",
      "Dbname" : "SLAVE_DBNAME",
      "Timezone" : "Asia/Kolkata",
      "MaxOpenCon" : 0,
      "MaxIdleCon" : 0
    }
  },
  "ScSlaveConfig":{
      "DriverName":"mysql",
      "Username":"SC_USERNAME",
      "Password" : "SC_PASSWORD",
      "Host" : "SC_HOST",
      "Port" : "3306",
      "Dbname" : "SC_DBNAME",
      "Timezone" : "Asia/Kolkata",
      "MaxOpenCon" : 0,
      "MaxIdleCon" : 0
    },
  "MongoDbConfig": {
    "Url":"MONGO_CONNECTION_URL",
    "Dbname":"MONGO_DBNAME"
  },
  "JabongBus": {
    "Url": "JABONG_BUS_CONNECTION_URL"
  }
}
`

func initTestConfig() {
	cm := new(service.ConfigManager)
	cm.InitializeGlobalConfigFromJson(globalConfigJson)
	cm.InitializeAppConfigFromJson(appConfigJson)
	cm.UpdateConfigFromEnv(config.ApplicationConfig, "application")
	cm.UpdateConfigFromEnv(config.GlobalAppConfig, "global")
}
