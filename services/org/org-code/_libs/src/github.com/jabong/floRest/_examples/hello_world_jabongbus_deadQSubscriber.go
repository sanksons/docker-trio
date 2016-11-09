package examples

import (
	"fmt"
	"github.com/jabong/floRest/src/common/jabongbus"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
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

type dclient struct { // dead queue subscriber
	Sub jabongbus.DeadQSubscriber
}

func (obj *dclient) Process(msg *jabongbus.Message, err error) {
	if err != nil { // some error, decide if you want to stop the subscriber or continue
		fmt.Println(err.Error())
		obj.Sub.StopSub()
	} else {
		// process the message
		fmt.Println(msg)
		// ack the message
		if ackErr := obj.Sub.Ack(); ackErr != nil {
			fmt.Println("ack failed,error:" + ackErr.Error())
		}
	}
}

func (a HelloWorld) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	var err error
	// subscribe to this message
	sconf := new(jabongbus.Subscriberconfig)
	sconf.Cluster = true
	sconf.Persistent = true
	sconf.Publisher = "rajcomics"
	sconf.RoutingKey = "chachachaudhri"
	sconf.RedisCon = "localhost:7001,localhost:7002"
	myClient := new(dclient)
	myClient.Sub, err = jabongbus.GetDeadQSubscriber(sconf)
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
