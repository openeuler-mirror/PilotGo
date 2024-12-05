/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugins licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 24 10:02:04 2024 +0800
 */
package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

var plugin_client *client.Client

// 注册event事件监听
func ListenEvent(eventTypes []int, callbacks common.EventCallback) error {
	var eventtypes []string
	for _, i := range eventTypes {
		eventtypes = append(eventtypes, strconv.Itoa(i))
	}

	eventServer, err := eventPluginServer()
	if err != nil {
		return err
	}

	url := eventServer + "/plugin/event/listener/register?eventTypes=" + strings.Join(eventtypes, ",")
	r, err := httputils.Put(url, &httputils.Params{
		Body: plugin_client.PluginInfo,
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
		registerEventCallback(eventType, callbacks)
	}
	return nil
}

// 取消注册event事件监听
func UnListenEvent(eventTypes []int) error {
	var eventtypes []string
	for _, i := range eventTypes {
		eventtypes = append(eventtypes, strconv.Itoa(i))
	}
	eventServer, err := eventPluginServer()
	if err != nil {
		return err
	}

	url := eventServer + "/plugin/event/listener/unregister?eventTypes=" + strings.Join(eventtypes, ",")
	r, err := httputils.Delete(url, &httputils.Params{
		Body: plugin_client.PluginInfo,
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
		unregisterEventCallback(eventType)
	}
	return nil
}

// 插件服务退出，取消注册所有本插件的event事件监听
func UnPluginListenEvent() error {
	eventServer, err := eventPluginServer()
	if err != nil {
		return err
	}

	url := eventServer + "/plugin/event/listener/unpluginRegister"
	r, err := httputils.Delete(url, &httputils.Params{
		Body: plugin_client.PluginInfo,
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
		EventType []int  `json:"eventType"`
		Status    string `json:"status"`
	}{}
	if err := resp.ParseData(data); err != nil || data.Status != "ok" {
		return err
	}
	for _, eventType := range data.EventType {
		unregisterEventCallback(eventType)
	}
	return nil
}

// 发布event事件
func PublishEvent(msg common.EventMessage) error {
	eventServer, err := eventPluginServer()
	if err != nil {
		return err
	}
	url := eventServer + "/plugin/event/publishEvent"
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
