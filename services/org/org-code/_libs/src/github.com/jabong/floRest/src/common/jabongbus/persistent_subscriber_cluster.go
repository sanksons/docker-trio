package jabongbus

import (
	"encoding/json"
	"errors"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
	"time"
)

// persistent subscriber with redis cluster
type persistSubCluster struct {
	clusterClient  *redis.ClusterClient
	id             int64
	messageList    string
	bufList        string
	exceptionList  string
	statusList     string
	getReq         *redis.StringCmd
	ackReq         *redis.StringCmd
	ackPending     bool
	usedNoAck      bool
	processMessage ProcessMessage
	exit           bool
	msgRetryCount  int
	retryCount     int
	statusQuit     chan bool
}

// initialize with config values
func (ps *persistSubCluster) init(conf *Subscriberconfig) error {
	// set the client
	ps.clusterClient = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:    strings.Split(conf.RedisCon, ","),
			PoolSize: MAX_ACTIVE_REDIS_CON,
		})
	ps.messageList = conf.Publisher + "-" + conf.RoutingKey
	if err := ps.setId(); err != nil { // failed in setting sub id
		return err
	} // set the id
	ps.bufList = "{" + ps.messageList + "}" + SUBSCRIBER + strconv.FormatInt(ps.id, 10) //e.g "catalog-rkey1-sub4"
	ps.exceptionList = "{" + ps.messageList + "}" + EXCEPTION
	ps.statusList = "{" + ps.messageList + "}" + SUBSTATUS
	ps.statusQuit = make(chan bool)
	ps.setAliveStatus()
	// return now
	return nil
}

// SetProcessMsg set client process message implementation
func (ps *persistSubCluster) SetProcessMsg(pm ProcessMessage) {
	ps.processMessage = pm
}

// set the subscriber id
func (ps *persistSubCluster) setId() error {
	// increment the subcount in redis. It will also create subcount key if not present
	req := ps.clusterClient.Incr("{" + ps.messageList + "}" + SUBSCRIBER_COUNT)
	if req.Err() != nil {
		logger.Error(req.Err().Error())
		return ErrGettingSubId
	} else {
		// set id
		ps.id = req.Val()
	}
	// return
	return nil
}

// Get get the message for client
func (ps *persistSubCluster) Get(timeout int) {
	var msg string
	var err error
	var retMsg *Message
	for {
		if ps.exit { // call client function to process message
			break // client wants to exit
		}
		msg = ""
		if ps.usedNoAck {
			ps.usedNoAck = false // mark it false, last message will be redelievered again
			msg = ps.getReq.Val()
		} else if ps.ackPending { // ack is pending send error
			err = errors.New(MSG_NOT_ACKED)
		} else { // get new message
			ps.getReq = ps.clusterClient.BRPopLPush(ps.messageList, ps.bufList, time.Duration(timeout)*time.Second)
			if ps.getReq.Err() == nil {
				ps.ackPending = true
				msg = ps.getReq.Val()
				ps.retryCount = DEFAULT_RETRY_COUNT
			} else {
				logger.Error(ps.getReq.Err().Error())
				err = ErrGettingPersMsg
			}
		}
		// check if any error
		if err == nil {
			retMsg, err = stringToMessage(msg)
			ps.msgRetryCount = retMsg.RetryCount // set msg retry count
			ps.processMessage.Process(retMsg, err)
		} else {
			ps.processMessage.Process(nil, err)
		}
	}
}

// Ack ack by client
func (ps *persistSubCluster) Ack() error {
	// ack the message
	ps.ackLastMsg()
	if ps.ackReq.Err() == nil {
		ps.ackPending = false
		ps.usedNoAck = false
	} else {
		logger.Error(ps.ackReq.Err().Error())
	}
	// return error
	return ps.ackReq.Err()
}

// ack the last message
func (ps *persistSubCluster) ackLastMsg() {
	ps.ackReq = ps.clusterClient.RPop(ps.bufList)
}

// NoAck no ack for last message, this message will be redelivered again
func (ps *persistSubCluster) NoAck() (err error) {
	// check if full retry reached
	if ps.retryCount < ps.msgRetryCount {
		ps.retryCount++
		ps.ackPending = false
		ps.usedNoAck = true
	} else { // max retry reached, move the message from buffer to exception list
		req := ps.clusterClient.RPopLPush(ps.bufList, ps.exceptionList)
		if req.Err() == nil {
			ps.ackPending = false
			ps.usedNoAck = false
		} else { // error in moving to exception list
			logger.Error(req.Err().Error())
			err = ErrNoAck
		}
	}
	// return now
	return err
}

// set the alive status
func (ps *persistSubCluster) setAliveStatus() {
	go func() {
		firstTime := true
		data := new(statusData)
		var lastStatus string
		for {
			data.Id = ps.id
			data.Time = time.Now().Unix() // time in seconds
			bufAry, err := json.Marshal(data)
			if err != nil {
				logger.Error(err.Error())
			} else {
				if firstTime {
					firstTime = false
				} else { // remove previous value
					ps.clusterClient.LRem(ps.statusList, 0, lastStatus)
				}
				// add alive status
				if req := ps.clusterClient.RPush(ps.statusList, string(bufAry)); req.Err() != nil {
					logger.Error(req.Err().Error())
				} else {
					lastStatus = string(bufAry) // save the last status, needs to be deleted
				}
				time.Sleep(time.Duration(SUB_STATUS_DURATION) * time.Second)
				if _, ok := <-ps.statusQuit; !ok { // subscriber has stopped
					return
				}
			}
		}
	}()
}

// StopSub stop the subscriber
func (ps *persistSubCluster) StopSub() {
	ps.exit = true
	ps.clusterClient.Close()
	close(ps.statusQuit)
}
