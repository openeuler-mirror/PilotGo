package eventbus

import (
	"net/http"
	"strconv"
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"k8s.io/klog/v2"
)

type Listener struct {
	Name string
	URL  string
}

type EventBus struct {
	sync.Mutex
	listeners []*Listener
	stop      chan struct{}
	event     chan *common.EventMessage
}

var eventTypeMap map[int][]Listener

// 添加监听事件
func (e *EventBus) AddListener(l *Listener) {
	e.Lock()
	defer e.Unlock()
	e.listeners = append(e.listeners, l)
}

// 删除监听事件
func (e *EventBus) RemoveListener(l *Listener) {
	e.Lock()
	defer e.Unlock()

	for index, v := range e.listeners {
		if v.Name == l.Name && v.URL == l.URL {
			if index == len(e.listeners)-1 {
				e.listeners = e.listeners[:index]
			} else {
				e.listeners = append(e.listeners[:index], e.listeners[index+1:]...)
			}
			break
		}
	}
}

// 添加event事件
func (e *EventBus) AddEventMap(eventtpye int, l *Listener) {
	e.Lock()
	defer e.Unlock()
	eventTypeMap[eventtpye] = append(eventTypeMap[eventtpye], *l)
}

// 删除event事件
func (e *EventBus) RemoveEventMap(eventtpye int, l *Listener) {
	e.Lock()
	defer e.Unlock()
	for i, v := range eventTypeMap[eventtpye] {
		if (v.Name == l.Name) && (v.URL == l.URL) {
			if i == len(eventTypeMap[eventtpye])-1 {
				eventTypeMap[eventtpye] = eventTypeMap[eventtpye][:i]
			} else {
				eventTypeMap[eventtpye] = append(eventTypeMap[eventtpye][:i], eventTypeMap[eventtpye][i+1:]...)
			}
			break
		}
	}
}

// 判断监听是否存在
func (e *EventBus) IsExitEventMap(l *Listener) bool {
	e.Lock()
	defer e.Unlock()
	for _, value := range eventTypeMap {
		for _, v := range value {
			if (v.Name == l.Name) && (v.URL == l.URL) {
				return true
			}
		}
	}
	return false
}

func (e *EventBus) Run(stopCh <-chan struct{}) {
	go func() {
		<-stopCh
		e.Stop()
		klog.Warningln("EventBus prepare stop")
	}()
	go func(e *EventBus) {
		for {
			select {
			case <-e.stop:
				klog.Warningln("EventBus success exit ")
				return
			case m := <-e.event:
				e.broadcast(m)
			}
		}
	}(e)

}

func (e *EventBus) Stop() {
	e.stop <- struct{}{}
}

func (e *EventBus) publish(m *common.EventMessage) {
	e.event <- m
}

func (e *EventBus) broadcast(msg *common.EventMessage) {
	listeners, ok := eventTypeMap[msg.MessageType]
	if ok {
		for _, listener := range listeners {
			r, err := httputils.Post(listener.URL+"/plugin_manage/api/v1/event", &httputils.Params{
				Body: msg,
			})
			if err != nil {
				logger.Error(listener.Name + "plugin process error:" + err.Error())
			}
			if r.StatusCode != http.StatusOK {
				logger.Error(listener.Name + "plugin process error:" + strconv.Itoa(r.StatusCode))
			}
		}
	}
}

var globalEventBus *EventBus

func Init(stopCh <-chan struct{}) {
	eventTypeMap = make(map[int][]Listener)
	globalEventBus = &EventBus{
		event: make(chan *common.EventMessage, 20),
		stop:  make(chan struct{}),
	}
	globalEventBus.Run(stopCh)
}

func Stop() {
	globalEventBus.Stop()
}

func AddListener(l *Listener) {
	globalEventBus.AddListener(l)
}

func RemoveListener(l *Listener) {
	globalEventBus.RemoveListener(l)
}

func PublishEvent(m *common.EventMessage) {
	globalEventBus.publish(m)
}

func AddEventMap(eventtype int, l *Listener) {
	globalEventBus.AddEventMap(eventtype, l)
}

func RemoveEventMap(eventtype int, l *Listener) {
	globalEventBus.RemoveEventMap(eventtype, l)
}

func IsExitEventMap(l *Listener) bool {
	return globalEventBus.IsExitEventMap(l)
}
