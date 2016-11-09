package servicetest

import (
	"github.com/jabong/floRest/src/common/logger"
)

func initTestLogger() {
	configLogger := `
	{
	    "ProfilerEnabled": false,
	    "LogLevel": 1,
	    "DefaultLogType": "dummyLoggerDefault",
	    "AppName" : "newApp-test-",
	    "FileLogger": [
	        {
	            "Key": "fileLoggerDefault",
	            "Path": "/tmp/",
	            "FileNamePrefix": "search-suggest-"            
	        }
	    ],
	    "DummyLogger": [
	        {
	            "Key": "dummyLoggerDefault"            
	        }        
    	]
	}`

	logger.InitialiseFromJson(configLogger)
}
