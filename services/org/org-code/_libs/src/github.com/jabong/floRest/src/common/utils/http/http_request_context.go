package http

import (
	"fmt"
)

/*
The Request Execution Context is used for tracking a particular request processing
*/
type RequestContext struct {
	AppName       string
	UserId        string
	SessionId     string
	RequestId     string
	TransactionId string
	URI           string
	ClientAppId   string
	TokenId       string
}

//Implements the Stringer interface
func (t RequestContext) String() string {
	format := "[AppName : %s, UserID : %s, SessionID : %s, RequestID : %s, TransactionID : %s, TokenId : %s, URI : %s, ClientAppId : %s]"
	return fmt.Sprintf(format,
		t.AppName,
		t.UserId,
		t.SessionId,
		t.RequestId,
		t.TransactionId,
		t.TokenId,
		t.URI,
		t.ClientAppId,
	)
}
