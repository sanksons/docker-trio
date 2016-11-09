package servicetest

import (
	"github.com/jabong/floRest/src/service"
)

var globalConfigJson = `
    {
	   "AppName":"newApp",
	   "AppVersion":"1.0.0",
	   "ServerPort":"8080",
	   "LogConfFile":"conf/logger.json",
	   "MonitorConfig":{  
	      "AppName":"newApp",
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
	   "Hello":{  
	      "ResponseHeaders":{  
	         "CacheControl":{  
	            "ResponseType":"public",
	            "NoCache":false,
	            "NoStore":false,
	            "MaxAgeInSeconds":300
	         }
	      }
	   }
	}
`

func initTestConfig() {
	cm := new(service.ConfigManager)
	cm.InitializeGlobalConfigFromJson(globalConfigJson)
	cm.InitializeAppConfigFromJson(appConfigJson)
}
