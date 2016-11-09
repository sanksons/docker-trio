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

type deadQSubCluster struct {
	clusterClient  *redis.ClusterClient
	id             int64
	messageList    string
	bufList        string
	exceptionList  string
	statusList     string
	getReq         *redis.StringCmd
	ackReq         *redis.StringCmd
	ackPending     bool
	processMessage ProcessMessage
	exit           bool
	statusQuit     chan bool
}

// initialize with config values
func (ds *deadQSubCluster) init(conf *Subscriberconfig) error {
	// set the client
	ds.clusterClient = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:    strings.Split(conf.RedisCon, ","),
			PoolSize: MAX_ACTIVE_REDIS_CON,
		})
	ds.messageList = conf.Publisher + "-" + conf.RoutingKey
	if err := ds.setId(); err != nil { // failed in setting sub id
		return err
	} // set the id
	ds.bufList = "{" + ds.messageList + "}" + DEADQ_SUBSCRIBER + strconv.FormatInt(ds.id, 10) //e.g "catalog-rkey1-sub4"
	ds.exceptionList = "{" + ds.messageList + "}" + EXCEPTION
	ds.statusList = "{" + ds.messageList + "}" + DEADQ_SUBSTATUS
	ds.statusQuit = make(chan bool)
	ds.setAliveStatus()
	// return now
	return nil
}

// SetProcessMsg set client process message implementation
func (ds *deadQSubCluster) SetProcessMsg(pm ProcessMessage) {
	ds.processMessage = pm
}

// set the subscriber id
func (ds *deadQSubCluster) setId() error {
	// increment the subcount in redis. It will also create subcount key if not present
	req := ds.clusterClient.Incr("{" + ds.messageList + "}" + DEADQ_SUBSCRIBER_COUNT)
	if req.Err() != nil {
		logger.Error(req.Err().Error())
		return ErrGettingSubId
	} else {
		// set id
		ds.id = req.Val()
	}
	// return
	return nil
}

// Get get the message for client
func (ds *deadQSubCluster) Get(timeout int) {
	var msg string
	var err error
	var retMsg *Message
	for {
		if ds.exit { // call client function to process message
			break // client wants to exit
		}
		msg = ""
		if ds.ackPending { // ack is pending send error
			err = errors.New(MSG_NOT_ACKED)
		} else { // get new message
			ds.getReq = ds.clusterClient.BRPopLPush(ds.exceptionList, ds.bufList, time.Duration(timeout)*time.Second)
			if ds.getReq.Err() == nil {
				ds.ackPending = true
				msg = ds.getReq.Val()
			} else {
				logger.Error(ds.getReq.Err().Error())
				err = ErrGettingPersMsg
			}
		}
		// check if any error
		if err == nil {
			retMsg, err = stringToMessage(msg)
			ds.processMessage.Process(retMsg, err)
		} else {
			ds.processMessage.Process(nil, err)
		}
	}
}

// Ack ack by client
func (ds *deadQSubCluster) Ack() error {
	// ack the message
	ds.ackLastMsg()
	if ds.ackReq.Err() == nil {
		ds.ackPending = false
	} else {
		logger.Error(ds.ackReq.Err().Error())
	}
	// return error
	return ds.ackReq.Err()
}

// ack the last message
func (ds *deadQSubCluster) ackLastMsg() {
	ds.ackReq = ds.clusterClient.RPop(ds.bufList)
}

// set the alive status
func (ds *deadQSubCluster) setAliveStatus() {
	go func() {
		firstTime := true
		data := new(statusData)
		var lastStatus string
		for {
			data.Id = ds.id
			data.Time = time.Now().Unix() // time in seconds
			bufAry, err := json.Marshal(data)
			if err != nil {
				logger.Error(err.Error())
			} else {
				if firstTime {
					firstTime = false
				} else { // remove previous value
					ds.clusterClient.LRem(ds.statusList, 0, lastStatus)
				}
				// add alive status
				if req := ds.clusterClient.RPush(ds.statusList, string(bufAry)); req.Err() != nil {
					logger.Error(req.Err().Error())
				} else {
					lastStatus = string(bufAry) // save the last status, needs to be deleted
				}
				time.Sleep(time.Duration(SUB_STATUS_DURATION) * time.Second)
				if _, ok := <-ds.statusQuit; !ok { // subscriber has stopped
					return
				}
			}
		}
	}()
}

// StopSub stop the subscriber
func (ds *deadQSubCluster) StopSub() {
	ds.exit = true
	ds.clusterClient.Close()
	close(ds.statusQuit)
}
