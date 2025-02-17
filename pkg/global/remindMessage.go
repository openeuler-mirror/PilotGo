package global

import (
	"container/list"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
)

var WARN_MSG chan *WebsocketSendMsg = make(chan *WebsocketSendMsg, 100)

type WebsocketSendMsgType int

const (
	MachineSendMsg WebsocketSendMsgType = iota
	PluginSendMsg
	ServerSendMsg
)

type WebsocketSendMsg struct {
	MsgType WebsocketSendMsgType `json:"msgtype"`
	Time    string               `json:"time"`
	Msg     string               `json:"msg"`
}

func SendRemindMsg(msg_type WebsocketSendMsgType, msg string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("send remind message to closed channel WARN_MSG: %+v", r)
		}
	}()

	WARN_MSG <- &WebsocketSendMsg{
		MsgType: msg_type,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Msg:     msg,
	}
}

type LimitedList struct {
	capacity int
	data     *list.List
}

func NewLimitedList(capacity int) *LimitedList {
	return &LimitedList{
		capacity: capacity,
		data:     list.New(),
	}
}

func (l *LimitedList) Store(value interface{}) {
	if l.data.Len() >= l.capacity {
		l.data.Remove(l.data.Front())
	}
	l.data.PushBack(value)
}

func (l *LimitedList) GetAll() []interface{} {
	result := make([]interface{}, 0, l.data.Len())
	for e := l.data.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value)
	}
	return result
}

func (l *LimitedList) Len() int {
	return l.data.Len()
}
