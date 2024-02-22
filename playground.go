package main

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessageEvent struct {
	Text string `json:"text"`
}

func main() {

	sendMsg := SendMessageEvent{
		Text: "hello",
	}

	sendMsgJSON, _ := json.Marshal(sendMsg)

	e := Event{
		Type:    "send_message",
		Payload: sendMsgJSON,
	}

	eJSON, _ := json.Marshal(e)

	fmt.Println(e.Payload)
	fmt.Println(string(eJSON))

	var NewSendMsg SendMessageEvent
	json.Unmarshal(e.Payload, &NewSendMsg)

	fmt.Println(NewSendMsg.Text)

}
