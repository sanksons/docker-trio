package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

//FileLoggerImpl is a file logger structure
type FileLoggerImpl struct {
	//filename is the log file name with absolute path
	filename string

	//file is the file handle for log filename
	file *os.File
}

//Trace logs a trace
func (fileLogger *FileLoggerImpl) Trace(msg LogMsg) {
	if CanLog(TraceLevel) == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "trace"
	msg.StackTraces = getStackTrace()
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)
}

//Warning logs a warning
func (fileLogger *FileLoggerImpl) Warning(msg LogMsg) {
	if CanLog(WarningLevel) == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "warning"
	msg.StackTraces = getStackTrace()
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)
}

//Debug logs a debug
func (fileLogger *FileLoggerImpl) Debug(msg LogMsg) {
	if CanLog(DebugLevel) == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "debug"
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)
}

//Info logs a info
func (fileLogger *FileLoggerImpl) Info(msg LogMsg) {
	if CanLog(InfoLevel) == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "info"
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)
}

//Error logs an error
func (fileLogger *FileLoggerImpl) Error(msg LogMsg) {
	if CanLog(ErrLevel) == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "error"
	msg.StackTraces = getStackTrace()
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)

}

//Profile logs a profile specifying memory and time taken by an execution
//point
func (fileLogger *FileLoggerImpl) Profile(msg LogMsg) {
	if conf.ProfilerEnabled == false {
		return
	}
	fileLogger.SetOutput()
	msg.Level = "profile"
	msg.TimeStamp = time.Now().Local().Format(time.RFC3339)
	str := msg.GetFormattedLog()
	log.Println(str)
}

//SetOutput sets the file handle where to write the log
func (filelogger *FileLoggerImpl) SetOutput() {
	logName := getLogFileExt(filelogger)
	log.SetFlags(0)
	if _, err := os.Stat(logName); err == nil {
		log.SetOutput(filelogger.file)
		return
	}
	if filelogger.file != nil {
		filelogger.file.Close()
	}
	err := setLogFileHandle(filelogger)
	if err != nil {
		fmt.Printf("\nError In Settting log file handle %s: %+v\n", filelogger.filename, err)
		return
	}
	log.SetOutput(filelogger.file)
}

//newFileLogger returns a file logger with log file name fname, having configuration
//specified in conf and allowedLogLevel specifies the log level that are actually to
//be logged
func newFileLogger(fname string, conf FileLoggerConfig) (*FileLoggerImpl, error) {
	logger := new(FileLoggerImpl)
	logger.filename = fname

	err := setLogFileHandle(logger)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

//setLogFileHandle opens a file and assigns the file handle to file in fileLogger
func setLogFileHandle(fileLogger *FileLoggerImpl) error {
	logName := getLogFileExt(fileLogger)
	f, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\nError In Opening log file %s: %+v\n", fileLogger.filename, err)
		return err
	}
	fileLogger.file = f
	return nil
}

//getLogFileExt returns the log file ext name after appending
//today's date and ".log" to the filename
func getLogFileExt(fileLogger *FileLoggerImpl) string {
	t := time.Now().Local()
	tf := t.Format("2006-01-02")
	return fileLogger.filename + tf + ".log"
}
