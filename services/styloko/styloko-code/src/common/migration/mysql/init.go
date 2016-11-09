package mysql

import (
	_ "common/appconfig"
	_ "database/sql"
	_ "fmt"
	//dotSql "github.com/dotsql-master"
	_ "github.com/jabong/floRest/src/common/config"
	_ "net/url"
)

func InitDotSql() {
	/*
		conf := config.ApplicationConfig.(*appconfig.AppConfig)
		mysqlConf := conf.MySqlConfig.Master
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
			mysqlConf.Username,
			mysqlConf.Password,
			mysqlConf.Host,
			mysqlConf.Port,
			mysqlConf.Dbname,
			url.QueryEscape(mysqlConf.Timezone),
		))
		if err != nil {
			panic("initDotSql(): mySql init failed" + err.Error())
		}

		dotsql, err := dotSql.LoadFromFile("queries.sql")
		if err != nil {
			panic("initDotSql(): cannot load schema file" + err.Error())
		}
		dotsql.Exec(db, "mysql-sync-drop")
		_, err = dotsql.Exec(db, "mysql-sync")
		if err != nil {
			panic("initDotSql(): cannot exec [mysql-sync]" + err.Error())
		}
	*/
	return
}
