package sqldb

import (
	"testing"
)

func TestSqlDb(t *testing.T) {
	var err error
	var dbObj SqlDbInterface
	conf := new(Config)
	// fill invalid driver name
	conf.DriverName = "invalid" // set driver name as mysql
	_, err = Get(conf)
	if err == nil {
		FailNow(t, "invalid driver must throw error")
	}

	// fill valid driver, invalid db details
	conf.DriverName = MYSQL // set driver name as mysql
	conf.Username = "root"
	conf.Password = "root"
	conf.Host = "localhost"
	conf.Port = "3306"
	conf.Dbname = "invalid"
	conf.Timezone = "Local"
	conf.MaxOpenCon = "2"
	conf.MaxIdleCon = "1"
	dbObj, err = Get(conf)
	if err == nil {
		FailNow(t, "init db must fail, but it passed")
	}
	// As invalid db object, assert error for all methods
	_, err = dbObj.Query("invalid query")
	if err == nil {
		FailNow(t, "query must fail for this invalid db")
	}
	_, err = dbObj.Execute("invalid execute")
	if err == nil {
		FailNow(t, "execute must fail for this invalid db")
	}
	_, err = dbObj.GetTxnObj()
	if err == nil {
		FailNow(t, "get txn object must fail for this invalid db")
	}
	err = dbObj.Close()
	if err == nil {
		FailNow(t, "close must fail for this invalid db")
	}
}

// logs error and exits the test run
func FailNow(t *testing.T, err string) {
	t.Error(err)
	t.FailNow()
}
