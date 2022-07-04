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
 * Date: 2022-06-20 02:43:29
 * LastEditTime: 2022-06-20 16:51:51
 * Description: get agent repo info.
 ******************************************************************************/
package os

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/host"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

func GetRepoSource() (interface{}, error) {
	repos, err := GetFiles(global.RepoPath)
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
