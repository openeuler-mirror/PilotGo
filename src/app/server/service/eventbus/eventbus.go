package eventbus

import (
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
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
