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
 * LastEditTime: 2023-02-21 18:24:37
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
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取全部安装的rpm包列表
func (b *BaseOS) GetAllRpm() []string {
	listRpm := make([]string, 0)
	result, err := utils.RunCommand("rpm -qa")
	if err != nil && len(result) != 0 {
		logger.Error("failed to get list of installed RPM packages: %s", err)
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
func (b *BaseOS) GetRpmSource(rpm string) ([]common.RpmSrc, error) {
	Getlist := make([]common.RpmSrc, 0)
	listRpmSource := make([]string, 0)
	listRpmName := make([]string, 0)
	listRpmProvides := make([]string, 0)
	result, err := utils.RunCommand("yum provides " + rpm)
	if err != nil && len(result) != 0 {
		logger.Error("failed to get source software package name and source: %s", err)
		return []common.RpmSrc{}, fmt.Errorf("failed to get source software package name and source")
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
		tmp := common.RpmSrc{}
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
	return string(""), fmt.Errorf("failed to match struct properties")
}

func (b *BaseOS) GetRpmInfo(rpm string) (common.RpmInfo, error) {
	rpminfo := common.RpmInfo{}
	result, err := utils.RunCommand("rpm -qi " + rpm)
	//未安装该软件包情况
	if err != nil && len(result) != 0 {
		logger.Error(" %s's RPM package not installed", rpm)
		return common.RpmInfo{}, fmt.Errorf("%s's RPM package not installed: %s", rpm, err)
	}
	reader := strings.NewReader(result)
	str, err := readInfo(reader, `^Name.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package name properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package name properties: %s", err)
	}
	rpminfo.Name = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Version.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Version properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Version properties: %s", err)
	}
	rpminfo.Version = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Release.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Release properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Release properties: %s", err)
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
		logger.Error("failed to read RPM package InstallDate properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package InstallDate properties: %s", err)
	}
	rpminfo.InstallDate = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Size.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Size properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Size properties: %s", err)
	}
	rpminfo.Size = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^License.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package License properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package License properties: %s", err)
	}
	rpminfo.License = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Signature.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Signature properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Signature properties: %s", err)
	}
	rpminfo.Signature = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Packager.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Packager properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Packager properties: %s", err)
	}
	rpminfo.Packager = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Vendor.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Vendor properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Vendor properties: %s", err)
	}
	rpminfo.Vendor = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^URL.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package URL properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package URL properties: %s", err)
	}
	rpminfo.URL = str
	reader = strings.NewReader(result)
	str, err = readInfo(reader, `^Summary.*`)
	if err != nil && len(str) != 0 {
		logger.Error("failed to read RPM package Summary properties")
		return common.RpmInfo{}, fmt.Errorf("failed to read RPM package Summary properties:%s", err)
	}
	rpminfo.Summary = str
	return rpminfo, nil
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
		logger.Error("failed to run RPM package installation command: %s", err)
		return fmt.Errorf("failed to run RPM package installation command")
	}
	if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
		logger.Error("failed to run RPM package installation command due to package already being installed")
		return fmt.Errorf("this RPM package is already installed")
	} else if verifyRpmInstalled(strings.NewReader(result), `^Error: Unable to find a match:.*`) {
		logger.Error("failed to run RPM package installation command due to the package not being found in the source")
		return fmt.Errorf("the package cannot be found in the source")
	} else {
		logger.Info("successfully installed %s", rpm)
		return nil
	}

}

// 卸载rpm软件包
func (b *BaseOS) RemoveRpm(rpm string) error {
	result, err := utils.RunCommand("yum -y remove " + rpm)
	if err != nil {
		logger.Error("failed to run RPM package uninstallation command: %s", err)
		return fmt.Errorf("failed to execute RPM package uninstallation command")
	}
	if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
		logger.Error("failed to run RPM package uninstallation command due to the package not being found")
		return fmt.Errorf("this RPM package does not exist")
	} else {
		logger.Info("successfully uninstalled %s", rpm)
		return nil
	}
}
