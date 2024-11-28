/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func DiskUsageHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_use, err := agent.DiskUsage()
	if err != nil {
		response.Fail(c, nil, "获取磁盘的使用情况失败!")
		return
	}
	response.Success(c, gin.H{"disk_use": disk_use}, "Success")
}

func DiskInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_info, err := agent.DiskInfo()
	if err != nil {
		response.Fail(c, nil, "获取磁盘的IO信息失败!")
		return
	}
	response.Success(c, gin.H{"disk_info": disk_info}, "Success")
}
func DiskMountHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	sourceDisk := c.Query("source")
	destPath := c.Query("mountpath")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_mount, err := agent.DiskMount(sourceDisk, destPath)
	if disk_mount != "" || err != nil {
		response.Fail(c, gin.H{"error": disk_mount}, "挂载磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_mount": disk_mount}, "Success")
}
func DiskUMountHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	diskPath := c.Query("path")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_umount, err := agent.DiskUMount(diskPath)
	if disk_umount != "" || err != nil {
		response.Fail(c, gin.H{"error": disk_umount}, "卸载磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_umount": disk_umount}, "Success")
}
func DiskFormatHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	fileType := c.Query("type")
	diskPath := c.Query("path")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_format, err := agent.DiskFormat(fileType, diskPath)
	if disk_format == "" || err != nil {
		response.Fail(c, gin.H{"error": disk_format}, "格式化磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_format": disk_format}, "Success")
}
