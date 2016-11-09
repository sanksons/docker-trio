package examples

import (
	"fmt"
	"github.com/jabong/floRest/src/common/mongodb"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
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
	return "HelloWord"
}

// an example for mongo document
type employeeInfo struct {
	Id   string
	Type string
}

func (a HelloWorld) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	// fill the config
	conf := new(mongodb.Config)
	conf.Url = "mongodb://localhost:27017"
	conf.DbName = "florest"

	collection := "employee"
	// get db object
	db, err := mongodb.Get(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		var query map[string]interface{}
		// insert
		err = db.Insert(collection, &employeeInfo{Id: "123", Type: "Manager"})
		if err != nil {
			fmt.Println(err)
		}
		// update 
		query = make(map[string]interface{}, 1)
		query["id"] = "123"
		err = db.Update(collection, query, &employeeInfo{Id: "123", Type: "Director"})
		if err != nil {
			fmt.Println(err)
		}

		// find one
		_, err = db.FindOne(collection, query)
		if err != nil {
			fmt.Println(err)
		}

		// find all
		_, err = db.FindAll(collection, query)
		if err != nil {
			fmt.Println(err)
		}

		// remove
		err = db.Remove(collection, query)
		if err != nil {
			fmt.Println(err)
		}
	}
	//Business Logic
	return io, nil
}
