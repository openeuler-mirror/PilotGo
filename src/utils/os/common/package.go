package common

import (
	"fmt"
	"regexp"
	"strings"

	"gitee.com/PilotGo/PilotGo/global"
	"gitee.com/PilotGo/PilotGo/utils"
	"github.com/shirou/gopsutil/host"
)

// 形如	openssl-1:1.1.1f-4.oe1.x86_64
//
//	OS
//	openssl=1:1.1.1f-4.oe1
type RpmSrc struct {
	Name     string
	Repo     string
	Provides string
}

type RpmInfo struct {
	Name         string
	Version      string
	Release      string
	Architecture string
	InstallDate  string
	Size         string
	License      string
	Signature    string
	Packager     string
	Vendor       string
	URL          string
	Summary      string
}

type RepoSource struct {
	Name    string
	Baseurl string
}

// TODO: yum源文件在agent端打开的情况下调用该接口匹配内容出错
func GetRepoSource() ([]RepoSource, error) {
	repos, err := utils.GetFiles(global.RepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo source file: %s", err)
	}

	SysInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get system's native repo: %s", err)
	}
	SysPlatform := SysInfo.Platform

	var repo string
	for _, repo = range repos {
		if SysPlatform == "centos" {
			reg := regexp.MustCompile(`(?i)(` + SysPlatform + "-base" + `)`)
			ok := reg.MatchString(repo)
			if ok {
				break
			}
		}
		reg := regexp.MustCompile(`(?i)(` + SysPlatform + `)`)
		ok := reg.MatchString(repo)
		if ok {
			break
		}
	}

	text, err := utils.FileReadString(global.RepoPath + "/" + repo)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo source data: %s", err)
	}

	reg1 := regexp.MustCompile(`\[.*]`)
	textType := reg1.FindAllString(text, -1)

	var reg2 *regexp.Regexp
	var BaseURL []string
	reg2 = regexp.MustCompile(`mirrorlist=http.*`)
	BaseURL = reg2.FindAllString(text, -1)
	if len(BaseURL) == 0 {
		reg2 = regexp.MustCompile(`baseurl.*`)
		BaseURL = reg2.FindAllString(text, -1)
	}

	datas := make([]RepoSource, 0)
	for i := 0; i < len(textType); i++ {
		data := RepoSource{
			Name:    textType[i][1 : len(textType[i])-1],
			Baseurl: "http" + strings.Split(BaseURL[i], "http")[1],
		}
		datas = append(datas, data)
	}

	return datas, nil
}
