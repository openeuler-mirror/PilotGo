package eventbus

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type Listener struct {
	Name string
	URL  string
}

type EventBus struct {
	sync.Mutex
	listeners []*Listener
	stop      chan struct{}
	wait      sync.WaitGroup
	event     chan *common.EventMessage
}

var eventTypeMap map[int][]Listener

func (e *EventBus) AddListener(l *Listener) {
	e.Lock()
	defer e.Unlock()
	e.listeners = append(e.listeners, l)
}

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

func (e *EventBus) AddEventMap(eventtpye int, l *Listener) {
	e.Lock()
	defer e.Unlock()
	eventTypeMap[eventtpye] = append(eventTypeMap[eventtpye], *l)
}

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

func (e *EventBus) Run() {
	go func(e *EventBus) {
		for {
			select {
			case <-e.stop:
				logger.Info("event bus exit")
				e.wait.Done()
				break
			case <-e.event:
				e.broadcast()
			}
		}
	}(e)
}

func (e *EventBus) Stop() {
	e.wait.Add(1)
	e.stop <- struct{}{}
	e.wait.Wait()
}

func (e *EventBus) publish(m *common.EventMessage) {
	e.event <- m
}

func (e *EventBus) broadcast() {
	// TODO：将event message发送给监听的listener
	mes := <-e.event
	listeners, ok := eventTypeMap[mes.MessageType]
	if ok {
		for _, listener := range listeners {
			r, err := httputils.Post(listener.URL+"/plugin_manage/api/v1/event", &httputils.Params{
				Body: mes,
			})
			if err != nil {
				logger.Error(listener.Name + " plugin request error:" + err.Error())
			}
			if r.StatusCode != http.StatusOK {
				logger.Error("server process error:" + strconv.Itoa(r.StatusCode))
			}

			resp := &common.CommonResult{}
			if err := json.Unmarshal(r.Body, resp); err != nil {
				logger.Error(listener.Name + " plugin request error:" + err.Error())
			}
			if resp.Code != http.StatusOK {
				logger.Error(resp.Message)
			}

			data := &struct {
				Status string `json:"status"`
				Error  string `json:"error"`
			}{}
			if err := resp.ParseData(data); err != nil {
				logger.Error(listener.Name + " plugin request error:" + err.Error())
			}
		}
	}
}

var globalEventBus *EventBus

func Init() {
	eventTypeMap = make(map[int][]Listener)
	globalEventBus = &EventBus{
		event: make(chan *common.EventMessage, 20),
	}
	globalEventBus.Run()
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
