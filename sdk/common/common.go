package common

import (
	"encoding/json"
)

type EventMessage struct {
	MessageType int    `json:"msgType"`
	MessageData string `json:"msgData"`
}

type CommonResult struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

func (r *CommonResult) ParseData(d interface{}) error {
	return json.Unmarshal(r.Data, d)
}

// 将 MessageData json字符串转换成指定结构体的message消息数据
func ToMessage(d string, s interface{}) error {
	return json.Unmarshal([]byte(d), s)
}
func ToJSONString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type EventCallback func(e *EventMessage)
