package logger

import ()

//DummyLoggerImpl is a file logger structure
type DummyLoggerImpl struct {
}

//Trace logs a trace
func (d *DummyLoggerImpl) Trace(msg LogMsg) {

}

//Warning logs a warning
func (d *DummyLoggerImpl) Warning(msg LogMsg) {

}

//Info logs a debug
func (d *DummyLoggerImpl) Debug(msg LogMsg) {

}

//Info logs a info
func (d *DummyLoggerImpl) Info(msg LogMsg) {

}

//Error logs an error
func (d *DummyLoggerImpl) Error(msg LogMsg) {

}

//Profile logs a profile specifying memory and time taken by an execution
//point
func (d *DummyLoggerImpl) Profile(msg LogMsg) {

}

//SetOutput sets the file handle where to write the log
func (d *DummyLoggerImpl) SetOutput() {

}

//newFileLogger returns a file logger with log file name fname, having configuration
//specified in conf and allowedLogLevel specifies the log level that are actually to
//be logged
func newDummyLogger(conf DummyLoggerConfig) *DummyLoggerImpl {
	return new(DummyLoggerImpl)
}
