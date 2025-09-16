package agentmanager

import (
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
)

// 远程在agent上运行脚本文件
func (a *Agent) AgentRunScripts(scriptType string, script string, params string) (*utils.CmdResult, error) {
	data := struct {
		ScriptContent string
		Params        string
		ScriptType    string
	}{
		ScriptContent: script,
		ScriptType:    scriptType,
		Params:        params,
	}
	responseMessage, err := a.SendMessageWrapper(protocol.AgentRunScripts, data, "failed to run script on agent", 0, nil, "")
	result, ok := responseMessage.(*utils.CmdResult)
	if !ok {
		return &utils.CmdResult{}, err
	}
	return result, err
}
