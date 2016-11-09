package appconfig

import (
	"os"

	"github.com/jabong/floRest/src/common/cache"
	florest_config "github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/mongodb"
	"github.com/jabong/floRest/src/common/sqldb"
)

// AppConfig struct
type AppConfig struct {
	florest_config.Application
	MySqlConfig         RDBConfig
	Cache               cache.Config
	JBus                JBusConfig
	Bootstrap           BootstrapConfig
	MongoDbConfig       mongodb.Config
	Org                 HttpAPIConfig
	Boutique            BoutiqueConfig
	CCapi               CCAPIConfig
	Redis               RedisConfig
	DbAdapter           string
	ProductConLimit     int
	ProductSyncing      ProductSyncingConfig
	Datadog             DatadogAPI
	NotifAddr           string
	AutoSequencing      bool
	JudgeDaemon         JudgeDaemonConfig
	ProductSellerUpdate bool
	SellerSkuLimit      int
	Sellers             map[string]string
	BrandsProcTime      map[string]string
}

// DatadogAPI struct
type DatadogAPI struct {
	APIKey string
	AppKey string
}

//Product Bootstrap config
type BootstrapConfig struct {
	WorkerCount string
	QueueSize   string
}

type ProductSyncingConfig struct {
	MaxRoutines string
	SleepTime   string
	MarkQCDown  string
}

// RDBConfig struct
type RDBConfig struct {
	Master sqldb.Config
	Slave  sqldb.Config
}

//Http API config
type HttpAPIConfig struct {
	Host    string
	Path    string
	Timeout int
}

//Http API config
type CCAPIConfig struct {
	Host     string
	Path     string
	Username string
	Password string
	Timeout  int
}

// BoutiqueConfig struct
type BoutiqueConfig struct {
	External HttpAPIConfig
	Internal HttpAPIConfig
}

type JBusConfig struct {
	URL        string
	Publisher  string
	RoutingKey string
}

// RedisConfig struct
type RedisConfig struct {
	Styloko RedisPoolConfig
	Stock   RedisPoolConfig
}

type RedisPoolConfig struct {
	Host     string
	Password string
	PoolSize int
}

//Judge Daemon struct
type JudgeDaemonConfig struct {
	ApiUrl        string
	ChunkSize     int
	Timeout       string
	RunEveryJob   string
	RunCleanupJob string
	RunResetJob   string
}

// MapEnvVariables -> Returns map of config values to be replaced by environment variables
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
	overrideVar["MongoDbConfig.Url"] = "MONGO_CONNECTION_URL"
	overrideVar["MongoDbConfig.DbName"] = "MONGO_DBNAME"
	overrideVar["Redis.Styloko.Host"] = "REDIS_HOSTS"
	overrideVar["Redis.Styloko.PoolSize"] = "REDIS_POOLSIZE"
	overrideVar["Redis.Stock.Host"] = "REDIS_STOCK_HOSTS"
	overrideVar["Redis.Stock.PoolSize"] = "REDIS_STOCK_POOLSIZE"
	overrideVar["JBus.URL"] = "JBUS_URL"
	overrideVar["JBus.Publisher"] = "JBUS_PUB"
	overrideVar["JBus.RoutingKey"] = "JBUS_RKEY"
	overrideVar["Bootstrap.WorkerCount"] = "BOOTPOOL_WORKER_COUNT"
	overrideVar["ProductSyncing.MaxRoutines"] = "PRO_SYC_ROUTINES"
	overrideVar["ProductSyncing.SleepTime"] = "PRO_SYC_SLEEP"
	overrideVar["ProductSyncing.MarkQCDown"] = "PRO_MRK_QCDWN"
	overrideVar["NotifAddr"] = "NOTIF_ADDR"
	overrideVar["CCapi.Username"] = "CCAPI_USERNAME"
	overrideVar["CCapi.Password"] = "CCAPI_PASSWORD"
	overrideVar["CCapi.Host"] = "CCAPI_URL"
	overrideVar["ProductConLimit"] = "PRO_MIG_LIMIT"
	overrideVar["JudgeDaemon.ApiUrl"] = "DAEMON_API_URL"
	overrideVar["JudgeDaemon.ChunkSize"] = "DAEMON_CHUNK_SIZE"
	overrideVar["ProductSellerUpdate"] = "PRODUCT_SELLER_UPDATE"
	overrideVar["SellerSkuLimit"] = "SELLER_SKU_LIMIT"

	checkEnv(overrideVar)
	return overrideVar
}

func MapEnvGlobalVariables() map[string]string {
	overrideVar := make(map[string]string)
	overrideVar["MonitorConfig.Enabled"] = "DATADOG_ENABLE"
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
