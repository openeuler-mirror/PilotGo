/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Aug 06 17:35:12 2025 +0800
 */
package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/jwt"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) RunCommand(batch *common.Batch, cmd string) ([]*common.CmdResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/run_command", serverInfo.Address, serverInfo.Port)

	p := &common.CmdStruct{
		Batch:   batch,
		Command: base64.StdEncoding.EncodeToString([]byte(cmd)),
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	res := &struct {
		Code    int                 `json:"code"`
		Message string              `json:"msg"`
		Data    []*common.CmdResult `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

type ScriptStruct struct {
	Batch  *common.Batch `json:"batch"`
	Script string        `json:"script"`
	Params []string      `json:"params"`
}

func (c *Client) RunScript(batch *common.Batch, script string, params []string) ([]*common.CmdResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/run_script", serverInfo.Address, serverInfo.Port)

	p := &ScriptStruct{
		Batch:  batch,
		Script: base64.StdEncoding.EncodeToString([]byte(script)),
		Params: params,
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	res := &struct {
		Code    int                 `json:"code"`
		Message string              `json:"msg"`
		Data    []*common.CmdResult `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) RunCommandAsync(batch *common.Batch, cmd string, callback RunCommandCallback) error {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/run_command_async?plugin_name=", serverInfo.Address, serverInfo.Port)

	p := &common.CmdStruct{
		Batch:   batch,
		Command: base64.StdEncoding.EncodeToString([]byte(cmd)),
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return err
	}

	res := struct {
		Code int `json:"code"`
		Data struct {
			TaskID  string `json:"task_id"`
			TaskLen int    `json:"task_len"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &res); err != nil {
		return err
	}

	taskID := res.Data.TaskID
	TaskLen := res.Data.TaskLen
	c.registerCommandResultCallback(taskID, TaskLen, callback)

	return nil
}

func (c *Client) startCommandResultProcessor() {
	go func() {
		for {
			d := <-c.asyncCmdResultChan

			cb, ok := c.cmdProcessorCallbackMap[d.TaskID]
			if !ok {
				continue
			}

			// 注意：map并发安全
			cb.RunCommandCallback(d.Result)
			cb.TaskLen = cb.TaskLen - len(d.Result)
			if cb.TaskLen == 0 {
				delete(c.cmdProcessorCallbackMap, d.TaskID)
			}
		}
	}()
}

func (c *Client) registerCommandResultCallback(taskID string, taskLen int, callback RunCommandCallback) {
	rccb := c.cmdProcessorCallbackMap[taskID]
	rccb.RunCommandCallback = callback
	rccb.TaskLen = taskLen
	c.cmdProcessorCallbackMap[taskID] = rccb
}

func (c *Client) ProcessCommandResult(command_result *common.AsyncCmdResult) {
	c.asyncCmdResultChan <- command_result
}
