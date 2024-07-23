package common

import "encoding/json"

type EventMessage struct {
	MessageType int
	MessageData string
}
type CommonResult struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

func (r *CommonResult) ParseData(d interface{}) error {
	return json.Unmarshal(r.Data, d)
}
