package jabongbus

import (
	"encoding/json"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// Message - structure to hold information
type Message struct {
	PublisherName string      `json:"publisher_name"`
	RoutingKey    string      `json:"routing_key"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	RetryCount    int         `json:"retry_count"`
}

// mapToMessage maps the string to message structure
func stringToMessage(str string) (*Message, error) {
	ret := new(Message)
	if err := json.Unmarshal([]byte(str), ret); err == nil {
		return ret, nil
	} else {
		logger.Error(err.Error() + ",input string:" + str)
		return nil, ErrJsonUnMarshal
	}
}

// messageToString maps the message structure to json string
func messageToString(msg *Message) (string, error) {
	ary, jerr := json.Marshal(msg)
	if jerr != nil {
		return "", ErrJsonMarshal
	} else {
		return string(ary), nil
	}
}
