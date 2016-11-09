package jabongbus

import (
	"errors"
)

// Json errors
var (
	// errJsonMarshal means error in marshaling structure to json string
	ErrJsonMarshal = errors.New("error in json marshal")

	// errJsonUnMarshal means error in unmarshaling jstring to json structure
	ErrJsonUnMarshal = errors.New("error in json unmarshal")
)

// Publisher errors
var (
	ErrHttpPost = errors.New("error in http post call to bus")
)

// Subscriber errors
var (
	ErrGettingSubId      = errors.New("error in getting subscriber id")
	ErrNetworkConnection = errors.New("error in network connection")
	ErrGettingNonPersMsg = errors.New("error in getting non-persistent message")
	ErrGettingPersMsg    = errors.New("error in getting persistent message")
	ErrNoAck             = errors.New("error in no ack for message")
	ErrNoDeadQForNonPers = errors.New("deadQ subscriber not supported for non-persistent")
)
