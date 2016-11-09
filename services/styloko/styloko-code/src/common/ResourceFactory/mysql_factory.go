package ResourceFactory

import (
	"common/appconfig"
	"errors"
	_ "fmt"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/sqldb"
)

var MySqlMap = map[string]sqldb.SqlDbInterface{}
var MySqlMapSlave = map[string]sqldb.SqlDbInterface{}

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

func GetMySqlDriverSlave(adapterName string) (sqldb.SqlDbInterface, error) {
	if _, ok := MySqlMapSlave[adapterName]; !ok {
		pool, err := InitMySqlPool(GetSlaveConfig())
		if err != nil {
			return pool, err
		}
		MySqlMapSlave[adapterName] = pool
	}
	return MySqlMapSlave[adapterName], nil
}

func InitMySqlPool(conf sqldb.Config) (sqldb.SqlDbInterface, error) {
	db, err := sqldb.Get(&conf)
	if err != nil {
		return nil, errors.New("Unable to Initialize MySql Pool: " + err.DeveloperMessage)
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
	conf.MySqlConfig.Slave.DriverName = "mysql"
	return conf.MySqlConfig.Slave
}

func GetDefaultMysqlDriver() (sqldb.SqlDbInterface, error) {
	return GetMySqlDriver("DEFAULT")
}

func GetDefaultMysqlDriverSlave() (sqldb.SqlDbInterface, error) {
	return GetMySqlDriverSlave("DEFAULT")
}
