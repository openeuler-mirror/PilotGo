package pluginapi

import (
	"encoding/base64"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	batchservice "gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ApplyConfig(c *gin.Context) {

}

// 配置文件下发
func FileDeploy(c *gin.Context) {
	fd := &struct {
		DeployBatch    common.Batch `json:"deploybatch"`
		DeployPath     string       `json:"deploypath"`
		DeployFileName string       `json:"deployname"`
		DeployText     string       `json:"deployfile"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	path := fd.DeployPath
	filename := fd.DeployFileName
	text := fd.DeployText

	if len(path) == 0 {
		response.Fail(c, nil, "路径为空，请检查文件路径")
		return
	}
	if len(filename) == 0 {
		response.Fail(c, nil, "文件名为空，请检查文件名字")
		return
	}
	if len(text) == 0 {
		response.Fail(c, nil, "文件内容为空，请重新文件内容")
		return
	}
	f := func(uuid string) batchservice.R {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			return common.NodeResult{UUID: uuid,
				Error: "get agent failed"}
		}

		_, Err, err := agent.UpdateConfigFile(path, filename, text)
		if len(Err) != 0 || err != nil {
			return common.NodeResult{UUID: uuid,
				Error: Err + err.Error()}
		}
		return common.NodeResult{UUID: uuid}
	}

	result := batchservice.BatchProcess(&fd.DeployBatch, f, path, filename, text)
	response.Success(c, result, "文件下发结果")
}

func GetNodeFiles(c *gin.Context) {
	fd := &struct {
		DeployBatch common.Batch `json:"deploybatch"`
		Path        string       `json:"path"`
		FileName    string       `json:"filename"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	f := func(uuid string) batchservice.R {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			logger.Error("get agent failed, agent:%s", uuid)
			return common.NodeResult{UUID: uuid,
				Error: "get agent failed"}
		}
		// 查找文件
		cmd := "find " + fd.Path + " -type f -name \"*" + fd.FileName + "\""
		data, err := agent.RunCommand(base64.StdEncoding.EncodeToString([]byte(cmd)))
		if err != nil {
			logger.Error("run command error, agent:%s, command:%s", uuid, cmd)
			return common.NodeResult{UUID: uuid,
				Error: "run command error"}
		}
		if data.Stdout != "" && data.Stderr == "" {
			result := []common.File{}
			for _, v := range strings.Split(data.Stdout, "\n") {
				file, _, err := agent.ReadConfigFile(v)
				if err != nil {
					logger.Error("failed to read the file:%s", err.Error())
				}
				name := strings.Split(v, "/")[len(strings.Split(v, "/"))-1]
				result = append(result, common.File{Path: fd.Path,
					Name:    name,
					Content: file})
			}
			return common.NodeResult{UUID: uuid,
				Data: result}
		}
		return common.NodeResult{UUID: uuid,
			Error: "no such file or directory"}
	}

	rs := batchservice.BatchProcess(&fd.DeployBatch, f)
	response.Success(c, rs, "文件获取完成")
}
