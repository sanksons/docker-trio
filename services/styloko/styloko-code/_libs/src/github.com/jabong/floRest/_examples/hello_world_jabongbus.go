package examples

import (
	"fmt"
	"github.com/jabong/floRest/src/common/jabongbus"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"time"
)

type HelloWorld struct {
	id string
}

func (n *HelloWorld) SetID(id string) {
	n.id = id
}

func (n HelloWorld) GetID() (id string, err error) {
	return n.id, nil
}

func (a HelloWorld) Name() string {
	return "HelloWord"
}

type Data struct {
	Name     string
	Location string
}

type pclient struct { // persistent client
	Sub jabongbus.PeristentSubscriber
}

func (obj *pclient) Process(msg *jabongbus.Message, err error) {
	if err != nil { // some error, decide if you want to stop the subscriber or continue
		fmt.Println(err.Error())
		obj.Sub.StopSub()
	} else {
		// NOTE: Message has reached to client, now you can either Ack() if processed message successfuly
		//  or NoAck() if not able to process the message.
		// In case of NoAck() message will be redelivered until max_retry is reached. This value is set
		// during publish of message itself.

		fmt.Println(msg) // print the message
		// ack the message
		if ackErr := obj.Sub.Ack(); ackErr != nil {
			fmt.Println("ack failed,error:" + ackErr.Error())
		}
		/*
			//no ack the message:
			if noackErr := obj.Sub.NoAck(); noackErr != nil {
				fmt.Println("noack failed,error:" + noackErr.Error())
			}
		*/
	}
}

func (a HelloWorld) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	//publish one message
	pconf := new(jabongbus.PublisherConfig)
	pconf.Url = "http://localhost:8090/omnibus/"
	myPublisher := jabongbus.GetPublisher(pconf)
	pubReq := new(jabongbus.PubRequest)
	// fill the headers
	pubReq.Headers.RequestId = "r"
	pubReq.Headers.SessionId = "s"
	pubReq.Headers.TransactionId = "t"
	pubReq.Headers.UserId = "u"
	// fill the body
	pubReq.Body.PublisherName = "rajcomics"
	pubReq.Body.RoutingKey = "chachachaudhri"
	pubReq.Body.Type = "images"
	pubReq.Body.Data = Data{Name: "sabu", Location: "jupiter"}
	pubReq.Body.RetryCount = 2
	// publish message
	resp, err := myPublisher.PublishMessage(pubReq, time.Duration(4)*time.Second, false)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(resp.BusResponse))
	}

	// subscribe to this message
	sconf := new(jabongbus.Subscriberconfig)
	sconf.Cluster = true
	sconf.Persistent = true
	sconf.Publisher = "rajcomics"
	sconf.RoutingKey = "chachachaudhri"
	sconf.RedisCon = "localhost:7001,localhost:7002"
	myClient := new(pclient)
	myClient.Sub, err = jabongbus.GetPersistentSubscriber(sconf)
	if err == nil {
		myClient.Sub.SetProcessMsg(myClient) // define your message processing function
		// NOTE: We have used 2 seconds as timeout in this example, in case you want to receive message
		// as and when it comes, use 0 as timeout value.
		myClient.Sub.Get(2) // get the message with 2 seconds timeout
	} else { // subscriber not initialized
		fmt.Println(err)
	}
	//Business Logic
	return io, nil
}
