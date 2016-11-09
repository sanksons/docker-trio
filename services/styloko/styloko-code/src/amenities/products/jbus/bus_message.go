package jbus

import (
	"bytes"
	"common/appconfig"
	"common/utils"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jabong/floRest/src/common/config"
	"github.com/satori/go.uuid"
)

var httpClient = &http.Client{
	Timeout: time.Duration(5 * time.Minute),
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	}}

const (
	SYSTEM       = "cat"
	USER         = "styl"
	MODULE_NAME  = "products"
	MESSAGE_TYPE = "Products"

	TYPE_PRODUCT_DELETE = "product_delete"
	TYPE_PRODUCT_UPDATE = "product_update"
)

type ProductMessage struct {
	MessageId       string                 `json:"id"`
	Timestamp       int64                  `json:"timestamp"`
	TransactionId   *string                `json:"transaction_id,omitempty"`
	Publisher       string                 `json:"publisher_name"`
	RoutingKey      string                 `json:"routing_key"`
	Type            string                 `json:"type"`
	TypeOfChange    string                 `json:"type_of_change"`
	ListOfChanges   string                 `json:"list_of_changes"`
	ExtraParameters map[string]interface{} `json:"extra_parameters,omitempty"`
	Data            interface{}            `json:"data"`
}

func (pm *ProductMessage) Publish() error {
	conf := GetConfig()
	if conf.JBus.URL == "" {
		return errors.New("(pm *ProductMessage)#Publish(): jbus URL Not defined")
	}
	url := conf.JBus.URL
	b, err := utils.JSONMarshal(pm, true)
	if err != nil {
		return errors.New("(pm *ProductMessage)#Publish(): Marshalling Failed " + err.Error())
	}
	req, err := http.NewRequest(
		"POST", url, bytes.NewBuffer(b),
	)
	if err != nil {
		return errors.New("(pm *ProductMessage)#Publish(): " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.New("(pm *ProductMessage)#Publish(): " + err.Error())
	}
	req.Close = true
	resp.Body.Close()
	return nil
}

func GetNewProductMessage() ProductMessage {
	msg := ProductMessage{}
	msg.Timestamp = time.Now().Unix()
	msg.MessageId = GenerateBusMessageId(SYSTEM, USER)
	msg.Publisher = GetConfig().JBus.Publisher
	msg.RoutingKey = GetConfig().JBus.RoutingKey
	return msg
}

func GetConfig() *appconfig.AppConfig {
	return config.ApplicationConfig.(*appconfig.AppConfig)
}

func GenerateTransactionId() string {
	return fmt.Sprintf("%s", uuid.NewV4())
}

func GenerateBusMessageId(system string, user string) string {
	template := "%s-%s-%s"
	requestId := strconv.Itoa(time.Now().Year()) + randomNumber(8)
	return fmt.Sprintf(template, requestId, system, user)
}

//takes an integer as input
//makes a slice of bytes of size integer value
//fill the slice by generating random no between 48 and 57
// and returns it as a string
func randomNumber(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(48, 57))
	}
	return string(bytes)
}

//This function takes two integers as parameters
//seeding is done using the current time(converted to unix time in nanosecond)
//returns a non-negative pseudo-random number
//from [0,max-min) + min
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
