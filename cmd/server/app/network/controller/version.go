/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

var version = ""
var commit = ""
var goVersion = ""
var buildTime = ""

func VersionHandler(c *gin.Context) {
	response.Success(c, struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		GoVersion string `json:"go_version"`
		BuildTime string `json:"build_time"`
	}{
		Version:   version,
		Commit:    commit,
		GoVersion: goVersion,
		BuildTime: buildTime,
	}, "")
}
