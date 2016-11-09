package ResourceFactory

import (
	"common/appconfig"
	_ "fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/monitor"
	"github.com/jabong/floRest/src/common/sqldb"
)

var MySqlMap = map[string]sqldb.SqlDbInterface{}

func GetMySqlDriver(adapterName string) (sqldb.SqlDbInterface, error) {
	if _, ok := MySqlMap[adapterName]; !ok {
		pool, err := InitMySqlPool(GetMasterConfig())
		if err != nil {
			return pool, err
		}
		MySqlMap[adapterName] = pool
	}
	return MySqlMap[adapterName], nil
}

func InitMySqlPool(conf sqldb.Config) (sqldb.SqlDbInterface, error) {
	db, err := sqldb.Get(&conf)
	if err != nil {
		monitor.GetInstance().Error("PANIC:", "Unable to Initialize MySql Pool", nil)
		panic("Unable to Initialize MySql Pool")
	}
	return db, nil
}

func GetMasterConfig() sqldb.Config {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	conf.MySqlConfig.Master.DriverName = "mysql"
	return conf.MySqlConfig.Master
}

func GetSlaveConfig() sqldb.Config {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	conf.MySqlConfig.Master.DriverName = "mysql"
	return conf.MySqlConfig.Slave
}
