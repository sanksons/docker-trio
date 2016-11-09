package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongodb driver
type mongoDriver struct {
	conn    *mgo.Database
	session *mgo.Session
}

// init method
func (obj *mongoDriver) init(conf *Config) (aerr *MongodbError) {
	var err error
	var tmp *mgo.Session
	// set the connection
	if tmp, err = mgo.Dial(conf.Url); err == nil {
		obj.session = tmp
		obj.conn = tmp.DB(conf.DbName)
	}
	// set error if needed
	if err != nil {
		aerr = getErrObj(ERR_INITIALIZATION, err.Error()+"-connection url:"+conf.Url)
	}
	// return
	return aerr

}

// find one
func (obj *mongoDriver) FindOne(collection string, query map[string]interface{}) (ret interface{}, aerr *MongodbError) {
	obj.session.Refresh()
	err := obj.conn.C(collection).Find(bson.M(query)).One(&ret)
	if err != nil {
		return nil, getErrObj(ERR_FINDONE_FAILURE, err.Error())
	} else {
		return ret, nil
	}
}

// find all
func (obj *mongoDriver) FindAll(collection string, query map[string]interface{}) (ret []interface{}, aerr *MongodbError) {
	obj.session.Refresh()
	err := obj.conn.C(collection).Find(bson.M(query)).All(&ret)
	if err != nil {
		return nil, getErrObj(ERR_FINDALL_FAILURE, err.Error())
	} else {
		return ret, nil
	}
}

// insert
func (obj *mongoDriver) Insert(collection string, value interface{}) *MongodbError {
	obj.session.Refresh()
	err := obj.conn.C(collection).Insert(value)
	if err != nil {
		return getErrObj(ERR_INSERT_FAILURE, err.Error())
	} else {
		return nil
	}
}

// update
func (obj *mongoDriver) Update(collection string, query map[string]interface{}, value interface{}) *MongodbError {
	obj.session.Refresh()
	err := obj.conn.C(collection).Update(bson.M(query), value)
	if err != nil {
		return getErrObj(ERR_UPDATE_FAILURE, err.Error())
	} else {
		return nil
	}
}

// delete
func (obj *mongoDriver) Remove(collection string, query map[string]interface{}) *MongodbError {
	obj.session.Refresh()
	err := obj.conn.C(collection).Remove(bson.M(query))
	if err != nil {
		return getErrObj(ERR_REMOVE_FAILURE, err.Error())
	} else {
		return nil
	}
}
