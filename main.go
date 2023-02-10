

/*
* File: main.go
* Author: Joel
* Copyright: 2023, Smart Cities Peru.
* License: MIT
*
* Purpose:
* This is a file that MQTT uses to subscribe to a topic and send a post request to an API
* then wait for a response from the request and submit this response via another topic.
* jksdljf;sdlkfjsdlk;jsdk;fjsd;kljsd;jsdfkjdfkjfjfj
*
* API = http://192.168.71.150/pstprnt
* TOPIC_PUB = /printer01/status
* TOPIC_SUB = /printer01/print
* Last modified: 2023-01-27
*
* fksj;dflksjdfklsdjklckv;lkdfjskldflf
* fjshdflkjshfkjshfjksjbvxcvxnvdjfkfd
*
*
*
*
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"os"
	"time"
	//"time"
)

var PathUrl = os.Getenv("URL_PATH")
var MqttBroker = os.Getenv("MQTT_BROKER")

type RequestBody struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type StatusBody struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}
type MqttClient struct {
	broker   string
	clientID string
}

type MqttSubscription struct{
	client MQTT.Client
	topic string
	qos byte
	cb MQTT.MessageHandler
}
func handleIncomingMessage(client MQTT.Client, message MQTT.Message) {
	var msg RequestBody
	if err := json.Unmarshal(message.Payload(), &msg); err != nil {
		return
	}

	// send POST format TEXT
	status := 0
	statusTmp, err := postPrinter(msg.Id, msg.Data)
	if err == nil {
		fmt.Println(err)
		status = *statusTmp
	}


	statusBody, err := createStatusBody(msg.Id, status)
	if err != nil {
		fmt.Println(err)
		return
	}

	err=publishTopic(client,"/printer01/state",statusBody)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (s *StatusBody) createStatusBody2() ([]byte, error){
	statusBodyJson, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return statusBodyJson, nil
}
func createStatusBody(msgId int, status int) ([]byte, error) {
	statusBody := StatusBody{
		Id:     msgId,
		Status: status,
	}
	statusBodyJson, err := json.Marshal(statusBody)
	if err != nil {
		return nil, err
	}
	return statusBodyJson, nil
}
func publishTopic(client MQTT.Client, topic string, payload []byte) error{
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func postPrinter(msgId int, data string) (*int, error) {
	// send POST format TEXT
	res, err := http.Post(PathUrl+"/pstprnt", "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	status := 0
	if res != nil {
		status = res.StatusCode
		defer res.Body.Close()
	}
	return &status, nil
}


func (m *MqttClient) Connect() (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", m.broker))
	opts.SetClientID(m.clientID)
	opts.SetDefaultPublishHandler(handleIncomingMessage)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}
func (s *MqttSubscription) Subscribe() MQTT.Token{
	return s.client.Subscribe(s.topic,s.qos,s.cb)
}
func main() {

	mqttClient := MqttClient{MqttBroker,"clientId-i3zf3Qj"}
	client,err:=mqttClient.Connect()
	if err!=nil{
		panic(err)
	}
	fmt.Println(client.IsConnected())

	subscription := MqttSubscription{client,"/printer01/print",0,nil}
	if token := subscription.Subscribe(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

