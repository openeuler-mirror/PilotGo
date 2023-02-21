package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/host"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
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

func GetRepoSource() (interface{}, error) {
	repos, err := utils.GetFiles(global.RepoPath)
	if err != nil {
		return "", fmt.Errorf("获取repo源文件失败:%s", err)
	}

	SysInfo, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("获取系统的原生repo失败:%s", err)
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
		return "", fmt.Errorf("读取repo源数据失败:%s", err)
	}

	reg1 := regexp.MustCompile(`\[.*]`)
	textType := reg1.FindAllString(text, -1)

	var reg2 *regexp.Regexp
	var BaseURL []string
	reg2 = regexp.MustCompile(`mirrorlist=http.*`)
	BaseURL = reg2.FindAllString(text, -1)
	if len(BaseURL) == 0 {
		reg2 = regexp.MustCompile(`baseurl=.*`)
		BaseURL = reg2.FindAllString(text, -1)
	}

	datas := make([]map[string]string, 0)
	for i := 0; i < len(textType); i++ {
		data := map[string]string{
			"name":    textType[i][1 : len(textType[i])-1],
			"baseurl": "http" + strings.Split(BaseURL[i], "http")[1],
		}
		datas = append(datas, data)
	}

	return datas, nil
}
