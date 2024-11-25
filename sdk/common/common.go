package common

import (
	"encoding/json"
	"time"
)

type MessageData struct {
	MsgType     int         `json:"msg_type_id"`
	MessageType string      `json:"msg_type"`
	TimeStamp   time.Time   `json:"timestamp"`
	Data        interface{} `json:"data"`
}

type EventMessage struct {
	MessageType int    `json:"msgType"`
	MessageData string `json:"msgData"`
}

type EventCallback func(e *EventMessage)

func (msgData *MessageData) ToMessageDataString() (string, error) {
	data, err := json.Marshal(msgData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type CommonResult struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

func (r *CommonResult) ParseData(d interface{}) error {
	return json.Unmarshal(r.Data, d)
}
