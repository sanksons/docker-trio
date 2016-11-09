package appconfig

import (
	"common/mongodb"
	"github.com/jabong/floRest/src/common/cache"
	florest_config "github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/sqldb"
	"os"
)

type AppConfig struct {
	florest_config.Application
	MySqlConfig   RDBConfig
	ScSlaveConfig sqldb.Config
	MongoDbConfig mongodb.Config
	Redis         RedisConfig
	Cache         cache.Config
	Erp           ErpConfig
	Styloko       StylokoConfig
	Datadog       DatadogAPI
	NotifAddr     string
}

type RDBConfig struct {
	Master sqldb.Config
	Slave  sqldb.Config
}

// RedisConfig struct
type RedisConfig struct {
	Host     string
	Password string
	PoolSize int
}

type ErpConfig struct {
	Url string
}

type StylokoConfig struct {
	Url     string
	Timeout string
}

// DatadogAPI struct
type DatadogAPI struct {
	APIKey string
	AppKey string
}

//creating map of all config variables to be overwritten by environment variables
func MapEnvVariables() map[string]string {
	overrideVar := make(map[string]string)
	overrideVar["MySqlConfig.Master.Username"] = "MASTER_USERNAME"
	overrideVar["MySqlConfig.Master.Password"] = "MASTER_PASSWORD"
	overrideVar["MySqlConfig.Master.Host"] = "MASTER_HOST"
	overrideVar["MySqlConfig.Master.Dbname"] = "MASTER_DBNAME"
	overrideVar["MySqlConfig.Master.MaxOpenCon"] = "MASTER_MAX_OPEN_CON"
	overrideVar["MySqlConfig.Master.MaxIdleCon"] = "MASTER_MAX_IDLE_CON"
	overrideVar["MySqlConfig.Slave.Username"] = "SLAVE_USERNAME"
	overrideVar["MySqlConfig.Slave.Password"] = "SLAVE_PASSWORD"
	overrideVar["MySqlConfig.Slave.Host"] = "SLAVE_HOST"
	overrideVar["MySqlConfig.Slave.Dbname"] = "SLAVE_DBNAME"
	overrideVar["MySqlConfig.Slave.MaxOpenCon"] = "SLAVE_MAX_OPEN_CON"
	overrideVar["MySqlConfig.Slave.MaxIdleCon"] = "SLAVE_MAX_IDLE_CON"
	overrideVar["ScSlaveConfig.Username"] = "SC_USERNAME"
	overrideVar["ScSlaveConfig.Password"] = "SC_PASSWORD"
	overrideVar["ScSlaveConfig.Host"] = "SC_HOST"
	overrideVar["ScSlaveConfig.Port"] = "SC_PORT"
	overrideVar["ScSlaveConfig.Dbname"] = "SC_DBNAME"
	overrideVar["ScSlaveConfig.MaxOpenCon"] = "SC_MAX_OPEN_CON"
	overrideVar["ScSlaveConfig.MaxIdleCon"] = "SC_MAX_IDLE_CON"
	overrideVar["MongoDbConfig.Url"] = "MONGO_CONNECTION_URL"
	overrideVar["MongoDbConfig.DbName"] = "MONGO_DBNAME"
	overrideVar["Redis.Host"] = "REDIS_HOSTS"
	overrideVar["Redis.PoolSize"] = "REDIS_POOLSIZE"
	overrideVar["Erp.Url"] = "ERP_URL"
	overrideVar["Styloko.Url"] = "STYLOKO_PRODUCT_URL"
	overrideVar["NotifAddr"] = "NOTIF_ADDR"
	checkEnv(overrideVar)
	return overrideVar
}

// checkEnv -> Checks environment variable availability in map, deletes entry if doesn't exist.
func checkEnv(override map[string]string) {
	for key, value := range override {
		if os.Getenv(value) == "" {
			delete(override, key)
		}
	}
}
