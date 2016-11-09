package logmonitor

import ()

// LogData is the log data which will be set in the logmonitor library. Application code doesn't need
// to set this explicitly
type LogData struct {

	// Level denotes the log-level. It can either be one of the following - ERROR,  INFO, WARN,
	// DEBUG, TRACE
	Level string

	// Time denotes the time the log is made in format yyyy-MM-dd HH:mm:ss,SSS TZ
	Time string

	// AppName denotes the application name
	AppName string

	// StackTrace is a list of the method calls that the application was in the middle of when the log was made
	StackTrace string
}

// AppLogData is the logging/monitoring data that application has to sent
type AppLogData struct {
	LogData

	// UserId
	UserId string

	// SessionId
	SessionId string

	// TransactionId
	TId string

	// RequestId
	ReqId string

	// Title denotes the title of a log message and event
	Title string

	// Body denotes the body of a log message and event
	Body string

	// Tags adds dimension to a event or log. This should be of the form <key:value> e.g env:prod
	Tags map[string]string

	// SentEvent controls whether the data has to be sent as event to some monitoring platform like datadog
	ToSendEvent bool

	// ToLog controls whether the data has to be logged in file as well
	ToLog bool

	//ToAsyncLog determines whether logging is to me made async or not
	ToAsyncLog bool
}
