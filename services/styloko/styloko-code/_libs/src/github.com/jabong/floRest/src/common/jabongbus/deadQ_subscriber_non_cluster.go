package jabongbus

import (
	"encoding/json"
	"errors"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/redis.v3"
	"strconv"
	"time"
)

// persistent subscriber with non-cluster redis
type deadQSubNonCluster struct {
	client         *redis.Client
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
	msgRetryCount  int
	statusQuit     chan bool
}

// initialize with config values
func (ds *deadQSubNonCluster) init(conf *Subscriberconfig) error {
	// set client
	ds.client = redis.NewClient(
		&redis.Options{
			Network:  "tcp",
			Addr:     conf.RedisCon,
			PoolSize: MAX_ACTIVE_REDIS_CON})

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
func (ds *deadQSubNonCluster) SetProcessMsg(pm ProcessMessage) {
	ds.processMessage = pm
}

// set the subscriber id
func (ds *deadQSubNonCluster) setId() error {
	// increment the subcount in redis. It will also create subcount key if not present
	req := ds.client.Incr("{" + ds.messageList + "}" + DEADQ_SUBSCRIBER_COUNT)
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

// get the message for client
func (ds *deadQSubNonCluster) Get(timeout int) {
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
			ds.getReq = ds.client.BRPopLPush(ds.exceptionList, ds.bufList, time.Duration(timeout)*time.Second)
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
func (ds *deadQSubNonCluster) Ack() error {
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
func (ds *deadQSubNonCluster) ackLastMsg() {
	ds.ackReq = ds.client.RPop(ds.bufList)
}

// set the alive status
func (ds *deadQSubNonCluster) setAliveStatus() {
	go func() {
		firstTime := true
		data := new(statusData)
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
					ds.client.LRem(ds.statusList, 0, string(bufAry))
				}
				// add alive status
				if req := ds.client.RPush(ds.statusList, string(bufAry)); req.Err() != nil {
					logger.Error(req.Err().Error())
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
func (ds *deadQSubNonCluster) StopSub() {
	ds.exit = true
	ds.client.Close()
	close(ds.statusQuit)
}
