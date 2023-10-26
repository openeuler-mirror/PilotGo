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
