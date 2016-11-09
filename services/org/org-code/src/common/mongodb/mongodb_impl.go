package mongodb

import (
	_ "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongodb driver
type MongoDriver struct {
	conn *mgo.Database
	sess *mgo.Session
}

var session *mgo.Session

// var MongoSession *MongoDriver

func GetInstance() *MongoDriver {
	tmp := new(MongoDriver)
	return tmp
}

type CounterInfo struct {
	Id  string `bson:"_id" mapstructure:"_id"`
	Seq int    `bson:"seqId" mapstructure:"seqId"`
}

type Query struct {
	Criteria interface{}
	Limit    *int
	Offset   *int
	Sort     []string
}

// init method
func (obj *MongoDriver) Initialize(conf *Config) (aerr *MongodbError) {
	var err error
	if session != nil {
		tmp := session.Copy()
		obj.conn = tmp.DB(conf.DbName)
		obj.sess = tmp
		return nil
	}
	session, err = mgo.Dial(conf.Url)
	if err == nil {
		obj.conn = session.DB(conf.DbName)
	} else {
		aerr = getErrObj(ERR_INITIALIZATION, err.Error()+"-connection url:"+conf.Url)
	}
	// return
	return aerr
}

func (driver *MongoDriver) Ping() bool {

	if driver.sess == nil {
		return false
	}
	err := driver.sess.Ping()
	if err != nil {
		return false
	}
	return true
}

func (driver *MongoDriver) GetConn() *mgo.Database {
	return driver.conn
}

func (obj *MongoDriver) Close() {
	obj.sess.Close()
}

func (obj *MongoDriver) Refresh() {
	obj.sess.Refresh()
}

func (obj *MongoDriver) CollectionExists() ([]string, error) {
	names, err := obj.conn.CollectionNames()
	return names, err
}

// find one
func (obj *MongoDriver) FindOne(collection string, query map[string]interface{}, ret interface{}) *MongodbError {
	err := obj.conn.C(collection).Find(bson.M(query)).One(&ret)
	if err != nil {
		return getErrObj(ERR_FINDONE_FAILURE, err.Error())
	} else {
		return nil
	}
}

// find all
func (obj *MongoDriver) FindAll(collection string, query Query, ret interface{}) *MongodbError {
	q := obj.conn.C(collection).Find(query.Criteria)
	if query.Sort != nil {
		q.Sort(query.Sort...)
	}
	if query.Limit != nil {
		q.Limit(*query.Limit)
	}
	if query.Offset != nil {
		q.Skip(*query.Offset)
	}
	err := q.All(ret)
	if err != nil {
		return getErrObj(ERR_FINDALL_FAILURE, err.Error())
	}
	return nil
}

// insert
func (obj *MongoDriver) Insert(collection string, value interface{}) *MongodbError {
	err := obj.conn.C(collection).Insert(value)
	if err != nil {
		return getErrObj(ERR_INSERT_FAILURE, err.Error())
	} else {
		return nil
	}
}

// update
func (obj *MongoDriver) Update(collection string, query map[string]interface{}, value interface{}) *MongodbError {
	err := obj.conn.C(collection).Update(bson.M(query), value)
	if err != nil {
		return getErrObj(ERR_UPDATE_FAILURE, err.Error())
	} else {
		return nil
	}
}

// delete
func (obj *MongoDriver) Remove(collection string, query map[string]interface{}) *MongodbError {
	err := obj.conn.C(collection).Remove(bson.M(query))
	if err != nil {
		return getErrObj(ERR_REMOVE_FAILURE, err.Error())
	} else {
		return nil
	}
}

//find and modify
func (obj *MongoDriver) FindAndModify(collectionName string, updatecriteria map[string]interface{}, findcriteria map[string]interface{}, upsertVal bool, deleteVal bool, returnNewVal bool, ret interface{}) *MongodbError {
	change := mgo.Change{Update: bson.M(updatecriteria), Upsert: upsertVal, ReturnNew: returnNewVal, Remove: deleteVal}
	_, err := obj.conn.C(collectionName).Find(bson.M(findcriteria)).Apply(change, &ret)
	if err != nil {
		return getErrObj(ERR_FINDMODIFY_FAILURE, err.Error())
	} else {
		return nil
	}
}

//set collection
func (obj *MongoDriver) SetCollection(collectionName string) *mgo.Collection {
	return obj.conn.C(collectionName)
}

func (obj *MongoDriver) GetNextSequence(collectionName string) int {
	retInfo := new(CounterInfo)
	change := mgo.Change{Update: bson.M{"$inc": bson.M{"seqId": 1}}, ReturnNew: true}
	_, err := obj.conn.C("counters").Find(bson.M{"_id": collectionName}).Apply(change, &retInfo)
	if err != nil {
		return 0
	}
	return retInfo.Seq
}
func (obj *MongoDriver) SetCollectionInCounter(collectionName string, seq int) *MongodbError {
	//create row in counter table
	counter := new(CounterInfo)
	counter.Id = collectionName
	counter.Seq = seq
	err := obj.Insert("counters", counter)
	if err != nil {
		return getErrObj(ERR_SETCOLINCOUNTER_FAILURE, err.DeveloperMessage)
	}
	return nil
}

func (obj *MongoDriver) DeleteDatabase() error {
	return obj.conn.DropDatabase()
}
