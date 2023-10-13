package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type EventCallback func(e *common.EventMessage)

// 注册event事件监听
func (c *Client) ListenEvent(eventType int, callback EventCallback) error {
	url := c.Server + "/api/v1/pluginapi/listener"
	r, err := httputils.Put(url, nil)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	c.registerEventCallback(eventType, callback)

	return nil
}

// 取消注册event事件监听
func (c *Client) UnListenEvent(eventType int) error {
	url := c.Server + "/api/v1/pluginapi/listener"
	r, err := httputils.Delete(url, nil)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	// TODO: unregister event handler here
	c.unregisterEventCallback(eventType)
	return nil
}

// 发布event事件
func (c *Client) PublishEvent(msg common.EventMessage) error {
	// TODO:
	return errors.New("not implemented")
}

func (c *Client) registerEventCallback(eventType int, callback EventCallback) {
	c.eventCallbackMap[eventType] = callback
}

func (c *Client) unregisterEventCallback(eventType int) {
	delete(c.eventCallbackMap, eventType)
}

func (c *Client) ProcessEvent(event *common.EventMessage) {
	c.eventChan <- event
}

func (c *Client) startEventProcessor() {
	for {
		e := <-c.eventChan

		// TODO: process event message
		cb, ok := c.eventCallbackMap[e.MessageType]
		if ok {
			cb(e)
		}
	}
}
