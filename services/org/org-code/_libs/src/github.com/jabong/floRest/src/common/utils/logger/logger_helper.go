package logger

import (
	"fmt"
	"github.com/jabong/floRest/src/common/logger"
	"github.com/jabong/floRest/src/common/monitor"
	utilHttp "github.com/jabong/floRest/src/common/utils/http"
)

//formatString is a log format string
var formatString string = "%v"

//Initialise initialises a logger
func Initialise(confFile string) error {
	return logger.Initialise(confFile)
}

//setFormatString sets a format string that is to be used while logging. The
//format has to be same as the standard fmt
func setFormatString(fmtString string) {
	if formatString != "" {
		formatString = fmtString
	}
}

//DebugSpecific logs a debug to a specific log handle
func DebugSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.DebugLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Debug Message. " + err.Error())
		return
	}
	loggerHandle.Debug(Convert(a...))
	//logmonitorUtils.PushEvent("INFO:", formatLogMsg(a), nil, logmonitor.INFO)
}

//Debug logs debug to a default log handle
func Debug(a ...interface{}) {
	DebugSpecific(GetDefaultLoggerType(), a...)
}

//InfoSpecific logs a info to a specific log handle
func InfoSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.InfoLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Info Message. " + err.Error())
		return
	}
	loggerHandle.Info(Convert(a...))
	//logmonitorUtils.PushEvent("INFO:", formatLogMsg(a), nil, logmonitor.INFO)
}

//Info logs info to a default log handle
func Info(a ...interface{}) {
	InfoSpecific(GetDefaultLoggerType(), a...)
}

//TraceSpecific logs a trace to a specific log handle
func TraceSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.TraceLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Trace Message. " + err.Error())
		return
	}
	loggerHandle.Trace(Convert(a...))
}

//Trace logs Trace to default log handle
func Trace(a ...interface{}) {
	TraceSpecific(GetDefaultLoggerType(), a...)
}

//WarningSpecific logs a warning to a specific log handle
func WarningSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.WarningLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Warning Message : " + err.Error())
		return
	}
	loggerHandle.Warning(Convert(a...))
	monitor.GetInstance().Warning("WARNING:", formatLogMsg(a), nil)
}

//Warning logs warning to default log handle
func Warning(a ...interface{}) {
	WarningSpecific(GetDefaultLoggerType(), a...)
}

//ErrorSpecific logs an error to a specific logger handle
func ErrorSpecific(logType string, a ...interface{}) {
	if !logger.CanLog(logger.ErrLevel) {
		return
	}
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Error Message. " + err.Error())
		return
	}
	loggerHandle.Error(Convert(a...))
	monitor.GetInstance().Error("ERROR:", formatLogMsg(a), nil)
}

//Error logs error to default log handle
func Error(a ...interface{}) {
	ErrorSpecific(GetDefaultLoggerType(), a...)
}

//ProfileSpecific logs a profile specifying memory and time taken by an execution
//point
func ProfileSpecific(logType string, a ...interface{}) {
	loggerHandle, err := logger.GetLoggerHandle(logType)
	if err != nil {
		fmt.Println("Skipping Log Profile Message : " + err.Error())
		return
	}
	loggerHandle.Profile(Convert(a...))
}

//Profile logs a profile specifying memory and time taken by an execution
//point to a default log handle
func Profile(a ...interface{}) {
	ProfileSpecific(GetDefaultLoggerType(), a...)
}

//GetDefaultLoggerType gets the key of a default logger type
func GetDefaultLoggerType() string {
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
func formatLogMsg(a ...interface{}) string {
	return fmt.Sprintf(formatString, a...)
}

//func converts application log to LogMsg format suitable for Logger
func Convert(a ...interface{}) logger.LogMsg {
	paramLength := len(a)
	if paramLength == 0 {
		return logger.LogMsg{
			Message: "Empty log param",
		}
	}
	if paramLength == 1 {
		//Only Log Message string is passed
		return logger.LogMsg{
			Message: fmt.Sprintf("%s", a[0]),
		}
	}

	//First param is message string; Second param is request context
	vMsg, msgOk := a[0].(string)
	vRc, rcOk := a[1].(utilHttp.RequestContext)

	if !msgOk || !rcOk {

		return logger.LogMsg{
			Message: fmt.Sprintf("Erorr in parsing logging params for %v", a),
		}
	}
	return logger.LogMsg{
		Message:       vMsg,
		TransactionId: vRc.TransactionId,
		SessionId:     vRc.SessionId,
		RequestId:     vRc.RequestId,
		AppId:         vRc.AppName,
		UserId:        vRc.UserId,
		Uri:           vRc.URI,
	}
}
