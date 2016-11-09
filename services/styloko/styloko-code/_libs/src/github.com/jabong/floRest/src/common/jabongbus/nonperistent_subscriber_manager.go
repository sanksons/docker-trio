package jabongbus

import (
	"gopkg.in/redis.v3"
	"io"
	"net"
	"strings"
	"time"
)

// nonPersistSubscriber non persistent subscriber structure
type nonPersistSubscriber struct {
	client         *redis.Client
	messageList    string
	processMessage ProcessMessage
	exit           bool
}

// GetNonPersistentSubscriber get non persistent subscriber implementation
func GetNonPersistentSubscriber(conf *Subscriberconfig) NonPersistentSubscriber {
	obj := new(nonPersistSubscriber)
	obj.init(conf)
	return obj
}

// initialize with config values
func (nps *nonPersistSubscriber) init(conf *Subscriberconfig) {
	nps.client = redis.NewClient(
		&redis.Options{
			Network: "tcp",
			Addr:    conf.RedisCon})

	nps.messageList = conf.Publisher + "-" + conf.RoutingKey
}

// Get get the message for client
func (nps *nonPersistSubscriber) Get(timeout int) {
	var retMsg *Message
	var err error
	var sub *redis.PubSub

	for {
		// check if client wants to exit
		if nps.exit {
			break
		}
		// check if sub is ok
		if sub == nil {
			if sub, err = nps.client.Subscribe(nps.messageList); err != nil {
				sub = nil
				nps.processMessage.Process(nil, ErrNetworkConnection)
			}
		} else {
			if msgi, merr := sub.ReceiveTimeout(time.Duration(timeout) * time.Second); merr != nil {
				if !isNetworkError(merr) {
					nps.processMessage.Process(nil, ErrGettingNonPersMsg)
				} else { // network error
					sub = nil // set sub as nil
				}
			} else { // no error read the message
				if msg, ok := msgi.(*redis.Message); ok { // check message format
					retMsg, err = stringToMessage(msg.Payload)
					nps.processMessage.Process(retMsg, err)
				}
			}
		}
	}
}

// SetProcessMsg set client process message implemenation
func (nps *nonPersistSubscriber) SetProcessMsg(pm ProcessMessage) {
	nps.processMessage = pm
}

// close the client
func (nps *nonPersistSubscriber) closeClient() {
	nps.client.Close()
}

// stop the subscriber
func (nps *nonPersistSubscriber) StopSub() {
	nps.exit = true
	nps.closeClient()
}

// is network error
func isNetworkError(err error) bool {
	if err == io.EOF {
		return true
	} else if strings.Contains(err.Error(), "i/o timeout") {
		return false
	}
	_, ok := err.(net.Error)
	return ok
}
