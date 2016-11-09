package examples

import (
	"fmt"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/sqldb"
)

type HelloWorld struct {
	id string
}

func (n *HelloWorld) SetID(id string) {
	n.id = id
}

func (n HelloWorld) GetID() (id string, err error) {
	return n.id, nil
}

func (a HelloWorld) Name() string {
	return "HelloWorld"
}

func (a HelloWorld) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	// fill sql config
	conf := new(sqldb.Config)
	conf.DriverName = sqldb.MYSQL // set driver name as mysql
	conf.Username = "root"
	conf.Password = "root"
	conf.Host = "localhost"
	conf.Port = "3306"
	conf.Dbname = "bobalice"
	conf.Timezone = "Local"
	conf.MaxOpenCon = 2
	conf.MaxIdleCon = 1
	// get db object
	db, err := sqldb.Get(conf) // It should be called only once and can be shared across go routines
	if err != nil {
		fmt.Println(err)
	} else { // got db object, try methods
		err = db.Ping() // try pinging db
		if err != nil {
			fmt.Println(err)
		}
		// execute a statement: create one table
		_, err = db.Execute("create table florest_employee (name varchar(255))")
		if err != nil {
			fmt.Println(err)
		}
		// raw query: query on table
		rows, qerr := db.Query("SELECT * from florest_employee")
		if qerr != nil {
			fmt.Println(qerr)
		} else {
			rows.Close()
		}
		name := "rajcomics"
		// query with run-time arguments: query on table
		rows, qerr = db.Query("SELECT * from florest_employee where name=?", &name)
		if qerr != nil {
			fmt.Println(qerr)
		} else {
			rows.Close()
		}
		// start and commit one txn: insert one row in table
		if txObj, terr := db.GetTxnObj(); terr == nil {
			txObj.Exec("insert into florest_employee (name) values('abc')")
			cerr := txObj.Commit()
			if cerr != nil {
				fmt.Println(cerr)
			}
		} else {
			fmt.Println(terr)
		}
		// start and rollback one txn
		if txObj, terr := db.GetTxnObj(); terr == nil {
			txObj.Exec("select * from florest_employee")
			txObj.Rollback()
		} else {
			fmt.Println(terr)
		}
		// delete the created table
		_, err = db.Execute("drop table florest_employee")
		if err != nil {
			fmt.Println(err)
		}
		// close the connection
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	//Business Logic
	return io, nil
}
