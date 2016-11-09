package sqldb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

type mysqlDriver struct {
	db *sql.DB
}

func (obj *mysqlDriver) init(conf *Config) (aerr *SqlDbError) {
	var err error
	// open connection
	obj.db, err = sql.Open(MYSQL, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Dbname,
		url.QueryEscape(conf.Timezone),
	))
	if err == nil {
		// set max open
		obj.db.SetMaxOpenConns(conf.MaxOpenCon)
		// set max idle
		obj.db.SetMaxIdleConns(conf.MaxIdleCon)
		// try a ping
		err = obj.db.Ping()
	}
	// set error if needed
	if err != nil {
		aerr = getErrObj(ERR_INITIALIZATION, err.Error())
	}
	// return
	return aerr
}

func (obj *mysqlDriver) Query(query string, args ...interface{}) (*sql.Rows, *SqlDbError) {
	rows, err := obj.db.Query(query, args...)
	if err != nil {
		return nil, getErrObj(ERR_QUERY_FAILURE, err.Error())
	} else {
		return rows, nil
	}
}

func (obj *mysqlDriver) Execute(query string, args ...interface{}) (sql.Result, *SqlDbError) {
	res, err := obj.db.Exec(query, args...)
	if err != nil {
		return nil, getErrObj(ERR_EXECUTE_FAILURE, err.Error())
	} else {
		return res, nil
	}
}

func (obj *mysqlDriver) GetTxnObj() (*sql.Tx, *SqlDbError) {
	txn, err := obj.db.Begin()
	if err != nil {
		return nil, getErrObj(ERR_GETTXN_FAILURE, err.Error())
	} else {
		return txn, nil
	}
}

func (obj *mysqlDriver) Ping() *SqlDbError {
	err := obj.db.Ping()
	if err != nil {
		return getErrObj(ERR_PING_FAILURE, err.Error())
	} else {
		return nil
	}
}

func (obj *mysqlDriver) Close() *SqlDbError {
	err := obj.db.Close()
	if err != nil {
		return getErrObj(ERR_CLOSE_FAILURE, err.Error())
	} else {
		return nil
	}
}
