package mongodb

import ()

// Get() - Creates, initializes and returns the mongo instance based on the given config
func Get(conf *Config) (ret MongodbInterface, err *MongodbError) {
	ret = new(mongoDriver)
	err = ret.init(conf)
	// return
	return ret, err
}
