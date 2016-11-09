package ResourceFactory

import (
	"common/appconfig"
	_ "fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/monitor"
	"github.com/jabong/floRest/src/common/sqldb"
)

var MySqlMapSc = map[string]sqldb.SqlDbInterface{}

func GetMySqlDriverSC(adapterName string) (sqldb.SqlDbInterface, error) {
	if _, ok := MySqlMapSc[adapterName]; !ok {
		pool, err := InitMySqlPoolSC(GetSlaveConfigSC())
		if err != nil {
			return pool, err
		}
		MySqlMapSc[adapterName] = pool
	}
	return MySqlMapSc[adapterName], nil
}

func InitMySqlPoolSC(conf sqldb.Config) (sqldb.SqlDbInterface, error) {
	db, err := sqldb.Get(&conf)
	if err != nil {
		monitor.GetInstance().Error("PANIC:", "Unable to Initialize MySql Pool", nil)
		// panic("Unable to Initialize MySql Pool")
	}
	return db, nil
}

func GetSlaveConfigSC() sqldb.Config {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	conf.ScSlaveConfig.DriverName = "mysql"
	return conf.ScSlaveConfig
}
