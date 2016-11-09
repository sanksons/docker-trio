package sqldb

import (
	"database/sql"
)

type SqlDbInterface interface {
	// init initialize db instance
	init(conf *Config) *SqlDbError
	// Query should be used for select purpose
	Query(string, ...interface{}) (*sql.Rows, *SqlDbError)
	// Execute should be used for data changes
	Execute(string, ...interface{}) (sql.Result, *SqlDbError)
	// Ping checks the connection
	Ping() *SqlDbError
	// Close close connection properly
	Close() *SqlDbError
	// GetTxnObj get transaction object
	GetTxnObj() (*sql.Tx, *SqlDbError)
}
