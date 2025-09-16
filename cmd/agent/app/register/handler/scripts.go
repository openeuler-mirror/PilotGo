package handler

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func getScriptInfo(scriptType string) (suffix string, interpreter string, err error) {
	switch scriptType {
	case "Shell":
		return ".sh", "/bin/bash", nil
	case "Python":
		return ".py", "/usr/bin/python3", nil
	case "Perl":
		return ".pl", "/usr/bin/perl", nil
	case "Ruby":
		return ".rb", "/usr/bin/ruby", nil
	case "PHP":
		return ".php", "/usr/bin/php", nil
	default:
		err = fmt.Errorf("不支持的脚本类型: %s", scriptType)
	}
	return
}

type ScriptsRun struct {
	Batch         *common.Batch `json:"batch"`
	ScriptType    string        `json:"script_type"`
	ScriptContent string        `json:"script_content"`
	Params        string        `json:"params"`
}

func AgentRunScriptsHandler(c *network.SocketClient, msg *protocol.Message) error {
	var workDir = "/tmp/scripts/"
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}

	d := &ScriptsRun{}
	err := msg.BindData(d)
	if err != nil {
		errorInfo := "parse data error:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	suffix, interpreter, err := getScriptInfo(d.ScriptType)
	if err != nil {
		errorInfo := "parse data error:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	decoded_script, err := base64.StdEncoding.DecodeString(d.ScriptContent)
	if err != nil {
		errorInfo := fmt.Sprintf("Err decoding base64: %v, %s", d.ScriptContent, err.Error())
		logger.Error("Err decoding base64: %v, %s", d.ScriptContent, err.Error())
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	// 创建临时脚本文件
	tmpFile, err := os.CreateTemp(workDir, "script_*"+suffix)
	if err != nil {
		errorInfo := "create temp file error:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}
	defer os.Remove(tmpFile.Name()) // 删除临时脚本文件

	_, err = tmpFile.Write(decoded_script)
	if err != nil {
		errorInfo := "write script to temp file error:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}
	tmpFile.Close()

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		errorInfo := "chmod temp file error:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	logger.Debug("process run script command: %s %v", interpreter, d.Params)

	args := append([]string{tmpFile.Name()}, []string{d.Params}...)
	cmd := exec.Command(interpreter, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errorInfo := fmt.Sprintf("run script error: %s, output: %s", err.Error(), string(output))
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	result := &utils.CmdResult{
		Stdout:  string(output),
		Stderr:  "",
		RetCode: cmd.ProcessState.ExitCode(),
	}

	resp_msg.Status = 0
	resp_msg.Data = result
	return c.Send(resp_msg)
}
