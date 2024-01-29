package pluginapi

import (
	"encoding/base64"
	"strings"

	"gitee.com/openeuler/PilotGo/utils"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	batchservice "gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/app/server/service/depart"
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
		DeployBatchIds  []int    `json:"deploybatches"`
		DeployDepartIds []int    `json:"deploydeparts"`
		DeployNodeUUIds []string `json:"deploynodes"`
		DeployPath      string   `json:"deploypath"`
		DeployFileName  string   `json:"deployname"`
		DeployText      string   `json:"deployfile"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	batchIds := fd.DeployBatchIds
	UUIDs := batchservice.BatchIds2UUIDs(batchIds)
	//添加部门机器数据，获取子部门，查找机器
	for _, v := range fd.DeployDepartIds {
		machinelist, err := depart.MachineList(v)
		if err != nil {
			response.Fail(c, nil, err.Error())
		}

		for _, m := range machinelist {
			UUIDs = append(UUIDs, m.UUID)
		}
	}

	UUIDs = append(UUIDs, fd.DeployNodeUUIds...)
	//uuid数组去重
	uuids := utils.RemoveRepByMap(UUIDs)
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

	result := []string{}
	for _, uuid := range uuids {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {

			result = append(result, uuid)
			continue
		}

		_, Err, err := agent.UpdateConfigFile(path, filename, text)
		if len(Err) != 0 || err != nil {

			result = append(result, uuid)
			continue
		}
	}

	switch len(result) > 0 {
	case true:
		response.Success(c, result, "部分文件下发失败")
	case false:
		response.Success(c, nil, "文件下发完成")
	}
}

func GetNodeFiles(c *gin.Context) {
	fd := &struct {
		NodeUUIds []string `json:"nodes"`
		Path      string   `json:"path"`
		FileName  string   `json:"filename"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	type FileResult struct {
		UUID  string
		Error string
		Data  interface{}
	}
	r := []FileResult{}
	for _, uuid := range fd.NodeUUIds {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			r = append(r, FileResult{UUID: uuid,
				Error: "get agent failed"})
			logger.Error("get agent failed, agent:%s", uuid)
			continue
		}

		//查找文件
		cmd := "find " + fd.Path + " -type f -name \"*" + fd.FileName + "\""
		data, err := agent.RunCommand(base64.StdEncoding.EncodeToString([]byte(cmd)))
		if err != nil {
			r = append(r, FileResult{UUID: uuid,
				Error: "run command error"})
			logger.Error("run command error, agent:%s, command:%s", uuid, cmd)
			continue
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
			r = append(r, FileResult{UUID: uuid,
				Data: result})
		} else {
			r = append(r, FileResult{UUID: uuid,
				Error: "no such file or directory"})
		}
	}
	response.Success(c, r, "文件获取完成")
}
