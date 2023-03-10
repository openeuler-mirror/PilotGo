package eventbus

import (
	"sync"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

type Listener struct {
	Name string
	URL  string
}

const (
	// 主机安装软件包
	MsgPackageInstall = 0
	// 主机升级软件包
	MsgPackageUpdate = 1
	// 主机卸载软件包
	MsgPackageUninstall = 2
	// 主机ip变更
	MsgIPChange = 3

	// 平台新增主机
	MsgHostAdd = 10
	// 平台移除主机
	MsgHostRemove = 11

	// 插件添加
	MsgPluginAdd = 20
	// 插件卸载
	MsgPluginRemove = 21
)

type EventMessage struct {
	MessageType int
	MessageData interface{}
}

type EventBus struct {
	sync.Mutex
	listeners []*Listener
	stop      chan struct{}
	wait      sync.WaitGroup
	event     chan *EventMessage
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
		if v.Name == l.Name && v.URL == v.URL {
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

func (e *EventBus) publish(m *EventMessage) {
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

func PublishEvent(m *EventMessage) {
	globalEventBus.publish(m)
}
