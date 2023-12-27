package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type CallbackHandler struct {
	RunCommandCallback RunCommandCallback
	TaskLen            int
}

type RunCommandCallback func([]*common.CmdResult)

func (c *Client) RunCommand(batch *common.Batch, cmd string) ([]*common.CmdResult, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}

	url := "http://" + c.Server() + "/api/v1/pluginapi/run_command"

	p := &common.CmdStruct{
		Batch:   batch,
		Command: base64.StdEncoding.EncodeToString([]byte(cmd)),
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
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
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/run_script"

	p := &ScriptStruct{
		Batch:  batch,
		Script: base64.StdEncoding.EncodeToString([]byte(script)),
		Params: params,
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
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
	if !c.IsBind() {
		return errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/run_command_async?plugin_name=" + c.PluginInfo.Name

	p := &common.CmdStruct{
		Batch:   batch,
		Command: base64.StdEncoding.EncodeToString([]byte(cmd)),
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
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
