package client

import (
	"encoding/base64"
	"encoding/json"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type TaskCmdResult struct {
	TaskId string    `json:"taskId"`
	Result CmdResult `json:"result"`
}
type CmdResult struct {
	MachineUUID string `json:"machine_uuid"`
	MachineIP   string `json:"machine_ip"`
	RetCode     int    `json:"retcode"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
}

type CmdStruct struct {
	Batch   *common.Batch `json:"batch"`
	Command string        `json:"command"`
	TaskId  string        `json:"taskId"`
}

func (c *Client) RunCommand(batch *common.Batch, cmd string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_command"

	p := &CmdStruct{
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
		Code    int          `json:"code"`
		Message string       `json:"msg"`
		Data    []*CmdResult `json:"data"`
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

func (c *Client) RunScript(batch *common.Batch, script string, params []string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_script"

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
		Code    int          `json:"code"`
		Message string       `json:"msg"`
		Data    []*CmdResult `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
