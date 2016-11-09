package mongodb

import ()

// config from mongo db
type Config struct {
	Url    string // e.g. mongodb://localhost:27017
	DbName string
}
