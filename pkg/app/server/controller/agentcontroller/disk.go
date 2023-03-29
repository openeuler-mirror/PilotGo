/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-02-16 15:13:25
 * LastEditTime: 2022-02-24 15:51:43
 * Description: provide agent disk manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
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
	if disk_mount != nil || err != nil {
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
	if disk_umount != nil || err != nil {
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
	if disk_format != nil || err != nil {
		response.Fail(c, gin.H{"error": disk_format}, "格式化磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_format": disk_format}, "Success")
}
