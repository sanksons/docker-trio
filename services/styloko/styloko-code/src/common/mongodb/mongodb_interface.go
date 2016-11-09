package mongodb

import ()

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
}
