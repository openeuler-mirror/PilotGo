package client

import (
	"encoding/json"
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/jwt"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type ScriptsRun struct {
	Batch         *common.Batch `json:"batch"`
	ScriptType    string        `json:"script_type"`
	ScriptContent string        `json:"script_content"`
	Params        string        `json:"params"`
	TimeOutSec    int           `json:"timeoutSec"`
}

func (c *Client) AgentRunScripts(batch *common.Batch, scriptType string, script string, params string, timeoutSec int) ([]*common.CmdResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/runScripts", serverInfo.Address, serverInfo.Port)

	scriptRun := &ScriptsRun{
		Batch:         batch,
		ScriptType:    scriptType,
		ScriptContent: script,
		Params:        params,
		TimeOutSec:    timeoutSec,
	}
	r, err := httputils.Post(url, &httputils.Params{
		Body: scriptRun,
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
