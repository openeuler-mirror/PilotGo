package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type EventCallback func(e *common.EventMessage)

// 注册event事件监听
func (c *Client) ListenEvent(eventTypes []int, callbacks []EventCallback) error {
	var eventtypes []string
	for _, i := range eventTypes {
		eventtypes = append(eventtypes, strconv.Itoa(i))
	}

	url := c.Server() + "/api/v1/pluginapi/listener?eventTypes=" + strings.Join(eventtypes, ",")
	r, err := httputils.Put(url, &httputils.Params{
		Body: c.PluginInfo,
	})
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("server process error:" + strconv.Itoa(r.StatusCode))
	}

	resp := &common.CommonResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(resp.Message)
	}

	data := &struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := resp.ParseData(data); err != nil {
		return err
	}
	for i, eventType := range eventTypes {
		c.registerEventCallback(eventType, callbacks[i])
	}
	return nil
}

// 取消注册event事件监听
func (c *Client) UnListenEvent(eventTypes []int) error {
	var eventtypes []string
	for _, i := range eventTypes {
		eventtypes = append(eventtypes, strconv.Itoa(i))
	}

	url := c.Server() + "/api/v1/pluginapi/listener?eventTypes=" + strings.Join(eventtypes, ",")
	r, err := httputils.Delete(url, &httputils.Params{
		Body: c.PluginInfo,
	})
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("server process error:" + strconv.Itoa(r.StatusCode))
	}

	resp := &common.CommonResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(resp.Message)
	}

	data := &struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := resp.ParseData(data); err != nil {
		return err
	}

	for _, eventType := range eventTypes {
		c.unregisterEventCallback(eventType)
	}
	return nil
}

// 发布event事件
func (c *Client) PublishEvent(msg common.EventMessage) error {
	url := c.Server() + "/api/v1/pluginapi/publish_event"
	r, err := httputils.Put(url, &httputils.Params{
		Body: &msg,
	})
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("server process error:" + strconv.Itoa(r.StatusCode))
	}

	resp := &common.CommonResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(resp.Message)
	}

	data := &struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := resp.ParseData(data); err != nil {
		return err
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
