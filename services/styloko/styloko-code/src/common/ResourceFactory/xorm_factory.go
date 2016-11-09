package ResourceFactory

import (
	"common/appconfig"
	"common/xorm/mysql"
	"github.com/jabong/floRest/src/common/config"
)

func initXormDb() {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	instance := mysql.GetInstance()
	instance.Init(conf.MySqlConfig.Master, conf.MySqlConfig.Slave)
}
