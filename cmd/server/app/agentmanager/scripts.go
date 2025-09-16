package agentmanager

import (
	"encoding/base64"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
)

// 远程在agent上运行脚本文件
func (a *Agent) AgentRunScripts(scriptType string, script string, params string) (*utils.CmdResult, error) {
	encoded_script := base64.StdEncoding.EncodeToString([]byte(script))
	data := struct {
		Script     string
		Params     string
		ScriptType string
	}{
		Script:     encoded_script,
		ScriptType: scriptType,
		Params:     params,
	}
	responseMessage, err := a.SendMessageWrapper(protocol.AgentRunScripts, data, "failed to run script on agent", 0, nil, "")
	result, ok := responseMessage.(*utils.CmdResult)
	if !ok {
		return &utils.CmdResult{}, err
	}
	return result, err
}
