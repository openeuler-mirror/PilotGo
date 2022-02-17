package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func AllRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_all, err := agent.AllRpm()
	if err != nil {
		response.Fail(c, nil, "获取已安装rpm包列表失败!")
		return
	}
	response.Success(c, gin.H{"rpm_all": rpm_all}, "Success")
}
func RpmSourceHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_source, err := agent.RpmSource(rpmname)
	if err != nil {
		response.Fail(c, nil, "获取源软件包名以及源失败!")
		return
	}
	response.Success(c, gin.H{"rpm_source": rpm_source}, "Success")
}
func RpmInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_info, err := agent.RpmInfo(rpmname)
	if err != nil {
		response.Fail(c, nil, "获取源软件包信息失败!")
		return
	}
	response.Success(c, gin.H{"rpm_info": rpm_info}, "Success")
}
func InstallRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_install, err := agent.InstallRpm(rpmname)
	if err != nil {
		response.Fail(c, nil, "rpm包安装命令执行失败!")
		return
	}
	response.Success(c, gin.H{"rpm_install": rpm_install}, "Success")
}
func RemoveRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_remove, err := agent.RemoveRpm(rpmname)
	if err != nil {
		response.Fail(c, nil, "rpm包卸载命令执行失败!")
		return
	}
	if rpm_remove != nil {
		response.Fail(c, nil, "软件包卸载失败!")
		return
	}
	response.Success(c, gin.H{"rpm_remove": "卸载成功!"}, "Success")
}
