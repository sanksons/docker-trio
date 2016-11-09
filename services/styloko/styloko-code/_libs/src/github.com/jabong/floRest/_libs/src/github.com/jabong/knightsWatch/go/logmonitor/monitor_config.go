package logmonitor

import ()

type Config struct {

	//AppName acts as a namespace which is prefixed with every custom metric
	AppName string

	// ApiKey is an optional parameter used for authenticating a write request in Datadog server.
	// This is needed only for the Http Apis ( which are not supported yet )
	ApiKey string

	// AppKey is an optional parameter used for authenticating a read request from Datadog server
	// along with ApiKey. This is needed only for the Http Apis ( which are not supported yet )
	AppKey string

	// ServerAddr is the montoring server ip and port
	ServerAddr string

	// Platform specifies monitoring platform that is being used. For now only
	// Datadog is supported in agent mode
	Platform string

	// Verbose option if set to true prints down some information for debugging purpose
	Verbose bool

	// Enabled determines whether to disable or enable sending metrics to a monitoring server
	Enabled bool

	// LogFileName holds the name of the log file along with the path
	LogFileName string
}
