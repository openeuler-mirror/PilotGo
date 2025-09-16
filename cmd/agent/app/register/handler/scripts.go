package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
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

func createTempScriptFile(workDir, suffix, encodedScript string) (string, error) {
	decodedScript, err := base64.StdEncoding.DecodeString(encodedScript)
	if err != nil {
		return "", fmt.Errorf("脚本内容base64解码失败: %s", err.Error())
	}

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		if err := os.MkdirAll(workDir, 0755); err != nil {
			return "", fmt.Errorf("创建临时工作目录失败: %s", err.Error())
		}
	}

	tmpFile, err := os.CreateTemp(workDir, "script_*"+suffix)
	if err != nil {
		return "", fmt.Errorf("创建临时脚本文件失败: %s", err.Error())
	}

	if _, err = tmpFile.Write(decodedScript); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("内容写入到脚本文件失败: %s", err.Error())
	}
	tmpFile.Close()

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("设置脚本可执行权限失败: %s", err.Error())
	}

	return tmpFile.Name(), nil
}

func runScript(interpreter, scriptPath string, params string, timeoutSec int) (string, string, int, error) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	args := append([]string{scriptPath}, params)
	cmd := exec.CommandContext(ctx, interpreter, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), "脚本执行超时", -1, fmt.Errorf("脚本执行超时")
	}

	if err != nil {
		retcode := -1
		if cmd.ProcessState != nil {
			retcode = cmd.ProcessState.ExitCode()
		}
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return stdout.String(), errMsg, retcode, err
	}

	return stdout.String(), stderr.String(), 0, nil
}

type ScriptsRun struct {
	ScriptType    string
	ScriptContent string
	Params        string
	TimeOut       int
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
		errorInfo := "反序列化失败:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	suffix, interpreter, err := getScriptInfo(d.ScriptType)
	if err != nil {
		errorInfo := "获取脚本类型失败:" + err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}

	scriptPath, err := createTempScriptFile(workDir, suffix, d.ScriptContent)
	if err != nil {
		errorInfo := err.Error()
		logger.Error("%s", errorInfo)
		resp_msg.Status = -1
		resp_msg.Error = errorInfo
		return c.Send(resp_msg)
	}
	defer os.Remove(scriptPath)

	logger.Debug("run script timeout: %v", d.TimeOut)
	logger.Debug("process run script command: %s %s %v", interpreter, scriptPath, d.Params)

	stdout, stderr, execCode, err := runScript(interpreter, scriptPath, d.Params, d.TimeOut)
	result := &utils.CmdResult{
		Stdout:  stdout,
		Stderr:  stderr,
		RetCode: execCode,
	}
	if err != nil {
		resp_msg.Status = -1
		resp_msg.Error = fmt.Sprintf("脚本执行错误: %s", err.Error())
		resp_msg.Data = result
		return c.Send(resp_msg)
	}

	resp_msg.Status = 0
	resp_msg.Data = result
	return c.Send(resp_msg)
}
