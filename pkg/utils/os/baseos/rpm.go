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
 * Date: 2022-01-17 02:43:29
 * LastEditTime: 2022-03-02 18:35:12
 * Description: provide agent rpm manager functions.
 ******************************************************************************/
package baseos

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
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

// 获取全部安装的rpm包列表
func (b *BaseOS) GetAllRpm() []string {
	listRpm := make([]string, 0)
	result, err := utils.RunCommand("rpm -qa")
	if err != nil && len(result) != 0 {
		logger.Error("获取已安装rpm包列表失败", err)
	}
	reader := strings.NewReader(result)
	scanner := bufio.NewScanner(reader)
	for {

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		listRpm = append(listRpm, line)
	}
	return listRpm
}

// 获取源软件包名以及源
func (b *BaseOS) GetRpmSource(rpm string) ([]RpmSrc, error) {
	Getlist := make([]RpmSrc, 0)
	listRpmSource := make([]string, 0)
	listRpmName := make([]string, 0)
	listRpmProvides := make([]string, 0)
	result, err := utils.RunCommand("yum provides " + rpm)
	if err != nil && len(result) != 0 {
		logger.Error("获取源软件包名以及源失败", err)
		return []RpmSrc{}, fmt.Errorf("获取源软件包名以及源失败")
	}
	reader := strings.NewReader(result)
	scanner := bufio.NewScanner(reader)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg2 := regexp.MustCompile(`^[R].*`)
		x := reg2.FindAllString(line, -1)
		if x == nil {
			continue
		}
		str2 := strings.Fields(x[0])
		listRpmSource = append(listRpmSource, str2[2])
	}
	reader = strings.NewReader(result)
	scanner = bufio.NewScanner(reader)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg1 := regexp.MustCompile(`^.*[.]+.*:`)
		x := reg1.FindAllString(line, -1)
		if x == nil {
			continue
		}
		str1 := strings.Fields(x[0])
		listRpmName = append(listRpmName, str1[0])
	}
	reader = strings.NewReader(result)
	scanner = bufio.NewScanner(reader)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg := regexp.MustCompile(`Provide.*:.*`)
		x := reg.FindAllString(line, -1)
		if x == nil {
			continue
		}
		str3 := strings.Fields(x[0])
		str := str3[2] + str3[3] + str3[4]
		listRpmProvides = append(listRpmProvides, str)
	}
	for key, value := range listRpmSource {
		tmp := RpmSrc{}
		tmp.Name = listRpmName[key]
		tmp.Provides = listRpmProvides[key]
		tmp.Repo = value
		Getlist = append(Getlist, tmp)
	}
	return Getlist, nil
}

// 按行使用正则语言查找结构体的属性信息
func readInfo(reader *strings.Reader, reg string) (string, error) {
	scanner := bufio.NewScanner(reader)
	var result string
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg := regexp.MustCompile(reg)
		x := reg.FindAllString(line, -1)
		if x == nil {
			continue
		}
		str := strings.Fields(x[0])
		length := len(str)
		if length < 3 {
			continue
		} else if length == 3 {
			result = str[2]
			return result, nil
		} else {
			i := 3
			result = str[2]
			for {
				if i == length {
					break
				}
				result = result + " " + str[i]
				i += 1

			}
			return result, nil
		}
	}
	return string(""), fmt.Errorf("匹配结构体属性失败")
}

func (b *BaseOS) GetRpmInfo(rpm string) (RpmInfo, error, error) {
	rpminfo := RpmInfo{}
	result, err := utils.RunCommand("rpm -qi " + rpm)
	//未安装该软件包情况
	if err != nil && len(result) != 0 {
		logger.Error(" %s的rpm包未安装", rpm)
		return RpmInfo{}, fmt.Errorf("%s的rpm包未安装", rpm), err
	}
	reader := strings.NewReader(result)
	str, err := readInfo(reader, `^Name.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包名属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包名属性失败"), err
	}
	rpminfo.Name = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Version.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Version属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Version属性失败"), err
	}
	rpminfo.Version = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Release.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Release属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Release属性失败"), err
	}
	rpminfo.Release = str
	// reader = strings.NewReader(result)
	// str, err = readInfo(reader, `^Architecture.*`)
	// if err != nil && len(str) != 0 {
	// 	logger.Error("读取rpm包Architecture属性失败")
	// 	return RpmInfo{}, fmt.Errorf("读取rpm包Architecture属性失败"), err
	// }
	rpminfo.Architecture = strings.Split(rpm, ".")[len(strings.Split(rpm, "."))-1]
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Install Date.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包InstallDate属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包InstallDate属性失败"), err
	}
	rpminfo.InstallDate = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Size.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Size属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Size属性失败"), err
	}
	rpminfo.Size = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^License.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包License属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包License属性失败"), err
	}
	rpminfo.License = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Signature.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Signature属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Signature属性失败"), err
	}
	rpminfo.Signature = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Packager.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Packager属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Packager属性失败"), err
	}
	rpminfo.Packager = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Vendor.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Vendor属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包Vendor属性失败"), err
	}
	rpminfo.Vendor = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^URL.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包URL属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包URL属性失败"), err
	}
	rpminfo.URL = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Summary.*`)
	if err != nil && len(str) != 0 {
		logger.Error("读取rpm包Summary属性失败")
		return RpmInfo{}, fmt.Errorf("读取rpm包URL属性失败"), err
	}
	rpminfo.Summary = str
	return rpminfo, nil, nil
}

// 判断rpm软件包是否安装/卸载成功
func verifyRpmInstalled(reader *strings.Reader, reg string) bool {
	scanner := bufio.NewScanner(reader)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg2 := regexp.MustCompile(reg)
		x := reg2.FindAllString(line, -1)
		if x != nil {
			return true
		}
	}
	return false
}

// 安装rpm软件包
func (b *BaseOS) InstallRpm(rpm string) error {
	result, err := utils.RunCommand("yum -y install " + rpm)
	if err != nil {
		logger.Error("rpm包安装命令运行失败: ", err)
		return fmt.Errorf("rpm包安装命令执行失败")
	}
	if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
		logger.Error("rpm包安装命令由于rpm包已安装而运行失败")
		return fmt.Errorf("该rpm包已安装")
	} else if verifyRpmInstalled(strings.NewReader(result), `^Error: Unable to find a match:.*`) {
		logger.Error("rpm包安装命令由于源内匹配不到该rpm包而运行失败")
		return fmt.Errorf("源内匹配不到该rpm包")
	} else {
		logger.Info("%s安装成功", rpm)
		return nil
	}

}

// 卸载rpm软件包
func (b *BaseOS) RemoveRpm(rpm string) error {
	result, err := utils.RunCommand("yum -y remove " + rpm)
	if err != nil {
		logger.Error("rpm包卸载命令运行失败: ", err)
		return fmt.Errorf("rpm包卸载命令执行失败")
	}
	if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
		logger.Error("rpm包卸载命令由于rpm包不存在而运行失败")
		return fmt.Errorf("该rpm包不存在")
	} else {
		logger.Info("%s卸载成功", rpm)
		return nil
	}
}
