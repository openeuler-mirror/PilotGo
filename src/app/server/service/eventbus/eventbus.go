package eventbus

import (
	"encoding/json"
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
	//EeventTypeMap[1] = append(eventTypeMap[1], l)
}

func (e *EventBus) RemoveListener(l *Listener) {
	e.Lock()
	defer e.Unlock()

	for index, v := range e.listeners {
		if v.Name == l.Name && v.URL == l.URL {
			e.listeners = append(e.listeners[:index], e.listeners[index+1:]...)
		}
	}
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
			r, err := httputils.Put(listener.URL+"/plugin_manage/api/v1/event", &httputils.Params{
				Body: mes,
			})
			if err != nil {
				logger.Error(listener.Name + " plugin request error:" + err.Error())
			}
			resp := &struct {
				Status string
				Error  string
			}{}
			if err := json.Unmarshal(r.Body, resp); err != nil {
				logger.Error(listener.Name + " plugin response error:" + err.Error())
			}
			if resp.Status != "ok" {
				logger.Error(listener.Name + " plugin response status:" + resp.Error)
			}
		}
	}
}

var globalEventBus *EventBus

func Init() {
	globalEventBus = &EventBus{}
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
