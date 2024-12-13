/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"regexp"

	scriptservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/script"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 存储脚本文件
func AddScriptHandler(c *gin.Context) {
	var script scriptservice.Script
	err := scriptservice.AddScript(&script)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "脚本文件添加失败")
		return
	}
	response.Success(c, nil, "脚本文件添加成功")
}

// 高危命令检测
func FindDangerousCommandsPos(content string) ([][]int, []string) {
	var positions [][]int
	var matchedCommands []string

	for _, pattern := range DangerousCommandsList {
		re, err := regexp.Compile(pattern)
		if err != nil {
			// TOODO: info remind
			continue
		}
		matches := re.FindAllStringIndex(content, -1)
		for _, match := range matches {
			start, end := match[0], match[1]-1
			positions = append(positions, []int{start, end})
			matchedCommands = append(matchedCommands, content[start:end+1])
		}
	}
	return positions, matchedCommands
}

var DangerousCommandsList = []string{
	`.*rm\s+-[r,f,rf].*`,
	`.*lvremove\s+-f.*`,
	`.*poweroff.*`,
	`.*shutdown\s+-[f,F,h,k,n,r,t,C].*`,
	`.*pvremove\s+-f.*`,
	`.*vgremove\s+-f.*`,
	`.*exportfs\s+-[a,u].*`,
	`.*umount.nfs+.*.+-[r,f,rf].*`,
	`.*mv+.*.+/dev/null.*`,
	`.*reboot.*`,
	`.*rmmod\s+-[a,s,v,f,w].*`,
	`.*dpkg-divert+.*.+-remove.*`,
	`.*dd.*`,
	`.*mkfs.*`,
	`.*vmo.*`,
	`.*init.*`,
	`.*halt.*`,
	`.*fasthalt.*`,
	`.*fastboot.*`,
	`.*startsrc.*`,
	`.*stopsrc.*`,
	`.*chkconfig.*`,
	`.*off.*`,
	`.*refresh.*`,
	`.*umount.*`,
	`.*rmdev.*`,
	`.*chdev.*`,
	`.*extendvg.*`,
	`.*reducevg.*`,
	`.*importvg.*`,
	`.*exportvg.*`,
	`.*mklv.*`,
	`.*rmlv.*`,
	`.*rmfs.*`,
	`.*chfs.*`,
	`.*installp.*`,
	`.*instfix.*`,
	`.*crontab.*`,
	`.*cfgmgr.*`,
	`.*mknod.*`,
}
