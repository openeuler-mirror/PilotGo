package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type EventCallback func(e *common.EventMessage)

// 注册event事件监听
func (c *Client) ListenEvent(eventType int, callback EventCallback) error {
	url := c.Server + "/api/v1/pluginapi/listener?eventType=" + strconv.Itoa(eventType)

	r, err := httputils.Put(url, &httputils.Params{
		Body: c.PluginInfo,
	})
	if err != nil {
		return err
	}

	resp := &common.RespResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK || resp.Code != http.StatusOK {
		return errors.New(resp.Msg)
	}

	c.registerEventCallback(eventType, callback)

	return nil
}

// 取消注册event事件监听
func (c *Client) UnListenEvent(eventType int) error {
	url := c.Server + "/api/v1/pluginapi/listener?eventType=" + strconv.Itoa(eventType)
	r, err := httputils.Delete(url, &httputils.Params{
		Body: c.PluginInfo,
	})
	if err != nil {
		return err
	}

	resp := &common.RespResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK || resp.Code != http.StatusOK {
		return errors.New(resp.Msg)
	}

	// TODO: unregister event handler here
	c.unregisterEventCallback(eventType)
	return nil
}

// 发布event事件
func (c *Client) PublishEvent(msg common.EventMessage) error {
	url := c.Server + "/api/v1/pluginapi/publish_event"
	r, err := httputils.Post(url, &httputils.Params{
		Body: &msg,
	})
	if err != nil {
		return err
	}

	resp := &common.RespResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK || resp.Code != http.StatusOK {
		return errors.New(resp.Msg)
	}

	return nil
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
	go func() {
		for {
			e := <-c.eventChan

			// TODO: process event message
			cb, ok := c.eventCallbackMap[e.MessageType]
			if ok {
				cb(e)
			}
		}
	}()

}
