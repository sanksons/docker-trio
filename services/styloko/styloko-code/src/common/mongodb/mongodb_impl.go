package mongodb

import (
	"os"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDriver is a the basic mongo struct returned to the user
type MongoDriver struct {
	conn *mgo.Database
	sess *mgo.Session
}

var UseSafeMode bool

var session *mgo.Session

// GetConn return connection object
func (driver *MongoDriver) GetConn() *mgo.Database {
	return driver.conn
}

// Ping return bool by pinging Mongo Servers
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

// GetInstance returns a new instance MongoDriver struct
func GetInstance() *MongoDriver {
	tmp := new(MongoDriver)
	return tmp
}

// CounterInfo stores info about counters
type CounterInfo struct {
	Id    string `bson:"_id"`
	SeqId int    `bson:"seqId"`
}

type Seq struct {
	SeqId int `bson:"seqId"`
}

// Initialize generates the connection bindings for mongo, calls the dial method.
func (driver *MongoDriver) Initialize(conf *Config, dbName string) (aerr *MongodbError) {
	var err error
	if session != nil {
		tmp := session.Copy()
		if dbName != "" {
			driver.conn = tmp.DB(dbName)
		} else {
			driver.conn = tmp.DB(conf.DbName)
		}
		driver.sess = tmp
		return nil
	}
	session, err = mgo.Dial(conf.Url)
	if err == nil {
		if dbName != "" {
			driver.conn = session.DB(dbName)
		} else {
			driver.conn = session.DB(conf.DbName)
		}
	} else {
		aerr = getErrObj(ERR_INITIALIZATION, err.Error()+"-connection url:"+conf.Url)
	}
	//check if we need to use safe mode.
	env := os.Getenv("ENVIRON_NAME")
	if (strings.ToLower(env) == "sc-perf") || (strings.ToLower(env) == "sc-prod") {
		UseSafeMode = true
	}
	return aerr
}

//Set Safe session to ensure write consistency.
func (driver *MongoDriver) SetSafe() {
	driver.sess.SetSafe(&mgo.Safe{WMode: "majority"})
}

//Ensure safe
func (driver *MongoDriver) EnsureSafe() {
	driver.sess.EnsureSafe(&mgo.Safe{W: 2, FSync: true})
}

// Close shuts down the current session.
func (driver *MongoDriver) Close() {
	driver.sess.Close()
	// Done to release sockets.
	driver.sess.Refresh()
}

// Refresh retries connection if lost and returns unused sockets back to the pool
func (driver *MongoDriver) Refresh() {
	driver.sess.Refresh()
}

// FindOne return one document
func (driver *MongoDriver) FindOne(collection string, query map[string]interface{}, ret interface{}) error {
	err := driver.conn.C(collection).Find(query).One(ret)
	if err != nil && err != mgo.ErrNotFound {
		return getErrObj(ERR_FINDONE_FAILURE, err.Error())
	}
	if err == mgo.ErrNotFound {
		return ErrNotFound
	}
	return nil
}

// Query struct
type Query struct {
	Criteria interface{}
	Limit    *int
	Offset   *int
	Sort     []string
}

// FindAll returns multiple document in an []interface{}
func (driver *MongoDriver) FindAll(collection string, query Query, ret interface{}) error {
	q := driver.conn.C(collection).Find(query.Criteria)
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

// Insert creates a new document
func (driver *MongoDriver) Insert(collection string, value interface{}) *MongodbError {
	//value.seqId := driver.GetNextSequence(collection)
	err := driver.conn.C(collection).Insert(value)
	if err != nil {
		return getErrObj(ERR_INSERT_FAILURE, err.Error())
	}
	return nil
}

// Update modifies a previous document. Query is required.
func (driver *MongoDriver) Update(collection string, query map[string]interface{}, value interface{}) *MongodbError {
	err := driver.conn.C(collection).Update(bson.M(query), value)
	if err != nil {
		return getErrObj(ERR_UPDATE_FAILURE, err.Error())
	}
	return nil
}

// Remove removes document(s) based on query
func (driver *MongoDriver) Remove(collection string, query map[string]interface{}) *MongodbError {
	err := driver.conn.C(collection).Remove(bson.M(query))
	if err != nil {
		return getErrObj(ERR_REMOVE_FAILURE, err.Error())
	}
	return nil

}

//FindAndModify queries and updates one document with given values
func (driver *MongoDriver) FindAndModify(collectionName string, updatecriteria map[string]interface{}, findcriteria map[string]interface{}, upsertFlag bool) (ret interface{}, err error) {
	change := mgo.Change{Update: bson.M(updatecriteria), ReturnNew: true, Upsert: upsertFlag}
	info, err := driver.conn.C(collectionName).Find(bson.M(findcriteria)).Apply(change, &ret)
	return info, err
}

//SetCollection sets the collection name for the session
func (driver *MongoDriver) SetCollection(collectionName string) *mgo.Collection {
	return driver.conn.C(collectionName)
}

// GetNextSequence returns the next sequence ID for MySQL and seqID for Mongo
func (driver *MongoDriver) GetNextSequence(collectionName string) int {
	retInfo := new(CounterInfo)
	change := mgo.Change{Update: bson.M{"$inc": bson.M{"seqId": 1}}, ReturnNew: true}
	_, err := driver.conn.C("counters").Find(bson.M{"_id": collectionName}).Apply(change, &retInfo)
	if err != nil {
		return 0
	}
	return retInfo.SeqId
}

// SetCollectionInCounter sets collection in counter
func (driver *MongoDriver) SetCollectionInCounter(collectionName string, seqId int) *MongodbError {
	//create row in counter table
	var seq Seq
	if seqId > 0 {
		seq.SeqId = seqId
	} else {
		err := driver.conn.C(collectionName).Find(nil).Sort("-seqId").Limit(1).One(&seq)
		if err != nil {
			return getErrObj(ERR_FINDONE_FAILURE, err.Error())
		}
	}
	if seq.SeqId > 0 {
		counter := new(CounterInfo)
		change := mgo.Change{Update: bson.M{"_id": collectionName, "seqId": seq.SeqId}, Upsert: true}
		_, err := driver.conn.C("counters").Find(bson.M{"_id": collectionName}).Apply(change, &counter)
		if err != nil {
			return getErrObj(ERR_INSERT_FAILURE, err.Error())
		}
	}
	return nil
}

func (obj *MongoDriver) CollectionExists() ([]string, error) {
	names, err := obj.conn.CollectionNames()
	return names, err
}
