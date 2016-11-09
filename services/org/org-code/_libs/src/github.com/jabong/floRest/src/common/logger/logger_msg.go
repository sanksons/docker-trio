package logger

import (
	"encoding/json"
	"fmt"
)

//formatString is format of the log in string formattype configuration
var formatString string = "[level : %s, message : %s, tId : %s, reqId : %s, appId : %s, sessionId : %s, userId : %s, stackTraces : %s, timestamp : %s, uri : %s]"

/*
LogMsg struct is used for bundling message/requestcontext attributes
for dumping into log
*/
type LogMsg struct {
	Level         string   `json:"level"`
	Message       string   `json:"message"`
	TransactionId string   `json:"tId,omitempty"`
	RequestId     string   `json:"reqId,omitempty"`
	AppId         string   `json:"appId,omitempty"`
	SessionId     string   `json:"sessionId,omitempty"`
	UserId        string   `json:"userId,omitempty"`
	StackTraces   []string `json:"stackTraces,omitempty"`
	TimeStamp     string   `json:"timestamp"`
	Uri           string   `json:"uri,omitempty"`
}

//Get log to dump from LogMsg struct
func (msg *LogMsg) GetFormattedLog() string {
	var str string
	if logFormatter.IsJson {
		str = msg.getJsonLog()
	} else {
		str = msg.getStringLog()
	}
	return str
}

//Get Json string from LogMsg struct
func (msg *LogMsg) getJsonLog() string {
	jMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("\nError In converting to json %+v\n", msg)
		return msg.getStringLog()
	}
	return string(jMsg)
}

//formatLogMsg specifies the log as specified in the format string
func (msg *LogMsg) getStringLog() string {
	return fmt.Sprintf(formatString, msg.Level,
		msg.Message,
		msg.TransactionId,
		msg.RequestId,
		msg.AppId,
		msg.SessionId,
		msg.UserId,
		msg.StackTraces,
		msg.TimeStamp,
		msg.Uri)
}
