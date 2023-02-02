/*
 * File: main.go
 * Author: User
 * Copyright: 2023, Smart Cities Peru.
 * License: MIT
 *
 * Purpose:
 * This is the entry point for the application.
 *
 * Last Modified: 2023-02-02
 */



package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"os"
)

var PathUrl =os.Getenv("URL_PATH")
var MqttBroker = os.Getenv("BROKER")

func main(){
	// Connect to MQTT broker
	var broker = MqttBroker
	//var port = 1883
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s",broker))
	opts.SetClientID("clientId-i3zf3Qj")
	opts.SetDefaultPublishHandler(onMessageReceived)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to topic
	if token := client.Subscribe("/printer01/print", 0, nil);
	token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	//select{}
	http.ListenAndServe(":8000", nil)
}
type RequestBody struct {
	Data string `json:"data"`
}
type StatusBody struct {
	Id float64 `json:"id"`
	Status int `json:"status"`
}
func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var msg map[string]interface{}
	json.Unmarshal(message.Payload(), &msg)

	// send POST format TEXT
	if data, ok := msg["data"].(string); ok {
		//jsonData, err := json.Marshal(RequestBody{Data:data})

		//fmt.Println(data)
		res, err := http.Post(PathUrl+"/pstprnt", "application/json", bytes.NewBuffer([]byte(data)))
		if err != nil {
			fmt.Println(err)
			//token := client.Publish("/printer/state", 0, false, "NOOOOOO")
			//fmt.Println("SUPERERROR")
			//token.Wait()
		}
		status:=res.StatusCode
		statusBody := StatusBody{
			Id:     msg["id"].(float64),
			Status: status,
		}
		statusBodyJson, err := json.Marshal(statusBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		//defer res.Body.Close()
		fmt.Println("Status:", res.StatusCode)
		if status==200 {
			token := client.Publish("/printer01/state", 0, false, statusBodyJson)
			token.Wait()
			fmt.Println(data)
			if token.Error() != nil {
				fmt.Println(token.Error())
			}
		} else{
			token := client.Publish("/printer01/state", 0, false, statusBodyJson)
			token.Wait()
			fmt.Println(data)
			if token.Error() != nil {
				fmt.Println(token.Error())
			}
		}


	} else {
		fmt.Println("data not found or is not a string")
	}

	//Format json
	//if data, ok := msg["data"].(string); ok {
	//	jsonData, err := json.Marshal(RequestBody{Data:data})
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	//fmt.Println(data)
	//	res, err := http.Post(PathUrl+"/pstprnt", "application/json", bytes.NewBuffer(jsonData))
	//	if err != nil {
	//		fmt.Println(err)
	//		//token := client.Publish("/printer/state", 0, false, "NOOOOOO")
	//		//fmt.Println("SUPERERROR")
	//		//token.Wait()
	//	}
	//	status:=res.StatusCode
	//	statusBody := StatusBody{
	//		Id:     msg["id"].(float64),
	//		Status: status,
	//	}
	//	statusBodyJson, err := json.Marshal(statusBody)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	//defer res.Body.Close()
	//	fmt.Println("Status:", res.StatusCode)
	//	if status==200 {
	//		token := client.Publish("/printer01/state", 0, false, statusBodyJson)
	//		token.Wait()
	//		fmt.Println(data)
	//		if token.Error() != nil {
	//			fmt.Println(token.Error())
	//		}
	//	} else{
	//		token := client.Publish("/printer01/state", 0, false, statusBodyJson)
	//		token.Wait()
	//		fmt.Println(data)
	//		if token.Error() != nil {
	//			fmt.Println(token.Error())
	//		}
	//	}
	//
	//
	//} else {
	//	fmt.Println("data not found or is not a string")
	//}
}
