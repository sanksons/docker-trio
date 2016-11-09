package mongodb

import (
	"gopkg.in/mgo.v2"
	"testing"
	//"gopkg.in/mgo.v2/bson"
)

func TestMongodb(t *testing.T) {
	var err error
	var dbObj MongodbInterface
	conf := new(Config)
	// fill invalid url
	conf.Url = "invalid"
	dbObj, err = Get(conf)
	if err == nil {
		FailNow(t, "invalid url must throw error")
	}
	// Test methods: valid url,db and collection
	conf.Url = "mongodb://localhost:27017"
	conf.DbName = "flashback"
	dbObj, err = Get(conf)
	// verify error is nil
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

// logs error and exits the test run
func FailNow(t *testing.T, err string) {
	t.Error(err)
	t.FailNow()
}
