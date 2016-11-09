package logger

//FileLoggerConfig specifies the configs for file logger
type FileLoggerConfig struct {
	//FileNamePrefix specifies the prefix for a log file name
	FileNamePrefix string `json:"FileNamePrefix"`

	//Key specifies a unique name for a paricular logger implementation
	Key string `json:"Key"`

	//Path specifies the file logger path
	Path string `json:"Path"`
}

//DummyLoggerConfig specifies the configs for file logger
type DummyLoggerConfig struct {
	//Key specifies a unique name for a paricular logger implementation
	Key string `json:"Key"`
}

//Config specifies the logger config
type Config struct {

	//FileLogger stores all the configs for file loggers
	FileLogger []FileLoggerConfig `json:"FileLogger"`

	//DummyLogger stores all the configs for file loggers
	DummyLogger []DummyLoggerConfig `json:"DummyLogger"`

	//ProfilerEnabled specifies if a profiler is enabled or not
	ProfilerEnabled bool `json:"ProfilerEnabled"`

	//LogLevel specifies the minimum log level to write (1 - Info, 2 - Trace,
	//3 - Warning, 4 - Error)
	LogLevel int `json:"LogLevel"`

	//DefaultLogType specifies the default logger type. This should match
	//one of keys specified in logger configs
	DefaultLogType string `json:"DefaultLogType"`

	//AppName is the application name for which logging is done. AppName is
	//prefixed with each log file name
	AppName string `json:"AppName"`

	//FormatType is the format in which logging is done.Currently It can be string or JSON.
	FormatType string `json:"FormatType"`
}
