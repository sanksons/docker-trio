package mongodb

import (
	"gopkg.in/mgo.v2"
)

// mongodbb interface
type MongodbInterface interface {
	// init initializes the mongodb instance
	init(*Config) *MongodbError
	// FindOne returns first matching item
	FindOne(string, map[string]interface{}) (interface{}, *MongodbError)
	// FindAll returns all matching items
	FindAll(string, map[string]interface{}) ([]interface{}, *MongodbError)
	// Insert add one item
	Insert(string, interface{}) *MongodbError
	// Update modify existing item
	Update(string, map[string]interface{}, interface{}) *MongodbError
	// Remove delete existing item
	Remove(string, map[string]interface{}) *MongodbError
	// Find And Modify existing item
	FindAndModify(string, map[string]interface{}, map[string]interface{}) (interface{}, *MongodbError)
	//Set Collection name
	SetCollection(string) *mgo.Collection
	//Generate Next Sequence Number
	GetNextSequence(string) int
	//Set Collection in Counter
	SetCollectionInCounter(string, int) *MongodbError
}
