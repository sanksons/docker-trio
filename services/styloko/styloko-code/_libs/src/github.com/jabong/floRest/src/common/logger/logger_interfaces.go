package logger

type LogInterface interface {
	Trace(msg LogMsg)
	Warning(msg LogMsg)
	Info(msg LogMsg)
	Error(msg LogMsg)
	Profile(msg LogMsg)
	Debug(msg LogMsg)
}
