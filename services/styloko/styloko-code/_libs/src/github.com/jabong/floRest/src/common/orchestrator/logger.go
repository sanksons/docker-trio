package orchestrator

import (
	"fmt"
	"github.com/jabong/floRest/src/common/logger"
)

//formatString is a log format string
var formatString string = "%v"

//InitLogger initialises a logger
func InitLogger(confFile string) error {
	return logger.Initialise(confFile)
}

//setFormatString sets a format string that is to be used while logging. The
//format has to be same as the standard fmt
func setFormatString(fmtString string) {
	if formatString != "" {
		formatString = fmtString
	}
}

//logInfoSpecific logs a info to a specific log handle
func logInfoSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.InfoLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Info Message. " + err.Error())
		return
	}
	loggerHandle.Info(formatLogMsg(a))
}

//logInfo logs info to a default log handle
func logInfo(a ...interface{}) {
	logInfoSpecific(getDefaultLoggerType(), a)
}

//logTraceSpecific logs a trace to a specific log handle
func logTraceSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.TraceLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Trace Message. " + err.Error())
		return
	}
	loggerHandle.Trace(formatLogMsg(a))
}

//logTrace logs Trace to default log handle
func logTrace(a ...interface{}) {
	logTraceSpecific(getDefaultLoggerType(), a)
}

//logWarningSpecific logs a warning to a specific log handle
func logWarningSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.WarningLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Warning Message : " + err.Error())
		return
	}
	loggerHandle.Warning(formatLogMsg(a))
}

//logWarning logs warning to default log handle
func logWarning(a ...interface{}) {
	logWarningSpecific(getDefaultLoggerType(), a)
}

//logErrorSpecific logs an error to a specific logger handle
func logErrorSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.ErrLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Error Message. " + err.Error())
		return
	}
	loggerHandle.Error(formatLogMsg(a))
}

//logError logs error to default log handle
func logError(a ...interface{}) {
	logErrorSpecific(getDefaultLoggerType(), a)
}

//logProfileSpecific logs a profile specifying memory and time taken by an execution
//point
func logProfileSpecific(logType string, a ...interface{}) {
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Profile Message : " + err.Error())
		return
	}
	loggerHandle.Profile(formatLogMsg(a))
}

//logProfile logs a profile specifying memory and time taken by an execution
//point to a default log handle
func logProfile(a ...interface{}) {
	logProfileSpecific(getDefaultLoggerType(), a)
}

//getDefaultLoggerType gets the key of a default logger type
func getDefaultLoggerType() string {
	return logger.GetDefaultLogTypeKey()
}

//getLoggerHandle gets the handle for a particular logger type
func getLoggerHandle(logType string) (logger.LogInterface, error) {
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println(err.Error())
		loggerHandle, err = logger.GetLoggerHandle(logger.GetDefaultLogTypeKey())
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return loggerHandle, nil
}

//formatLogMsg specifies the log as specified in the format string
func formatLogMsg(a ...interface{}) logger.LogMsg {
	return logger.LogMsg{
		Message: fmt.Sprintf(formatString, a...),
	}
}
