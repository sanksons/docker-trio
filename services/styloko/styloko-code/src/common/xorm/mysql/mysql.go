package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-xorm/xorm"
	"github.com/jabong/floRest/src/common/sqldb"
)

//
// type sqldb.Config struct {
// 	Username   string
// 	Password   string
// 	Host       string
// 	Port       string
// 	Dbname     string
// 	Timezone   string
// 	MaxOpenCon int
// 	MaxIdleCon int
// }

type MySqlDb struct {
	writeConn *myCon
	readConn  *myCon
	engines   map[string]*xorm.Engine
}

type myCon struct {
	conn    string
	maxOpen int
	maxIdle int
}

const (
	driver_name = "mysql"
)

var MysqlObj *MySqlDb

func GetInstance() *MySqlDb {
	if MysqlObj == nil {
		MysqlObj = new(MySqlDb)
	}
	return MysqlObj
}

func (obj *MySqlDb) Init(master sqldb.Config, slave sqldb.Config) (err error) {
	obj.writeConn = new(myCon)
	obj.readConn = new(myCon)
	if err = obj.writeConn.setConn(&master); err == nil {
		if err = obj.readConn.setConn(&slave); err == nil {
			obj.engines = make(map[string]*xorm.Engine)
			return nil
		}
	}
	return err
}

func (obj *myCon) setConn(info *sqldb.Config) (err error) {
	obj.conn = connString(info) // save connection string
	obj.maxOpen = info.MaxOpenCon
	obj.maxIdle = info.MaxIdleCon
	return err
}

// return dns connection string
func connString(db *sqldb.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Dbname,
		url.QueryEscape(db.Timezone),
	)
}

// choose which connection object to return
func (obj *MySqlDb) getConn(write bool) (ret *myCon) {
	if write || obj.readConn == nil {
		ret = obj.writeConn
	} else {
		ret = obj.readConn
	}
	// return now
	return ret
}

// get the orm engine
func (obj *MySqlDb) orm(write bool) (orm *xorm.Engine, err error) {
	var conObj = obj.getConn(write)
	var ok bool

	if orm, ok = obj.engines[conObj.conn]; ok {
		err = orm.Ping()
	} else {
		if orm, err = xorm.NewEngine(driver_name, conObj.conn); err == nil {
			orm.SetMaxOpenConns(conObj.maxOpen)
			orm.SetMaxIdleConns(conObj.maxIdle)
			obj.engines[conObj.conn] = orm
			err = orm.Ping()
		}
	}
	// return now
	return orm, err
}

// run query
func (obj *MySqlDb) Query(query string, write bool) (ret interface{}, err error) {
	var engine *xorm.Engine
	if engine, err = obj.orm(write); err == nil {
		xormRows, cerr := engine.DB().Query(query)
		if cerr == nil { // convert to sql rows
			ret = xormRows
		} else {
			err = errors.New("err in runquery:" + cerr.Error())
		}
	}
	// return
	return ret, err
}

// execute query
func (obj *MySqlDb) Execute(query string, write bool) (ret sql.Result, err error) {
	var engine *xorm.Engine
	if engine, err = obj.orm(write); err == nil {
		ret, err = engine.DB().Exec(query)
	}
	// return
	return ret, err
}

// ping master
func (obj *MySqlDb) PingMaster() (err error) {
	_, err = obj.orm(true) // ping write orm
	return err
}

// ping slave
func (obj *MySqlDb) PingSlave() (err error) {
	_, err = obj.orm(false) // ping read orm
	return err
}
