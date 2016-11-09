package ResourceFactory

import (
	"common/appconfig"
	"common/mongodb"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
)

/*
 * Mongo factory specific code below.
 */

// mongoInfo -> Global variable required for connection startup
var mongoInfo *mongodb.Config

// GetMongoSession -> Same as adapter, removes confusion.
func GetMongoSession(adapterName string) *mongodb.MongoDriver {
	logger.Info("Mongo session returned for ", adapterName)
	session := mongodb.GetInstance()
	session.Initialize(mongoInfo, "")
	session.Refresh()
	//check if safe mode needs to be enabled.
	if mongodb.UseSafeMode {
		//session.SetSafe()
	}
	return session
}

// GetMongoSessionWithPing -> Returns mongo session by Pinging it. If ping flag is false,
// then connection failure has occured.
func GetMongoSessionWithPing(adapterName string) (*mongodb.MongoDriver, bool) {
	session := GetMongoSession(adapterName)
	ping := session.Ping()
	return session, ping
}

// InitMongoDb -> Initializes Mongo DB with configs and sets database value.
func initMongoDb() {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	mongoInfo = new(mongodb.Config)
	mongoInfo.DbName = conf.MongoDbConfig.DbName
	mongoInfo.Url = conf.MongoDbConfig.Url
	err := mongodb.GetInstance().Initialize(mongoInfo, "")
	if err != nil {
		logger.Error(fmt.Sprintf("Error in initializing Mongo Instance %v", err.DeveloperMessage))
	}
}

// GetMongoSession -> Same as adapter but connects to the db passed in params, removes confusion.
func GetMongoSessionWithDb(adapterName string, dbName string) *mongodb.MongoDriver {
	logger.Info("Mongo session returned for ", adapterName)
	session := mongodb.GetInstance()
	session.Initialize(mongoInfo, dbName)
	session.Refresh()
	//check if safe mode needs to be enabled.
	if mongodb.UseSafeMode {
		//session.SetSafe()
	}
	return session
}
